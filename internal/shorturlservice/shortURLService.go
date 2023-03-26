package shorturlservice

import (
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Данные о пользователе
type CookiesAuthentication struct {
	NameUser  string `json:"NameUser"`
	ValueUser string `json:"ValueUser"`
}

// Данные о URL и  пользователе
type URLInfo struct {
	URL                   string `json:"url"`
	ShortURL              string `json:"shortURL"`
	CookiesAuthentication CookiesAuthentication
	Deleted               string `json:"deleted"`
}

var urlInfo = &URLInfo{}

type Generator interface {
	shortURL() string
}

type TestGenerator struct {
	Result []string
	Index  int
}

func (g *TestGenerator) shortURL() string {
	if g.Index >= len(g.Result) {
		return ""
	}
	ret := g.Result[g.Index]

	g.Index += 1
	return ret
}

type RandomGenerator struct {
}

// shortURL Генератор коротких ссылок
func (RandomGenerator) shortURL() string {
	a := make([]byte, 7)
	for i := range a {
		a[i] = letters[rand.Intn(len(letters))]
	}
	return string(a)
}

// SetStructURL запись данных в структуру и получение структуры
func SetStructURL(url string, short string) (info *URLInfo) {
	urlInfo.URL = url
	urlInfo.ShortURL = short
	urlInfo.Deleted = "false"
	// urlInfo.D = " "
	info = urlInfo
	return
}

// SetStructCookies запись данных в структуру и получение структуры
func SetStructCookies(nameUser string, value string) {
	urlInfo.CookiesAuthentication = CookiesAuthentication{nameUser, value}

}

// GetStructCookies получение данных о пользователе
func GetStructCookies() *CookiesAuthentication {
	return &urlInfo.CookiesAuthentication
}

// GetStructURL получение данных о пользователе и URL
func GetStructURL() *URLInfo {
	return urlInfo
}

// GetCookieValue получение расшифрованного пользователя из куки
func GetCookieValue(body []*http.Cookie) string {
	if len(body) == 0 {
		return ""
	}
	deHexCookies, err := hex.DecodeString(body[0].Value)
	if err != nil {
		log.Fatal(err)
	}
	deCryptoCookies, err := DeCryptoToken(deHexCookies)
	if err != nil {
		log.Fatal(err)
	}
	hexCookies := hex.EncodeToString(deCryptoCookies)
	return hexCookies
}
