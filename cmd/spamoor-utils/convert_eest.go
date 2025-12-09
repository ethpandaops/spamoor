package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/eestconv"
)

// NewConvertEESTCmd creates the convert-eest command
func NewConvertEESTCmd(logger logrus.FieldLogger) *cobra.Command {
	var outputFile string
	var verbosePayloads bool
	var testPattern string
	var excludePattern string

	cmd := &cobra.Command{
		Use:   "convert-eest <path>",
		Short: "Convert EEST fixtures to intermediate representation",
		Long: `Converts Ethereum Execution Spec Test (EEST) fixtures to an intermediate
representation that can be replayed on a normal network.

The intermediate representation includes:
- Unsigned transactions with placeholder addresses ($contract[1], $sender[1], etc.)
- Prerequisites (sender funding amounts)
- Post-execution checks (storage values, balances)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			return runConvertEEST(logger, inputPath, outputFile, testPattern, excludePattern, verbosePayloads)
		},
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().BoolVar(&verbosePayloads, "verbose-payloads", false, "Print each converted payload name")
	cmd.Flags().StringVar(&testPattern, "test-pattern", "", "Regex pattern to include tests by path/name")
	cmd.Flags().StringVar(&excludePattern, "exclude-pattern", "", "Regex pattern to exclude tests by path/name")

	return cmd
}

func runConvertEEST(logger logrus.FieldLogger, inputPath, outputFile, testPattern, excludePattern string, verbose bool) error {
	converter, err := eestconv.NewConverter(logger, eestconv.ConvertOptions{
		TestPattern:    testPattern,
		ExcludePattern: excludePattern,
		Verbose:        verbose,
	})
	if err != nil {
		return fmt.Errorf("failed to create converter: %w", err)
	}

	output, err := converter.ConvertDirectory(inputPath)
	if err != nil {
		return fmt.Errorf("failed to convert fixtures: %w", err)
	}

	yamlData, err := yaml.Marshal(output)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if outputFile != "" {
		err = os.WriteFile(outputFile, yamlData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		logger.WithField("file", outputFile).Info("wrote output file")
	} else {
		fmt.Println(string(yamlData))
	}

	return nil
}
