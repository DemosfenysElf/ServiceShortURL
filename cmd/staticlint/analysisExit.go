package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// osExitAnalyzer проверяет не используется ли os.Exit() в main
var osExitAnalyzer = &analysis.Analyzer{
	Name: "osExitCheck",
	Doc:  "check os.Exit in main",
	Run:  run,
}

// run osExitAnalyzer
func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			ast.Inspect(file, func(n ast.Node) bool {
				if c, ok := n.(*ast.CallExpr); ok {
					if s, ok := c.Fun.(*ast.SelectorExpr); ok {
						if s.X.(*ast.Ident).Name == "os" && s.Sel.Name == "Exit" {
							pass.Reportf(s.X.(*ast.Ident).NamePos, "calling os.Exit in main")
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
