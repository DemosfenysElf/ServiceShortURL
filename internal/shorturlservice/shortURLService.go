package shorturlservice

import "math/rand"

var urlmap = make(map[string]string)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetURL(short string) (url string) {
	return urlmap[short]
}

func SetURL(url string) (short string) {
	short = shortURL()

	for _, ok := urlmap[short]; ok; {
		short = shortURL()
	}
	urlmap[short] = url

	return
}

func shortURL() string {

	a := make([]byte, 7)
	for i := range a {
		a[i] = letters[rand.Intn(len(letters))]
	}
	return string(a)
}
