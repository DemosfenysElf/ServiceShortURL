package main

import "log"

var urlmap = make(map[string]string)

func main() {
	err := router()
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
