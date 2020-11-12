package keeper

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
)

var gzipIdent = []byte("\x1F\x8B\x08")

const maxSize = 400 * 1024

func UnCompress(src []byte) ([]byte, error) {
	if len(src) < 3 {
		return src, nil
	}

	if !bytes.Equal(gzipIdent, src[0:3]) {
		return src, nil
	}

	zr, err := gzip.NewReader(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}
	zr.Multistream(false)

	return ioutil.ReadAll(io.LimitReader(zr, maxSize))
}