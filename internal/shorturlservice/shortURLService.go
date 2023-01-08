package shorturlservice

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type CookiesAuthentication struct {
	NameUser  string `json:"NameUser"`
	ValueUser string `json:"ValueUser"`
}

type URLInfo struct {
	URL                   string `json:"url"`
	ShortURL              string `json:"shortURL"`
	CookiesAuthentication CookiesAuthentication
}

var urlInfo = &URLInfo{}

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

func GetStructURL() *URLInfo {
	return urlInfo
}
