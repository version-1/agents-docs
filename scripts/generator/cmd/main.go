package main

import (
	"flag"
	"fmt"
	"os"

	"generator/internal/app"
	"generator/internal/infra"
)

func main() {
	var (
		inRoot  = flag.String("input", "", "input root directory to scan (e.g., ./.)")
		outRoot = flag.String("output", "out/.codex/skills", "output directory (e.g., .out/codex/skills)")
	)
	flag.Parse()

	gen := app.NewGenerator(infra.OSFS{})
	if err := gen.Run(*inRoot, *outRoot); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
