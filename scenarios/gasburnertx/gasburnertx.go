package gasburnertx

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	geas "github.com/fjl/geas/asm"
	"gopkg.in/yaml.v3"

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
	TotalCount     uint64 `yaml:"total_count"`
	Throughput     uint64 `yaml:"throughput"`
	MaxPending     uint64 `yaml:"max_pending"`
	MaxWallets     uint64 `yaml:"max_wallets"`
	Rebroadcast    uint64 `yaml:"rebroadcast"`
	BaseFee        uint64 `yaml:"base_fee"`
	TipFee         uint64 `yaml:"tip_fee"`
	GasUnitsToBurn uint64 `yaml:"gas_units_to_burn"`
	GasRemainder   uint64 `yaml:"gas_remainder"`
	Timeout        string `yaml:"timeout"`
	OpcodesEas     string `yaml:"opcodes"`
	InitOpcodesEas string `yaml:"init_opcodes"`
	ClientGroup    string `yaml:"client_group"`
	LogTxs         bool   `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	gasBurnerContractAddr common.Address

	pendingWGroup sync.WaitGroup
}

var ScenarioName = "gasburnertx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:     0,
	Throughput:     10,
	MaxPending:     0,
	MaxWallets:     0,
	Rebroadcast:    1,
	BaseFee:        20,
	TipFee:         2,
	GasUnitsToBurn: 2000000,
	GasRemainder:   10000,
	Timeout:        "",
	OpcodesEas:     "",
	InitOpcodesEas: "",
	ClientGroup:    "",
	LogTxs:         false,
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
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in gasburner transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in gasburner transactions (in gwei)")
	flags.Uint64Var(&s.options.GasUnitsToBurn, "gas-units-to-burn", ScenarioDefaultOptions.GasUnitsToBurn, "The number of gas units for each tx to cost")
	flags.Uint64Var(&s.options.GasRemainder, "gas-remainder", ScenarioDefaultOptions.GasRemainder, "The minimum number of gas units that must be left to do another round of the gasburner loop")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.OpcodesEas, "opcodes", "", "EAS opcodes to use for burning gas in the gasburner contract")
	flags.StringVar(&s.options.InitOpcodesEas, "init-opcodes", "", "EAS opcodes to use for the init code of the gasburner contract")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")
	return nil
}

func (s *Scenario) Init(options *scenario.Options) error {
	s.walletPool = options.WalletPool

	if options.Config != "" {
		err := yaml.Unmarshal([]byte(options.Config), &s.options)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.TotalCount > 0 {
		if s.options.TotalCount < 1000 {
			s.walletPool.SetWalletCount(s.options.TotalCount)
		} else {
			s.walletPool.SetWalletCount(1000)
		}
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
		ProcessNextTxFn: func(ctx context.Context, txIdx uint64, onComplete func()) (func(), error) {
			logger := s.logger
			tx, client, wallet, err := s.sendTx(ctx, txIdx, onComplete)
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}

			return func() {
				if err != nil {
					logger.Warnf("could not send transaction: %v", err)
				} else if s.options.LogTxs {
					logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				} else {
					logger.Debugf("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
				}
			}, err
		},
	})

	s.pendingWGroup.Wait()
	s.logger.Infof("finished sending transactions, awaiting block inclusion...")

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
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	if client == nil {
		return nil, client, fmt.Errorf("no client available")
	}

	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return nil, client, err
		}
	}

	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
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

	var txReceipt *types.Receipt
	var txErr error
	txWg := sync.WaitGroup{}
	txWg.Add(1)

	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				txWg.Done()
			}()

			txErr = err
			txReceipt = receipt
		},
	})
	if err != nil {
		return nil, client, err
	}

	txWg.Wait()
	if txErr != nil {
		return nil, client, txErr
	}
	if txReceipt == nil {
		return nil, client, fmt.Errorf("deployment transaction receipt is nil")
	}
	return txReceipt, client, nil
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, int(txIdx), s.options.ClientGroup)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
	transactionSubmitted := false

	defer func() {
		if !transactionSubmitted {
			onComplete()
		}
	}()

	if client == nil {
		return nil, client, wallet, fmt.Errorf("no client available")
	}

	var feeCap *big.Int
	var tipCap *big.Int

	if s.options.BaseFee > 0 {
		feeCap = new(big.Int).Mul(big.NewInt(int64(s.options.BaseFee)), big.NewInt(1000000000))
	}
	if s.options.TipFee > 0 {
		tipCap = new(big.Int).Mul(big.NewInt(int64(s.options.TipFee)), big.NewInt(1000000000))
	}

	if feeCap == nil || tipCap == nil {
		var err error
		feeCap, tipCap, err = client.GetSuggestedFee(s.walletPool.GetContext())
		if err != nil {
			return nil, client, wallet, err
		}
	}

	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
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
			gasLimit = 30000000
		} else if gasLimit == 0 {
			// Final fallback to a reasonable default if no block gas limit is available
			gasLimit = 30000000
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
		return nil, nil, wallet, err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, nil, wallet, err
	}

	s.pendingWGroup.Add(1)
	transactionSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer func() {
				onComplete()
				s.pendingWGroup.Done()
			}()

			if err != nil {
				s.logger.WithField("rpc", client.GetName()).Warnf("tx %6d: await receipt failed: %v", txIdx+1, err)
				return
			}
			if receipt == nil {
				return
			}

			txFees := utils.GetTransactionFees(tx, receipt)
			s.logger.WithField("rpc", client.GetName()).Debugf(" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v", txIdx+1, receipt.BlockNumber.String(), txFees.TotalFeeGwei(), txFees.TxBaseFeeGwei(), len(receipt.Logs))
		},
		LogFn: func(client *spamoor.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName()).WithField("nonce", tx.Nonce())
			if retry == 0 && rebroadcast > 0 {
				logger.Infof("rebroadcasting tx %6d", txIdx+1)
			}
			if retry > 0 {
				logger = logger.WithField("retry", retry)
			}
			if rebroadcast > 0 {
				logger = logger.WithField("rebroadcast", rebroadcast)
			}
			if err != nil {
				logger.Debugf("failed sending tx %6d: %v", txIdx+1, err)
			} else if retry > 0 || rebroadcast > 0 {
				logger.Debugf("successfully sent tx %6d", txIdx+1)
			}
		},
	})
	if err != nil {
		// reset nonce if tx was not sent
		wallet.ResetPendingNonce(ctx, client)

		return nil, client, wallet, err
	}

	return tx, client, wallet, nil
}
