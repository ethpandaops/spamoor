package erc4337

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/erc4337/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

// Per-UserOperation gas allowances. Static (no per-op estimation) to keep the
// hot path cheap - the paymaster reimburses unused gas, so over-allocation only
// locks deposit briefly.
//
// A UserOperation that deploys a fresh account needs a large verificationGasLimit
// to cover the initCode CREATE2: under the Amsterdam fee schedule deploying
// contract code carries a state-creation surcharge, so a too-tight limit reverts
// handleOps with "Out of gas" inside the factory CREATE2 (same class of fix as
// the erc20/erc721 scenarios). A UserOperation that reuses an already-deployed
// account does no deployment and needs far less (just signature + nonce +
// paymaster validation).
const (
	newAccountVerificationGas   = uint64(1500000)
	reuseAccountVerificationGas = uint64(250000)

	opCallGasLimit       = uint64(200000)
	opPreVerificationGas = uint64(60000)
	opPaymasterVerifGas  = uint64(100000)
	opPaymasterPostOpGas = uint64(50000)

	// perOpOverhead is the EntryPoint per-op bookkeeping gas added on top of the
	// op's own limits when sizing the outer handleOps tx.
	perOpOverhead = uint64(60000)
	// handleOpsBaseGas is the fixed handleOps overhead independent of bundle size.
	handleOpsBaseGas = uint64(80000)
)

// opTxGas returns the outer handleOps gas budget attributed to a single bundled
// op, given whether that op deploys a new account.
func opTxGas(newAccount bool) uint64 {
	verif := reuseAccountVerificationGas
	if newAccount {
		verif = newAccountVerificationGas
	}
	return verif + opCallGasLimit + opPreVerificationGas + opPaymasterVerifGas + opPaymasterPostOpGas + perOpOverhead
}

type ScenarioOptions struct {
	TotalCount         uint64  `yaml:"total_count"`
	Throughput         uint64  `yaml:"throughput"`
	MaxPending         uint64  `yaml:"max_pending"`
	MaxWallets         uint64  `yaml:"max_wallets"`
	Rebroadcast        uint64  `yaml:"rebroadcast"`
	BaseFee            float64 `yaml:"base_fee"`
	TipFee             float64 `yaml:"tip_fee"`
	BaseFeeWei         string  `yaml:"base_fee_wei"`
	TipFeeWei          string  `yaml:"tip_fee_wei"`
	BundleSize         uint64  `yaml:"bundle_size"`
	NewAccountInterval uint64  `yaml:"new_account_interval"`
	PaymasterDeposit   float64 `yaml:"paymaster_deposit"`
	Timeout            string  `yaml:"timeout"`
	ClientGroup        string  `yaml:"client_group"`
	LogTxs             bool    `yaml:"log_txs"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	deploymentInfo *DeploymentInfo

	// accountPools tracks the reusable SimpleAccount pool per owner (child)
	// wallet, keyed by owner address -> *walletAccountPool.
	accountPools sync.Map
}

// accountEntry is one reusable SimpleAccount: its counterfactual address, the
// CREATE2 salt used to derive/deploy it, and the next EntryPoint nonce (key 0)
// to use for it.
type accountEntry struct {
	addr  common.Address
	salt  *big.Int
	nonce uint64
}

// walletAccountPool holds the accounts owned by a single child wallet. All of a
// wallet's accounts are only ever used as senders in handleOps bundled by that
// same wallet, so the mutex (held across nonce assignment and the BuildBoundTx
// that assigns the bundler EOA nonce) keeps each account's EntryPoint nonce
// ordered consistently with the bundler's EOA nonce.
type walletAccountPool struct {
	mu       sync.Mutex
	accounts []*accountEntry
	opCount  uint64
}

var ScenarioName = "erc4337"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:         0,
	Throughput:         10,
	MaxPending:         0,
	MaxWallets:         0,
	Rebroadcast:        1,
	BaseFee:            20,
	TipFee:             2,
	BundleSize:         1,
	NewAccountInterval: 1000,
	PaymasterDeposit:   10,
	Timeout:            "",
	ClientGroup:        "",
	LogTxs:             false,
}
var ScenarioDescriptor = scenario.Descriptor{
	Name:           ScenarioName,
	Aliases:        []string{"aa"},
	Description:    "Send ERC-4337 v0.7 UserOperations via EntryPoint.handleOps (account abstraction)",
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
	flags.Uint64VarP(&s.options.TotalCount, "count", "c", ScenarioDefaultOptions.TotalCount, "Total number of handleOps transactions to send")
	flags.Uint64VarP(&s.options.Throughput, "throughput", "t", ScenarioDefaultOptions.Throughput, "Number of handleOps transactions to send per slot")
	flags.Uint64Var(&s.options.MaxPending, "max-pending", ScenarioDefaultOptions.MaxPending, "Maximum number of pending transactions")
	flags.Uint64Var(&s.options.MaxWallets, "max-wallets", ScenarioDefaultOptions.MaxWallets, "Maximum number of child wallets to use")
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Enable reliable rebroadcast system")
	flags.Float64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transactions (in gwei)")
	flags.Float64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transactions (in gwei)")
	flags.StringVar(&s.options.BaseFeeWei, "basefee-wei", "", "Max fee per gas in wei (overrides --basefee for L2 sub-gwei fees)")
	flags.StringVar(&s.options.TipFeeWei, "tipfee-wei", "", "Max tip per gas in wei (overrides --tipfee for L2 sub-gwei fees)")
	flags.Uint64Var(&s.options.BundleSize, "bundle-size", ScenarioDefaultOptions.BundleSize, "Number of UserOperations to bundle into each handleOps transaction")
	flags.Uint64Var(&s.options.NewAccountInterval, "new-account-interval", ScenarioDefaultOptions.NewAccountInterval, "Deploy a new smart account every Nth UserOperation; reuse existing accounts in between (1 = new account every op, 0 = only ever create the first account per wallet)")
	flags.Float64Var(&s.options.PaymasterDeposit, "paymaster-deposit", ScenarioDefaultOptions.PaymasterDeposit, "EntryPoint deposit (in ETH) to keep funded on the sponsoring paymaster")
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

	if s.options.BundleSize == 0 {
		s.options.BundleSize = 1
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

	// deployer needs enough headroom to deploy the full stack in one batch.
	s.walletPool.AddWellKnownWallet(&spamoor.WellKnownWalletConfig{
		Name:          "deployer",
		RefillAmount:  uint256.NewInt(2000000000000000000), // 2 ETH
		RefillBalance: uint256.NewInt(1000000000000000000), // 1 ETH
	})

	if s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	// Worst-case bundle: every op deploys a new account. Used only for the
	// block-limit warning - actual per-tx gas is computed per bundle.
	worstCaseGas := handleOpsBaseGas + s.options.BundleSize*opTxGas(true)
	if blockLimit := s.walletPool.GetTxPool().GetCurrentGasLimit(); blockLimit > 0 && worstCaseGas > blockLimit {
		s.logger.Warnf("worst-case handleOps gas %d (bundle-size %d, all new accounts) exceeds block gas limit %d; reduce --bundle-size", worstCaseGas, s.options.BundleSize, blockLimit)
	}

	return nil
}

// gasConfigFor returns the per-op gas config for a UserOperation, sizing the
// verification gas to whether the op deploys a new account.
func gasConfigFor(newAccount bool) *userOpGasConfig {
	verif := reuseAccountVerificationGas
	if newAccount {
		verif = newAccountVerificationGas
	}
	return &userOpGasConfig{
		VerificationGasLimit: verif,
		CallGasLimit:         opCallGasLimit,
		PreVerificationGas:   opPreVerificationGas,
		PaymasterVerifGas:    opPaymasterVerifGas,
		PaymasterPostOpGas:   opPaymasterPostOpGas,
	}
}

func (s *Scenario) Run(ctx context.Context) error {
	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	// deploy the ERC-4337 stack
	deploymentInfo, err := s.DeployContracts(ctx, false)
	if err != nil {
		s.logger.Errorf("could not deploy ERC-4337 contracts: %v", err)
		return err
	}
	s.deploymentInfo = deploymentInfo
	s.logger.Infof("deployed EntryPoint=%v factory=%v paymaster=%v counter=%v",
		deploymentInfo.EntryPointAddr.Hex(), deploymentInfo.FactoryAddr.Hex(),
		deploymentInfo.PaymasterAddr.Hex(), deploymentInfo.CounterAddr.Hex())

	// fund the paymaster deposit that sponsors every UserOperation
	topUp := etherToWei(s.options.PaymasterDeposit)
	minBalance := new(big.Int).Div(topUp, big.NewInt(2))
	if err := s.ensurePaymasterDeposit(ctx, topUp, topUp); err != nil {
		return fmt.Errorf("could not fund paymaster deposit: %w", err)
	}

	// keep the paymaster deposit topped up while the scenario runs. The top-up
	// loop is driven by its own context so it is torn down when Run returns -
	// the parent ctx is only cancelled by the caller *after* Run returns, so
	// waiting on it directly here would deadlock once the tx loop completes.
	topUpCtx, topUpCancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.runDepositTopUp(topUpCtx, minBalance, topUp)
	}()
	defer func() {
		topUpCancel()
		wg.Wait()
	}()

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

			if _, err := receiptChan.Wait(ctx); err != nil {
				return err
			}

			return err
		},
	})

	return err
}

// runDepositTopUp periodically replenishes the paymaster deposit until the
// context is cancelled.
func (s *Scenario) runDepositTopUp(ctx context.Context, minBalance, topUpAmount *big.Int) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.ensurePaymasterDeposit(ctx, minBalance, topUpAmount); err != nil {
				s.logger.Warnf("could not top up paymaster deposit: %v", err)
			}
		}
	}
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
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	// Build the bundle from this wallet's account pool. The wallet is both the
	// bundler and the owner of every account it uses, and the pool lock is held
	// across nonce assignment and BuildBoundTx so each account's EntryPoint nonce
	// (key 0) stays ordered consistently with the bundler EOA nonce.
	pool := s.poolForWallet(wallet.GetAddress())
	pool.mu.Lock()

	ops := make([]contract.PackedUserOperation, 0, s.options.BundleSize)
	txGas := handleOpsBaseGas
	for j := uint64(0); j < s.options.BundleSize; j++ {
		entry, err := s.nextAccountLocked(ctx, pool, wallet)
		if err != nil {
			pool.mu.Unlock()
			return nil, nil, client, wallet, err
		}

		nonce := entry.nonce
		entry.nonce++
		withInitCode := nonce == 0 // an undeployed account is created on its first op

		op, err := buildSignedUserOp(
			gasConfigFor(withInitCode),
			s.deploymentInfo.EntryPointAddr,
			s.deploymentInfo.FactoryAddr,
			s.deploymentInfo.PaymasterAddr,
			s.deploymentInfo.CounterAddr,
			wallet.GetAddress(),
			wallet.GetPrivateKey(),
			entry.addr,
			entry.salt,
			nonce,
			withInitCode,
			wallet.GetChainId(),
			feeCap,
			tipCap,
		)
		if err != nil {
			pool.mu.Unlock()
			return nil, nil, client, wallet, fmt.Errorf("could not build userOp: %w", err)
		}
		ops = append(ops, op)
		txGas += opTxGas(withInitCode)
	}

	if txGas > utils.MaxGasLimitPerTx {
		txGas = utils.MaxGasLimitPerTx
	}

	tx, err := wallet.BuildBoundTx(ctx, &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       txGas,
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deploymentInfo.EntryPoint.HandleOps(transactOpts, ops, wallet.GetAddress())
	})
	pool.mu.Unlock()
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
				" transaction %d confirmed in block #%v. total fee: %v gwei (base: %v) ops: %v logs: %v",
				txIdx+1,
				receipt.BlockNumber.String(),
				txFees.TotalFeeGweiString(),
				txFees.TxBaseFeeGweiString(),
				s.options.BundleSize,
				len(receipt.Logs),
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

// poolForWallet returns the account pool for the given owner wallet, creating an
// empty one on first use.
func (s *Scenario) poolForWallet(owner common.Address) *walletAccountPool {
	if v, ok := s.accountPools.Load(owner); ok {
		return v.(*walletAccountPool)
	}
	v, _ := s.accountPools.LoadOrStore(owner, &walletAccountPool{})
	return v.(*walletAccountPool)
}

// nextAccountLocked returns the account to use for the next UserOperation in a
// bundle, applying the new-vs-reuse policy. It must be called with pool.mu held.
//
// A new account is created when the pool is empty or when the configured
// new-account interval is hit (every Nth op; interval 1 = every op, interval 0 =
// only the first account per wallet ever). Otherwise an existing account is
// reused round-robin. New accounts have their EntryPoint nonce synced from chain
// so restarts reuse the same deterministic-salt accounts instead of colliding.
func (s *Scenario) nextAccountLocked(ctx context.Context, pool *walletAccountPool, owner *spamoor.Wallet) (*accountEntry, error) {
	pool.opCount++

	needNew := len(pool.accounts) == 0
	if !needNew && s.options.NewAccountInterval >= 1 && pool.opCount%s.options.NewAccountInterval == 0 {
		needNew = true
	}

	if !needNew {
		return pool.accounts[pool.opCount%uint64(len(pool.accounts))], nil
	}

	// Deterministic salt = index within this wallet's pool, so the same accounts
	// are reproduced across restarts.
	salt := big.NewInt(int64(len(pool.accounts)))
	addr, err := s.deploymentInfo.Factory.GetAddress(&bind.CallOpts{Context: ctx}, owner.GetAddress(), salt)
	if err != nil {
		return nil, fmt.Errorf("could not compute account address: %w", err)
	}
	nonce, err := s.deploymentInfo.EntryPoint.GetNonce(&bind.CallOpts{Context: ctx}, addr, big.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("could not read account nonce: %w", err)
	}

	entry := &accountEntry{addr: addr, salt: salt, nonce: nonce.Uint64()}
	pool.accounts = append(pool.accounts, entry)
	return entry, nil
}

// etherToWei converts a (possibly fractional) ether amount to wei.
func etherToWei(eth float64) *big.Int {
	wei := new(big.Float).Mul(big.NewFloat(eth), big.NewFloat(1e18))
	out, _ := wei.Int(nil)
	return out
}
