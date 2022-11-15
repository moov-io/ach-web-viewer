package web

import (
	"io"
	"strings"

	"github.com/moov-io/ach"
	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/moov-io/ach/cmd/achcli/describe"
)

func File(w io.Writer, cfg *service.DisplayConfig, file *ach.File) error {
	switch strings.ToLower(cfg.Format) {
	case "human-readable":
		return achcliDescribe(w, cfg.Masking, file)
	}
	return nil
}

func achcliDescribe(w io.Writer, cfg service.MaskingConfig, file *ach.File) error {
	describe.File(w, file, &describe.Opts{
		MaskAccountNumbers: cfg.AccountNumbers,
		MaskCorrectedData:  cfg.CorrectedData,
		MaskNames:          cfg.Names,
		PrettyAmounts:      cfg.PrettyAmounts,
	})
	return nil
}
