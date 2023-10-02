package filelist

import (
	"io"
	"strings"

	"github.com/moov-io/ach"
	"github.com/moov-io/ach-web-viewer/pkg/service"
)

func readFile(r io.Reader, cfg service.DisplayConfig) (*ach.File, error) {
	reader := ach.NewReader(r)
	reader.SetValidation(&ach.ValidateOpts{
		AllowMissingBatchHeader: cfg.AllowMissingBatchHeader,
	})
	file, err := reader.Read()
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
