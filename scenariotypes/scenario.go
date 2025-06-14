package scenariotypes

import (
	"context"

	"github.com/ethpandaops/spamoor/spamoortypes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// ScenarioDescriptor describes a scenario.
type ScenarioDescriptor struct {
	Name           string
	Description    string
	DefaultOptions interface{}
	NewScenario    func(logger logrus.FieldLogger) Scenario
}

// ScenarioOptions contains the options for the scenario initialization.
type ScenarioOptions struct {
	WalletPool spamoortypes.WalletPool
	Config     string
	GlobalCfg  map[string]interface{}
}

type Scenario interface {
	// Flags registers the scenario's flags with the given flag set.
	Flags(flags *pflag.FlagSet) error
	// Init initializes the scenario with the given options.
	Init(options *ScenarioOptions) error
	// Run runs the scenario.
	Run(ctx context.Context) error
}
