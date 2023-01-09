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

func (p *producer) WriteURL(short *URLInfo) error {
	return p.encoder.Encode(&short)
}

func (p *producer) WriteUser(user *CookiesAuthentication) error {
	return p.encoder.Encode(&user)
}

func (p *producer) Close() error {
	return p.file.Close()
}

// Consumer:
type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

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
func (c *consumer) ReadURLInfo() (*URLInfo, error) {
	urli := &URLInfo{}
	if err := c.decoder.Decode(&urli); err != nil {
		return nil, err
	}
	return urli, nil
}

func (c *consumer) ReadUser() (*CookiesAuthentication, error) {
	user := &CookiesAuthentication{}
	if err := c.decoder.Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}
