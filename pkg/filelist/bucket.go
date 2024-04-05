package filelist

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/moov-io/ach-web-viewer/pkg/yyyymmdd"
	"github.com/moov-io/cryptfs"

	"cloud.google.com/go/storage"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob"
	"golang.org/x/sync/errgroup"
)

type bucketLister struct {
	sourceID string
	buck     *blob.Bucket
	paths    []string
	cryptors []*cryptfs.FS
}

func newBucketLister(sourceID string, cfg service.Source) (*bucketLister, error) {
	if cfg.Bucket == nil {
		return nil, errors.New("missing BucketConfig")
	}
	buck, err := blob.OpenBucket(context.Background(), cfg.Bucket.URL)
	if err != nil {
		return nil, err
	}
	ls := &bucketLister{
		sourceID: sourceID,
		buck:     buck,
		paths:    cfg.Bucket.Paths,
	}
	if cfg.Encryption != nil && cfg.Encryption.GPG != nil {
		conf := cfg.Encryption.GPG
		// Decrypt with either config
		if conf.KeyFile != "" {
			cc, err := cryptfs.FromCryptor(cryptfs.NewGPGDecryptorFile(conf.KeyFile, []byte(conf.KeyPassword)))
			if err != nil {
				return nil, fmt.Errorf("problem reading %s: %v", conf.KeyFile, err)
			}
			if cc != nil {
				ls.cryptors = append(ls.cryptors, cc)
			}
		}
		// Read from .Files and add those
		for i := range conf.Files {
			cc, err := cryptfs.FromCryptor(cryptfs.NewGPGDecryptorFile(conf.Files[i].KeyFile, []byte(conf.Files[i].KeyPassword)))
			if err != nil {
				return nil, fmt.Errorf("problem reading %s: %v", conf.Files[i].KeyFile, err)
			}
			if cc != nil {
				ls.cryptors = append(ls.cryptors, cc)
			}
		}
	}
	return ls, nil
}

func (ls *bucketLister) SourceID() string {
	return ls.sourceID
}

func (ls *bucketLister) GetFiles(opts ListOpts) (Files, error) {
	out := Files{
		SourceID:   ls.sourceID,
		SourceType: "Bucket",
	}
	for i := range ls.paths {
		files, err := ls.listFiles(opts, ls.paths[i])
		if err != nil {
			return out, fmt.Errorf("error reading %s bucket path: %v", ls.paths[i], err)
		}
		out.Files = append(out.Files, files...)
	}
	return out, nil
}

func (ls *bucketLister) GetFile(path string) (*File, error) {
	rdr, err := ls.buck.NewReader(context.Background(), path, nil)
	if err != nil {
		return nil, err
	}

	bs, err := ls.maybeDecrypt(rdr)
	if err != nil {
		rdr.Close()
		return nil, err
	}
	rdr.Close()

	_, name := filepath.Split(path)

	file, err := readFile(bytes.NewReader(bs))

	return &File{
		Name:        name,
		StoragePath: path,
		Contents:    file,
		CreatedAt:   rdr.ModTime(),
		Size:        int64(len(bs)),
	}, err
}

func (ls *bucketLister) maybeDecrypt(r io.Reader) ([]byte, error) {
	initial, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	for i := range ls.cryptors {
		bs, err := ls.cryptors[i].Reveal(initial)
		if len(bs) > 0 && err == nil {
			return bs, err
		}
	}
	return initial, err
}

func (ls *bucketLister) listFiles(opts ListOpts, pathPrefix string) ([]File, error) {
	// Different underlying storage engines will let us scan/glob parts of the bucket differently.
	var gcsBucket *storage.Client
	if ls.buck.As(&gcsBucket) {
		return ls.listFilesFromGCSBucket(opts, pathPrefix)
	}
	return ls.listFilesFromCDKBucket(opts, pathPrefix)
}

func (ls *bucketLister) listFilesFromGCSBucket(opts ListOpts, pathPrefix string) ([]File, error) {
	var g errgroup.Group
	datePrefixes := yyyymmdd.Prefixes(opts.StartDate, opts.EndDate)

	discoveredFiles := make(chan []File)

	for _, datePrefix := range datePrefixes {
		g.Go(func() error {
			beforeList := func(as func(interface{}) bool) error {
				var q *storage.Query
				if as(&q) {
					// If pathPrefix contains a "/" then it's formatted as "folder/hostname"
					folder, hostname := pathPrefix, "*"

					idx := strings.Index(pathPrefix, "/")
					if idx > 0 {
						folder, hostname = pathPrefix[:idx], pathPrefix[idx:]

						hostname = strings.ReplaceAll(hostname, "/", "")
					}

					// achgateway stores files under two path schemes:
					//
					//   /odfi     /sftp.bank.com:22 /Returned   /2024-04-05
					//   /outbound /sftp.bank.com:22 /2022-10-17
					//
					// Glob for the optional "Returned" directory
					q.MatchGlob = fmt.Sprintf("%s/%s/**/%s*/**", folder, hostname, datePrefix)
				}
				return nil
			}

			listOptions := &blob.ListOptions{
				Prefix:     pathPrefix, //  + "/",
				BeforeList: beforeList,
			}

			files, err := ls.listFilesFromCursor(opts, ls.buck.List(listOptions))
			if len(files) > 0 {
				go func() {
					discoveredFiles <- files
				}()
			}
			return err
		})
	}

	err := g.Wait()
	go func() {
		discoveredFiles <- nil
	}()
	if err != nil {
		return nil, err
	}

	var out []File
	for {
		files := <-discoveredFiles
		if len(files) == 0 {
			break
		}
		out = append(out, files...)
	}
	return out, nil

}

func (ls *bucketLister) listFilesFromCDKBucket(opts ListOpts, pathPrefix string) ([]File, error) {
	return ls.listFilesFromCursor(opts, ls.buck.List(&blob.ListOptions{
		Delimiter: "/",
		Prefix:    pathPrefix,
	}))
}

func (ls *bucketLister) listFilesFromCursor(opts ListOpts, cur *blob.ListIterator) ([]File, error) {
	var out []File
	for {
		obj, err := cur.Next(context.Background())
		if err != nil {
			if err == io.EOF {
				break // finished with cursor
			}
			return nil, err
		}

		// Skip this file if it's outside of our query params
		if !opts.Inside(obj.ModTime) {
			continue
		}

		dir, name := filepath.Split(obj.Key)
		out = append(out, File{
			Name:        name,
			StoragePath: dir,
			CreatedAt:   obj.ModTime,
		})
	}
	return out, nil
}
