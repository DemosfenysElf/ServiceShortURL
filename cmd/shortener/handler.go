package main

import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postURLToShort(w, r)
	case "GET":
		getShortToURL(w, r)
	}
}
