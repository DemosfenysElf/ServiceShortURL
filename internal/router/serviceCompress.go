package router

import (
	"bytes"
	"compress/gzip"
	"log"
)

func serviceCompress(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write(b)
	if err != nil {
		log.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}

	return buf.Bytes(), nil
}
