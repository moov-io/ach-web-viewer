package filelist

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/moov-io/ach-web-viewer/internal/gpgx"

	"golang.org/x/crypto/openpgp"
)

func maybeDecrypt(r io.Reader, gpgKey openpgp.EntityList) (io.Reader, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(gpgKey) > 0 {
		bs, err = gpgx.Decrypt(bs, gpgKey)
	}
	return bytes.NewReader(bs), err
}
