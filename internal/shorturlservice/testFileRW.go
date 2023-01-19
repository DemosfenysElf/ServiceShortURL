package shorturlservice

import (
	"encoding/json"
	"os"
)

// Producer:
type nRW struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func (fs *FileStorage) newRW(filename string) (*nRW, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0777)
	if err != nil {
		return nil, err
	}

	return &nRW{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (p *nRW) WriteURLInfo(short *URLInfo) error {
	return p.encoder.Encode(&short)
}

func (p *nRW) ReadURLInfo() (*URLInfo, error) {
	urli := &URLInfo{}
	if err := p.decoder.Decode(&urli); err != nil {
		return nil, err
	}
	return urli, nil
}

func (p *nRW) Close() error {
	return p.file.Close()
}
