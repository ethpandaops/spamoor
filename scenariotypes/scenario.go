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

type Scenario interface {
	Flags(flags *pflag.FlagSet) error
	Init(walletPool *spamoor.WalletPool, config string) error
	Config() string
	Run(ctx context.Context) error
}
