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

// NewProducer открываем файл для записи и чтения
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

// WriteURLInfo сохраняем в файл короткий url, оригинальный, пользователя
func (p *nRW) WriteURLInfo(short *URLInfo) error {
	return p.encoder.Encode(&short)
}

// WriteDelet перезапись метки об удалении
func (p *nRW) WriteDelet(i int64) error {
	del := []byte("delet")
	_, err := p.file.WriteAt(del, i)
	if err != nil {
		return err
	}
	return nil
}

// ReadURLInfo получение данных о сохраненной ссылке
func (p *nRW) ReadURLInfo() (*URLInfo, int, error) {
	urli := &URLInfo{}
	if err := p.decoder.Decode(&urli); err != nil {
		return nil, 0, err
	}
	jString, _ := json.Marshal(urli)
	lenInfo := len(jString)
	return urli, lenInfo, nil
}

// Close закрываем файл
func (p *nRW) Close() error {
	return p.file.Close()
}
