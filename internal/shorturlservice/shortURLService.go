package shorturlservice

import (
	"log"
	"math/rand"
)

var urlmap = make(map[string]string)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetURL(short string, storage string) (url string) {

	cons, err := NewConsumer(storage)
	if err != nil {
		return urlmap[short]
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

func SetURL(url string, storage string) (short string) {
	short = shortURL()
	for _, ok := urlmap[short]; ok; {
		short = shortURL()
	}
	urlmap[short] = url
	////////// дублирование в файл

	urli := &URLInfo{URL: url, ShortURL: short}

	prod, err := NewProducer(storage)
	if err != nil {
		return
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
