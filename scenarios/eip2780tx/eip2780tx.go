package eip2780tx

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/eip2780tx/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

// mode constants for EIP-2780 gas path exercise
const (
	modeNewAccount   = "new-account"
	modePrecompile   = "precompile"
	modeAccessList   = "access-list"
	modeContractCall = "contract-call"
	modeValueCall    = "value-call"
	modeValueCallNew = "value-call-new"
	modeSelfdestruct = "selfdestruct"
	modeAll          = "all"
)

// allModes lists the individual modes cycled through by "all"
var allModes = []string{
	modeNewAccount,
	modePrecompile,
	modeAccessList,
	modeContractCall,
	modeValueCall,
	modeValueCallNew,
	modeSelfdestruct,
}

// default gas limits per mode (used when --gaslimit is 0)
// These defaults are set high enough to work on both pre-EIP-2780 and EIP-2780 networks.
var defaultGasLimits = map[string]uint64{
	modeNewAccount:   50000,
	modePrecompile:   30000,
	modeAccessList:   30000,
	modeContractCall: 50000,
	modeValueCall:    60000,
	modeValueCallNew: 80000,
	modeSelfdestruct: 200000,
}

type ScenarioOptions struct {
	TotalCount        uint64  `yaml:"total_count"`
	Throughput        uint64  `yaml:"throughput"`
	MaxPending        uint64  `yaml:"max_pending"`
	MaxWallets        uint64  `yaml:"max_wallets"`
	Rebroadcast       uint64  `yaml:"rebroadcast"`
	BaseFee           float64 `yaml:"base_fee"`
	TipFee            float64 `yaml:"tip_fee"`
	BaseFeeWei        string  `yaml:"base_fee_wei"`
	TipFeeWei         string  `yaml:"tip_fee_wei"`
	Amount            uint64  `yaml:"amount"`
	GasLimit          uint64  `yaml:"gas_limit"`
	RandomAmount      bool    `yaml:"random_amount"`
	Mode              string  `yaml:"mode"`
	Timeout           string  `yaml:"timeout"`
	ClientGroup       string  `yaml:"client_group"`
	DeployClientGroup string  `yaml:"deploy_client_group"`
	LogTxs            bool    `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	contractAddr common.Address
	needContract bool
}

var ScenarioName = "eip2780tx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        10,
	MaxPending:        0,
	MaxWallets:        0,
	Rebroadcast:       1,
	BaseFee:           20,
	TipFee:            2,
	Amount:            20,
	GasLimit:          0,
	RandomAmount:      false,
	Mode:              modeAll,
	Timeout:           "",
	ClientGroup:       "",
	DeployClientGroup: "",
	LogTxs:            false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Send transactions exercising EIP-2780 gas accounting paths",
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
	flags.Uint64Var(&s.options.Amount, "amount", ScenarioDefaultOptions.Amount, "Transfer amount per transaction (in gwei)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit per transaction (0 = use mode-specific defaults)")
	flags.BoolVar(&s.options.RandomAmount, "random-amount", ScenarioDefaultOptions.RandomAmount, "Use random amounts for transactions (with --amount as limit)")
	flags.StringVar(&s.options.Mode, "mode", ScenarioDefaultOptions.Mode, "EIP-2780 mode: new-account, precompile, access-list, contract-call, value-call, value-call-new, selfdestruct, all")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.DeployClientGroup, "deploy-client-group", ScenarioDefaultOptions.DeployClientGroup, "Client group to use for contract deployment")
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

	// Validate mode
	switch s.options.Mode {
	case modeNewAccount, modePrecompile, modeAccessList,
		modeContractCall, modeValueCall, modeValueCallNew,
		modeSelfdestruct, modeAll:
	default:
		return fmt.Errorf("unknown mode %q: must be one of new-account, precompile, access-list, contract-call, value-call, value-call-new, selfdestruct, all", s.options.Mode)
	}

	// Determine if contract deployment is needed
	s.needContract = modeNeedsContract(s.options.Mode)

	// Configure wallet count
	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		maxWallets := s.options.TotalCount / 50
		if maxWallets < 10 {
			maxWallets = 10
		} else if maxWallets > 1000 {
			maxWallets = 1000
		}
		s.walletPool.SetWalletCount(maxWallets)
	} else {
		if s.options.Throughput*10 < 1000 {
			s.walletPool.SetWalletCount(s.options.Throughput * 10)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
	}

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	if s.options.GasLimit > utils.MaxGasLimitPerTx {
		s.logger.Warnf("Gas limit %d exceeds %d and will most likely be dropped by the execution layer client", s.options.GasLimit, utils.MaxGasLimitPerTx)
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s (mode: %s)", ScenarioName, s.options.Mode)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Deploy contract if any mode needs it
	if s.needContract {
		contractReceipt, _, err := s.sendDeploymentTx(ctx)
		if err != nil {
			s.logger.Errorf("could not deploy EIP2780Helper contract: %v", err)
			return err
		}
		if contractReceipt == nil {
			return fmt.Errorf("could not deploy EIP2780Helper contract")
		}
		s.contractAddr = contractReceipt.ContractAddress
		s.logger.Infof("deployed EIP2780Helper contract: %v (confirmed in block #%v)", s.contractAddr.String(), contractReceipt.BlockNumber.String())
	}

	// Configure max pending
	maxPending := s.options.MaxPending
	if maxPending == 0 {
		maxPending = s.options.Throughput * 10
		if maxPending == 0 {
			maxPending = 4000
		}
		if maxPending > s.walletPool.GetConfiguredWalletCount()*10 {
			maxPending = s.walletPool.GetConfiguredWalletCount() * 10
		}
	}

	// Parse timeout
	var timeout time.Duration
	if s.options.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
	}

	err := scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount:                  s.options.TotalCount,
		Throughput:                  s.options.Throughput,
		MaxPending:                  maxPending,
		ThroughputIncrementInterval: 0,
		Timeout:                     timeout,
		WalletPool:                  s.walletPool,

		Logger: s.logger,
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			logger := s.logger
			receiptChan, tx, client, wallet, err := s.sendTx(ctx, params.TxIdx)
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

			if _, err := receiptChan.Wait(ctx); err != nil {
				return err
			}

			return err
		},
	})

	return err
}

func (s *Scenario) sendDeploymentTx(ctx context.Context) (*types.Receipt, *spamoor.Client, error) {
	deployClientGroup := s.options.DeployClientGroup
	if deployClientGroup == "" {
		deployClientGroup = s.options.ClientGroup
	}

	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(deployClientGroup),
	)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	if client == nil {
		return nil, nil, scenario.ErrNoClients
	}
	if wallet == nil {
		return nil, nil, scenario.ErrNoWallet
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, client, err
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       2000000,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		_, deployTx, _, err := contract.DeployEIP2780Helper(transactOpts, client.GetEthClient())
		return deployTx, err
	})
	if err != nil {
		return nil, nil, err
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: deployClientGroup,
		Rebroadcast: true,
	})
	if err != nil {
		return nil, client, err
	}

	return receipt, client, nil
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))

	if client == nil {
		return nil, nil, client, wallet, scenario.ErrNoClients
	}
	if wallet == nil {
		return nil, nil, client, wallet, scenario.ErrNoWallet
	}

	if err := wallet.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, nil, client, wallet, err
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	// Resolve mode for this transaction
	mode := s.resolveMode(txIdx)

	// Build mode-specific transaction
	tx, err := s.buildModeTx(ctx, mode, txIdx, wallet, client, feeCap, tipCap)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	receiptChan := make(scenario.ReceiptChan, 1)
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			receiptChan <- receipt
		},
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt) {
			txFees := utils.GetTransactionFees(tx, receipt)
			s.logger.WithField("rpc", client.GetName()).Debugf(
				" transaction %d [%s] confirmed in block #%v. gas: %v, total fee: %v gwei (base: %v)",
				txIdx+1,
				mode,
				receipt.BlockNumber.String(),
				receipt.GasUsed,
				txFees.TotalFeeGweiString(),
				txFees.TxBaseFeeGweiString(),
			)
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())
		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}

// resolveMode returns the concrete mode for the given txIdx.
// In "all" mode, cycles through allModes round-robin.
func (s *Scenario) resolveMode(txIdx uint64) string {
	if s.options.Mode != modeAll {
		return s.options.Mode
	}
	return allModes[txIdx%uint64(len(allModes))]
}

// gasLimitForMode returns the gas limit, respecting user override via --gaslimit.
func (s *Scenario) gasLimitForMode(mode string) uint64 {
	if s.options.GasLimit > 0 {
		return s.options.GasLimit
	}
	if gl, ok := defaultGasLimits[mode]; ok {
		return gl
	}
	return 50000
}

// resolveAmount calculates the transfer amount (in wei), optionally randomized.
func (s *Scenario) resolveAmount() *uint256.Int {
	amount := new(uint256.Int).Mul(uint256.NewInt(s.options.Amount), uint256.NewInt(1000000000))
	if s.options.RandomAmount {
		n, err := rand.Int(rand.Reader, amount.ToBig())
		if err == nil {
			amount = uint256.MustFromBig(n)
		}
	}
	return amount
}

// modeNeedsContract returns true if the mode (or any mode in "all") requires the helper contract.
func modeNeedsContract(mode string) bool {
	switch mode {
	case modeContractCall, modeValueCall, modeValueCallNew, modeSelfdestruct, modeAll:
		return true
	default:
		return false
	}
}

// buildModeTx dispatches to the correct transaction builder based on mode.
func (s *Scenario) buildModeTx(
	ctx context.Context,
	mode string,
	txIdx uint64,
	wallet *spamoor.Wallet,
	client *spamoor.Client,
	feeCap, tipCap *big.Int,
) (*types.Transaction, error) {
	gasLimit := s.gasLimitForMode(mode)
	amount := s.resolveAmount()

	switch mode {
	case modeNewAccount:
		return s.buildNewAccountTx(wallet, feeCap, tipCap, gasLimit, amount)
	case modePrecompile:
		return s.buildPrecompileTx(wallet, txIdx, feeCap, tipCap, gasLimit, amount)
	case modeAccessList:
		return s.buildAccessListTx(wallet, txIdx, feeCap, tipCap, gasLimit, amount)
	case modeContractCall:
		return s.buildContractCallTx(ctx, wallet, client, feeCap, tipCap, gasLimit)
	case modeValueCall:
		return s.buildValueCallTx(ctx, wallet, client, txIdx, feeCap, tipCap, gasLimit, amount)
	case modeValueCallNew:
		return s.buildValueCallNewTx(ctx, wallet, client, feeCap, tipCap, gasLimit, amount)
	case modeSelfdestruct:
		return s.buildSelfdestructTx(ctx, wallet, client, feeCap, tipCap, gasLimit, amount)
	default:
		return nil, fmt.Errorf("unknown mode: %s", mode)
	}
}

// buildNewAccountTx sends value to a random fresh address (exercises GAS_NEW_ACCOUNT surcharge).
func (s *Scenario) buildNewAccountTx(
	wallet *spamoor.Wallet,
	feeCap, tipCap *big.Int,
	gasLimit uint64,
	amount *uint256.Int,
) (*types.Transaction, error) {
	addrBytes := make([]byte, 20)
	rand.Read(addrBytes)
	toAddr := common.BytesToAddress(addrBytes)

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     amount,
	})
	if err != nil {
		return nil, err
	}
	return wallet.BuildDynamicFeeTx(txData)
}

// buildPrecompileTx sends value to precompile addresses 0x01..0x09 (warm at tx start).
func (s *Scenario) buildPrecompileTx(
	wallet *spamoor.Wallet,
	txIdx uint64,
	feeCap, tipCap *big.Int,
	gasLimit uint64,
	amount *uint256.Int,
) (*types.Transaction, error) {
	// Cycle through precompiles 0x01 to 0x09
	precompileIdx := (txIdx % 9) + 1
	toAddr := common.BigToAddress(big.NewInt(int64(precompileIdx)))

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     amount,
	})
	if err != nil {
		return nil, err
	}
	return wallet.BuildDynamicFeeTx(txData)
}

// buildAccessListTx sends value to a child wallet with recipient in access list (WARM vs COLD).
func (s *Scenario) buildAccessListTx(
	wallet *spamoor.Wallet,
	txIdx uint64,
	feeCap, tipCap *big.Int,
	gasLimit uint64,
	amount *uint256.Int,
) (*types.Transaction, error) {
	// Send to another child wallet, pre-warming it via access list
	targetWallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx)+1)
	toAddr := targetWallet.GetAddress()

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		To:        &toAddr,
		Value:     amount,
		AccessList: types.AccessList{
			{Address: toAddr},
		},
	})
	if err != nil {
		return nil, err
	}
	return wallet.BuildDynamicFeeTx(txData)
}

// buildContractCallTx calls contract.Nop() with no value (exercises COLD_ACCOUNT_COST_CODE).
func (s *Scenario) buildContractCallTx(
	ctx context.Context,
	wallet *spamoor.Wallet,
	client *spamoor.Client,
	feeCap, tipCap *big.Int,
	gasLimit uint64,
) (*types.Transaction, error) {
	helperContract, err := contract.NewEIP2780Helper(s.contractAddr, client.GetEthClient())
	if err != nil {
		return nil, err
	}

	return wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return helperContract.Nop(transactOpts)
	})
}

// buildValueCallTx calls contract.ForwardValue(existingWallet) (exercises CALL_VALUE_COST).
func (s *Scenario) buildValueCallTx(
	ctx context.Context,
	wallet *spamoor.Wallet,
	client *spamoor.Client,
	txIdx uint64,
	feeCap, tipCap *big.Int,
	gasLimit uint64,
	amount *uint256.Int,
) (*types.Transaction, error) {
	helperContract, err := contract.NewEIP2780Helper(s.contractAddr, client.GetEthClient())
	if err != nil {
		return nil, err
	}

	// Forward to an existing child wallet
	targetWallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx)+1)
	targetAddr := targetWallet.GetAddress()

	return wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		Value:     uint256.MustFromBig(amount.ToBig()),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return helperContract.ForwardValue(transactOpts, targetAddr)
	})
}

// buildValueCallNewTx calls contract.ForwardValue(freshAddress) (exercises CALL_VALUE_COST + GAS_NEW_ACCOUNT).
func (s *Scenario) buildValueCallNewTx(
	ctx context.Context,
	wallet *spamoor.Wallet,
	client *spamoor.Client,
	feeCap, tipCap *big.Int,
	gasLimit uint64,
	amount *uint256.Int,
) (*types.Transaction, error) {
	helperContract, err := contract.NewEIP2780Helper(s.contractAddr, client.GetEthClient())
	if err != nil {
		return nil, err
	}

	// Forward to a random fresh address
	addrBytes := make([]byte, 20)
	rand.Read(addrBytes)
	targetAddr := common.BytesToAddress(addrBytes)

	return wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		Value:     uint256.MustFromBig(amount.ToBig()),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return helperContract.ForwardValue(transactOpts, targetAddr)
	})
}

// buildSelfdestructTx calls contract.CreateAndDestroy() (exercises EIP-6780 create+selfdestruct).
func (s *Scenario) buildSelfdestructTx(
	ctx context.Context,
	wallet *spamoor.Wallet,
	client *spamoor.Client,
	feeCap, tipCap *big.Int,
	gasLimit uint64,
	amount *uint256.Int,
) (*types.Transaction, error) {
	helperContract, err := contract.NewEIP2780Helper(s.contractAddr, client.GetEthClient())
	if err != nil {
		return nil, err
	}

	return wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		Value:     uint256.MustFromBig(amount.ToBig()),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return helperContract.CreateAndDestroy(transactOpts)
	})
}
