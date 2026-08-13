package filelist

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadAckFile(t *testing.T) {
	report := readTestdataReport(t, "ACHFAHK673960043AIN202605261654134.ack")
	require.Equal(t, "ack", report.Kind())
	require.Greater(t, report.RecordCount(), int64(0))
	require.NotEmpty(t, report.Body)
	require.Equal(t, "1", report.Metadata()["Batches"])
	require.Equal(t, "1", report.Metadata()["Entries"])

	require.Len(t, report.Sections, 2)
	require.Equal(t, "FILE ERRORS:", report.Sections[0].Heading)
	require.Equal(t, "BATCH ERRORS:", report.Sections[1].Heading)

	err := report.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "TH145-SENDING ELECTRONIC CONNECTION OWNER")
	require.Contains(t, err.Error(), "FH239-INVALID SENDING POINT NOT AUTHORIZED")
	require.Contains(t, err.Error(), "BH232-INVALID ORIGINATING DFI IDENTIFICATION")
}

func TestReadAckFile_FileLevelErrors(t *testing.T) {
	report := readTestdataReport(t, "ACHFAHK673960043AIN202608121530761.ack")
	require.Equal(t, "1", report.Metadata()["Batches"])
	require.Equal(t, "2", report.Metadata()["Entries"])
	require.Len(t, report.Sections, 1)
	require.Equal(t, "FILE ERRORS:", report.Sections[0].Heading)

	err := report.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "FH238-INVALID IMMEDIATE ORIGIN NOT AUTHORIZED AS A SENDING POINT")
	require.Contains(t, err.Error(), "IMMEDIATE ORIGIN = 073923156")
}

func TestReadAckFile_AcceptedNoErrors(t *testing.T) {
	report := readTestdataReport(t, "ACHFAHK673960043AIN202608121534803.ack")
	require.Empty(t, report.Sections)
	require.NoError(t, report.Validate())
}

func TestReadAckFile_FileLevelExample(t *testing.T) {
	report := readTestdataReport(t, "file-level-example.ack")
	require.NotEmpty(t, report.Body)
	require.Equal(t, "AJ001A01A08052", report.Body[0])
	require.Contains(t, report.Body, "****** ACKNOWLEDGEMENT OF ACH FILE DEPOSITS ******")
}

func TestReadFile_ACHStillWorks(t *testing.T) {
	fd, err := os.Open(filepath.Join("..", "..", "testdata", "ppd-debit.ach"))
	require.NoError(t, err)
	t.Cleanup(func() { fd.Close() })

	parsed, err := Read(filepath.Base(fd.Name()), fd)
	require.NoError(t, err)
	achFile, ok := parsed.Document.(*ACHFile)
	require.True(t, ok)
	require.NotNil(t, achFile.File)
	require.NoError(t, parsed.Validate())
}

func TestReadFile_UnknownExtension(t *testing.T) {
	_, err := Read("notes.md", nil)
	require.EqualError(t, err, "unknown file extension for notes.md")
}

func TestReadAckFile_NotFedACH(t *testing.T) {
	fd, err := os.Open(filepath.Join("..", "..", "testdata", "ppd-debit.ach"))
	require.NoError(t, err)
	t.Cleanup(func() { fd.Close() })

	parsed, err := Read("not-fedach.ack", fd)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no FedACH records found")
	require.Nil(t, parsed.Document)
}

func TestReadFile_PendingFedACHFormat(t *testing.T) {
	parsed, err := Read("balance.eod", strings.NewReader("not parsed yet"))
	require.NoError(t, err)
	require.Equal(t, "eod", parsed.Document.Kind())
	require.ErrorIs(t, parsed.Validate(), ErrUnsupported{Kind: "eod", Title: "FedACH end-of-day balance report", Ext: ".eod"})
}

func TestRegisterFormat_NewExtension(t *testing.T) {
	RegisterFormat(Format{
		Kind:       "demo",
		Title:      "Demo Report",
		Extensions: []string{".demo"},
		Parse: func(name string, r io.Reader) (Document, error) {
			return &TextReport{
				ReportKind:  "demo",
				ReportTitle: "Demo Report",
				Body:        []string{"hello from " + name},
			}, nil
		},
	})

	parsed, err := Read("sample.demo", strings.NewReader("unused"))
	require.NoError(t, err)
	report, ok := parsed.Document.(*TextReport)
	require.True(t, ok)
	require.Equal(t, "demo", report.Kind())
	require.Equal(t, []string{"hello from sample.demo"}, report.Body)
}

func readTestdataReport(t *testing.T, name string) *TextReport {
	t.Helper()

	fd, err := os.Open(filepath.Join("..", "..", "testdata", "fedach", name))
	require.NoError(t, err)
	t.Cleanup(func() { fd.Close() })

	parsed, err := Read(name, fd)
	require.NoError(t, err)
	report, ok := parsed.Document.(*TextReport)
	require.True(t, ok)
	return report
}
