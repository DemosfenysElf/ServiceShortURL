package main

import (
	"log"
	"router"
)

func main() {

	err := router.Router()
	if err != nil {
		log.Fatal("Router:", err)
	}
}
