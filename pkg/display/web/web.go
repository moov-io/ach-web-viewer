package web

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/moov-io/ach"
	"github.com/moov-io/ach-web-viewer/pkg/display"
	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/moov-io/ach/cmd/achcli/describe"
)

type DisplayParams struct {
	Format string
	SortBy string
	Order  string
}

func File(w io.Writer, displayParams DisplayParams, cfg service.MaskingConfig, file *ach.File) error {
	switch strings.ToLower(displayParams.Format) {
	case "human-readable":
		return achcliDescribe(w, cfg, file)
	case "table":
		return tableDescribe(w, displayParams, cfg, file)
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

func tableDescribe(w io.Writer, displayParams DisplayParams, cfg service.MaskingConfig, file *ach.File) error {
	fmt.Printf("display params: %#v", displayParams)

	switch strings.ToLower(displayParams.SortBy) {
	// Batch Sorting
	case "companyname":
		sort.Slice(file.Batches, func(i, j int) bool {
			bh1 := file.Batches[i].GetHeader()
			bh2 := file.Batches[j].GetHeader()
			if asc(displayParams.Order) {
				return bh1.CompanyName < bh2.CompanyName
			}
			return bh1.CompanyName > bh2.CompanyName
		})

	case "companyidentification":

	// EntryDetail Sorting
	case "amount":
		for i := range file.Batches {
			sort.Slice(file.Batches[i], func(p, m int) bool {
				entries := file.Batches[i].GetEntries()
				if asc(displayParams.Order) {
					return entries[p].Amount < entries[m].Amount
				}
				return entries[p].Amount > entries[m].Amount
			})
		}

	case "individualname":

	case "tracenumber":

	}

	tableTemplate.Execute(w, file)
	return nil
}

func asc(order string) bool {
	return strings.ToLower(order) == "asc"
}
