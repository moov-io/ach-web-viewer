package filelist

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/moov-io/ach"
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
	ls := &filesystemLister{
		sourceID: sourceID,
	}
	for i := range cfg.Paths {
		if fd, err := os.Stat(cfg.Paths[i]); err == nil && fd.IsDir() {
			ls.dirs = append(ls.dirs, cfg.Paths[i])
		} else {
			return nil, fmt.Errorf("error reading %s: %v", cfg.Paths[i], err)
		}
	}
	return ls, nil
}

func (ls *filesystemLister) SourceID() string {
	return ls.sourceID
}

func (ls *filesystemLister) GetFiles() (Files, error) {
	out := Files{
		SourceID:   ls.sourceID,
		SourceType: "Filesystem",
	}
	for i := range ls.dirs {
		err := filepath.Walk(ls.dirs[i], func(path string, info fs.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			file, err := ls.GetFile(path)
			if err != nil {
				return fmt.Errorf("error reading %s: %v", path, err)
			}
			dir, _ := filepath.Split(path)
			out.Files = append(out.Files, File{
				Name:        filepath.Base(path),
				StoragePath: dir,
				File:        file,
			})
			return nil
		})
		if err != nil {
			return out, fmt.Errorf("error reading %s: %v", ls.dirs[i], err)
		}
	}
	return out, nil
}

func (ls *filesystemLister) GetFile(path string) (*ach.File, error) {
	return ach.ReadFile(path)
}
