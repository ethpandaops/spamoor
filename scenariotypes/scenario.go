package scenariotypes

import (
	"context"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type ScenarioDescriptor struct {
	Name           string
	Description    string
	DefaultOptions interface{}
	NewScenario    func(logger logrus.FieldLogger) Scenario
}

type ScenarioOptions struct {
	WalletPool *spamoor.WalletPool
	Config     string
	GlobalCfg  map[string]interface{}
}

type Scenario interface {
	Flags(flags *pflag.FlagSet) error
	Init(options *ScenarioOptions) error
	Config() string
	Run(ctx context.Context) error
}
