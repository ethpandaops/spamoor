package geastx

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	geas "github.com/fjl/geas/asm"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenariotypes"
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
	Amount         uint64 `yaml:"amount"`
	BaseFee        uint64 `yaml:"base_fee"`
	TipFee         uint64 `yaml:"tip_fee"`
	GasLimit       uint64 `yaml:"gas_limit"`
	DeployGasLimit uint64 `yaml:"deploy_gas_limit"`
	GeasFile       string `yaml:"geas_file"`
	GeasCode       string `yaml:"geas_code"`
	ClientGroup    string `yaml:"client_group"`
	Timeout        string `yaml:"timeout"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	geasContractAddr common.Address

	pendingWGroup sync.WaitGroup
}

var ScenarioName = "geastx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:     0,
	Throughput:     100,
	MaxPending:     0,
	MaxWallets:     0,
	Rebroadcast:    1,
	Amount:         0,
	BaseFee:        20,
	TipFee:         2,
	GasLimit:       1000000,
	DeployGasLimit: 1000000,
	GeasFile:       "",
	GeasCode:       "",
	ClientGroup:    "",
	Timeout:        "",
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Send transactions that execute custom geas code with different configurations",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of geas transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of geas transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Uint64Var(&s.options.Amount, "amount", ScenarioDefaultOptions.Amount, "Amount to send in geas transactions")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in geas transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in geas transactions (in gwei)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Max gas limit to use in geas transactions")
	flags.Uint64Var(&s.options.DeployGasLimit, "deploy-gaslimit", ScenarioDefaultOptions.GasLimit, "Max gas limit to use in deployment transaction")
	flags.StringVar(&s.options.GeasFile, "geasfile", "", "Path to the geas file to use for execution")
	flags.StringVar(&s.options.GeasCode, "geascode", "", "Geas code to use for execution")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
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

	geasCode := s.options.GeasCode
	if geasCode == "" && s.options.GeasFile != "" {
		if strings.HasPrefix(s.options.GeasFile, "https://") || strings.HasPrefix(s.options.GeasFile, "http://") {
			resp, err := http.Get(s.options.GeasFile)
			if err != nil {
				return fmt.Errorf("failed to download geas file: %w", err)
			}
			defer resp.Body.Close()
			geasBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read geas file response: %w", err)
			}
			geasCode = string(geasBytes)
		} else {
			_, err := os.Stat(s.options.GeasFile)
			if err != nil {
				return fmt.Errorf("failed to check if geas file exists: %w", err)
			}
			geasBytes, err := os.ReadFile(s.options.GeasFile)
			if err != nil {
				return fmt.Errorf("failed to read geas file: %w", err)
			}
			geasCode = string(geasBytes)
		}
	}

	if geasCode == "" {
		return fmt.Errorf("no geas code or file provided")
	}

	receipt, _, err := s.sendDeploymentTx(ctx, s.trimGeasOpcodes(geasCode))
	if err != nil {
		return err
	}

	s.geasContractAddr = receipt.ContractAddress

	s.logger.Infof("deployed contract with geas code at %v", s.geasContractAddr.String())

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

	err = utils.RunTransactionScenario(ctx, utils.TransactionScenarioOptions{
		TotalCount:                  s.options.TotalCount,
		Throughput:                  s.options.Throughput,
		MaxPending:                  maxPending,
		ThroughputIncrementInterval: 0,
		Timeout:                     timeout,

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
				} else {
					logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
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
	compiler := geas.NewCompiler(nil)

	initcode := compiler.CompileString(initcodeGeas)
	if initcode == nil {
		return nil, client, fmt.Errorf("failed to compile initcode")
	}

	var deployData []byte

	if len(opcodesGeas) > 0 && strings.HasPrefix(opcodesGeas, "0x") {
		// raw bytecode format
		contractCodeBytes, err := hex.DecodeString(strings.ReplaceAll(opcodesGeas, "0x", ""))
		if err != nil {
			return nil, client, fmt.Errorf("failed to decode contract code: %w", err)
		}

		deployData = contractCodeBytes
	} else {
		// opcodes in geas format
		deployData = compiler.CompileString(opcodesGeas)
		if deployData == nil {
			return nil, client, fmt.Errorf("failed to compile template contract code")
		}
	}

	deployData = append(initcode, deployData...)

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.DeployGasLimit,
		To:        nil,
		Value:     uint256.NewInt(0),
		Data:      deployData,
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
		return nil, client, err
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

	amount := uint256.NewInt(s.options.Amount)
	amount = amount.Mul(amount, uint256.NewInt(1000000000))

	txIdBytes := make([]byte, 4)
	txIdBytes[0] = byte(txIdx >> 24)
	txIdBytes[1] = byte(txIdx >> 16)
	txIdBytes[2] = byte(txIdx >> 8)
	txIdBytes[3] = byte(txIdx)

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		To:        &s.geasContractAddr,
		Value:     amount,
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

			effectiveGasPrice := receipt.EffectiveGasPrice
			if effectiveGasPrice == nil {
				effectiveGasPrice = big.NewInt(0)
			}
			feeAmount := new(big.Int).Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
			totalAmount := new(big.Int).Add(tx.Value(), feeAmount)
			wallet.SubBalance(totalAmount)

			gweiTotalFee := new(big.Int).Div(feeAmount, big.NewInt(1000000000))
			gweiBaseFee := new(big.Int).Div(effectiveGasPrice, big.NewInt(1000000000))

			s.logger.WithField("rpc", client.GetName()).Debugf(" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) logs: %v", txIdx+1, receipt.BlockNumber.String(), gweiTotalFee, gweiBaseFee, len(receipt.Logs))
		},
		LogFn: func(client *spamoor.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName())
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
