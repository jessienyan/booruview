package api

import (
	"bytes"
	"compress/flate"
	"io"
)

func CompressData(dst io.Writer, data []byte) error {
	w, err := flate.NewWriter(dst, flate.BestCompression)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write(data)
	return err
}

func DecompressData(data []byte) []byte {
	r := flate.NewReader(bytes.NewReader(data))
	defer r.Close()
	buf := bytes.Buffer{}
	io.Copy(&buf, r)
	return buf.Bytes()
}
