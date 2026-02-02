package main

import (
	"fmt"
	"os"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/loader"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// ValidateCommand validates a dynamic scenario file.
func ValidateCommand(args []string) {
	flags := pflag.NewFlagSet("validate-scenario", pflag.ExitOnError)
	verbose := flags.BoolP("verbose", "v", false, "Show verbose output")
	flags.Parse(args)

	if flags.NArg() < 1 {
		fmt.Println("Usage: spamoor validate-scenario [options] <path>")
		fmt.Println()
		fmt.Println("Validates a dynamic scenario file by:")
		fmt.Println("  1. Loading it via Yaegi interpreter")
		fmt.Println("  2. Verifying ScenarioDescriptor fields")
		fmt.Println("  3. Instantiating the scenario")
		fmt.Println("  4. Checking registered flags")
		fmt.Println()
		fmt.Println("Options:")
		flags.PrintDefaults()
		os.Exit(1)
	}

	path := flags.Args()[0]

	// Configure logger
	logger := logrus.New()
	if *verbose {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.WarnLevel)
	}

	fmt.Printf("Validating scenario: %s\n", path)
	fmt.Println()

	// Load the scenario
	l := loader.NewScenarioLoader(logger)
	desc, err := l.LoadFromFile(path)
	if err != nil {
		fmt.Printf("✗ Failed to load scenario\n")
		fmt.Printf("  Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Scenario loaded successfully\n")
	fmt.Println()

	// Validate descriptor fields
	issues := validateDescriptor(desc)
	hasError := false

	fmt.Println("Descriptor:")
	fmt.Printf("  Name:        %s", desc.Name)
	if desc.Name == "" {
		fmt.Printf(" ⚠ (empty)")
		hasError = true
	}
	fmt.Println()

	fmt.Printf("  Description: %s", desc.Description)
	if desc.Description == "" {
		fmt.Printf(" ⚠ (empty)")
	}
	fmt.Println()

	if len(desc.Aliases) > 0 {
		fmt.Printf("  Aliases:     %v\n", desc.Aliases)
	}

	fmt.Printf("  NewScenario: ")
	if desc.NewScenario != nil {
		fmt.Printf("✓ defined\n")
	} else {
		fmt.Printf("✗ nil\n")
		hasError = true
	}

	fmt.Printf("  DefaultOpts: ")
	if desc.DefaultOptions != nil {
		fmt.Printf("✓ defined (%T)\n", desc.DefaultOptions)
	} else {
		fmt.Printf("⚠ nil (no default configuration)\n")
	}

	fmt.Println()

	// Try to instantiate
	if desc.NewScenario != nil {
		fmt.Println("Instantiation:")
		instance := desc.NewScenario(logger)
		if instance == nil {
			fmt.Printf("  ✗ NewScenario returned nil\n")
			hasError = true
		} else {
			fmt.Printf("  ✓ Instance created successfully\n")

			// Count flags
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			if err := instance.Flags(flagSet); err != nil {
				fmt.Printf("  ⚠ Flags() returned error: %v\n", err)
			} else {
				flagCount := 0
				flagSet.VisitAll(func(f *pflag.Flag) {
					flagCount++
				})
				fmt.Printf("  ✓ Flags registered: %d\n", flagCount)

				if *verbose && flagCount > 0 {
					fmt.Println()
					fmt.Println("  Available flags:")
					flagSet.VisitAll(func(f *pflag.Flag) {
						fmt.Printf("    --%s: %s (default: %s)\n", f.Name, f.Usage, f.DefValue)
					})
				}
			}
		}
		fmt.Println()
	}

	// Print warnings
	if len(issues) > 0 {
		fmt.Println("Warnings:")
		for _, issue := range issues {
			fmt.Printf("  ⚠ %s\n", issue)
		}
		fmt.Println()
	}

	// Final status
	if hasError {
		fmt.Printf("Result: ✗ Scenario has errors and may not work correctly\n")
		os.Exit(1)
	} else if len(issues) > 0 {
		fmt.Printf("Result: ⚠ Scenario loaded with warnings\n")
	} else {
		fmt.Printf("Result: ✓ Scenario '%s' is valid\n", desc.Name)
	}
}

// validateDescriptor checks for common issues in a scenario descriptor.
func validateDescriptor(desc *scenario.Descriptor) []string {
	var issues []string

	if desc.Name == "" {
		issues = append(issues, "Name is empty - scenario won't be findable by name")
	}

	if desc.Description == "" {
		issues = append(issues, "Description is empty - users won't know what this scenario does")
	}

	if desc.NewScenario == nil {
		issues = append(issues, "NewScenario is nil - scenario cannot be instantiated")
	}

	return issues
}
