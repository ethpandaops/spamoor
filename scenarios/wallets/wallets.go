package wallets

import (
	"context"
	"fmt"

	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	Wallets uint64 `yaml:"wallets"`
	Reclaim bool   `yaml:"reclaim"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool
}

var ScenarioName = "wallets"
var ScenarioDefaultOptions = ScenarioOptions{
	Wallets: 0,
	Reclaim: false,
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Show wallet balances",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		options: ScenarioDefaultOptions,
		logger:  logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.Wallets, "max-wallets", "w", ScenarioDefaultOptions.Wallets, "Maximum number of child wallets to use")
	flags.BoolVarP(&s.options.Reclaim, "reclaim", "r", ScenarioDefaultOptions.Reclaim, "Reclaim funds from wallets")
	return nil
}

func (s *Scenario) Init(options *scenariotypes.ScenarioOptions) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		err := yaml.Unmarshal([]byte(options.Config), &s.options)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	if s.options.Wallets > 0 {
		s.walletPool.SetWalletCount(s.options.Wallets)
	} else {
		s.walletPool.SetWalletCount(1000)
	}

	// skip funding for this scenario
	s.walletPool.SetRunFundings(false)

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	wallet := s.walletPool.GetRootWallet().GetWallet()
	s.logger.Infof("Root Wallet  %v  nonce: %6d  balance: %v ETH", wallet.GetAddress().String(), wallet.GetNonce(), utils.WeiToEther(uint256.MustFromBig(wallet.GetBalance())))
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, "")

	if client == nil {
		return fmt.Errorf("no client available")
	}

	if s.options.Reclaim {
		s.logger.Infof("Reclaiming funds from wallets")
		err := s.walletPool.ReclaimFunds(ctx, client)
		if err != nil {
			return err
		}
	}

	for i := 0; i < int(s.walletPool.GetWalletCount()); i++ {
		wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, i)
		pendingNonce, _ := client.GetPendingNonceAt(ctx, wallet.GetAddress())

		s.logger.Infof("Child Wallet %4d  %v  nonce: %6d (%6d)  balance: %v ETH", i+1, wallet.GetAddress().String(), wallet.GetNonce(), pendingNonce, utils.WeiToEther(uint256.MustFromBig(wallet.GetBalance())))
	}

	return nil
}
