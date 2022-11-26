package main

import (
	"net/http"
)

func router() error {

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}
	return nil
}
