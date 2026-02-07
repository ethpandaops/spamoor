package gasburnertx

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	geas "github.com/fjl/geas/asm"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

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
	GasUnitsToBurn    uint64  `yaml:"gas_units_to_burn"`
	GasRemainder      uint64  `yaml:"gas_remainder"`
	Timeout           string  `yaml:"timeout"`
	OpcodesEas        string  `yaml:"opcodes"`
	InitOpcodesEas    string  `yaml:"init_opcodes"`
	ClientGroup       string  `yaml:"client_group"`
	DeployClientGroup string  `yaml:"deploy_client_group"`
	LogTxs            bool    `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	gasBurnerContractAddr common.Address
}

var ScenarioName = "gasburnertx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        10,
	MaxPending:        0,
	MaxWallets:        0,
	Rebroadcast:       1,
	BaseFee:           20,
	TipFee:            2,
	GasUnitsToBurn:    2000000,
	GasRemainder:      10000,
	Timeout:           "",
	OpcodesEas:        "",
	InitOpcodesEas:    "",
	ClientGroup:       "",
	DeployClientGroup: "",
	LogTxs:            false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Send gasburner transactions with different configurations",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of gasburner transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of gasburner transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in gasburner transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in gasburner transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.Uint64Var(&s.options.GasUnitsToBurn, "gas-units-to-burn", ScenarioDefaultOptions.GasUnitsToBurn, "The number of gas units for each tx to cost")
	flags.Uint64Var(&s.options.GasRemainder, "gas-remainder", ScenarioDefaultOptions.GasRemainder, "The minimum number of gas units that must be left to do another round of the gasburner loop")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.OpcodesEas, "opcodes", "", "EAS opcodes to use for burning gas in the gasburner contract")
	flags.StringVar(&s.options.InitOpcodesEas, "init-opcodes", "", "EAS opcodes to use for the init code of the gasburner contract")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.DeployClientGroup, "deploy-client-group", ScenarioDefaultOptions.DeployClientGroup, "Client group to use for deployments")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		// Use the generalized config validation and parsing helper
		err := scenario.ParseAndValidateConfig(&ScenarioDescriptor, options.Config, &s.options, s.logger)
		if err != nil {
			return err
		}
	}

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

	if s.options.GasUnitsToBurn > utils.MaxGasLimitPerTx {
		s.logger.Warnf("Gas units to burn %d exceeds %d and will most likely be dropped by the execution layer client", s.options.GasUnitsToBurn, utils.MaxGasLimitPerTx)
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// deploy gas burner contract
	receipt, _, err := s.sendDeploymentTx(ctx, s.trimGeasOpcodes(s.options.OpcodesEas))
	if err != nil {
		return err
	}

	s.gasBurnerContractAddr = receipt.ContractAddress

	s.logger.Infof("deployed gas burner contract at %v", s.gasBurnerContractAddr.String())

	// send transactions
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

	// Parse timeout duration
	var timeout time.Duration
	if s.options.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout format '%s': %w", s.options.Timeout, err)
		}
	}

	err = scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
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

			// wait for receipt
			if _, err := receiptChan.Wait(ctx); err != nil {
				return err
			}

			return err
		},
	})

	return err
}

func (s *Scenario) trimGeasOpcodes(opcodesGeas string) string {
	if strings.Contains(opcodesGeas, "\n") {
		return opcodesGeas
	}

	opcodesGeas = strings.ReplaceAll(opcodesGeas, ";", "\n")
	return opcodesGeas
}

func (s *Scenario) sendDeploymentTx(ctx context.Context, opcodesGeas string) (*types.Receipt, *spamoor.Client, error) {
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
		return nil, client, scenario.ErrNoClients
	}

	if wallet == nil {
		return nil, client, scenario.ErrNoWallet
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, client, err
	}

	// build the worker code
	initcodeGeas := `
	;; Init code
	push @.start
	codesize
	sub
	dup1
	push @.start
	push0
	codecopy
	push0
	return
	
	.start:
	`
	defaultOpcodesGeas := `
    push 0x1337
	pop
    `
	contractInitOpcodesGeas := `
	push 0                ;; [custom]
	`
	contractGeasTpl := `
	%s
	gas                   ;; [gas, custom]
	push 0                ;; [loop_counter, gas, custom]
	jump @loop

	exit:
		push 0            ;; [0, loop_counter, gas, custom]
        mstore            ;; [gas, custom]
        push 32           ;; [32, gas, custom]
        push 0            ;; [0, 32, gas, custom]
		log1              ;; [custom]
        stop              ;; [custom]

	loop:
		push %d           ;; [gas_remainder, loop_counter, gas, custom]
		gas               ;; [gas, 10000, loop_counter, gas, custom]
		lt                ;; [gas < 10000, loop_counter, gas, custom]
		jumpi @exit       ;; [loop_counter, gas, custom]

		;; increase loop_counter
		push 1            ;; [1, loop_counter, gas, custom]
		add               ;; [loop_counter+1, gas, custom]

		;; dummy opcodes to burn gas
		%s

		jump @loop
	`

	if s.options.InitOpcodesEas != "" {
		contractInitOpcodesGeas = s.trimGeasOpcodes(s.options.InitOpcodesEas)
	}

	compiler := geas.NewCompiler(nil)

	initcode := compiler.CompileString(initcodeGeas)
	if initcode == nil {
		return nil, client, fmt.Errorf("failed to compile initcode")
	}

	var workerCodeBytes []byte

	if len(opcodesGeas) > 0 && strings.HasPrefix(opcodesGeas, "0x") {
		// opcodes in bytecode format
		contractCode := compiler.CompileString(fmt.Sprintf(contractGeasTpl, contractInitOpcodesGeas, s.options.GasRemainder, defaultOpcodesGeas))
		if contractCode == nil {
			return nil, client, fmt.Errorf("failed to compile template contract code")
		}

		defaultOpcodes := compiler.CompileString(defaultOpcodesGeas)
		if defaultOpcodes == nil {
			return nil, client, fmt.Errorf("failed to compile default opcodes")
		}

		// replace default opcodes with provided opcodes
		contractCodeHex := strings.Replace(hex.EncodeToString(contractCode), hex.EncodeToString(defaultOpcodes), strings.ReplaceAll(opcodesGeas, "0x", ""), 1)
		contractCodeBytes, err := hex.DecodeString(contractCodeHex)
		if err != nil {
			return nil, client, fmt.Errorf("failed to decode contract code: %w", err)
		}

		workerCodeBytes = contractCodeBytes
	} else if len(opcodesGeas) > 0 {
		// opcodes in geas format
		workerCodeBytes = compiler.CompileString(fmt.Sprintf(contractGeasTpl, contractInitOpcodesGeas, s.options.GasRemainder, opcodesGeas))
		if workerCodeBytes == nil {
			return nil, client, fmt.Errorf("failed to compile template contract code")
		}
	} else {
		workerCodeBytes = compiler.CompileString(fmt.Sprintf(contractGeasTpl, contractInitOpcodesGeas, s.options.GasRemainder, defaultOpcodesGeas))
		if workerCodeBytes == nil {
			return nil, client, fmt.Errorf("failed to compile default contract code")
		}
	}

	workerCodeBytes = append(initcode, workerCodeBytes...)

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       2000000,
		To:        nil,
		Value:     uint256.NewInt(0),
		Data:      workerCodeBytes,
	})
	if err != nil {
		return nil, client, err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, client, err
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: deployClientGroup,
		Rebroadcast: true,
	})
	if err != nil {
		return nil, client, err
	}

	if receipt == nil {
		return nil, client, fmt.Errorf("deployment transaction receipt is nil")
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

	txIdBytes := make([]byte, 4)
	txIdBytes[0] = byte(txIdx >> 24)
	txIdBytes[1] = byte(txIdx >> 16)
	txIdBytes[2] = byte(txIdx >> 8)
	txIdBytes[3] = byte(txIdx)

	// Determine gas limit: use block gas limit if GasUnitsToBurn is 0
	gasLimit := s.options.GasUnitsToBurn
	if gasLimit == 0 {
		var err error
		gasLimit, err = s.walletPool.GetTxPool().GetCurrentGasLimitWithInit()
		if err != nil {
			s.logger.Warnf("tx %6d: failed to fetch current gas limit: %v, using fallback", txIdx+1, err)
			gasLimit = utils.MaxGasLimitPerTx
		} else if gasLimit == 0 {
			// Final fallback to a reasonable default if no block gas limit is available
			gasLimit = utils.MaxGasLimitPerTx
			s.logger.Warnf("tx %6d: no gas limit available, using fallback %v", txIdx+1, gasLimit)
		} else {
			s.logger.Debugf("tx %6d: using block gas limit %v", txIdx+1, gasLimit)
		}
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		To:        &s.gasBurnerContractAddr,
		Value:     uint256.NewInt(0),
		Data:      txIdBytes,
	})
	if err != nil {
		return nil, nil, client, wallet, err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
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
				" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v",
				txIdx+1,
				receipt.BlockNumber.String(),
				txFees.TotalFeeGweiString(),
				txFees.TxBaseFeeGweiString(),
				len(receipt.Logs),
			)
		},
		LogFn: spamoor.GetDefaultLogFn(s.logger, "", fmt.Sprintf("%6d", txIdx+1), tx),
	})
	if err != nil {
		// mark nonce as skipped if tx was not sent
		wallet.MarkSkippedNonce(tx.Nonce())

		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}
