package filelist

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/moov-io/ach"
	"github.com/moov-io/ach-web-viewer/internal/gpgx"
	"github.com/moov-io/ach-web-viewer/pkg/service"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob"
	"golang.org/x/crypto/openpgp"
)

type bucketLister struct {
	sourceID string
	buck     *blob.Bucket
	paths    []string
	gpgKey   openpgp.EntityList
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
		ls.gpgKey, err = gpgx.ReadPrivateKeyFile(conf.KeyFile, []byte(conf.KeyPassword))
		if err != nil {
			return nil, fmt.Errorf("problem reading %s: %v", conf.KeyFile, err)
		}
	}
	return ls, nil
}

func (ls *bucketLister) SourceID() string {
	return ls.sourceID
}

func (ls *bucketLister) GetFiles() (Files, error) {
	out := Files{
		SourceID:   ls.sourceID,
		SourceType: "Bucket",
	}
	for i := range ls.paths {
		files, err := ls.listFiles(ls.buck.List(&blob.ListOptions{
			Prefix: ls.paths[i],
		}))
		if err != nil {
			return out, fmt.Errorf("error reading %s bucket path: %v", ls.paths[i], err)
		}
		out.Files = append(out.Files, files...)
	}
	return out, nil
}

func (ls *bucketLister) GetFile(path string) (*ach.File, error) {
	rdr, err := ls.buck.NewReader(context.Background(), path, nil)
	if err != nil {
		return nil, err
	}
	bs, err := maybeDecrypt(rdr, ls.gpgKey)
	if err != nil {
		rdr.Close()
		return nil, err
	}
	rdr.Close()
	file, err := ach.NewReader(bs).Read()
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (ls *bucketLister) listFiles(cur *blob.ListIterator) ([]File, error) {
	var out []File
	for {
		obj, err := cur.Next(context.Background())
		if err != nil {
			if err == io.EOF {
				break // finished with cursor
			}
			return nil, err
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
