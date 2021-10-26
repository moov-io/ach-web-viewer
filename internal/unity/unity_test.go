package unity

import (
	"path/filepath"
	"testing"

	"github.com/moov-io/ach-web-viewer/internal/gpgproton"
	"github.com/moov-io/ach-web-viewer/internal/gpgx"

	"github.com/stretchr/testify/require"
)

var (
	password = []byte("password")

	privateKeyPath = filepath.Join("..", "gpgx", "testdata", "moov.key")
	publicKeyPath  = filepath.Join("..", "gpgx", "testdata", "moov.pub")
)

func TestGPG__GoToProton(t *testing.T) {
	// Encrypt
	pubKey, err := gpgx.ReadArmoredKeyFile(publicKeyPath)
	require.NoError(t, err)
	msg, err := gpgx.Encrypt([]byte("hello, world"), pubKey)
	require.NoError(t, err)
	if len(msg) == 0 {
		t.Error("empty encrypted message")
	}

	// Decrypt
	privKey, err := gpgproton.ReadPrivateKeyFile(privateKeyPath, password)
	require.NoError(t, err)
	require.NoError(t, err)
	out, err := gpgproton.Decrypt(msg, privKey)
	require.NoError(t, err)

	if v := string(out); v != "hello, world" {
		t.Errorf("got %q", v)
	}
}

func TestGPG__ProtonToGo(t *testing.T) {
	// Encrypt
	pubKey, err := gpgproton.ReadArmoredKeyFile(publicKeyPath)
	require.NoError(t, err)
	msg, err := gpgproton.Encrypt([]byte("hello, world"), pubKey)
	require.NoError(t, err)
	if len(msg) == 0 {
		t.Error("empty encrypted message")
	}

	// Decrypt
	privKey, err := gpgx.ReadPrivateKeyFile(privateKeyPath, password)
	require.NoError(t, err)
	require.NoError(t, err)
	out, err := gpgx.Decrypt(msg, privKey)
	require.NoError(t, err)

	if v := string(out); v != "hello, world" {
		t.Errorf("got %q", v)
	}
}
