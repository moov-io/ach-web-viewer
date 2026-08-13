package web

import (
	"io"
	"strings"

	"github.com/moov-io/ach"
	"github.com/moov-io/ach-web-viewer/pkg/filelist"
	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/moov-io/ach/cmd/achcli/describe"
	"github.com/moov-io/ach/cmd/achcli/describe/mask"
)

func Contents(w io.Writer, cfg *service.DisplayConfig, file *filelist.File) error {
	if file == nil || file.Document == nil {
		return nil
	}
	if format(cfg) != "human-readable" {
		return nil
	}

	switch doc := file.Document.(type) {
	case *filelist.ACHFile:
		return File(w, cfg, doc.File)
	case filelist.HumanWriter:
		return doc.WriteHuman(w)
	}
	return nil
}

func File(w io.Writer, cfg *service.DisplayConfig, file *ach.File) error {
	if format(cfg) == "human-readable" {
		return achcliDescribe(w, cfg.Masking, file)
	}
	return nil
}

func format(cfg *service.DisplayConfig) string {
	if cfg == nil || cfg.Format == "" {
		return "human-readable"
	}
	return strings.ToLower(cfg.Format)
}

func achcliDescribe(w io.Writer, cfg service.MaskingConfig, file *ach.File) error {
	describe.File(w, file, &describe.Opts{
		Options: mask.Options{
			MaskAccountNumbers: cfg.AccountNumbers,
			MaskCorrectedData:  cfg.CorrectedData,
			MaskNames:          cfg.Names,
		},
		PrettyAmounts: cfg.PrettyAmounts,
	})
	return nil
}
