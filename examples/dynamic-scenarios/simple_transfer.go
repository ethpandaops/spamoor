// Example dynamic scenario: Minimal example
//
// This file demonstrates how to write a dynamic scenario that can be loaded
// at runtime using Yaegi without recompiling spamoor.
//
// Usage:
//
//	spamoor --scenario-file examples/dynamic-scenarios/simple_transfer.go dyn-transfer -h localhost:8545
package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
)

// ScenarioOptions defines the configuration for this scenario.
type ScenarioOptions struct {
	TotalCount uint64 `yaml:"total_count"`
	Throughput uint64 `yaml:"throughput"`
	MaxWallets uint64 `yaml:"max_wallets"`
}

// Default option values
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount: 10,
	Throughput: 5,
	MaxWallets: 5,
}

// Scenario implements the scenario.Scenario interface
type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool
}

// ScenarioDescriptor is the exported variable that spamoor looks for.
var ScenarioDescriptor = scenario.Descriptor{
	Name:           "dyn-transfer",
	Description:    "Dynamic scenario example",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options: ScenarioDefaultOptions,
		logger:  logger.WithField("scenario", "dyn-transfer"),
	}
}

// Flags registers CLI flags for this scenario.
func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of transactions")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Transactions per slot")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum wallets")
	return nil
}

// Init initializes the scenario.
func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	}

	s.logger.Infof("initialized: count=%d, throughput=%d", s.options.TotalCount, s.options.Throughput)
	return nil
}

// Run executes the scenario.
func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Info("starting dynamic scenario")

	for i := uint64(0); i < s.options.TotalCount; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(i))
			if wallet != nil {
				s.logger.Infof("tx %d: would send from wallet %s", i+1, wallet.GetAddress().Hex())
			} else {
				return fmt.Errorf("no wallet available for tx %d", i+1)
			}
		}
	}

	s.logger.Info("scenario complete")
	return nil
}

// main is a placeholder required for package main compilation.
// This file is loaded dynamically by Yaegi at runtime, not executed directly.
func main() {}
