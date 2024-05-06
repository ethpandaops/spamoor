package scenariotypes

import (
	"github.com/ethpandaops/spamoor/tester"
	"github.com/spf13/pflag"
)

type Scenario interface {
	Flags(flags *pflag.FlagSet) error
	Init(testerCfg *tester.TesterConfig) error
	Run(tester *tester.Tester) error
}
