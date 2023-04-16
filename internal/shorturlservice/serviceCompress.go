package shorturlservice

import (
	"bytes"
	"compress/gzip"
)

func ServiceCompress(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	_, err := zw.Write(b)
	if err != nil {
		return nil, err
	}

	if err = zw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

//func serviceDeCompress(data []byte) ([]byte, error) {
//	r, _ := gzip.NewReader(bytes.NewReader(data))
//	defer r.Close()
//	var b bytes.Buffer
//	_, err := b.ReadFrom(r)
//	if err != nil {
//		return nil, fmt.Errorf("failed decompress data: %v", err)
//	}
//	return b.Bytes(), nil
//}
