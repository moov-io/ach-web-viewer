package filelist

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/moov-io/cryptfs"
	"github.com/stretchr/testify/require"
)

func TestMaybeDecrypt_Plaintext(t *testing.T) {
	ls := &bucketLister{}
	bs, err := ls.maybeDecrypt(strings.NewReader("        AJ001A01A08052"))
	require.NoError(t, err)
	require.Equal(t, "        AJ001A01A08052", string(bs))
}

func TestMaybeDecrypt_EncryptedWithoutKeys(t *testing.T) {
	ls := &bucketLister{}
	_, err := ls.maybeDecrypt(strings.NewReader("-----BEGIN PGP MESSAGE-----\n\nwV4D\n-----END PGP MESSAGE-----\n"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "no decryption keys")
}

func TestMaybeDecrypt_EncryptedWrongKey(t *testing.T) {
	pub := filepath.Join(os.Getenv("GOPATH"), "pkg", "mod", "github.com", "moov-io", "cryptfs@v0.11.0", "internal", "gpgx", "testdata", "key.pub")
	priv := filepath.Join(os.Getenv("GOPATH"), "pkg", "mod", "github.com", "moov-io", "cryptfs@v0.11.0", "internal", "gpgx", "testdata", "key.priv")
	if os.Getenv("GOPATH") == "" {
		home, err := os.UserHomeDir()
		require.NoError(t, err)
		pub = filepath.Join(home, "go", "pkg", "mod", "github.com", "moov-io", "cryptfs@v0.11.0", "internal", "gpgx", "testdata", "key.pub")
		priv = filepath.Join(home, "go", "pkg", "mod", "github.com", "moov-io", "cryptfs@v0.11.0", "internal", "gpgx", "testdata", "key.priv")
	}

	enc, err := cryptfs.FromCryptor(cryptfs.NewGPGEncryptorFile(pub))
	require.NoError(t, err)
	cipher, err := enc.Disfigure([]byte("        AJ001A01A08052"))
	require.NoError(t, err)
	require.True(t, looksEncrypted(cipher))

	dec, err := cryptfs.FromCryptor(cryptfs.NewGPGDecryptorFile(priv, []byte("wrong-password")))
	// A wrong password may fail at key-load time; either way Reveal must not
	// silently return the ciphertext.
	if err != nil {
		ls := &bucketLister{}
		_, revealErr := ls.maybeDecrypt(bytes.NewReader(cipher))
		require.Error(t, revealErr)
		return
	}

	ls := &bucketLister{cryptors: []*cryptfs.FS{dec}}
	_, err = ls.maybeDecrypt(bytes.NewReader(cipher))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to decrypt")
}

func TestMaybeDecrypt_EncryptedMatchingKey(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	dir := filepath.Join(home, "go", "pkg", "mod", "github.com", "moov-io", "cryptfs@v0.11.0", "internal", "gpgx", "testdata")
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		dir = filepath.Join(gopath, "pkg", "mod", "github.com", "moov-io", "cryptfs@v0.11.0", "internal", "gpgx", "testdata")
	}

	enc, err := cryptfs.FromCryptor(cryptfs.NewGPGEncryptorFile(filepath.Join(dir, "key.pub")))
	require.NoError(t, err)
	cipher, err := enc.Disfigure([]byte("        AJ001A01A08052"))
	require.NoError(t, err)

	dec, err := cryptfs.FromCryptor(cryptfs.NewGPGDecryptorFile(filepath.Join(dir, "key.priv"), []byte("password")))
	require.NoError(t, err)

	ls := &bucketLister{cryptors: []*cryptfs.FS{dec}}
	bs, err := ls.maybeDecrypt(bytes.NewReader(cipher))
	require.NoError(t, err)
	require.Equal(t, "        AJ001A01A08052", string(bs))
}
