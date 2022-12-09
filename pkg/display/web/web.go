package web

import (
	"io"
	"strings"

	"github.com/moov-io/ach"
	"github.com/moov-io/ach-web-viewer/pkg/display"
	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/moov-io/ach/cmd/achcli/describe"
)

func File(w io.Writer, format string, cfg service.MaskingConfig, file *ach.File) error {
	switch strings.ToLower(format) {
	case "human-readable":
		return achcliDescribe(w, cfg, file)
	case "table":
		return tableDescribe(w, cfg, file)
	}
	return nil
}

func achcliDescribe(w io.Writer, cfg service.MaskingConfig, file *ach.File) error {
	w.Write([]byte("<pre>"))
	describe.File(w, file, &describe.Opts{
		MaskAccountNumbers: cfg.AccountNumbers,
		MaskCorrectedData:  cfg.CorrectedData,
		MaskNames:          cfg.Names,
		PrettyAmounts:      cfg.PrettyAmounts,
	})
	w.Write([]byte("</pre>"))
	return nil
}

var tableTemplate = display.InitTemplate("describe-file-table", "/webui/describe-file-table.html.tpl")

func tableDescribe(w io.Writer, cfg service.MaskingConfig, file *ach.File) error {
	tableTemplate.Execute(w, file)
	return nil
}
