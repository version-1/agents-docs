package main

import (
	"flag"
	"fmt"
	"os"

	"deploy/internal/deploy"
)

func main() {
	var (
		configPath = flag.String("config", "", "deployment config file path")
		dryRun     = flag.Bool("dry-run", false, "print planned copies without writing files")
	)
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "error: -config is required")
		os.Exit(1)
	}

	runner := deploy.NewRunner(os.Stdout)
	if err := runner.Run(*configPath, deploy.Options{DryRun: *dryRun}); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
