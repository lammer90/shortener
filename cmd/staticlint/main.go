package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/unitchecker"
	"golang.org/x/tools/go/ssa"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/unused"
)

func main() {
	mychecks := []*analysis.Analyzer{
		asmdecl.Analyzer,
		atomic.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		httpresponse.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
	}

	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	simpleChecks := map[string]bool{
		"S1001": true,
		"S1002": true,
		"S1003": true,
	}

	for _, v := range simple.Analyzers {
		if simpleChecks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	quickfixChecks := map[string]bool{
		"QF1003": true,
		"QF1004": true,
	}

	for _, v := range simple.Analyzers {
		if quickfixChecks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	mychecks = append(mychecks, unused.Analyzer.Analyzer)
	mychecks = append(mychecks, exitCallerAnalyzer)

	unitchecker.Main(
		mychecks...,
	)
}

var exitCallerAnalyzer = &analysis.Analyzer{
	Name: "exitCaller",
	Doc:  "reports direct calls to os.Exit in main function",
	Run:  runExitCallerAnalyzer,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

func runExitCallerAnalyzer(pass *analysis.Pass) (interface{}, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)

	for _, fn := range s.SrcFuncs {
		if fn.Synthetic == "package initializer" {
			continue
		}
		if fn.Pkg.Pkg.Name() == "main" && fn.Name() == "main" {
			for _, b := range fn.Blocks {
				for _, i := range b.Instrs {
					call, ok := i.(*ssa.Call)
					if !ok {
						continue
					}
					obj := call.Common().StaticCallee()
					if obj == nil {
						continue
					}
					if obj.Pkg != nil && obj.Pkg.Pkg.Path() == "os" && obj.Name() == "Exit" {
						pass.Reportf(call.Pos(), "direct call to os.Exit detected in main function")
					}
				}
			}
		}
	}

	return nil, nil
}
