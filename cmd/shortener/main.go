package main

import (
	"log"

	"ServiceShortURL/internal/router"
)

func main() {
	rout := router.URLServer{}
	err := rout.Router()
	if err != nil {
		log.Fatal("Router:", err)
	}
}
