package shorturlservice

import (
	"github.com/caarlos0/env"
	"log"
	"math/rand"
)

var urlmap = make(map[string]string)

type FileStorage struct {
	Storage string `env:"FILE_STORAGE_PATH"`
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetURL(short string) (url string) {

	f := FileStorage{}
	errConfig := env.Parse(&f)
	if errConfig != nil {
		return urlmap[short]
	}

	cons, err := NewConsumer(f.Storage)
	if err != nil {
		log.Fatal(err)
	}
	defer cons.Close()

	for {
		readURL, err := cons.ReadURL()
		if err != nil {
			break
		}
		if readURL.ShortURL == short {
			return urlmap[readURL.ShortURL]
		}
	}

	return urlmap[short]
}

func SetURL(url string) (short string) {
	short = shortURL()
	for _, ok := urlmap[short]; ok; {
		short = shortURL()
	}
	urlmap[short] = url
	////////// дублирование в файл
	f := FileStorage{}
	errConfig := env.Parse(&f)
	if errConfig != nil {
		log.Fatal(errConfig)
	}
	urli := &URLInfo{URL: url, ShortURL: short}

	//defer os.Remove(fileName)
	prod, err := NewProducer(f.Storage)
	if err != nil {
		log.Fatal(err)
	}
	defer prod.Close()
	if err := prod.WriteURL(urli); err != nil {
		log.Fatal(err)
	}
	return
}

func shortURL() string {

	a := make([]byte, 7)
	for i := range a {
		a[i] = letters[rand.Intn(len(letters))]
	}
	return string(a)
}
