package main

import (
	"fmt"
	"os"

	"github.com/ethpandaops/spamoor/plugin"
	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// ValidateCommand validates a plugin tar.gz file or local directory.
func ValidateCommand(args []string) {
	flags := pflag.NewFlagSet("validate-plugin", pflag.ExitOnError)
	verbose := flags.BoolP("verbose", "v", false, "Show verbose output")
	flags.Parse(args)

	if flags.NArg() < 1 {
		fmt.Println("Usage: spamoor validate-plugin [options] <plugin.tar.gz | plugin-dir>")
		fmt.Println()
		fmt.Println("Validates a plugin archive or directory by:")
		fmt.Println("  1. Loading it via Yaegi interpreter")
		fmt.Println("  2. Verifying PluginDescriptor fields")
		fmt.Println("  3. Checking registered scenarios")
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

	fmt.Printf("Validating plugin: %s\n", path)
	fmt.Println()

	// Initialize registries for plugin loading
	pluginRegistry, scenarioRegistry := scenarios.InitRegistries()

	// Load the plugin - detect source type
	l := plugin.NewPluginLoader(logger, pluginRegistry, scenarioRegistry)

	var loaded *plugin.LoadedPlugin
	var err error

	if isDirectory(path) {
		loaded, err = l.LoadFromLocalPath(path)
	} else {
		loaded, err = l.LoadFromFile(path)
	}

	if err != nil {
		fmt.Printf("✗ Failed to load plugin\n")
		fmt.Printf("  Error: %v\n", err)
		os.Exit(1)
	}

	// Access the descriptor from loaded plugin
	desc := loaded.Descriptor

	fmt.Printf("✓ Plugin loaded successfully\n")
	fmt.Println()

	// Validate descriptor fields
	hasError := false

	fmt.Println("Plugin Descriptor:")
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

	fmt.Printf("  Scenarios:   %d\n", len(desc.Scenarios))
	fmt.Println()

	fmt.Println("Plugin Metadata (from plugin.yaml):")
	if loaded.Metadata != nil {
		fmt.Printf("  Name:        %s\n", loaded.Metadata.Name)
		fmt.Printf("  BuildTime:   %s\n", loaded.Metadata.BuildTime)
		fmt.Printf("  GitVersion:  %s\n", loaded.Metadata.GitVersion)
	} else {
		fmt.Printf("  (no metadata available)\n")
	}
	fmt.Println()

	fmt.Println("Loading Details:")
	fmt.Printf("  TempDir:     %s\n", loaded.TempDir)
	fmt.Printf("  PluginPath:  %s\n", loaded.PluginPath)
	fmt.Printf("  SourceType:  %s\n", loaded.SourceType)
	fmt.Println()

	// Check each scenario
	if len(desc.Scenarios) == 0 {
		fmt.Printf("  ⚠ No scenarios defined in plugin\n")
	} else {
		fmt.Println("Scenarios:")
		for i, scenarioDesc := range desc.Scenarios {
			fmt.Printf("  [%d] %s\n", i, scenarioDesc.Name)
			if scenarioDesc.Name == "" {
				fmt.Printf("      ⚠ Name is empty\n")
				hasError = true
			}
			if scenarioDesc.NewScenario == nil {
				fmt.Printf("      ✗ NewScenario is nil\n")
				hasError = true
			} else {
				fmt.Printf("      ✓ NewScenario defined\n")

				// Try to instantiate
				instance := scenarioDesc.NewScenario(logger)
				if instance == nil {
					fmt.Printf("      ✗ NewScenario returned nil\n")
					hasError = true
				} else {
					fmt.Printf("      ✓ Instance created successfully\n")

					// Count flags
					flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
					if err := instance.Flags(flagSet); err != nil {
						fmt.Printf("      ⚠ Flags() returned error: %v\n", err)
					} else {
						flagCount := 0
						flagSet.VisitAll(func(f *pflag.Flag) {
							flagCount++
						})
						fmt.Printf("      ✓ Flags registered: %d\n", flagCount)
					}
				}
			}
			if *verbose && scenarioDesc.Description != "" {
				fmt.Printf("      Description: %s\n", scenarioDesc.Description)
			}
			fmt.Println()
		}
	}

	// Final status
	if hasError {
		fmt.Printf("Result: ✗ Plugin has errors and may not work correctly\n")
		os.Exit(1)
	} else if len(desc.Scenarios) == 0 {
		fmt.Printf("Result: ⚠ Plugin loaded but has no scenarios\n")
	} else {
		fmt.Printf("Result: ✓ Plugin '%s' is valid with %d scenario(s)\n", desc.Name, len(desc.Scenarios))
	}

	// Cleanup temp directory
	if err := l.CleanupPlugin(loaded); err != nil {
		fmt.Printf("Warning: Failed to cleanup temp directory: %v\n", err)
	}
}
