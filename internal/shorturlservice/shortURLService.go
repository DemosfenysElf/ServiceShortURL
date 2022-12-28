package shorturlservice

import (
	"log"
	"math/rand"
)

var urlmap = make(map[string]string)
var urlInfo = &URLInfo{}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetURL(short string, storage string) (url string) {

	consumerURL, err := NewConsumer(storage)
	if err != nil {
		return urlmap[short]
	}
	defer consumerURL.Close()

	for {
		readURL, err := consumerURL.ReadURL()
		if err != nil {
			break
		}
		if readURL.ShortURL == short {
			return urlmap[readURL.ShortURL]
		}
	}
	return urlmap[short]
}

func SetURL(url string, storageURL string) (short string) {
	short = shortURL()
	for _, ok := urlmap[short]; ok; {
		short = shortURL()
	}
	urlmap[short] = url
	////////// дублирование в файл

	urli := SetStructURL(url, short)

	producerURL, err := NewProducer(storageURL)
	if err != nil {
		return
	}
	defer producerURL.Close()
	if err := producerURL.WriteURL(urli); err != nil {
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

func SetStructURL(url string, short string) (info *URLInfo) {
	urlInfo.URL = url
	urlInfo.ShortURL = short
	info = urlInfo
	return
}

func SetStructCoockies(nameUser string, value string) {
	urlInfo.CookiesAuthentication = CookiesAuthentication{nameUser, value}

}

func GetStructCoockies() *CookiesAuthentication {
	return &urlInfo.CookiesAuthentication
}
