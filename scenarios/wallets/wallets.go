package wallets

import (
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/tester"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	Wallets uint64
}

type Scenario struct {
	options ScenarioOptions
	logger  *logrus.Entry
	wallets uint64
}

func NewScenario() scenariotypes.Scenario {
	return &Scenario{
		logger: logrus.WithField("scenario", "wallets"),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.Wallets, "max-wallets", "w", 0, "Maximum number of child wallets to use")
	return nil
}

func (s *Scenario) Init(testerCfg *tester.TesterConfig) error {
	if s.options.Wallets > 0 {
		testerCfg.WalletCount = s.options.Wallets
	} else {
		testerCfg.WalletCount = 1000
	}
	s.wallets = testerCfg.WalletCount
	return nil
}

func (s *Scenario) Run(t *tester.Tester) error {
	wallet := t.GetRootWallet()
	s.logger.Infof("Root Wallet  %v  nonce: %6d  balance: %v ETH", wallet.GetAddress().String(), wallet.GetNonce(), utils.WeiToEther(uint256.MustFromBig(wallet.GetBalance())))
	client := t.GetClient(tester.SelectByIndex, 0)

	for i := 0; i < int(s.wallets); i++ {
		wallet := t.GetWallet(tester.SelectByIndex, i)
		pendingNonce, _ := client.GetPendingNonceAt(wallet.GetAddress())

		s.logger.Infof("Child Wallet %4d  %v  nonce: %6d (%6d)  balance: %v ETH", i+1, wallet.GetAddress().String(), wallet.GetNonce(), pendingNonce, utils.WeiToEther(uint256.MustFromBig(wallet.GetBalance())))
	}

	return nil
}
