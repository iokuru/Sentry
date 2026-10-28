package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/iokuru/sentry/internal/costguard"
)

func main() {
	planFile := flag.String("plan", "", "path to terraform plan json file (reads stdin if -)")
	dir := flag.String("dir", "infra", "path to terraform configuration directory")
	flag.Parse()

	var violations []costguard.Violation
	var err error

	if *planFile != "" {
		var planBytes []byte
		if *planFile == "-" {
			planBytes, err = io.ReadAll(os.Stdin)
		} else {
			planBytes, err = os.ReadFile(*planFile)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read plan: %v\n", err)
			os.Exit(1)
		}
		violations, err = costguard.ValidatePlanJSON(planBytes)
	} else {
		violations, err = costguard.ValidateSourceDirectory(*dir)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "cost guard check failed: %v\n", err)
		os.Exit(1)
	}

	if len(violations) > 0 {
		fmt.Fprintf(os.Stderr, "[COST VIOLATION] Found %d forbidden non-free resource(s):\n", len(violations))
		for _, v := range violations {
			fmt.Fprintf(os.Stderr, "  - %s: %s\n", v.Resource, v.Reason)
		}
		os.Exit(1)
	}

	fmt.Println("[COST GUARD] All resources comply with zero dollar free tier constraints.")
}
