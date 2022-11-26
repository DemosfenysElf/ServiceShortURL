package main

import (
	"fmt"
	"io"
	"net/http"
)

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
