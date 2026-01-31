package api

import (
	"bytes"
	"compress/flate"
	"io"

	"github.com/rs/zerolog/log"
)

func CompressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	w, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}
	defer w.Close()

	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Flush()

	compressed := buf.Bytes()
	before := len(data)
	after := len(compressed)
	log.Info().Msgf("compressed from %d bytes to %d bytes (%d bytes saved, %.f%% reduction)", before, after, before-after, 100-float64(after)/float64(before)*100)

	return compressed, err
}

func DecompressData(data []byte) []byte {
	r := flate.NewReader(bytes.NewReader(data))
	defer r.Close()
	buf, _ := io.ReadAll(r)
	return buf
}
