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
			strings.Contains(message, "*errors.errorString"),
			strings.Contains(message, "none or more than one"):
			return &file, nil
		}
	}
	return &file, err
}
