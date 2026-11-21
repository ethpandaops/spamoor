package bloatedeoa

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/txtypes"
	"github.com/ethpandaops/spamoor/utils"
)

const (
	ModeSload  = "sload"
	ModeSstore = "sstore"

	// maxTxGasLimit is the EIP-7825 per-transaction gas limit cap.
	maxTxGasLimit = 1 << 24

	// sstoreGasPerSlot is a conservative upper bound on the regular gas a single
	// cold SSTORE in the attack loop costs: 2,100 cold slot access + 20,000
	// first-write to a fresh slot + ~400 loop overhead. Existing-slot rewrites
	// and no-op zero writes are cheaper, so this over-estimates and keeps a tx
	// from running out of gas mid-range. Override the derived count with
	// --slots-per-tx to pack more aggressively.
	sstoreGasPerSlot = 22_500

	// stateCreationBytesPerSlot mirrors the EIP-8037 STORAGE_CREATION_SIZE
	// (key + value) charged for every fresh storage slot at the chain's cost
	// per state byte.
	stateCreationBytesPerSlot = 64

	// sstoreTxOverheadGas covers intrinsic gas, the 64-byte range calldata and a
	// per-tx safety margin subtracted before dividing the gas limit into slots.
	sstoreTxOverheadGas = 50_000
)

// nonExistingStartSlot mirrors EELS START_SLOT: keccak256("random") masked
// to 160 bits, a slot range that is assumed to be empty in the bloated EOA.
var nonExistingStartSlot = new(big.Int).Mod(
	common.HexToHash("0xa4896a3f93bf4bf58378e579f3cf193bb4af1022af7d2089f37d8bae7157b85f").Big(),
	new(big.Int).Lsh(big.NewInt(1), 160),
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
	GasLimit          uint64  `yaml:"gas_limit"`
	Mode              string  `yaml:"mode"`
	AuthorityPrivkey  string  `yaml:"authority_privkey"`
	ExistingSlots     bool    `yaml:"existing_slots"`
	WriteNewValue     bool    `yaml:"write_new_value"`
	StartSlot         string  `yaml:"start_slot"`
	SlotsPerTx        uint64  `yaml:"slots_per_tx"`
	KeepCounter       bool    `yaml:"keep_counter"`
	LoopGasThreshold  uint64  `yaml:"loop_gas_threshold"`
	Timeout           string  `yaml:"timeout"`
	ClientGroup       string  `yaml:"client_group"`
	DeployClientGroup string  `yaml:"deploy_client_group"`
	LogTxs            bool    `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	authority    *spamoor.Wallet
	attackAddr   common.Address
	txGasLimit   uint64
	startCounter *big.Int
	slotsPerTx   uint64
}

var ScenarioName = "bloated-eoa"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        1,
	MaxPending:        0,
	MaxWallets:        0,
	Rebroadcast:       1,
	BaseFee:           20,
	TipFee:            2,
	GasLimit:          0,
	Mode:              ModeSload,
	AuthorityPrivkey:  "",
	ExistingSlots:     true,
	WriteNewValue:     false,
	StartSlot:         "",
	SlotsPerTx:        0,
	KeepCounter:       false,
	LoopGasThreshold:  0xFFFF,
	Timeout:           "",
	ClientGroup:       "",
	DeployClientGroup: "",
	LogTxs:            false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Description:    "Delegate a storage-bloated EOA to an SLOAD/SSTORE attack contract and hammer its storage",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of attack transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of attack transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in attack transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in attack transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit per attack transaction (0 = min(block gas limit, 2^24))")
	flags.StringVar(&s.options.Mode, "mode", ScenarioDefaultOptions.Mode, "Attack mode: 'sload' (test_sload_bloated) or 'sstore' (test_sstore_bloated)")
	flags.StringVar(&s.options.AuthorityPrivkey, "authority-privkey", ScenarioDefaultOptions.AuthorityPrivkey, "Private key of the storage-bloated EOA to delegate to the attack contract")
	flags.BoolVar(&s.options.ExistingSlots, "existing-slots", ScenarioDefaultOptions.ExistingSlots, "Target existing slots (start at slot 1) instead of empty slots (start at keccak256('random') % 2^160)")
	flags.BoolVar(&s.options.WriteNewValue, "write-new-value", ScenarioDefaultOptions.WriteNewValue, "sstore mode: write slot+1 instead of the current value (slot for existing slots, 0 for empty slots)")
	flags.StringVar(&s.options.StartSlot, "start-slot", ScenarioDefaultOptions.StartSlot, "Override the starting slot (decimal or 0x hex, default depends on --existing-slots)")
	flags.Uint64Var(&s.options.SlotsPerTx, "slots-per-tx", ScenarioDefaultOptions.SlotsPerTx, "sstore mode: slots to write per transaction (0 = derive from the gas limit)")
	flags.BoolVar(&s.options.KeepCounter, "keep-counter", ScenarioDefaultOptions.KeepCounter, "sload mode: continue from the slot counter stored in slot 0 of the EOA instead of resetting it")
	flags.Uint64Var(&s.options.LoopGasThreshold, "loop-gas-threshold", ScenarioDefaultOptions.LoopGasThreshold, "sload mode: remaining gas at which the attack loop stops and persists the counter")
	flags.StringVar(&s.options.Timeout, "timeout", ScenarioDefaultOptions.Timeout, "Timeout for the scenario (e.g. '1h', '30m', '5s') - empty means no timeout")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	flags.StringVar(&s.options.DeployClientGroup, "deploy-client-group", ScenarioDefaultOptions.DeployClientGroup, "Client group to use for deployment and delegation transactions")
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

	s.options.Mode = strings.ToLower(s.options.Mode)
	if s.options.Mode != ModeSload && s.options.Mode != ModeSstore {
		return fmt.Errorf("invalid mode %q, must be %q or %q", s.options.Mode, ModeSload, ModeSstore)
	}

	if s.options.Mode != ModeSstore {
		if s.options.WriteNewValue {
			s.logger.Warnf("--write-new-value only applies to sstore mode, ignoring")
		}
		if s.options.SlotsPerTx > 0 {
			s.logger.Warnf("--slots-per-tx only applies to sstore mode, ignoring")
		}
	}
	if s.options.Mode != ModeSload && s.options.KeepCounter {
		s.logger.Warnf("--keep-counter only applies to sload mode, ignoring")
	}

	if s.options.AuthorityPrivkey == "" {
		return fmt.Errorf("--authority-privkey is required (private key of the storage-bloated EOA)")
	}

	privkey, address, err := spamoor.LoadPrivateKey(s.options.AuthorityPrivkey)
	if err != nil {
		return fmt.Errorf("invalid authority private key: %w", err)
	}

	s.authority = spamoor.NewWallet(privkey, address)

	if s.options.StartSlot != "" {
		startSlot, ok := new(big.Int).SetString(s.options.StartSlot, 0)
		if !ok || startSlot.Sign() <= 0 || startSlot.BitLen() > 256 {
			return fmt.Errorf("invalid start slot %q", s.options.StartSlot)
		}

		s.startCounter = startSlot
	} else if s.options.ExistingSlots {
		s.startCounter = big.NewInt(1)
	} else {
		s.startCounter = new(big.Int).Set(nonExistingStartSlot)
	}

	s.txGasLimit = s.options.GasLimit
	blockLimit := s.walletPool.GetTxPool().GetCurrentGasLimit()
	if s.txGasLimit == 0 {
		s.txGasLimit = maxTxGasLimit
		if blockLimit > 0 && blockLimit < s.txGasLimit {
			s.txGasLimit = blockLimit
		}
	} else if blockLimit > 0 && s.txGasLimit > blockLimit {
		s.logger.Warnf("Gas limit %d exceeds block gas limit %d and will most likely be dropped by the execution layer client", s.txGasLimit, blockLimit)
	}

	if s.options.Mode == ModeSstore {
		s.slotsPerTx = s.computeSlotsPerTx()
	}

	return nil
}

// computeSlotsPerTx returns how many slots each sstore transaction should write.
// It honours an explicit --slots-per-tx, otherwise it divides the usable tx gas
// by a conservative per-slot estimate (including EIP-8037 state-creation gas
// where the chain charges it) so a transaction does not run out of gas
// mid-range.
func (s *Scenario) computeSlotsPerTx() uint64 {
	if s.options.SlotsPerTx > 0 {
		return s.options.SlotsPerTx
	}

	perSlot := uint64(sstoreGasPerSlot)
	if cpsb := s.walletPool.GetTxPool().GetCostPerStateByte(); cpsb > 0 {
		perSlot += stateCreationBytesPerSlot * cpsb
	}

	if s.txGasLimit <= sstoreTxOverheadGas {
		return 1
	}

	slots := (s.txGasLimit - sstoreTxOverheadGas) / perSlot
	if slots == 0 {
		return 1
	}

	return slots
}

// sstoreRange returns the [start, end) slot range a given sstore transaction
// writes. Ranges are derived from the transaction index so they never overlap,
// matching the disjoint per-tx ranges of EELS test_sstore_bloated.
func (s *Scenario) sstoreRange(txIdx uint64) (start, end *big.Int) {
	offset := new(big.Int).Mul(new(big.Int).SetUint64(txIdx), new(big.Int).SetUint64(s.slotsPerTx))
	start = new(big.Int).Add(s.startCounter, offset)
	end = new(big.Int).Add(start, new(big.Int).SetUint64(s.slotsPerTx))
	return start, end
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	s.logger.Infof("bloated EOA: %v, mode: %v, existing slots: %v", s.authority.GetAddress().String(), s.options.Mode, s.options.ExistingSlots)

	receipt, err := s.deployAttackContract(ctx)
	if err != nil {
		return fmt.Errorf("could not deploy attack contract: %w", err)
	}

	s.attackAddr = receipt.ContractAddress
	s.logger.Infof("deployed attack contract at %v", s.attackAddr.String())

	counter, err := s.delegateAuthority(ctx)
	if err != nil {
		return fmt.Errorf("could not delegate bloated EOA: %w", err)
	}

	if s.options.Mode == ModeSstore {
		s.logger.Infof("delegated %v to attack contract, start slot: 0x%x, %d slots per tx", s.authority.GetAddress().String(), counter, s.slotsPerTx)
	} else {
		s.logger.Infof("delegated %v to attack contract, slot counter: 0x%x", s.authority.GetAddress().String(), counter)
	}

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
		timeout, err = time.ParseDuration(s.options.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout value: %v", err)
		}
		s.logger.Infof("Timeout set to %v", timeout)
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

func (s *Scenario) deployAttackContract(ctx context.Context) (*txtypes.Receipt, error) {
	client, wallet, feeCap, tipCap, err := s.getSetupTxContext()
	if err != nil {
		return nil, err
	}

	initcode, err := compileInitcode()
	if err != nil {
		return nil, err
	}

	runtimeCode, err := compileAttackContract(&s.options)
	if err != nil {
		return nil, err
	}

	deployData := append(initcode, runtimeCode...)

	deployGas := s.walletPool.EstimateDeployGas(ctx, client, wallet.GetAddress(), uint256.NewInt(0), deployData)
	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       deployGas,
		To:        nil,
		Value:     uint256.NewInt(0),
		Data:      deployData,
	})
	if err != nil {
		return nil, err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, err
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.getDeployClientGroup(),
		Rebroadcast: true,
	})
	if err != nil {
		return nil, err
	}

	if receipt == nil {
		return nil, fmt.Errorf("deployment transaction receipt is nil")
	}

	if receipt.Status != txtypes.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("deployment transaction reverted: %v", tx.Hash().String())
	}

	return receipt, nil
}

// delegateAuthority sends a setcode tx from a child wallet that delegates the
// bloated EOA to the attack contract. Returns the slot the attack starts from.
//
// In sload mode the same tx also initializes the slot-0 counter via calldata.
// In sstore mode each benchmark tx carries its own [start, end) range, so no
// counter is kept on chain and the delegation tx is sent with empty calldata
// (a no-op through the contract's calldata guard).
func (s *Scenario) delegateAuthority(ctx context.Context) (*big.Int, error) {
	client, wallet, feeCap, tipCap, err := s.getSetupTxContext()
	if err != nil {
		return nil, err
	}

	authorityAddr := s.authority.GetAddress()

	counter := s.startCounter
	var calldata []byte
	if s.options.Mode == ModeSload {
		if s.options.KeepCounter {
			stored, err := client.GetStorageAt(ctx, authorityAddr, common.Hash{})
			if err != nil {
				return nil, fmt.Errorf("could not read slot counter: %w", err)
			}

			if stored.Big().Sign() == 0 {
				s.logger.Warnf("slot counter of %v is empty, starting at 0x%x", authorityAddr.String(), s.startCounter)
			} else {
				counter = stored.Big()
			}
		}

		calldata = common.BigToHash(counter).Bytes()
	}

	authorityNonce, err := client.GetPendingNonceAt(ctx, authorityAddr)
	if err != nil {
		return nil, fmt.Errorf("could not get nonce of bloated EOA: %w", err)
	}

	authorization, err := txtypes.SignAuthorization(txtypes.SetCodeAuthorization{
		ChainID: *uint256.MustFromBig(s.walletPool.GetChainId()),
		Address: s.attackAddr,
		Nonce:   authorityNonce,
	}, s.authority.GetPrivateKey())
	if err != nil {
		return nil, fmt.Errorf("could not sign set code authorization: %w", err)
	}

	// intrinsic + authorization + (sload mode) one SSTORE to slot 0, plus the
	// EIP-8037 state gas for the authorization and a potentially new slot.
	gasLimit := uint64(200_000)
	if cpsb := s.walletPool.GetTxPool().GetCostPerStateByte(); cpsb > 0 {
		gasLimit += (spamoor.AuthorizationCreationSize + spamoor.AccountCreationSize + 64) * cpsb
	}

	txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       gasLimit,
		To:        &authorityAddr,
		Value:     uint256.NewInt(0),
		Data:      calldata,
		AuthList:  []txtypes.SetCodeAuthorization{authorization},
	})
	if err != nil {
		return nil, err
	}

	tx, err := wallet.BuildSetCodeTx(txData)
	if err != nil {
		return nil, err
	}

	receipt, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.getDeployClientGroup(),
		Rebroadcast: true,
	})
	if err != nil {
		return nil, err
	}

	if receipt == nil {
		return nil, fmt.Errorf("delegation transaction receipt is nil")
	}

	if receipt.Status != txtypes.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("delegation transaction reverted: %v", tx.Hash().String())
	}

	code, err := client.GetCodeAt(ctx, authorityAddr)
	if err != nil {
		return nil, fmt.Errorf("could not verify delegation: %w", err)
	}

	expectedCode := append([]byte{0xef, 0x01, 0x00}, s.attackAddr.Bytes()...)
	if string(code) != string(expectedCode) {
		return nil, fmt.Errorf("bloated EOA code is 0x%x after delegation, expected 0x%x (authorization nonce mismatch?)", code, expectedCode)
	}

	return counter, nil
}

func (s *Scenario) getDeployClientGroup() string {
	if s.options.DeployClientGroup != "" {
		return s.options.DeployClientGroup
	}

	return s.options.ClientGroup
}

func (s *Scenario) getSetupTxContext() (*spamoor.Client, *spamoor.Wallet, *big.Int, *big.Int, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.getDeployClientGroup()),
	)
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	if client == nil {
		return nil, nil, nil, nil, scenario.ErrNoClients
	}

	if wallet == nil {
		return nil, nil, nil, nil, scenario.ErrNoWallet
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return client, wallet, feeCap, tipCap, nil
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64) (scenario.ReceiptChan, *txtypes.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
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
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	// In sstore mode each tx carries its own disjoint [start, end) range as
	// calldata (32-byte start || 32-byte end). In sload mode the loop reads its
	// counter from slot 0 and the tx carries no calldata.
	var calldata []byte
	if s.options.Mode == ModeSstore {
		start, end := s.sstoreRange(txIdx)
		calldata = append(common.BigToHash(start).Bytes(), common.BigToHash(end).Bytes()...)
	}

	authorityAddr := s.authority.GetAddress()
	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.txGasLimit,
		To:        &authorityAddr,
		Value:     uint256.NewInt(0),
		Data:      calldata,
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
		OnComplete: func(tx *txtypes.Transaction, receipt *txtypes.Receipt, err error) {
			receiptChan <- receipt
		},
		OnConfirm: func(tx *txtypes.Transaction, receipt *txtypes.Receipt) {
			logger := s.logger.WithField("rpc", client.GetName())
			if receipt.Status != txtypes.ReceiptStatusSuccessful {
				logger.Warnf("transaction %d reverted in block #%v (gas used: %v); if this is sstore mode, lower --slots-per-tx", txIdx+1, receipt.BlockNumber.String(), receipt.GasUsed)
				return
			}

			txFees := utils.GetTransactionFees(tx, receipt)
			logger.Debugf(
				" transaction %d confirmed in block #%v. gas used: %v, total fee: %v gwei (base: %v)",
				txIdx+1,
				receipt.BlockNumber.String(),
				receipt.GasUsed,
				txFees.TotalFeeGweiString(),
				txFees.TxBaseFeeGweiString(),
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
