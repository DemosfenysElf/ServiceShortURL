package main

import (
	"strings"

	analyzer1 "github.com/go-critic/go-critic/checkers/analyzer"
	analyzer2 "github.com/sashamelentyev/usestdlibvars/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	// создаём пустой массив проверок
	var mychecks []*analysis.Analyzer
	// проходимся по staticcheck.Analyzers
	for _, v := range staticcheck.Analyzers {
		// и добавляем из него в массив все проверки класса "SA"
		if strings.Contains(v.Analyzer.Name, "SA") {
			mychecks = append(mychecks, v.Analyzer)
		}

	}

	// проходимся по quickfix.Analyzers
	for _, v := range quickfix.Analyzers {
		// и добавляем из него в массив проверку "QF1003"
		if v.Analyzer.Name == "QF1003" {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	// выборочно добавляем проверки в массив из:
	// "golang.org/x/tools/go/analysis/"
	mychecks = append(mychecks, printf.Analyzer, shadow.Analyzer, structtag.Analyzer)
	// "github.com/go-critic/go-critic/"
	mychecks = append(mychecks, analyzer1.Analyzer)
	// где-то тут я поломался и не понял что дальше творю
	mychecks = append(mychecks, analyzer2.New())
	mychecks = append(mychecks, asmdecl.Analyzer)

	multichecker.Main(
		mychecks...,
	)

}
