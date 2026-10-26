package main

import (
	"deploy/internal/plugins"
	"flag"
	"fmt"
	"os"
)

func main() {
	groups := flag.String("groups", "", "group definition file")
	external := flag.String("external-skills", "", "external lock file")
	output := flag.String("output", "", "marketplace output directory")
	dry := flag.Bool("dry-run", false, "validate and print without writing")
	flag.Parse()
	sum, err := plugins.Build(plugins.Options{GroupsPath: *groups, ExternalPath: *external, Output: *output, DryRun: *dry})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("groups=%v skills=%v managed_plugins=%v targets=%v\n", sum.Groups, sum.Skills, sum.ManagedPlugins, sum.Targets)
}
