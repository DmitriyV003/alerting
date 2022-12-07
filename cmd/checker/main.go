package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	osExitChecker := &analysis.Analyzer{
		Name: "osexitchecker",
		Doc:  "Check Os.Exit",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			for _, p := range pass.Files {
				ast.Inspect(p, func(node ast.Node) bool {
					switch x := node.(type) {
					case *ast.File:
						if x.Name.Name != "main" {
							return false
						}
					case *ast.SelectorExpr:
						if x.Sel.Name == "Exit" {
							pass.Reportf(x.Pos(), "Called os.Exit() in main")
						}
					}
					return true
				})
			}
			return nil, nil
		},
	}

	checks := map[string]bool{
		"SA*": true,
	}
	var mychecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if checks[v.Name] {
			mychecks = append(mychecks, v)
		}
	}
	mychecks = append(mychecks, printf.Analyzer, shadow.Analyzer, structtag.Analyzer, osExitChecker)
	multichecker.Main(mychecks...)
}
