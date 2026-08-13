package web

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moov-io/ach-web-viewer/pkg/filelist"
	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/stretchr/testify/require"
)

func TestContents_AckReport(t *testing.T) {
	parsed := readTestdata(t, "fedach", "ACHFAHK673960043AIN202605261654134.ack")

	var buf bytes.Buffer
	err := Contents(&buf, &service.DisplayConfig{Format: "human-readable"}, parsed)
	require.NoError(t, err)

	out := buf.String()
	require.Contains(t, out, "Type:     ACK (FAHK - Acknowledgement of ACH File Deposits)")
	require.Contains(t, out, "ORIGINAL FILE")
	require.Contains(t, out, "1 batches, 1 entries")
	require.Contains(t, out, "0.00 debits, 321.45 credits")
	require.Contains(t, out, "FILE ERRORS:")
	require.Contains(t, out, "TH145-SENDING ELECTRONIC CONNECTION OWNER")
	require.Contains(t, out, "FH239-INVALID SENDING POINT NOT AUTHORIZED")
	require.Contains(t, out, "BATCH ERRORS:")
	require.Contains(t, out, "BH232-INVALID ORIGINATING DFI IDENTIFICATION")
	require.Contains(t, out, "=== Reconstructed Visual Lines ===")
	require.Contains(t, out, "ACKNOWLEDGEMENT OF ACH FILE DEPOSITS")
}

func TestContents_AckReport_Accepted(t *testing.T) {
	parsed := readTestdata(t, "fedach", "ACHFAHK673960043AIN202608121534803.ack")

	var buf bytes.Buffer
	err := Contents(&buf, &service.DisplayConfig{}, parsed)
	require.NoError(t, err)

	out := buf.String()
	require.Contains(t, out, "ACK (FAHK - Acknowledgement of ACH File Deposits)")
	require.NotContains(t, out, "FILE ERRORS:")
	require.NotContains(t, out, "BATCH ERRORS:")
	require.Contains(t, out, "=== Reconstructed Visual Lines ===")
}

func TestContents_ACHFile(t *testing.T) {
	parsed := readTestdata(t, "ppd-debit.ach")

	var buf bytes.Buffer
	err := Contents(&buf, &service.DisplayConfig{Format: "human-readable"}, parsed)
	require.NoError(t, err)
	require.Contains(t, buf.String(), "Origin")
	require.Contains(t, buf.String(), "BatchNumber")
}

func TestContents_PendingFedACHFormat(t *testing.T) {
	parsed, err := filelist.Read("balance.eod", strings.NewReader("unused"))
	require.NoError(t, err)

	var buf bytes.Buffer
	err = Contents(&buf, &service.DisplayConfig{Format: "human-readable"}, parsed)
	require.NoError(t, err)
	require.Contains(t, buf.String(), "FedACH end-of-day balance report")
	require.Contains(t, buf.String(), "not yet supported")
}

func readTestdata(t *testing.T, parts ...string) *filelist.File {
	t.Helper()

	path := filepath.Join(append([]string{"..", "..", "..", "testdata"}, parts...)...)
	fd, err := os.Open(path)
	require.NoError(t, err)
	t.Cleanup(func() { fd.Close() })

	parsed, err := filelist.Read(filepath.Base(path), fd)
	require.NoError(t, err)
	return parsed
}
