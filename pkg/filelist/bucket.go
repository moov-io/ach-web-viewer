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
	"github.com/moov-io/base/telemetry"
	"github.com/moov-io/cryptfs"

	"cloud.google.com/go/storage"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

func (ls *bucketLister) GetFiles(ctx context.Context, opts ListOpts) (Files, error) {
	out := Files{
		SourceID:   ls.sourceID,
		SourceType: "Bucket",
	}

	g, ctx := errgroup.WithContext(ctx)
	fileChan := make(chan []File)

	for _, path := range ls.paths {
		path := path // capture range variable
		g.Go(func() error {
			files, err := ls.listFiles(ctx, opts, path)
			if err != nil {
				return fmt.Errorf("error reading %s bucket path: %v", path, err)
			}
			select {
			case fileChan <- files:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}

	// Start a goroutine to collect files
	go func() {
		for files := range fileChan {
			out.Files = append(out.Files, files...)
		}
	}()

	// Wait for all goroutines to complete or first error
	if err := g.Wait(); err != nil {
		return out, err
	}

	// Close channel after all goroutines are done
	close(fileChan)

	return out, nil
}

func (ls *bucketLister) GetFile(ctx context.Context, path string) (*File, error) {
	ctx, span := telemetry.StartSpan(ctx, "filelist-bucket-getfile", trace.WithAttributes(
		attribute.String("search.path", path),
	))
	defer span.End()

	rdr, err := ls.buck.NewReader(ctx, path, nil)
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

	file, err := readFile(name, bytes.NewReader(bs))

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

func (ls *bucketLister) listFiles(ctx context.Context, opts ListOpts, pathPrefix string) ([]File, error) {
	// Different underlying storage engines will let us scan/glob parts of the bucket differently.
	var gcsBucket *storage.Client
	if ls.buck.As(&gcsBucket) {
		return ls.listFilesFromGCSBucket(ctx, opts, pathPrefix)
	}
	return ls.listFilesFromCDKBucket(ctx, opts, pathPrefix)
}

func (ls *bucketLister) listFilesFromGCSBucket(ctx context.Context, opts ListOpts, pathPrefix string) ([]File, error) {
	ctx, span := telemetry.StartSpan(ctx, "filelist-bucket-list-files-gcs-bucket", trace.WithAttributes(
		attribute.String("search.path_prefix", pathPrefix),
	))
	defer span.End()

	var g errgroup.Group
	datePrefixes := yyyymmdd.Prefixes(opts.StartDate, opts.EndDate)
	span.SetAttributes(
		attribute.StringSlice("search.yymmdd_prefixes", datePrefixes),
	)

	discoveredFiles := make(chan []File, len(datePrefixes))

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

			files, err := ls.listFilesFromCursor(ctx, opts, ls.buck.List(listOptions))
			if len(files) > 0 {
				discoveredFiles <- files
			}
			return err
		})
	}

	err := g.Wait()
	if err != nil {
		return nil, err
	}
	close(discoveredFiles)

	var out []File
	for files := range discoveredFiles {
		out = append(out, files...)
	}
	return out, nil

}

func (ls *bucketLister) listFilesFromCDKBucket(ctx context.Context, opts ListOpts, pathPrefix string) ([]File, error) {
	ctx, span := telemetry.StartSpan(ctx, "filelist-bucket-list-files-cdk-bucket", trace.WithAttributes(
		attribute.String("search.path_prefix", pathPrefix),
	))
	defer span.End()

	return ls.listFilesFromCursor(ctx, opts, ls.buck.List(&blob.ListOptions{
		Delimiter: "/",
		Prefix:    pathPrefix,
	}))
}

func (ls *bucketLister) listFilesFromCursor(ctx context.Context, opts ListOpts, cur *blob.ListIterator) ([]File, error) {
	var out []File
	for {
		obj, err := cur.Next(ctx)
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
