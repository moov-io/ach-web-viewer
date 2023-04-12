package filelist

import (
	"fmt"
	"time"

	"github.com/moov-io/ach"
	"github.com/moov-io/ach-web-viewer/pkg/service"
)

type Files struct {
	SourceID   string
	SourceType string
	Files      []File
}

type File struct {
	Name        string
	StoragePath string
	Contents    *ach.File

	CreatedAt   time.Time
	Size        int64
	RecordCount int64
}

type Lister interface {
	SourceID() string

	GetFile(path string) (*File, error)
	GetFiles(opts ListOpts) (Files, error)
}

type Listers []Lister

func (ls Listers) GetFiles(opts ListOpts) (map[string]Files, error) {
	out := make(map[string]Files)
	for i := range ls {
		files, err := ls[i].GetFiles(opts)
		if err != nil {
			return out, err
		}
		out[ls[i].SourceID()] = files
	}
	return out, nil
}

func NewListers(ss service.Sources) (Listers, error) {
	var out Listers
	for i := range ss {
		ls, err := createLister(ss[i])
		if err != nil {
			return out, err
		}
		out = append(out, ls)
	}
	return out, nil
}

func createLister(src service.Source) (Lister, error) {
	switch {
	case src.ACHGateway != nil:
		return newACHGatewayLister(src.ID, *src.ACHGateway)

	case src.Bucket != nil:
		return newBucketLister(src.ID, src)

	case src.Filesystem != nil:
		return newFilesystemLister(src.ID, src.Filesystem)
	}
	return nil, fmt.Errorf("unknown source: %#v", src)
}

func (ls Listers) GetFile(sourceID, path string) (*File, error) {
	for i := range ls {
		if ls[i].SourceID() == sourceID {
			return ls[i].GetFile(path)
		}
	}
	return nil, fmt.Errorf("%s not found for sourceID=%s", path, sourceID)
}
