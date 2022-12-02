package main

import (
	"ServiceShortURL/internal/router"
	"log"
)

func main() {
	err := router.Router()
	if err != nil {
		log.Fatal("Router:", err)
	}
}
