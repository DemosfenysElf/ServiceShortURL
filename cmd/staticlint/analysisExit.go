package main

import (
	"golang.org/x/tools/go/analysis"
)

var ExitAnalyzer = &analysis.Analyzer{
	Name: "os.Exit",
	Doc:  "check os.Exit in main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	//
	//for i, file := range pass.Files {
	//	if file.Name.Name == "main" {
	//
	//	}
	//}

	return nil, nil
}
