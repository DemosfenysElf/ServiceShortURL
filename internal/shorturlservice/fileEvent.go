package shorturlservice

import (
	"encoding/json"
	"os"
)

// Producer:
type producer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewProducer открываем файл для записи
func NewProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteURL сохраняем в файл короткий url, оригинальный, пользователя
func (p *producer) WriteURL(short *URLInfo) error {
	return p.encoder.Encode(&short)
}

// WriteUser сохраняем в файл пользователя
func (p *producer) WriteUser(user *CookiesAuthentication) error {
	return p.encoder.Encode(&user)
}

// Close producer.
func (p *producer) Close() error {
	return p.file.Close()
}

// Consumer:
type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

// NewConsumer открываем файл для чтения
func NewConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// ReadURLInfo получаем из файла строку данных о сохранённом URL в виде структуры URLInfo.
func (c *consumer) ReadURLInfo() (*URLInfo, error) {
	urli := &URLInfo{}
	if err := c.decoder.Decode(&urli); err != nil {
		return nil, err
	}
	return urli, nil
}

// ReadUser получаем из файла строку данных о пользователе в виде структуры CookiesAuthentication.
func (c *consumer) ReadUser() (*CookiesAuthentication, error) {
	user := &CookiesAuthentication{}
	if err := c.decoder.Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

// Close сonsumer.
func (c *consumer) Close() error {
	return c.file.Close()
}
