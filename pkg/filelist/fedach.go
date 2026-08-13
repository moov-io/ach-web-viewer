package filelist

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/moov-io/fedach/pkg/ack"
	"golang.org/x/text/encoding/charmap"
)

func init() {
	RegisterFormat(Format{
		Kind:       "ack",
		Title:      "ACK (FAHK - Acknowledgement of ACH File Deposits)",
		Extensions: []string{".ack"},
		Parse:      parseACK,
	})

	// Known FedPayments Reporter / FedACH types that achgateway already
	// distinguishes. Replace these with real parsers as moov-io/fedach grows.
	RegisterFormat(pendingFormat("adv", "FedACH advice", ".adv"))
	RegisterFormat(pendingFormat("crf", "FedACH CRF", ".crf"))
	RegisterFormat(pendingFormat("eod", "FedACH end-of-day balance report", ".eod"))
	RegisterFormat(pendingFormat("xml", "FedEDI Plus", ".xml"))
}

// pendingFormat registers a known report type that we cannot parse yet.
// Opening one still produces a Document so the viewer can explain the gap
// instead of treating the file as unknown.
func pendingFormat(kind, title string, exts ...string) Format {
	return Format{
		Kind:       kind,
		Title:      title,
		Extensions: exts,
		Parse: func(name string, r io.Reader) (Document, error) {
			ext := normalizeExt(filepath.Ext(name))
			err := ErrUnsupported{Kind: kind, Title: title, Ext: ext}
			return &TextReport{
				ReportKind:  kind,
				ReportTitle: title,
				Summary:     []string{err.Error()},
				Err:         err,
			}, nil
		},
	}
}

func parseACK(_ string, r io.Reader) (Document, error) {
	bs, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading ack file: %w", err)
	}

	records := ack.Split(bs)
	if !isAckReport(bs, records) {
		if looksEncrypted(bs) {
			return nil, errors.New("file is still GPG-encrypted; decryption failed or no key matched")
		}
		// Incoming IBM Connect:Direct FedACH files are often still
		// EBCDIC CP037 after GPG unwrap (80-byte blocked records).
		if decoded, ok := decodeEBCDICACK(bs); ok {
			bs = decoded
			records = ack.Split(bs)
		}
	}
	if !isAckReport(bs, records) {
		return nil, errors.New("no FedACH records found")
	}

	totals, _ := ack.ParseFileTotals(records)
	fileErrors, batchErrors := ack.FindErrorBlocks(records)
	lines := ack.SplitLines(bs)

	meta := map[string]string{}
	if totals.Batches > 0 || totals.Entries > 0 {
		meta["Batches"] = fmt.Sprintf("%d", totals.Batches)
		meta["Entries"] = fmt.Sprintf("%d", totals.Entries)
	}

	report := &TextReport{
		ReportKind:     "ack",
		ReportTitle:    "ACK (FAHK - Acknowledgement of ACH File Deposits)",
		Header:         []string{fmt.Sprintf("Records:  %d tagged logical records", len(records)), fmt.Sprintf("Lines:    %d reconstructed visual lines", len(lines))},
		SummaryHeading: "ORIGINAL FILE",
		Summary: []string{
			fmt.Sprintf("%d batches, %d entries", totals.Batches, totals.Entries),
			fmt.Sprintf("%.2f debits, %.2f credits", float64(totals.DebitTotal)/100.0, float64(totals.CreditTotal)/100.0),
		},
		BodyHeading: "=== Reconstructed Visual Lines ===",
		Body:        lines,
		NumberBody:  true,
		Records:     int64(len(records)),
		Meta:        meta,
		Err:         ackValidationError(fileErrors, batchErrors),
	}

	if len(fileErrors) > 0 {
		report.Sections = append(report.Sections, formatAckSection("FILE ERRORS:", fileErrors))
	}
	if len(batchErrors) > 0 {
		report.Sections = append(report.Sections, formatAckSection("BATCH ERRORS:", batchErrors))
	}

	return report, nil
}

func formatAckSection(heading string, blocks [][]ack.Record) ReportSection {
	section := ReportSection{Heading: heading}
	for _, block := range blocks {
		if msg := ack.FormatErrorBlock(block); msg != "" {
			section.Lines = append(section.Lines, "  "+msg)
		}
	}
	return section
}

func ackValidationError(fileErrors, batchErrors [][]ack.Record) error {
	var parts []string
	for _, block := range fileErrors {
		if msg := ack.FormatErrorBlock(block); msg != "" {
			parts = append(parts, msg)
		}
	}
	for _, block := range batchErrors {
		if msg := ack.FormatErrorBlock(block); msg != "" {
			parts = append(parts, msg)
		}
	}
	if len(parts) == 0 {
		return nil
	}
	return errors.New(strings.Join(parts, "; "))
}

// decodeEBCDICACK converts IBM-037 (US EBCDIC) to UTF-8 when the result
// looks like a FedACH FAHK acknowledgement.
func decodeEBCDICACK(bs []byte) ([]byte, bool) {
	decoded, err := charmap.CodePage037.NewDecoder().Bytes(bs)
	if err != nil || len(decoded) == 0 {
		return nil, false
	}
	if !isAckReport(decoded, ack.Split(decoded)) {
		return nil, false
	}
	return decoded, true
}

// isAckReport reports whether the bytes (and any records Split produced)
// look like a FedACH FAHK acknowledgement rather than an ACH payment file
// that happens to contain uppercase letters.
func isAckReport(raw []byte, records []ack.Record) bool {
	if len(records) == 0 {
		return false
	}
	if bytes.Contains(raw, []byte("AJ001A01A08052")) {
		return true
	}
	if bytes.Contains(bytes.ToUpper(raw), []byte("ACKNOWLEDGEMENT OF ACH FILE DEPOSITS")) {
		return true
	}
	return records[0].Prefix == 'A' && bytes.Contains(records[0].Content, []byte("AJ001A01A08052"))
}
