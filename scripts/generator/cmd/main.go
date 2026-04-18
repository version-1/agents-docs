package main

import (
	"flag"
	"fmt"
	"os"

	"generator/internal/app"
	"generator/internal/domain"
	"generator/internal/infra"
)

func main() {
	var (
		inRoot  = flag.String("input", "", "input root directory to scan (e.g., ./.)")
		outRoot = flag.String("output", "out/.codex/skills", "output directory (e.g., .out/codex/skills)")
		mode    = flag.String("mode", "codex", "output mode: codex (nested skills) or claude (flat skills)")
	)
	flag.Parse()

	var skills domain.SkillGenerator
	switch *mode {
	case "codex":
		skills = domain.NewPathRespectSkillGenerator()
	case "claude":
		skills = domain.NewFlatSkillGenerator()
	default:
		fmt.Fprintf(os.Stderr, "error: unknown mode %q (use \"codex\" or \"claude\")\n", *mode)
		os.Exit(1)
	}

	gen := app.NewGenerator(infra.OSFS{}, skills)
	if err := gen.Run(*inRoot, *outRoot); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
