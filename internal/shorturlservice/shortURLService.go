package shorturlservice

import (
	"fmt"
	"log"
	"math/rand"
)

var urlmap = make(map[string]string)
var urlInfo = &URLInfo{}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetURL(short string, storage string) (url string, err error) {

	consumerURL, err := NewConsumer(storage)
	if err != nil {
		return urlmap[short], nil
	}
	defer consumerURL.Close()

	for {
		readURL, err := consumerURL.ReadURLInfo()
		if err != nil {
			break
		}
		fmt.Println(readURL)
		if readURL.ShortURL == short {
			return readURL.URL, nil
		}
	}

	return "", fmt.Errorf("no found url")
}

func SetURL(url string, storageURL string) (short string) {
	short = shortURL()
	for {
		ssss, _ := GetURL(short, storageURL)
		if ssss == "" {
			break
		}
		short = shortURL()
	}

	//fmt.Println("-----File-Short: ", short)

	////////// дублирование в файл

	urli := SetStructURL(url, short)
	//fmt.Println("-----File-urli: ", urli)
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

func SetStructCookies(nameUser string, value string) {
	urlInfo.CookiesAuthentication = CookiesAuthentication{nameUser, value}

}

func GetStructCookies() *CookiesAuthentication {
	return &urlInfo.CookiesAuthentication
}

//func GetStructURL() *URLInfo {
//	return urlInfo
//}
