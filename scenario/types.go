package scenario

import (
	"context"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Descriptor describes a scenario.
type Descriptor struct {
	Name           string
	Aliases        []string
	Description    string
	DefaultOptions any
	NewScenario    func(logger logrus.FieldLogger) Scenario
}

// Options contains the options for the scenario initialization.
type Options struct {
	WalletPool *spamoor.WalletPool
	Config     string
	GlobalCfg  map[string]any
	PluginPath string // Path to plugin resources (empty for native scenarios)
}

type Scenario interface {
	// Flags registers the scenario's flags with the given flag set.
	Flags(flags *pflag.FlagSet) error
	// Init initializes the scenario with the given options.
	Init(options *Options) error
	// Run runs the scenario.
	Run(ctx context.Context) error
}
