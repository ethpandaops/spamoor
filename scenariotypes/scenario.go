package scenariotypes

import (
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/spf13/pflag"
)

type Scenario interface {
	Flags(flags *pflag.FlagSet) error
	Init(walletPool *spamoor.WalletPool) error
	Run() error
}
