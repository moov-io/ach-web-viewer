package filelist

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/moov-io/ach-web-viewer/pkg/service"
)

type filesystemLister struct {
	sourceID string
	dirs     []string
}

func newFilesystemLister(sourceID string, cfg *service.FilesystemConfig) (*filesystemLister, error) {
	if cfg == nil {
		return nil, errors.New("missing FilesystemConfig")
	}
	return &filesystemLister{
		sourceID: sourceID,
		dirs:     cfg.Paths,
	}, nil
}

func (ls *filesystemLister) SourceID() string {
	return ls.sourceID
}

func (ls *filesystemLister) GetFiles(opts ListOpts) (Files, error) {
	out := Files{
		SourceID:   ls.sourceID,
		SourceType: "Filesystem",
	}
	for i := range ls.dirs {
		err := filepath.Walk(ls.dirs[i], func(path string, info fs.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			// Skip this file if it's outside of our query params
			if !opts.Inside(info.ModTime()) {
				return nil
			}

			dir, _ := filepath.Split(path)
			out.Files = append(out.Files, File{
				Name:        filepath.Base(path),
				StoragePath: dir,
				CreatedAt:   info.ModTime(),
			})
			return nil
		})
		if err != nil {
			return out, fmt.Errorf("error reading %s: %v", ls.dirs[i], err)
		}
	}
	return out, nil
}

func (ls *filesystemLister) GetFile(path string, cfg service.DisplayConfig) (*File, error) {
	path = filepath.Clean(path)

	if strings.Contains(path, "..") || strings.HasPrefix(path, "/") {
		return nil, errors.New("invalid path")
	}

	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s: %v", path, err)
	}

	_, name := filepath.Split(fd.Name())

	file, err := readFile(fd, cfg)

	var stat fs.FileInfo
	if fd != nil {
		stat, _ = fd.Stat()
	}

	return &File{
		Name:        name,
		StoragePath: fd.Name(),
		Contents:    file,
		CreatedAt:   stat.ModTime(),
		Size:        stat.Size(),
	}, err
}
