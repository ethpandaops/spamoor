package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ethpandaops/spamoor/utils"
)

func main() {
	logger := logrus.StandardLogger()

	rootCmd := &cobra.Command{
		Use:   "spamoor-utils",
		Short: "Spamoor utilities",
		Long:  `A collection of utility commands for spamoor.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				logger.Error(err)
			}
		},
	}

	rootCmd.AddCommand(NewConvertEESTCmd(logger))
	rootCmd.AddCommand(NewValidatePluginCmd(logger))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().Bool("trace", false, "Trace output")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		trace, _ := cmd.Flags().GetBool("trace")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if trace {
			logrus.SetLevel(logrus.TraceLevel)
		} else if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}
	}

	logger.WithFields(logrus.Fields{
		"version":   utils.GetBuildVersion(),
		"buildtime": utils.BuildTime,
	}).Info("starting spamoor-utils")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
