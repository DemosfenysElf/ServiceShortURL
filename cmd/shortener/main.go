package main

import (
	"log"
)

func main() {
	err := router()
	if err != nil {
		log.Fatal("Router:", err)
	}
}
