package main

import (
	"fmt"
	"os"

	"github.com/ethpandaops/spamoor/plugin"
	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewValidatePluginCmd creates the validate-plugin command.
func NewValidatePluginCmd(logger logrus.FieldLogger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-plugin <plugin.tar.gz | plugin-dir>",
		Short: "Validate a plugin archive or directory",
		Long: `Validates a plugin archive or directory by:
  1. Loading it via Yaegi interpreter
  2. Verifying PluginDescriptor fields
  3. Checking registered scenarios`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidatePlugin(logger, args[0])
		},
	}

	return cmd
}

func runValidatePlugin(logger logrus.FieldLogger, path string) error {
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
		return fmt.Errorf("failed to load plugin: %w", err)
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
						flagSet.VisitAll(func(_ *pflag.Flag) {
							flagCount++
						})
						fmt.Printf("      ✓ Flags registered: %d\n", flagCount)
					}
				}
			}
			if scenarioDesc.Description != "" {
				fmt.Printf("      Description: %s\n", scenarioDesc.Description)
			}
			fmt.Println()
		}
	}

	// Final status
	if hasError {
		fmt.Printf("Result: ✗ Plugin has errors and may not work correctly\n")
		return fmt.Errorf("plugin has errors")
	} else if len(desc.Scenarios) == 0 {
		fmt.Printf("Result: ⚠ Plugin loaded but has no scenarios\n")
	} else {
		fmt.Printf("Result: ✓ Plugin '%s' is valid with %d scenario(s)\n", desc.Name, len(desc.Scenarios))
	}

	// Cleanup temp directory
	if err := l.CleanupPlugin(loaded); err != nil {
		fmt.Printf("Warning: Failed to cleanup temp directory: %v\n", err)
	}

	return nil
}

// isDirectory checks if the given path is an existing directory.
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
