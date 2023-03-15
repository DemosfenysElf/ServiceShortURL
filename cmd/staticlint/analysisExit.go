package main

import (
	"golang.org/x/tools/go/analysis"
)

// exitAnalyzer don't work
var exitAnalyzer = &analysis.Analyzer{
	Name: "os.Exit",
	Doc:  "check os.Exit in main",
	Run:  run,
}

// rundon't work
func run(pass *analysis.Pass) (interface{}, error) {
	//
	//for i, file := range pass.Files {
	//	if file.Name.Name == "main" {
	//
	//	}
	//}

	return nil, nil
}
