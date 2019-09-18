package main

import (
	"flag"
	"go/ast"
	"log"
	"os"
	"strings"

	"github.com/ashanbrown/nolintlint/nolintlint"
	"golang.org/x/tools/go/packages"
)

func main() {
	log.SetFlags(0) // remove log timestamp

	setExitStatus := flag.Bool("set_exit_status", false, "Set exit status to 1 if any issues are found")
	explain := flag.Bool("explain", true, "Require explanation for nolint directives")
	specific := flag.Bool("specific", true, "Require specific linters for nolint directives")
	machine := flag.Bool("machine", false, "Require machine-readable directives")
	nolint := flag.String("nolint", "nolint", "comma-separated list of nolint directives")

	flag.Parse()

	cfg := packages.Config{
		Mode: packages.NeedSyntax |
			packages.NeedName |
			packages.NeedTypes,
	}
	pkgs, err := packages.Load(&cfg, flag.Args()...)
	if err != nil {
		log.Fatalf("Could not load packages: %s", err)
	}
	var needs nolintlint.Needs
	if *explain {
		needs |= nolintlint.NeedsExplanation
	}
	if *specific {
		needs |= nolintlint.NeedsSpecific
	}
	if *machine {
		needs |= nolintlint.NeedsMachine
	}
	linter := nolintlint.NewLinter(strings.Split(*nolint, ","), needs)

	var issues []nolintlint.Issue //nolint:prealloc // don't know how many there will be
	for _, p := range pkgs {
		nodes := make([]ast.Node, 0, len(p.Syntax))
		for _, n := range p.Syntax {
			nodes = append(nodes, n)
		}
		newIssues, err := linter.Run(p.Fset, nodes...)
		if err != nil {
			log.Fatalf("failed: %s", err)
		}
		if err != nil {
			log.Fatalf("failed: %s", err)
		}
		issues = append(issues, newIssues...)
	}

	for _, issue := range issues {
		log.Println(issue)
	}

	if *setExitStatus && len(issues) > 0 {
		os.Exit(1)
	}
}
