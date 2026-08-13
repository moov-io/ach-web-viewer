package filelist

import (
	"fmt"
	"io"
	"strings"

	"github.com/moov-io/ach"
)

func init() {
	RegisterFormat(Format{
		Kind:       "ach",
		Title:      "ACH File",
		Extensions: []string{".ach", ".txt"},
		Parse:      parseACH,
	})
	RegisterFormat(Format{
		Kind:       "ach-json",
		Title:      "ACH File",
		Extensions: []string{".json"},
		Parse:      parseACHJSON,
	})
}

// ACHFile is a parsed NACHA ACH file.
type ACHFile struct {
	File *ach.File
}

func (a *ACHFile) Kind() string  { return "ach" }
func (a *ACHFile) Title() string { return "ACH File" }

func (a *ACHFile) Validate() error {
	if a == nil || a.File == nil {
		return errMissingFile
	}
	return a.File.Validate()
}

func (a *ACHFile) RecordCount() int64 { return 0 }

func (a *ACHFile) Metadata() map[string]string { return nil }

func parseACH(_ string, r io.Reader) (Document, error) {
	file, err := readACHFile(r)
	if file == nil {
		return nil, err
	}
	return &ACHFile{File: file}, err
}

func parseACHJSON(_ string, r io.Reader) (Document, error) {
	file, err := readJSONFile(r)
	if file == nil {
		return nil, err
	}
	return &ACHFile{File: file}, err
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
