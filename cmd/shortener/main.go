package main

import (
	"ServiceShortURL/internal/router"
	"log"
)

func main() {
	rout := router.Server{}
	err := rout.Router()
	if err != nil {
		log.Fatal("Router:", err)
	}
}
