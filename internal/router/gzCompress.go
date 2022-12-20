package router

import (
	"bytes"
	"compress/flate"
	"fmt"
)

func serviceCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		fmt.Println(">>>>>>>>>>A_5")
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}

	_, err = w.Write(data)
	if err != nil {
		fmt.Println(">>>>>>>>>>A_6")
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	err = w.Close()
	if err != nil {
		fmt.Println(">>>>>>>>>>A_7")
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	fmt.Println(">>>>>>>>>>A_8")
	return b.Bytes(), nil
}
