package txfuzz

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
)

type ScenarioOptions struct {
	TotalCount  uint64  `yaml:"total_count"`
	Throughput  uint64  `yaml:"throughput"`
	MaxPending  uint64  `yaml:"max_pending"`
	MaxWallets  uint64  `yaml:"max_wallets"`
	Rebroadcast uint64  `yaml:"rebroadcast"`
	BaseFee     float64 `yaml:"base_fee"`
	TipFee      float64 `yaml:"tip_fee"`
	BaseFeeWei  string  `yaml:"base_fee_wei"`
	TipFeeWei   string  `yaml:"tip_fee_wei"`
	GasLimit    uint64  `yaml:"gas_limit"`
	Timeout     string  `yaml:"timeout"`
	ClientGroup string  `yaml:"client_group"`
	LogTxs      bool    `yaml:"log_txs"`

	// Fuzzing specific options
	TxTypes       string `yaml:"tx_types"`        // comma list: legacy,accesslist,dynfee,blob,setcode (or "all")
	PayloadSeed   string `yaml:"payload_seed"`    // optional hex seed for reproducible fuzzing
	TxIdOffset    uint64 `yaml:"tx_id_offset"`    // start fuzzing from a specific txID
	MaxCallData   uint64 `yaml:"max_call_data"`   // maximum calldata/initcode size in bytes
	MaxAccessList uint64 `yaml:"max_access_list"` // maximum access list entries / storage keys
	MaxAuthList   uint64 `yaml:"max_auth_list"`   // maximum EIP-7702 authorizations per tx
	MaxBlobs      uint64 `yaml:"max_blobs"`       // maximum blob sidecars per blob tx
}

type Scenario struct {
	options    ScenarioOptions
	logger     logrus.FieldLogger
	walletPool *spamoor.WalletPool
	fuzzer     *fuzzer
	seed       string
}

var ScenarioName = "tx-fuzz"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:  0,
	Throughput:  50,
	MaxPending:  100,
	MaxWallets:  0,
	Rebroadcast: 30,
	BaseFee:     20,
	TipFee:      2,
	GasLimit:    500000,
	Timeout:     "",
	ClientGroup: "",
	LogTxs:      false,

	TxTypes:       "all",
	PayloadSeed:   "",
	TxIdOffset:    0,
	MaxCallData:   1024,
	MaxAccessList: 5,
	MaxAuthList:   5,
	MaxBlobs:      3,
}

var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Fuzzes the transaction layer by sending well-formed txs across all types (legacy/2930/1559/4844/7702) with randomized calldata, access lists, authorizations, blobs and targets",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of fuzzed transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of fuzzed transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast with unlimited retries and exponential backoff")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit to use in transactions")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.BoolVar(&s.options.LogTxs, "log-txs", ScenarioDefaultOptions.LogTxs, "Log all submitted transactions")

	flags.StringVar(&s.options.TxTypes, "tx-types", ScenarioDefaultOptions.TxTypes, "Comma-separated tx types to fuzz: legacy,accesslist,dynfee,blob,setcode (or 'all')")
	flags.StringVar(&s.options.PayloadSeed, "payload-seed", ScenarioDefaultOptions.PayloadSeed, "Custom hex seed for reproducible fuzzing (e.g. 0x1234abcd, empty means random)")
	flags.Uint64Var(&s.options.TxIdOffset, "tx-id-offset", ScenarioDefaultOptions.TxIdOffset, "Start fuzzing from a specific transaction ID")
	flags.Uint64Var(&s.options.MaxCallData, "max-call-data", ScenarioDefaultOptions.MaxCallData, "Maximum calldata/initcode size in bytes")
	flags.Uint64Var(&s.options.MaxAccessList, "max-access-list", ScenarioDefaultOptions.MaxAccessList, "Maximum access list entries and storage keys per entry")
	flags.Uint64Var(&s.options.MaxAuthList, "max-auth-list", ScenarioDefaultOptions.MaxAuthList, "Maximum EIP-7702 authorizations per setcode tx")
	flags.Uint64Var(&s.options.MaxBlobs, "max-blobs", ScenarioDefaultOptions.MaxBlobs, "Maximum blob sidecars per blob tx")

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

	enabledKinds, err := parseTxTypes(s.options.TxTypes)
	if err != nil {
		return err
	}

	if blockLimit := s.walletPool.GetTxPool().GetCurrentGasLimit(); blockLimit > 0 && s.options.GasLimit > blockLimit {
		s.logger.Warnf("Gas limit %d exceeds block gas limit %d and will most likely be dropped by the execution layer client", s.options.GasLimit, blockLimit)
	}

	if s.options.PayloadSeed != "" {
		if err := s.validateSeed(s.options.PayloadSeed); err != nil {
			return fmt.Errorf("invalid payload seed: %v", err)
		}
	}

	s.fuzzer = &fuzzer{
		chainID:      s.walletPool.GetChainId().Uint64(),
		enabledKinds: enabledKinds,
		maxCallData:  int(s.options.MaxCallData),
		maxAccessLen: int(s.options.MaxAccessList),
		maxAuthList:  int(s.options.MaxAuthList),
		maxBlobs:     int(s.options.MaxBlobs),
		poolAddrs:    s.poolAddresses,
	}

	return nil
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// Generate seed once at scenario start if not provided
	s.seed = s.options.PayloadSeed
	if s.seed == "" {
		randomBytes := make([]byte, 32)
		rand.Read(randomBytes)
		s.seed = hex.EncodeToString(randomBytes)
		s.logger.Infof("Generated random seed for this run: 0x%s", s.seed)
	} else {
		s.logger.Infof("Using provided seed: %s", s.seed)
	}
	s.fuzzer.seed = []byte(s.seed)

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

	var timeout time.Duration
	var err error
	if s.options.Timeout != "" {
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
	}

	s.logger.WithFields(logrus.Fields{
		"total":      s.options.TotalCount,
		"throughput": s.options.Throughput,
		"maxPending": maxPending,
		"txTypes":    s.options.TxTypes,
	}).Info("Starting transaction-layer fuzzer scenario")

	err = scenario.RunTransactionScenario(ctx, scenario.TransactionScenarioOptions{
		TotalCount: s.options.TotalCount,
		Throughput: s.options.Throughput,
		MaxPending: maxPending,
		Timeout:    timeout,
		WalletPool: s.walletPool,
		Logger:     s.logger.(*logrus.Entry),
		ProcessNextTxFn: func(ctx context.Context, params *scenario.ProcessNextTxParams) error {
			logger := s.logger
			receiptChan, tx, client, wallet, kind, err := s.sendFuzzedTx(ctx, params.TxIdx)
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
					logger.Warnf("fuzz tx %6d.0 (%s) failed: %v", params.TxIdx+1, kind, err)
				} else if s.options.LogTxs {
					logger.Infof("fuzz tx %6d.0 (%s) sent:  %v (%d bytes)", params.TxIdx+1, kind, tx.Hash().String(), len(tx.Data()))
				} else {
					logger.Debugf("fuzz tx %6d.0 (%s) sent:  %v (%d bytes)", params.TxIdx+1, kind, tx.Hash().String(), len(tx.Data()))
				}
			})

			if receiptChan != nil {
				if _, err := receiptChan.Wait(ctx); err != nil {
					return err
				}
			}

			return err
		},
	})

	return err
}

func (s *Scenario) sendFuzzedTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, string, error) {
	ftx := s.fuzzer.generate(txIdx + s.options.TxIdOffset)

	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByPendingTxCount, int(txIdx))
	if wallet == nil {
		return nil, nil, nil, nil, ftx.kind.String(), fmt.Errorf("no wallet available")
	}

	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return nil, nil, nil, wallet, ftx.kind.String(), fmt.Errorf("no client available")
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, ftx.kind.String(), err
	}

	tx, err := s.buildTx(wallet, ftx, feeCap, tipCap)
	if err != nil {
		return nil, nil, client, wallet, ftx.kind.String(), err
	}

	receiptChan := make(scenario.ReceiptChan, 1)
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			receiptChan <- receipt
		},
	})
	if err != nil {
		return nil, tx, client, wallet, ftx.kind.String(), err
	}

	return receiptChan, tx, client, wallet, ftx.kind.String(), nil
}

// buildTx turns a fuzzed envelope into a concrete signed transaction of the
// matching type using the wallet's managed-nonce build methods.
func (s *Scenario) buildTx(wallet *spamoor.Wallet, ftx *fuzzedTx, feeCap, tipCap *big.Int) (*types.Transaction, error) {
	fee := uint256.MustFromBig(feeCap)
	tip := uint256.MustFromBig(tipCap)
	meta := &txbuilder.TxMetadata{
		GasFeeCap:  fee,
		GasTipCap:  tip,
		Gas:        s.options.GasLimit,
		To:         ftx.to,
		Value:      ftx.value,
		Data:       ftx.data,
		AccessList: ftx.accessList,
		AuthList:   ftx.authList,
	}

	switch ftx.kind {
	case kindLegacy:
		txData, err := txbuilder.LegacyTx(meta)
		if err != nil {
			return nil, err
		}
		return wallet.BuildLegacyTx(txData)
	case kindAccessList:
		txData, err := txbuilder.AccessListTx(meta)
		if err != nil {
			return nil, err
		}
		return wallet.BuildAccessListTx(txData)
	case kindDynFee:
		txData, err := txbuilder.DynFeeTx(meta)
		if err != nil {
			return nil, err
		}
		return wallet.BuildDynamicFeeTx(txData)
	case kindSetCode:
		txData, err := txbuilder.SetCodeTx(meta)
		if err != nil {
			return nil, err
		}
		return wallet.BuildSetCodeTx(txData)
	case kindBlob:
		// blob fee cap: reuse the (already bumped) gas fee cap. Blob gas is priced
		// separately and is far cheaper than this on test networks, so this keeps
		// blob txs comfortably includable without ballooning the tx cost.
		meta.BlobFeeCap = uint256.MustFromBig(feeCap)
		txData, err := txbuilder.BuildBlobTx(meta, ftx.blobRefs)
		if err != nil {
			return nil, err
		}
		return wallet.BuildBlobTx(txData)
	default:
		return nil, fmt.Errorf("unknown tx kind: %s", ftx.kind)
	}
}

func (s *Scenario) poolAddresses() []common.Address {
	wallets := s.walletPool.GetAllWallets()
	addrs := make([]common.Address, 0, len(wallets))
	for _, w := range wallets {
		if w != nil {
			addrs = append(addrs, w.GetAddress())
		}
	}
	return addrs
}

func (s *Scenario) validateSeed(seed string) error {
	cleanSeed := strings.TrimPrefix(seed, "0x")
	if _, err := hex.DecodeString(cleanSeed); err != nil {
		return fmt.Errorf("seed must be a valid hex string (with or without 0x prefix): %v", err)
	}
	return nil
}

func parseTxTypes(s string) ([]txKind, error) {
	if s == "" || s == "all" {
		return allKinds, nil
	}

	byName := map[string]txKind{
		"legacy":     kindLegacy,
		"accesslist": kindAccessList,
		"dynfee":     kindDynFee,
		"blob":       kindBlob,
		"setcode":    kindSetCode,
	}

	var kinds []txKind
	seen := map[txKind]bool{}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(strings.ToLower(part))
		if part == "" {
			continue
		}
		k, ok := byName[part]
		if !ok {
			return nil, fmt.Errorf("unknown tx type %q (valid: legacy, accesslist, dynfee, blob, setcode, all)", part)
		}
		if !seen[k] {
			kinds = append(kinds, k)
			seen[k] = true
		}
	}
	if len(kinds) == 0 {
		return nil, fmt.Errorf("no valid tx types selected")
	}
	return kinds, nil
}
