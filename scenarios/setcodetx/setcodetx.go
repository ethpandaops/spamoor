package setcodetx

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/ethpandaops/spamoor/scenariotypes"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

type ScenarioOptions struct {
	TotalCount        uint64 `yaml:"total_count"`
	Throughput        uint64 `yaml:"throughput"`
	MaxPending        uint64 `yaml:"max_pending"`
	MaxWallets        uint64 `yaml:"max_wallets"`
	MinAuthorizations uint64 `yaml:"min_authorizations"`
	MaxAuthorizations uint64 `yaml:"max_authorizations"`
	MaxDelegators     uint64 `yaml:"max_delegators"`
	Rebroadcast       uint64 `yaml:"rebroadcast"`
	BaseFee           uint64 `yaml:"base_fee"`
	TipFee            uint64 `yaml:"tip_fee"`
	GasLimit          uint64 `yaml:"gas_limit"`
	Amount            uint64 `yaml:"amount"`
	Data              string `yaml:"data"`
	CodeAddr          string `yaml:"code_addr"`
	RandomAmount      bool   `yaml:"random_amount"`
	RandomTarget      bool   `yaml:"random_target"`
	RandomCodeAddr    bool   `yaml:"random_code_addr"`
	MaxBloating       bool   `yaml:"max_bloating"`
	ClientGroup       string `yaml:"client_group"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup
	delegatorSeed []byte
	delegators    []*txbuilder.Wallet
}

var ScenarioName = "setcodetx"
var ScenarioDefaultOptions = ScenarioOptions{
	TotalCount:        0,
	Throughput:        0,
	MaxPending:        0,
	MaxWallets:        0,
	MinAuthorizations: 1,
	MaxAuthorizations: 10,
	MaxDelegators:     0,
	Rebroadcast:       120,
	BaseFee:           20,
	TipFee:            2,
	GasLimit:          200000,
	Amount:            20,
	Data:              "",
	CodeAddr:          "",
	RandomAmount:      false,
	RandomTarget:      false,
	RandomCodeAddr:    false,
	MaxBloating:       false,
	ClientGroup:       "",
}
var ScenarioDescriptor = scenariotypes.ScenarioDescriptor{
	Name:           ScenarioName,
	Description:    "Send setcode transactions with different configurations",
	DefaultOptions: ScenarioDefaultOptions,
	NewScenario:    newScenario,
}

func newScenario(logger logrus.FieldLogger) scenariotypes.Scenario {
	return &Scenario{
		logger: logger.WithField("scenario", ScenarioName),
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
	flags.Uint64Var(&s.options.Rebroadcast, "rebroadcast", ScenarioDefaultOptions.Rebroadcast, "Number of seconds to wait before re-broadcasting a transaction")
	flags.Uint64Var(&s.options.BaseFee, "basefee", ScenarioDefaultOptions.BaseFee, "Max fee per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.TipFee, "tipfee", ScenarioDefaultOptions.TipFee, "Max tip per gas to use in transfer transactions (in gwei)")
	flags.Uint64Var(&s.options.GasLimit, "gaslimit", ScenarioDefaultOptions.GasLimit, "Gas limit to use in transactions")
	flags.Uint64Var(&s.options.Amount, "amount", ScenarioDefaultOptions.Amount, "Transfer amount per transaction (in gwei)")
	flags.StringVar(&s.options.Data, "data", ScenarioDefaultOptions.Data, "Transaction call data to send")
	flags.StringVar(&s.options.CodeAddr, "code-addr", ScenarioDefaultOptions.CodeAddr, "Code delegation target address to use for transactions")
	flags.BoolVar(&s.options.RandomAmount, "random-amount", ScenarioDefaultOptions.RandomAmount, "Use random amounts for transactions (with --amount as limit)")
	flags.BoolVar(&s.options.RandomTarget, "random-target", ScenarioDefaultOptions.RandomTarget, "Use random to addresses for transactions")
	flags.BoolVar(&s.options.RandomCodeAddr, "random-code-addr", ScenarioDefaultOptions.RandomCodeAddr, "Use random delegation target for transactions")
	flags.BoolVar(&s.options.MaxBloating, "max-bloating", ScenarioDefaultOptions.MaxBloating, "Enable maximum state bloating mode: creates ~960 EOA delegations in a single block-filling transaction")
	flags.StringVar(&s.options.ClientGroup, "client-group", ScenarioDefaultOptions.ClientGroup, "Client group to use for sending transactions")
	return nil
}

func (s *Scenario) Init(walletPool *spamoor.WalletPool, config string) error {
	s.walletPool = walletPool

	if config != "" {
		err := yaml.Unmarshal([]byte(config), &s.options)
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

	if s.options.MaxPending > 0 {
		s.pendingChan = make(chan bool, s.options.MaxPending)
	}

	s.delegatorSeed = make([]byte, 32)
	rand.Read(s.delegatorSeed)

	if s.options.MaxDelegators > 0 {
		s.delegators = make([]*txbuilder.Wallet, 0, s.options.MaxDelegators)
	}

	// In max-bloating mode, we handle throughput automatically, so skip this validation
	if !s.options.MaxBloating && s.options.TotalCount == 0 && s.options.Throughput == 0 {
		return fmt.Errorf("neither total count nor throughput limit set, must define at least one of them (see --help for list of all flags)")
	}

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

func (s *Scenario) Run(ctx context.Context) error {
	// Check if max bloating mode is enabled
	if s.options.MaxBloating {
		return s.runMaxBloatingMode(ctx)
	}

	// Original implementation for backward compatibility
	txIdxCounter := uint64(0)
	pendingCount := atomic.Int64{}
	txCount := atomic.Uint64{}

	s.logger.Infof("starting scenario: %s", ScenarioName)
	defer s.logger.Infof("scenario %s finished.", ScenarioName)

	var lastChan chan bool

	initialRate := rate.Limit(float64(s.options.Throughput) / float64(utils.SecondsPerSlot))
	if initialRate == 0 {
		initialRate = rate.Inf
	}
	limiter := rate.NewLimiter(initialRate, 1)

	for {
		if err := limiter.Wait(ctx); err != nil {
			if ctx.Err() != nil {
				break
			}

			s.logger.Debugf("rate limited: %s", err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
		}

		txIdx := txIdxCounter
		txIdxCounter++

		if s.pendingChan != nil {
			// await pending transactions
			s.pendingChan <- true
		}
		pendingCount.Add(1)

		currentChan := make(chan bool, 1)

		go func(txIdx uint64, lastChan, currentChan chan bool) {
			defer func() {
				pendingCount.Add(-1)
				currentChan <- true
			}()

			logger := s.logger
			tx, client, wallet, err := s.sendTx(ctx, txIdx, func() {
				if s.pendingChan != nil {
					time.Sleep(100 * time.Millisecond)
					<-s.pendingChan
				}
			})
			if client != nil {
				logger = logger.WithField("rpc", client.GetName())
			}
			if tx != nil {
				logger = logger.WithField("nonce", tx.Nonce())
			}
			if wallet != nil {
				logger = logger.WithField("wallet", s.walletPool.GetWalletName(wallet.GetAddress()))
			}
			if lastChan != nil {
				<-lastChan
				close(lastChan)
			}
			if err != nil {
				logger.Warnf("could not send transaction: %v", err)
				return
			}

			txCount.Add(1)
			logger.Infof("sent tx #%6d: %v", txIdx+1, tx.Hash().String())
		}(txIdx, lastChan, currentChan)

		lastChan = currentChan

		count := txCount.Load() + uint64(pendingCount.Load())
		if s.options.TotalCount > 0 && count >= s.options.TotalCount {
			break
		}
	}

	<-lastChan
	close(lastChan)

	s.logger.Infof("finished sending transactions, awaiting block inclusion...")
	s.pendingWGroup.Wait()

	return nil
}

func (s *Scenario) sendTx(ctx context.Context, txIdx uint64, onComplete func()) (*types.Transaction, *txbuilder.Client, *txbuilder.Wallet, error) {
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

	rebroadcast := 0
	if s.options.Rebroadcast > 0 {
		rebroadcast = 10
	}

	s.pendingWGroup.Add(1)
	transactionSubmitted = true
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     rebroadcast,
		RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
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
		LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName())
			if retry > 0 {
				logger = logger.WithField("retry", retry)
			}
			if rebroadcast > 0 {
				logger = logger.WithField("rebroadcast", rebroadcast)
			}
			if err != nil {
				logger.Warnf("failed sending tx %6d: %v", txIdx+1, err)
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

		var delegator *txbuilder.Wallet
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
			// Use a fixed delegate contract address for maximum efficiency
			// In max bloating mode, we want all EOAs to delegate to the same existing contract
			// to benefit from reduced gas costs (PER_AUTH_BASE_COST vs PER_EMPTY_ACCOUNT_COST)
			// Precompiles are ideal as they're guaranteed to exist with code on all networks
			codeAddr = common.HexToAddress("0x0000000000000000000000000000000000000001")
			s.logger.Infof("no --code-addr specified, using ecrecover precompile as delegate: %s", codeAddr.Hex())
		}

		authorization := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(s.walletPool.GetRootWallet().GetChainId().Uint64()),
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

func (s *Scenario) prepareDelegator(delegatorIndex uint64) (*txbuilder.Wallet, error) {
	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes, delegatorIndex)
	if s.options.MaxDelegators > 0 {
		seedBytes := []byte(s.delegatorSeed)
		idxBytes = append(idxBytes, seedBytes...)
	}
	childKey := sha256.Sum256(append(common.FromHex(s.walletPool.GetRootWallet().GetAddress().Hex()), idxBytes...))
	return txbuilder.NewWallet(fmt.Sprintf("%x", childKey))
}

func (s *Scenario) buildMaxBloatingAuthorizations(targetCount int, iteration int) []types.SetCodeAuthorization {
	authorizations := make([]types.SetCodeAuthorization, 0, targetCount)

	s.logger.Infof("building %d unique EOA authorizations for maximum state bloat (iteration %d)", targetCount, iteration)

	// Use a fixed delegate contract address for maximum efficiency
	// In max bloating mode, we want all EOAs to delegate to the same existing contract
	// to benefit from reduced gas costs (PER_AUTH_BASE_COST vs PER_EMPTY_ACCOUNT_COST)
	// Precompiles are ideal as they're guaranteed to exist with code on all networks
	var codeAddr common.Address
	if s.options.CodeAddr != "" {
		codeAddr = common.HexToAddress(s.options.CodeAddr)
		s.logger.Infof("using configured delegate address: %s", codeAddr.Hex())
	} else {
		// Default to using the ecrecover precompile (0x1) as delegate target
		// This is perfect for max-bloating mode as it's guaranteed to exist with code
		codeAddr = common.HexToAddress("0x0000000000000000000000000000000000000001")
		s.logger.Infof("no --code-addr specified, using ecrecover precompile as delegate: %s", codeAddr.Hex())
	}

	chainId := s.walletPool.GetRootWallet().GetChainId().Uint64()

	for i := 0; i < targetCount; i++ {
		// Create a unique delegator for each authorization
		// Include iteration counter to ensure different addresses for each iteration
		delegatorIndex := uint64(iteration*targetCount + i)

		delegator, err := s.prepareDelegator(delegatorIndex)
		if err != nil {
			s.logger.Errorf("could not prepare delegator %v: %v", delegatorIndex, err)
			continue
		}

		// Each EOA uses auth_nonce = 0 (assuming first EIP-7702 operation)
		// This creates maximum new state as each EOA gets its first delegation
		authorization := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(chainId),
			Address: codeAddr,
			Nonce:   0, // First delegation for each EOA
		}

		// Sign the authorization with the delegator's private key
		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), authorization)
		if err != nil {
			s.logger.Errorf("could not sign set code authorization for delegator %v: %v", delegatorIndex, err)
			continue
		}

		authorizations = append(authorizations, signedAuth)

		// Log progress every 100 authorizations
		if (i+1)%100 == 0 {
			s.logger.Debugf("generated %d/%d authorizations", i+1, targetCount)
		}
	}

	s.logger.Infof("successfully generated %d authorizations (target: %d)", len(authorizations), targetCount)
	return authorizations
}

func (s *Scenario) runMaxBloatingMode(ctx context.Context) error {
	s.logger.Infof("starting max bloating mode: targeting ~960 EOA delegations per block, continuous operation")
	defer s.logger.Infof("max bloating mode finished")

	// For max bloating, we use exactly one transaction with ~960 authorizations per iteration
	// This targets the theoretical maximum state bloat efficiency
	const targetAuthorizations = 960
	const targetGasLimit = 29000000 // Near block gas limit, leaving room for base transaction costs

	blockCounter := 0

	for {
		select {
		case <-ctx.Done():
			s.logger.Infof("max bloating mode stopping due to context cancellation")
			return ctx.Err()
		default:
		}

		blockCounter++
		s.logger.Infof("=== Starting Max Bloating Iteration #%d ===", blockCounter)

		// Fund delegator accounts first
		err := s.fundMaxBloatingDelegators(ctx, targetAuthorizations, blockCounter)
		if err != nil {
			s.logger.Errorf("failed to fund delegators for iteration %d: %v", blockCounter, err)
			time.Sleep(5 * time.Second) // Wait before retry
			continue
		}

		// Send the max bloating transaction
		err = s.sendMaxBloatingTransaction(ctx, targetAuthorizations, targetGasLimit, blockCounter)
		if err != nil {
			s.logger.Errorf("failed to send max bloating transaction for iteration %d: %v", blockCounter, err)
			time.Sleep(5 * time.Second) // Wait before retry
			continue
		}

		s.logger.Infof("=== Completed Max Bloating Iteration #%d ===", blockCounter)

		// Small delay between iterations to allow for block inclusion
		time.Sleep(2 * time.Second)
	}
}

func (s *Scenario) fundMaxBloatingDelegators(ctx context.Context, targetCount int, iteration int) error {
	s.logger.Infof("funding %d unique delegator accounts with 1 wei each (iteration %d)", targetCount, iteration)

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return fmt.Errorf("no client available for funding delegators")
	}

	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	// Get suggested fees for funding transactions
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
			return fmt.Errorf("failed to get suggested fees for funding: %w", err)
		}
	}

	// Minimum gas prices
	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
	}

	// Fund with 1 wei as requested by user
	fundingAmount := uint256.NewInt(1)

	successCount := 0

	// Fund delegators sequentially to avoid nonce conflicts
	for i := 0; i < targetCount; i++ {
		// Generate unique delegator address for this iteration
		delegatorIndex := uint64(iteration*targetCount + i)
		delegator, err := s.prepareDelegator(delegatorIndex)
		if err != nil {
			s.logger.Errorf("could not prepare delegator %v for funding: %v", delegatorIndex, err)
			continue
		}

		// Build funding transaction
		delegatorAddr := delegator.GetAddress()
		txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       21000, // Standard ETH transfer gas
			To:        &delegatorAddr,
			Value:     fundingAmount,
			Data:      []byte{},
		})
		if err != nil {
			s.logger.Errorf("failed to build funding tx for delegator %d: %v", i, err)
			continue
		}

		tx, err := wallet.BuildDynamicFeeTx(txData)
		if err != nil {
			s.logger.Errorf("failed to build funding transaction for delegator %d: %v", i, err)
			continue
		}

		// Send funding transaction with no retries to avoid duplicates
		err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
			Client:          client,
			MaxRebroadcasts: 0, // No retries to avoid duplicates
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					s.logger.Debugf("funding tx failed for delegator %d: %v", i, err)
					return
				}
				if receipt != nil && receipt.Status == 1 {
					successCount++
				}
			},
			LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
				// Only log actual failures
				if err != nil {
					s.logger.Debugf("funding tx failed for delegator %d: %v", i, err)
				}
			},
		})

		if err != nil {
			s.logger.Debugf("failed to send funding transaction for delegator %d: %v", i, err)
			continue
		}

		// Log progress every 100 transactions
		if (i+1)%100 == 0 {
			s.logger.Infof("sent %d/%d funding transactions", i+1, targetCount)
		}

		// Small delay between transactions to ensure proper nonce ordering
		time.Sleep(10 * time.Millisecond)
	}

	s.logger.Infof("funding completed for iteration %d - sent %d funding transactions", iteration, targetCount)

	// Wait for funding transactions to be included
	s.logger.Infof("waiting for funding transactions to be included...")
	time.Sleep(5 * time.Second)

	return nil
}

func (s *Scenario) sendMaxBloatingTransaction(ctx context.Context, targetAuthorizations int, targetGasLimit uint64, blockCounter int) error {
	s.logger.Infof("sending max bloating transaction with %d authorizations, gas limit: %d", targetAuthorizations, targetGasLimit)

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return fmt.Errorf("no client available for sending max bloating transaction")
	}

	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, 0)

	// Get suggested fees or use configured values
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
			return fmt.Errorf("failed to get suggested fees: %w", err)
		}
	}

	// Ensure minimum gas prices for inclusion
	if feeCap.Cmp(big.NewInt(1000000000)) < 0 {
		feeCap = big.NewInt(1000000000)
	}
	if tipCap.Cmp(big.NewInt(1000000000)) < 0 {
		tipCap = big.NewInt(1000000000)
	}

	// Use minimal amount for max bloating (focus on authorizations, not value transfer)
	amount := uint256.NewInt(0) // No value transfer needed

	// Target address - use our own wallet for simplicity
	toAddr := wallet.GetAddress()

	// Minimal call data
	txCallData := []byte{}
	if s.options.Data != "" {
		dataBytes, err := txbuilder.ParseBlobRefsBytes(strings.Split(s.options.Data, ","), nil)
		if err != nil {
			return fmt.Errorf("failed to parse call data: %w", err)
		}
		txCallData = dataBytes
	}

	// Build the authorizations for maximum state bloat
	authorizations := s.buildMaxBloatingAuthorizations(targetAuthorizations, blockCounter)
	s.logger.Infof("generated %d authorizations for max bloating", len(authorizations))

	txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       targetGasLimit,
		To:        &toAddr,
		Value:     amount,
		Data:      txCallData,
		AuthList:  authorizations,
	})
	if err != nil {
		return fmt.Errorf("failed to build transaction metadata: %w", err)
	}

	tx, err := wallet.BuildSetCodeTx(txData)
	if err != nil {
		return fmt.Errorf("failed to build transaction: %w", err)
	}

	// Use WaitGroup to wait for transaction completion
	var wg sync.WaitGroup
	wg.Add(1)

	// Send the transaction
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     10,
		RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			defer wg.Done()

			if err != nil {
				s.logger.WithField("rpc", client.GetName()).Errorf("max bloating tx failed: %v", err)
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

			// Calculate efficiency metrics
			authCount := len(authorizations)
			gasPerAuth := float64(receipt.GasUsed) / float64(authCount)
			estimatedBytesPerAuth := 135.0 // ~135 bytes state change per EOA delegation
			gasPerByte := gasPerAuth / estimatedBytesPerAuth

			s.logger.WithField("rpc", client.GetName()).Infof(
				"MAX BLOATING SUCCESS - Block #%s, Gas Used: %d, Authorizations: %d, Gas/Auth: %.1f, Gas/Byte: %.1f, Total Fee: %s gwei",
				receipt.BlockNumber.String(),
				receipt.GasUsed,
				authCount,
				gasPerAuth,
				gasPerByte,
				gweiTotalFee.String(),
			)
		},
		LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
			logger := s.logger.WithField("rpc", client.GetName())
			if retry > 0 {
				logger = logger.WithField("retry", retry)
			}
			if rebroadcast > 0 {
				logger = logger.WithField("rebroadcast", rebroadcast)
			}
			if err != nil {
				logger.Errorf("failed sending max bloating tx: %v", err)
			} else if retry > 0 || rebroadcast > 0 {
				logger.Infof("successfully sent max bloating tx")
			}
		},
	})

	if err != nil {
		wallet.ResetPendingNonce(ctx, client)
		return fmt.Errorf("failed to send max bloating transaction: %w", err)
	}

	s.logger.Infof("max bloating transaction submitted: %s", tx.Hash().String())
	s.logger.Infof("awaiting max bloating transaction confirmation...")
	wg.Wait()
	s.logger.Infof("max bloating transaction confirmed and processed")
	return nil
}
