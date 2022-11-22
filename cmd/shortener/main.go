package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

var urlmap = make(map[string]string)

func main() {

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postURLToShort(w, r)
	case "GET":
		getShortToURL(w, r)
	}
}

func postURLToShort(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	short := shortURL()
	for ; urlmap[short] != string(body); short = shortURL() {
		urlmap[short] = string(body)
	}
	write := []byte("http://localhost:8080/" + short)
	w.WriteHeader(http.StatusCreated)
	w.Write(write)
}

func getShortToURL(w http.ResponseWriter, r *http.Request) {
	short := r.URL.String()
	short = short[1:]
	if urlmap[short] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Location", urlmap[short])
	w.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Println(w.Header())
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func shortURL() string {
	a := make([]byte, 7)
	for i := range a {
		a[i] = letters[rand.Intn(len(letters))]
	}
	return string(a)
}
