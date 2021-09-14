package filelist

import (
	"io"
	"strings"

	"github.com/moov-io/ach"
)

func readFile(r io.Reader) (*ach.File, error) {
	file, err := ach.NewReader(r).Read()
	if err != nil {
		message := err.Error()
		switch {
		case strings.Contains(message, "*ach.BatchError"),
			strings.Contains(message, "*ach.FieldError"),
			strings.Contains(message, "*errors.errorString"):
			return &file, nil
		}
	}
	return &file, err
}
