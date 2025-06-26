package filelist

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/moov-io/ach"
)

func readFile(name string, r io.Reader) (*ach.File, error) {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".ach", ".txt":
		return readACHFile(r)

	case ".json":
		return readJSONFile(r)
	}
	return nil, fmt.Errorf("unknown file extension for %s", name)
}

func readACHFile(r io.Reader) (*ach.File, error) {
	file, err := ach.NewReader(r).Read()
	if err != nil {
		if isSkippableError(err) {
			return &file, nil
		}
	}
	return &file, err
}

func readJSONFile(r io.Reader) (*ach.File, error) {
	bs, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("draining reader: %w", err)
	}

	file, err := ach.FileFromJSON(bs)
	if err != nil {
		if isSkippableError(err) {
			return file, nil
		}
	}
	return file, err
}

func isSkippableError(err error) bool {
	message := err.Error()

	return strings.Contains(message, "*ach.BatchError") ||
		strings.Contains(message, "*ach.FieldError") ||
		strings.Contains(message, "*errors.errorString") ||
		strings.Contains(message, "none or more than one")
}
