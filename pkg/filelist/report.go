package filelist

import (
	"fmt"
	"io"
	"strings"
)

// TextReport is a parsed line-oriented report (FedACH ACK today; other
// FedPayments Reporter formats can reuse this as they are added).
type TextReport struct {
	ReportKind  string
	ReportTitle string

	Header []string

	SummaryHeading string
	Summary        []string

	Sections []ReportSection

	BodyHeading string
	Body        []string
	NumberBody  bool

	Records int64
	Meta    map[string]string
	Err     error
}

// ReportSection is a labeled block in a TextReport (errors, totals, etc.).
type ReportSection struct {
	Heading string
	Lines   []string
}

func (r *TextReport) Kind() string {
	if r == nil {
		return ""
	}
	return r.ReportKind
}

func (r *TextReport) Title() string {
	if r == nil {
		return ""
	}
	return r.ReportTitle
}

func (r *TextReport) Validate() error {
	if r == nil {
		return errMissingFile
	}
	return r.Err
}

func (r *TextReport) RecordCount() int64 {
	if r == nil {
		return 0
	}
	return r.Records
}

func (r *TextReport) Metadata() map[string]string {
	if r == nil {
		return nil
	}
	out := make(map[string]string, len(r.Meta)+1)
	if r.ReportKind != "" {
		out["Report Type"] = strings.ToUpper(r.ReportKind)
	}
	for k, v := range r.Meta {
		out[k] = v
	}
	return out
}

func (r *TextReport) WriteHuman(w io.Writer) error {
	if r == nil {
		return nil
	}

	if r.ReportTitle != "" {
		fmt.Fprintf(w, "Type:     %s\n", r.ReportTitle)
	}
	for _, line := range r.Header {
		fmt.Fprintln(w, line)
	}
	if len(r.Header) > 0 || r.ReportTitle != "" {
		fmt.Fprintln(w)
	}

	if r.SummaryHeading != "" {
		fmt.Fprintln(w, r.SummaryHeading)
	}
	for _, line := range r.Summary {
		fmt.Fprintln(w, line)
	}
	if r.SummaryHeading != "" || len(r.Summary) > 0 {
		fmt.Fprintln(w)
	}

	for _, section := range r.Sections {
		if section.Heading != "" {
			fmt.Fprintln(w, section.Heading)
		}
		for _, line := range section.Lines {
			fmt.Fprintln(w, line)
		}
		fmt.Fprintln(w)
	}

	if r.BodyHeading != "" {
		fmt.Fprintln(w, r.BodyHeading)
	}
	for i, line := range r.Body {
		if r.NumberBody {
			fmt.Fprintf(w, "%3d: %s\n", i+1, line)
		} else {
			fmt.Fprintln(w, line)
		}
	}
	return nil
}
