package filelist

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/stretchr/testify/require"
)

func TestFilesystemGetFileJail(t *testing.T) {
	ach, err := os.ReadFile(filepath.Join("..", "..", "testdata", "ppd-debit.ach"))
	require.NoError(t, err)

	dir := t.TempDir()
	root := filepath.Join(dir, "ach")
	require.NoError(t, os.Mkdir(root, 0750))
	require.NoError(t, os.WriteFile(filepath.Join(root, "ppd-debit.ach"), ach, 0600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "secret.txt"), []byte("pwn"), 0600))

	t.Chdir(dir)

	ls, err := newFilesystemLister("src", &service.FilesystemConfig{Paths: []string{"ach"}})
	require.NoError(t, err)

	f, err := ls.GetFile(context.Background(), filepath.Join("ach", "ppd-debit.ach"))
	require.NoError(t, err)
	require.Equal(t, "ppd-debit.ach", f.Name)

	_, err = ls.GetFile(context.Background(), "secret.txt")
	require.EqualError(t, err, "invalid path")

	_, err = ls.GetFile(context.Background(), "../secret.txt")
	require.EqualError(t, err, "invalid path")

	_, err = ls.GetFile(context.Background(), "/etc/passwd")
	require.EqualError(t, err, "invalid path")
}
