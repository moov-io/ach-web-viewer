package filelist

import (
	"fmt"

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
	File        *ach.File
}

type Lister interface {
	SourceID() string

	GetFile(path string) (*ach.File, error)
	GetFiles() (Files, error)
}

type Listers []Lister

func (ls Listers) GetFiles() (map[string]Files, error) {
	out := make(map[string]Files)
	for i := range ls {
		files, err := ls[i].GetFiles()
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
	case src.Bucket != nil:
		return newBucketLister(src.ID, src)

	case src.Filesystem != nil:
		return newFilesystemLister(src.ID, src.Filesystem)
	}
	return nil, fmt.Errorf("unknown source: %#v", src)
}

func (ls Listers) GetFile(sourceID, path string) (*ach.File, error) {
	for i := range ls {
		if ls[i].SourceID() == sourceID {
			return ls[i].GetFile(path)
		}
	}
	return nil, fmt.Errorf("%s not found for sourceID=%s", path, sourceID)
}
