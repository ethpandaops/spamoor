package example

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/_example1/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

// ScenarioOptions defines the configuration options for the example scenario
type ScenarioOptions struct {
	TotalCount      uint64  `yaml:"total_count"`
	Throughput      uint64  `yaml:"throughput"`
	MaxPending      uint64  `yaml:"max_pending"`
	MaxWallets      uint64  `yaml:"max_wallets"`
	Rebroadcast     uint64  `yaml:"rebroadcast"`
	BaseFee         float64 `yaml:"base_fee"`
	TipFee          float64 `yaml:"tip_fee"`
	BaseFeeWei      string  `yaml:"base_fee_wei"`
	TipFeeWei       string  `yaml:"tip_fee_wei"`
	InitialValue    uint64  `yaml:"initial_value"`
	MaxIncrement    uint64  `yaml:"max_increment"`
	RandomizeValues bool    `yaml:"randomize_values"`
	Timeout         string  `yaml:"timeout"`
	ClientGroup     string  `yaml:"client_group"`
	LogTxs          bool    `yaml:"log_txs"`
}

// Scenario represents the example scenario implementation
type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	contractAddr common.Address
}

var ScenarioName = "_example1"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:      0,
	Throughput:      10,
	MaxPending:      0,
	MaxWallets:      0,
	Rebroadcast:     1,
	BaseFee:         20,
	TipFee:          2,
	InitialValue:    100,
	MaxIncrement:    10,
	RandomizeValues: true,
	Timeout:         "",
	ClientGroup:     "",
	LogTxs:          false,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Example scenario demonstrating contract deployment and bound transactions",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.Uint64Var(&s.options.InitialValue, "initial-value", ScenarioDefaultOptions.InitialValue, "Initial value for the storage contract")
	flags.Uint64Var(&s.options.MaxIncrement, "max-increment", ScenarioDefaultOptions.MaxIncrement, "Maximum increment value for random increments")
	flags.BoolVar(&s.options.RandomizeValues, "randomize-values", ScenarioDefaultOptions.RandomizeValues, "Use random values for contract interactions")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger)
		if err != nil {
			return err
		}
	}

	// Configure wallet pool based on transaction volume
	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		// Rule of thumb: 1 wallet per 50 transactions
		maxWallets := s.options.TotalCount / 50
		if maxWallets < 10 {
			maxWallets = 10
		} else if maxWallets > 1000 {
			maxWallets = 1000
		}
		s.walletPool.SetWalletCount(maxWallets)
	} else {
		// For throughput-based scenarios
		if s.options.Throughput*10 < 1000 {
			s.walletPool.SetWalletCount(s.options.Throughput * 10)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
	}

	// Add a well-known wallet for contract deployment
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  utils.EtherToWei(uint256.NewInt(10)), // 10 ETH
		RefillBalance: utils.EtherToWei(uint256.NewInt(5)),  // 5 ETH threshold
		VeryWellKnown: false,                                // Scenario-specific
	})

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Step 1: Deploy the SimpleStorage contract
	contractReceipt, err := s.deployContract(ctx)
	if err != nil {
		return fmt.Errorf("failed to deploy contract: %w", err)
	}

	s.contractAddr = contractReceipt.ContractAddress
	s.logger.Infof("deployed SimpleStorage contract at %s (block #%v)",
		s.contractAddr.Hex(), contractReceipt.BlockNumber.Uint64())

	// Step 2: Configure pending transaction limits
	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.options.Throughput * 10
		if maxPending == 0 {
			maxPending = 4000
		}
		// Don't exceed 10 pending per wallet
		if maxPending > s.walletPool.GetConfiguredWalletCount()*10 {
			maxPending = s.walletPool.GetConfiguredWalletCount() * 10
		}
	}

	// Step 3: Parse timeout if specified
	var timeout time.Duration
	if s.options.Timeout != "" {
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("timeout set to %v", timeout)
	}

	// Step 4: Run the transaction scenario using the helper function
	return scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: s.options.TotalCount,
		Throughput: s.options.Throughput,
		MaxPending: maxPending,
		Timeout:    timeout,
		WalletPool: s.walletPool,
		Logger:     s.logger,
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			logger := s.logger
			receiptChan, tx, client, wallet, err := s.sendNextTransaction(ctx, params.TxIdx)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}

			params.NotifySubmitted()
			params.OrderedLogCb(func() {
				if err != nil {
					logger.Warnf("could not send transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent tx #%6d: %v", params.TxIdx+1, tx.Hash().String())
				}
			})

			// wait for receipt
			if _, err := receiptChan.Wait(ctx); err != nil {
				return err
			}

			return err
		},
	})
}

func (s *Scenario) deployContract(ctx context.Context) (*types.Receipt, error) {
	// Use the well-known deployer wallet for consistent deployment address
	deployerWallet := s.walletPool.GetWellKnownWallet("deployer")
	if deployerWallet == nil {
		return nil, fmt.Errorf("deployer wallet not found")
	}

	// Get a client for deployment
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return nil, scenario.ErrNoClients
	}

	// Get suggested fees
	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	s.logger.Infof("deploying SimpleStorage contract with initial value %d", s.options.InitialValue)

	// Build deployment transaction using BuildBoundTx pattern
	deploymentTx, err := deployerWallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       2000000, // Sufficient gas for deployment
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		// Use abigen-generated deployment function
		_, deployTx, _, err := contract.DeploySimpleStorage(
			transactOpts,
			client.GetEthClient(),
			big.NewInt(int64(s.options.InitialValue)),
		)
		return deployTx, err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to build deployment transaction: %w", err)
	}

	// Submit deployment transaction and wait for confirmation
	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, deployerWallet, deploymentTx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true, // Enable rebroadcast for important deployment
	})

	if err != nil {
		return nil, fmt.Errorf("deployment transaction failed: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("deployment transaction reverted")
	}

	return receipt, nil
}

func (s *Scenario) sendNextTransaction(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	// Select wallet and client for this transaction
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
	if wallet == nil {
		return nil, nil, nil, nil, scenario.ErrNoWallet
	}

	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return nil, nil, client, wallet, scenario.ErrNoClients
	}

	if err := wallet.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, nil, client, wallet, err
	}

	// Get suggested fees
	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, fmt.Errorf("failed to get suggested fees: %w", err)
	}

	// Create contract instance for interaction
	storageContract, err := contract.NewSimpleStorage(s.contractAddr, client.GetEthClient())
	if err != nil {
		return nil, nil, client, wallet, fmt.Errorf("failed to create contract instance: %w", err)
	}

	// Determine which operation to perform (alternating between setValue and increment)
	var tx *types.Transaction
	var opName string
	if txIdx%3 == 0 {
		// setValue operation
		opName = "setValue"
		value := s.options.InitialValue
		if s.options.RandomizeValues {
			value = uint64(rand.Intn(1000)) + 1
		}

		tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       100000, // Sufficient gas for setValue call
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return storageContract.SetValue(transactOpts, big.NewInt(int64(value)))
		})
	} else if txIdx%3 == 1 {
		// increment operation
		opName = "increment"
		tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       100000, // Sufficient gas for increment call
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return storageContract.Increment(transactOpts)
		})
	} else {
		// incrementBy operation
		opName = "incrementBy"
		increment := uint64(1)
		if s.options.RandomizeValues && s.options.MaxIncrement > 0 {
			increment = uint64(rand.Intn(int(s.options.MaxIncrement))) + 1
		}

		tx, err = wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       100000, // Sufficient gas for incrementBy call
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return storageContract.IncrementBy(transactOpts, big.NewInt(int64(increment)))
		})
	}

	if err != nil {
		return nil, nil, client, wallet, fmt.Errorf("failed to build transaction: %w", err)
	}

	receiptChan := make(scenario.ReceiptChan, 1)

	// Submit the transaction
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			receiptChan <- receipt
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			if receipt.Status == types.ReceiptStatusSuccessful {
				// Log successful contract interaction
				s.logger.WithFields(logrus.Fields{
					"txHash":    tx.Hash().Hex(),
					"block":     receipt.BlockNumber.Uint64(),
					"gasUsed":   receipt.GasUsed,
					"contract":  s.contractAddr.Hex(),
					"operation": opName,
				}).Debug("contract interaction confirmed")
			}
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "contract", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		// mark nonce as skipped if tx was not sent
		wallet.MarkSkippedNonce(tx.Nonce())

		return nil, nil, client, wallet, fmt.Errorf("failed to send transaction: %w", err)
	}

	return receiptChan, tx, client, wallet, nil
}
