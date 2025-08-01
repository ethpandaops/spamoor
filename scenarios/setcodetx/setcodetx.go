package setcodetx

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount        uint64  `yaml:"total_count"`
	Throughput        uint64  `yaml:"throughput"`
	MaxPending        uint64  `yaml:"max_pending"`
	MaxWallets        uint64  `yaml:"max_wallets"`
	MinAuthorizations uint64  `yaml:"min_authorizations"`
	MaxAuthorizations uint64  `yaml:"max_authorizations"`
	MaxDelegators     uint64  `yaml:"max_delegators"`
	Rebroadcast       uint64  `yaml:"rebroadcast"`
	BaseFee           float64 `yaml:"base_fee"`
	TipFee            float64 `yaml:"tip_fee"`
	GasLimit          uint64  `yaml:"gas_limit"`
	Amount            uint64  `yaml:"amount"`
	Data              string  `yaml:"data"`
	CodeAddr          string  `yaml:"code_addr"`
	RandomAmount      bool    `yaml:"random_amount"`
	RandomTarget      bool    `yaml:"random_target"`
	RandomCodeAddr    bool    `yaml:"random_code_addr"`
	Timeout           string  `yaml:"timeout"`
	ClientGroup       string  `yaml:"client_group"`
	LogTxs            bool    `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	delegatorSeed []byte
	delegators    []*spamoor.Wallet
}

var ScenarioName = "setcodetx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        100,
	MaxPending:        0,
	MaxWallets:        0,
	MinAuthorizations: 1,
	MaxAuthorizations: 10,
	MaxDelegators:     0,
	Rebroadcast:       1,
	BaseFee:           20,
	TipFee:            2,
	GasLimit:          200000,
	Amount:            20,
	Data:              "",
	CodeAddr:          "",
	RandomAmount:      false,
	RandomTarget:      false,
	RandomCodeAddr:    false,
	Timeout:           "",
	ClientGroup:       "",
	LogTxs:            false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Send setcode transactions with different configurations",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of transfer transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of transfer transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.MinAuthorizations, "min-authorizations", ScenarioDefaultOptions.MinAuthorizations, "Minimum number of authorizations to send per transaction")
	flags.Uint64Var(&s.options.MaxAuthorizations, "max-authorizations", ScenarioDefaultOptions.MaxAuthorizations, "Maximum number of authorizations to send per transaction")
	flags.Uint64Var(&s.options.MaxDelegators, "max-delegators", ScenarioDefaultOptions.MaxDelegators, "Maximum number of random delegators to use (0 = no delegator gets reused)")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transfer transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit to use in transactions")
	flags.Uint64Var(&s.options.Amount, "amount", ScenarioDefaultOptions.Amount, "Transfer amount per transaction (in gwei)")
	flags.StringVar(&s.options.Data, "data", ScenarioDefaultOptions.Data, "Transaction call data to send")
	flags.StringVar(&s.options.CodeAddr, "code-addr", ScenarioDefaultOptions.CodeAddr, "Code delegation target address to use for transactions")
	flags.BoolVar(&s.options.RandomAmount, "random-amount", ScenarioDefaultOptions.RandomAmount, "Use random amounts for transactions (with --amount as limit)")
	flags.BoolVar(&s.options.RandomTarget, "random-target", ScenarioDefaultOptions.RandomTarget, "Use random to addresses for transactions")
	flags.BoolVar(&s.options.RandomCodeAddr, "random-code-addr", ScenarioDefaultOptions.RandomCodeAddr, "Use random delegation target for transactions")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
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

	if s.options.GasLimit > utils.MaxGasLimitPerTx {
		s.logger.Warnf("Gas limit %d exceeds %d and will most likely be dropped by the execution layer client", s.options.GasLimit, utils.MaxGasLimitPerTx)
	}

	s.delegatorSeed = make([]byte, 32)
	rand.Read(s.delegatorSeed)

	if s.options.MaxDelegators > 0 {
		s.delegators = make([]*spamoor.Wallet, 0, s.options.MaxDelegators)
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

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

	err := scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
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

	return err
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

	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, s.options.BaseFee, s.options.TipFee)
	if err != nil {
		return nil, client, wallet, err
	}

	amount := uint256.NewInt(s.options.Amount)
	amount = amount.Mul(amount, uint256.NewInt(1000000000))
	if s.options.RandomAmount {
		n, err := rand.Int(rand.Reader, amount.ToBig())
		if err == nil {
			amount = uint256.MustFromBig(n)
		}
	}

	toAddr := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx)+1).GetAddress()
	if s.options.RandomTarget {
		addrBytes := make([]byte, 20)
		rand.Read(addrBytes)
		toAddr = common.Address(addrBytes)
	}

	txCallData := []byte{}

	if s.options.Data != "" {
		dataBytes, err := txbuilder.ParseBlobRefsBytes(strings.Split(s.options.Data, ","), nil)
		if err != nil {
			return nil, nil, wallet, err
		}

		txCallData = dataBytes
	}

	txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		To:        &toAddr,
		Value:     amount,
		Data:      txCallData,
		AuthList:  s.buildSetCodeAuthorizations(txIdx),
	})
	if err != nil {
		return nil, nil, wallet, err
	}

	tx, err := wallet.BuildSetCodeTx(txData)
	if err != nil {
		return nil, nil, wallet, err
	}

	transactionSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			onComplete()
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
		// reset nonce if tx was not sent
		wallet.ResetPendingNonce(ctx, client)

		return nil, client, wallet, err
	}

	return tx, client, wallet, nil
}

func (s *Scenario) buildSetCodeAuthorizations(txIdx uint64) []types.SetCodeAuthorization {
	authorizations := []types.SetCodeAuthorization{}

	if s.options.MaxAuthorizations == 0 {
		return authorizations
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(s.options.MaxAuthorizations-s.options.MinAuthorizations+1)))
	authorizationCount := int(n.Int64()) + int(s.options.MinAuthorizations)

	for i := 0; i < authorizationCount; i++ {
		delegatorIndex := (txIdx * s.options.MaxAuthorizations) + uint64(i)
		if s.options.MaxDelegators > 0 {
			delegatorIndex = delegatorIndex % s.options.MaxDelegators
		}

		var delegator *spamoor.Wallet
		if s.options.MaxDelegators > 0 && len(s.delegators) > int(delegatorIndex) {
			delegator = s.delegators[delegatorIndex]
		} else {
			d, err := s.prepareDelegator(delegatorIndex)
			if err != nil {
				s.logger.Errorf("could not prepare delegator %v: %v", delegatorIndex, err)
				continue
			}

			delegator = d

			if s.options.MaxDelegators > 0 {
				s.delegators = append(s.delegators, delegator)
			}
		}

		var codeAddr common.Address
		if s.options.RandomCodeAddr {
			codeAddr = common.Address(make([]byte, 20))
			rand.Read(codeAddr[:])
		} else if s.options.CodeAddr != "" {
			codeAddr = common.HexToAddress(s.options.CodeAddr)
		} else {
			codeAddr = s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx)).GetAddress()
		}

		authorization := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(s.walletPool.GetChainId().Uint64()),
			Address: codeAddr,
			Nonce:   delegator.GetNextNonce(),
		}

		authorization, err := types.SignSetCode(delegator.GetPrivateKey(), authorization)
		if err != nil {
			s.logger.Errorf("could not sign set code authorization: %v", err)
			continue
		}

		authorizations = append(authorizations, authorization)
	}

	return authorizations
}

func (s *Scenario) prepareDelegator(delegatorIndex uint64) (*spamoor.Wallet, error) {
	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes, delegatorIndex)
	if s.options.MaxDelegators > 0 {
		seedBytes := []byte(s.delegatorSeed)
		idxBytes = append(idxBytes, seedBytes...)
	}
	childKey := sha256.Sum256(append(common.FromHex(s.walletPool.GetRootWallet().GetWallet().GetAddress().Hex()), idxBytes...))
	return spamoor.NewWallet(fmt.Sprintf("%x", childKey))
}
