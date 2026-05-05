package main

import (
	"flag"
	"fmt"
	"os"

	"deploy/internal/deploy"
)

func main() {
	var (
		configPath         = flag.String("config", "", "deployment config file path")
		externalSkillsPath = flag.String("external-skills", "", "external skills config file path")
		dryRun             = flag.Bool("dry-run", false, "print planned copies without writing files")
		noColor            = flag.Bool("no-color", false, "disable ANSI color output")
	)
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "error: -config is required")
		os.Exit(1)
	}

	runner := deploy.NewRunner(os.Stdout)
	if err := runner.Run(*configPath, deploy.Options{DryRun: *dryRun, NoColor: *noColor, ExternalSkillsPath: *externalSkillsPath}); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
