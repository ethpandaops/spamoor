package codeadopt

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount  uint64  `yaml:"total_count"`
	Throughput  uint64  `yaml:"throughput"`
	MaxPending  uint64  `yaml:"max_pending"`
	MaxWallets  uint64  `yaml:"max_wallets"`
	BaseFee     float64 `yaml:"base_fee"`
	TipFee      float64 `yaml:"tip_fee"`
	BaseFeeWei  string  `yaml:"base_fee_wei"`
	TipFeeWei   string  `yaml:"tip_fee_wei"`
	Timeout     string  `yaml:"timeout"`
	ClientGroup string  `yaml:"client_group"`
	LogTxs      bool    `yaml:"log_txs"`

	Cases        string `yaml:"cases"`
	FundAmount   uint64 `yaml:"fund_amount"`
	SendAllValue uint64 `yaml:"sendall_value"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool
	cases      []caseFn
	caseNames  []string
	support    spamoor.FrameSupport

	valueSource    common.Address
	approverSource common.Address
	adopter        common.Address
	staticCaller   common.Address
	fixtureBlock   uint64

	casesRun     atomic.Uint64
	checksPassed atomic.Uint64
	checksFailed atomic.Uint64
}

var ScenarioName = "codeadopt"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:   0,
	Throughput:   2,
	MaxPending:   0,
	MaxWallets:   0,
	BaseFee:      20,
	TipFee:       2,
	Cases:        "all",
	FundAmount:   10_000_000,
	SendAllValue: 1_000,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Exercise EIP-8298 SETCODEFROM, EIP-8151 code-restricted ecRecover and EIP-4758 SENDALL, including the EOA migration from ECDSA to frame transactions approved by adopted code, and check the resulting state",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenario.Scenario {
	return &Scenario{
		options: ScenarioDefaultOptions,
		logger:  logger.WithField("scenario", ScenarioName),
	}
}

func (s *Scenario) Flags(flags *pflag.FlagSet) error {
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of cases to run")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of cases to start per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of cases in flight")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", ScenarioDefaultOptions.BaseFeeWei, "Max fee per gas in wei (overrides --basefee)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", ScenarioDefaultOptions.TipFeeWei, "Max tip per gas in wei (overrides --tipfee)")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m')")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log every case result")
	flags.StringVar(&s.options.Cases, "cases", ScenarioDefaultOptions.Cases, "Comma-separated cases to run ("+strings.Join(allCaseNames(), ", ")+"), or 'all'")
	flags.Uint64Var(&s.options.FundAmount, "fund-amount", ScenarioDefaultOptions.FundAmount, "Amount in gwei given to each migrated account to pay for its frame transaction")
	flags.Uint64Var(&s.options.SendAllValue, "sendall-value", ScenarioDefaultOptions.SendAllValue, "Amount in gwei each SENDALL case moves")

	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		if err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger); err != nil {
			return err
		}
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else {
		s.walletPool.SetWalletCount(max(10, min(s.options.Throughput*10, 200)))
	}

	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "codeadopt-deployer",
		RefillAmount:  utils.EtherToWei(uint256.NewInt(5)),
		RefillBalance: utils.EtherToWei(uint256.NewInt(2)),
	})

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	cases, names, err := selectCases(s, s.options.Cases)
	if err != nil {
		return err
	}

	s.cases = cases
	s.caseNames = names

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s (cases: %s)", ScenarioName, strings.Join(s.caseNames, ", "))
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	support, err := s.walletPool.GetTxPool().AwaitFrameSupport(ctx, s.logger)
	if err != nil {
		return err
	}

	s.support = support

	if err := s.deployFixtures(ctx); err != nil {
		return err
	}

	if err := s.awaitActivation(ctx); err != nil {
		return err
	}

	if err := s.checkSources(ctx); err != nil {
		return err
	}

	var timeout time.Duration
	if s.options.Timeout != "" {
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
	}

	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = max(s.options.Throughput*4, 4)
	}

	defer s.logSummary()

	return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: s.options.TotalCount,
		Throughput: s.options.Throughput,
		MaxPending: maxPending,
		Timeout:    timeout,
		WalletPool: s.walletPool,
		Logger:     s.logger,
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			idx := params.TxIdx % uint64(len(s.cases))
			name := s.caseNames[idx]
			run := s.newCaseRun(name, params.TxIdx)

			params.NotifySubmitted()

			err := s.cases[idx](ctx, run)

			params.OrderedLogCb(func() {
				switch {
				case err != nil:
					run.logger.Warnf("case %s did not complete: %v", name, err)
				case run.failed > 0:
					run.logger.Warnf("case %s: %d of %d checks failed", name, run.failed, run.failed+run.passed)
				case s.options.LogTxs:
					run.logger.Infof("case %s: %d checks passed", name, run.passed)
				default:
					run.logger.Debugf("case %s: %d checks passed", name, run.passed)
				}
			})

			if n := s.casesRun.Add(1); n%50 == 0 {
				s.logSummary()
			}

			return nil
		},
	})
}

func (s *Scenario) logSummary() {
	s.logger.Infof("codeadopt checks: %d cases, %d checks passed, %d failed",
		s.casesRun.Load(), s.checksPassed.Load(), s.checksFailed.Load())
}

func (s *Scenario) client() *spamoor.Client {
	return s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientRandom),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
}

func (s *Scenario) deployFixtures(ctx context.Context) error {
	deployer := s.walletPool.GetWellKnownWallet("codeadopt-deployer")
	if deployer == nil {
		return fmt.Errorf("deployer wallet not found")
	}

	client := s.client()
	if client == nil {
		return scenario.ErrNoClients
	}

	fixtures := []struct {
		name    string
		code    []byte
		address *common.Address
	}{
		{"value source", deployCode(valueSourceCode), &s.valueSource},
		{"approver source", deployCode(approverSourceCode), &s.approverSource},
		{"adopter", deployCode(adopterCode), &s.adopter},
	}

	for _, fixture := range fixtures {
		receipt, err := s.send(ctx, client, deployer, nil, nil, fixture.code, nil)
		if err != nil {
			return fmt.Errorf("deploying %s: %w", fixture.name, err)
		}

		*fixture.address = receipt.ContractAddress
	}

	receipt, err := s.send(ctx, client, deployer, nil, nil, deployCode(staticCallerCode(s.adopter)), nil)
	if err != nil {
		return fmt.Errorf("deploying static caller: %w", err)
	}

	s.staticCaller = receipt.ContractAddress
	s.fixtureBlock = receipt.BlockNumber.Uint64()

	s.logger.Infof("fixtures: value source %s, approver source %s, adopter %s, static caller %s",
		s.valueSource, s.approverSource, s.adopter, s.staticCaller)

	return nil
}

func (s *Scenario) awaitActivation(ctx context.Context) error {
	waiting := false

	for {
		client := s.client()
		if client == nil {
			return scenario.ErrNoClients
		}

		out, err := s.call(ctx, client, s.adopter, common.LeftPadBytes(s.adopter.Bytes(), 32))
		if err == nil && word(out) == 1 && s.awaitBlock(ctx, client, s.fixtureBlock) == nil {
			if waiting {
				s.logger.Infof("SETCODEFROM is active")
			}

			return nil
		}

		if !waiting {
			waiting = true
			s.logger.Warnf("SETCODEFROM is not active yet, waiting for the fork")
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(12 * time.Second):
		}
	}
}
