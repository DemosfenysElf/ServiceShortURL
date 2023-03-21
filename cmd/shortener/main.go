package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"

	"ServiceShortURL/internal/router"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	printBuildVersion()
	rout := router.InitServer()
	err := rout.Router()
	if err != nil {
		log.Fatal("Router:", err)
	}
}

func printBuildVersion() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
}
