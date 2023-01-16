package main

import (
	"log"

	"ServiceShortURL/internal/router"
)

func main() {
	rout := router.ServerShortener{}
	err := rout.Router()
	if err != nil {
		log.Fatal("Router:", err)
	}
}
