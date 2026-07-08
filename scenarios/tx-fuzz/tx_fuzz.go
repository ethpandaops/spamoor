package txfuzz

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"strings"
	"sync/atomic"
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
	UnstuckTime uint64  `yaml:"unstuck_time"` // seconds to wait for a tx before replacing it to free the nonce (0 = disabled)

	// Fuzzing specific options
	TxTypes       string `yaml:"tx_types"`        // comma list: legacy,accesslist,dynfee,blob,setcode (or "all")
	PayloadSeed   string `yaml:"payload_seed"`    // optional hex seed for reproducible fuzzing
	TxIdOffset    uint64 `yaml:"tx_id_offset"`    // start fuzzing from a specific txID
	MaxCallData   uint64 `yaml:"max_call_data"`   // maximum calldata/initcode size in bytes
	MaxAccessList uint64 `yaml:"max_access_list"` // maximum access list entries / storage keys
	MaxAuthList   uint64 `yaml:"max_auth_list"`   // maximum EIP-7702 authorizations per tx
	MaxBlobs      uint64 `yaml:"max_blobs"`       // maximum blob sidecars per blob tx

	// Blob format (EIP-4844 / EIP-7594) options
	BlobV1Percent  uint64                   `yaml:"blob_v1_percent"` // % of blob txs sent with the v1 (cell-proof) wrapper after Fulu
	FuluActivation utils.FlexibleJsonUInt64 `yaml:"fulu_activation"` // unix timestamp of the Fulu activation
}

type Scenario struct {
	options    ScenarioOptions
	logger     logrus.FieldLogger
	walletPool *spamoor.WalletPool
	fuzzer     *fuzzer
	seed       string
}

// unstuckMaxTries bounds how many escalating-fee replacement attempts are made
// to clear a single stuck nonce before giving up and letting the pool's own
// rebroadcast/gap-fill take over.
const unstuckMaxTries = 5

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
	UnstuckTime: 60,

	TxTypes:       "all",
	PayloadSeed:   "",
	TxIdOffset:    0,
	MaxCallData:   1024,
	MaxAccessList: 5,
	MaxAuthList:   5,
	MaxBlobs:      3,

	BlobV1Percent:  100,
	FuluActivation: 0, // 0 = Fulu active since genesis -> send v1 (cell-proof) blobs by default
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
	flags.Uint64Var(&s.options.UnstuckTime, "unstuck-time", ScenarioDefaultOptions.UnstuckTime, "Seconds to wait for a fuzzed tx before replacing it with a cancel tx to free the nonce (0 disables)")

	flags.StringVar(&s.options.TxTypes, "tx-types", ScenarioDefaultOptions.TxTypes, "Comma-separated tx types to fuzz: legacy,accesslist,dynfee,blob,setcode (or 'all')")
	flags.StringVar(&s.options.PayloadSeed, "payload-seed", ScenarioDefaultOptions.PayloadSeed, "Custom hex seed for reproducible fuzzing (e.g. 0x1234abcd, empty means random)")
	flags.Uint64Var(&s.options.TxIdOffset, "tx-id-offset", ScenarioDefaultOptions.TxIdOffset, "Start fuzzing from a specific transaction ID")
	flags.Uint64Var(&s.options.MaxCallData, "max-call-data", ScenarioDefaultOptions.MaxCallData, "Maximum calldata/initcode size in bytes")
	flags.Uint64Var(&s.options.MaxAccessList, "max-access-list", ScenarioDefaultOptions.MaxAccessList, "Maximum access list entries and storage keys per entry")
	flags.Uint64Var(&s.options.MaxAuthList, "max-auth-list", ScenarioDefaultOptions.MaxAuthList, "Maximum EIP-7702 authorizations per setcode tx")
	flags.Uint64Var(&s.options.MaxBlobs, "max-blobs", ScenarioDefaultOptions.MaxBlobs, "Maximum blob sidecars per blob tx")
	flags.Uint64Var(&s.options.BlobV1Percent, "blob-v1-percent", ScenarioDefaultOptions.BlobV1Percent, "Percentage of blob transactions to send with the v1 (cell-proof) wrapper format after Fulu")
	flags.Uint64Var((*uint64)(&s.options.FuluActivation), "fulu-activation", uint64(ScenarioDefaultOptions.FuluActivation), "Unix timestamp of the Fulu activation")

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

	// Pick up the network-wide Fulu activation timestamp from the daemon config
	// (unless explicitly overridden) so blob txs use the correct sidecar format.
	if options.GlobalCfg != nil {
		if v, ok := options.GlobalCfg["fulu_activation"]; ok && s.options.FuluActivation == ScenarioDefaultOptions.FuluActivation {
			s.options.FuluActivation = utils.FlexibleJsonUInt64(v.(uint64))
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
				// Bound the wait: a fuzzed tx may be accepted into a mempool but
				// never mined (invalid execution, pool eviction, blob/non-blob
				// account reservation conflicts, ...). The unstuck goroutine
				// spawned in sendFuzzedTx replaces such a tx at its nonce so the
				// wallet keeps moving; this deadline is a backstop so the slot is
				// freed even if every replacement also fails to land.
				waitCtx := ctx
				if s.options.UnstuckTime > 0 {
					var cancel context.CancelFunc
					deadline := time.Duration(s.options.UnstuckTime*(unstuckMaxTries+2)) * time.Second
					waitCtx, cancel = context.WithTimeout(ctx, deadline)
					defer cancel()
				}
				if _, werr := receiptChan.Wait(waitCtx); werr != nil {
					if ctx.Err() != nil {
						return werr
					}
					// per-tx deadline hit (scenario still running): give up on this
					// tx and free the slot rather than block the whole scenario.
					logger.Warnf("fuzz tx %6d.0 (%s) not confirmed within unstuck deadline, freeing slot", params.TxIdx+1, kind)
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
	confirmed := &atomic.Bool{}
	sendOpts := &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			confirmed.Store(true)
			receiptChan <- receipt
		},
	}

	// Blob txs: after Fulu the network expects the v1 (EIP-7594 cell-proof)
	// sidecar wrapper. Convert in OnEncode (once, idempotent across rebroadcasts)
	// so we submit a correctly-formatted blob tx. Before Fulu the v0 KZG-proof
	// sidecar built by txbuilder is used as-is.
	if ftx.kind == kindBlob {
		blobV1Converted := false
		sendOpts.OnEncode = func(tx *types.Transaction) ([]byte, error) {
			sendAsV1 := uint64(time.Now().Unix()) > uint64(s.options.FuluActivation) &&
				mathrand.Intn(100) < int(s.options.BlobV1Percent)
			if sendAsV1 && !blobV1Converted {
				if sidecar := tx.BlobTxSidecar(); sidecar != nil {
					if err := sidecar.ToV1(); err != nil {
						return nil, err
					}
				}
				blobV1Converted = true
			}
			return nil, nil
		}
	}

	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, sendOpts)
	if err != nil {
		return nil, tx, client, wallet, ftx.kind.String(), err
	}

	// The tx submitted fine, but fuzzed txs frequently never mine (invalid
	// execution, pool eviction, blob/non-blob account reservation conflicts).
	// Watch it and, once it looks stuck, replace it at its nonce with a cheap
	// cancel tx so the wallet's nonce advances regardless of the original's fate.
	if s.options.UnstuckTime > 0 {
		go s.unstuckLoop(ctx, wallet, tx, confirmed)
	}

	return receiptChan, tx, client, wallet, ftx.kind.String(), nil
}

// unstuckLoop waits for the fuzzed tx to confirm; if it doesn't within
// UnstuckTime, it repeatedly submits a fee-bumped cancel transaction at the same
// nonce until either the nonce clears or the try budget is exhausted. This is the
// key liveness guarantee for the fuzzer: because a fuzzed tx may be accepted by a
// node yet never be includable, the only way to keep a wallet usable is to
// forcibly replace the stuck nonce rather than wait on an outcome that never comes.
func (s *Scenario) unstuckLoop(ctx context.Context, wallet *spamoor.Wallet, stuckTx *types.Transaction, confirmed *atomic.Bool) {
	nonce := stuckTx.Nonce()
	for try := 0; try < unstuckMaxTries; try++ {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(s.options.UnstuckTime) * time.Second):
		}

		// The original (or a prior replacement) confirmed, or the wallet already
		// mined past this nonce: nothing to unstuck.
		if confirmed.Load() || wallet.GetConfirmedNonce() > nonce {
			return
		}

		if err := s.sendCancelTx(ctx, wallet, stuckTx, try); err != nil {
			s.logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress())).
				Debugf("unstuck attempt %d for nonce %d failed: %v", try+1, nonce, err)
		} else {
			s.logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress())).
				Debugf("submitted unstuck cancel tx for nonce %d (attempt %d)", nonce, try+1)
		}
	}
}

// sendCancelTx replaces the stuck tx at its nonce with a minimal, definitely-
// mineable transaction carrying aggressively bumped fees. A blob tx can only be
// replaced by another blob tx (go-ethereum keeps blob and non-blob txs in
// separate subpools), so the cancel matches the stuck tx's pool class.
func (s *Scenario) sendCancelTx(ctx context.Context, wallet *spamoor.Wallet, stuckTx *types.Transaction, try int) error {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientRandom),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return fmt.Errorf("no client available")
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return err
	}

	// Escalate fees each attempt so replacements clear go-ethereum's price-bump
	// requirement (>=10% for regular txs, >=100% for blob fees) even against a
	// previous replacement.
	bump := big.NewInt(int64(1) << uint(try+1)) // 2x, 4x, 8x, ...
	feeCap = new(big.Int).Mul(feeCap, bump)
	tipCap = new(big.Int).Mul(tipCap, bump)

	nonce := stuckTx.Nonce()
	to := wallet.GetAddress()

	var cancelTx *types.Transaction
	if stuckTx.Type() == types.BlobTxType {
		meta := &txbuilder.TxMetadata{
			GasFeeCap:  uint256.MustFromBig(feeCap),
			GasTipCap:  uint256.MustFromBig(tipCap),
			BlobFeeCap: uint256.MustFromBig(feeCap),
			Gas:        21000,
			To:         &to,
			Value:      uint256.NewInt(0),
		}
		blobTx, berr := txbuilder.BuildBlobTx(meta, [][]string{{"random:full"}})
		if berr != nil {
			return berr
		}
		cancelTx, err = wallet.ReplaceBlobTx(blobTx, nonce)
	} else {
		txData := &types.DynamicFeeTx{
			GasFeeCap: feeCap,
			GasTipCap: tipCap,
			Gas:       21000,
			To:        &to,
			Value:     big.NewInt(0),
		}
		cancelTx, err = wallet.ReplaceDynamicFeeTx(txData, nonce)
	}
	if err != nil {
		return err
	}

	return s.walletPool.GetTxPool().SendTransaction(ctx, wallet, cancelTx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: true,
	})
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
