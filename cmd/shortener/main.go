package main

import (
	"log"
	_ "net/http/pprof"

	"ServiceShortURL/internal/router"
)

func main() {
	rout := router.InitServer()
	err := rout.Router()
	if err != nil {
		log.Fatal("Router:", err)
	}
}
