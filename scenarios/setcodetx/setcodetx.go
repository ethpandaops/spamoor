package setcodetx

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
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

// EIP-7702 gas cost constants
const (
	// PER_AUTH_BASE_COST (EIP-7702) - upper bound on gas cost per authorization when delegating to existing contract in this scenario
	GasPerAuthorization = 26000
	// EstimatedBytesPerAuth - estimated state change in bytes per EOA delegation
	EstimatedBytesPerAuth = 135.0
	// DefaultTargetGasRatio - target percentage of block gas limit to use (99% for safety margin)
	DefaultTargetGasRatio = 0.99
	// FallbackBlockGasLimit - fallback gas limit if network query fails
	FallbackBlockGasLimit = 30000000
	// BaseTransferCost - gas cost for a standard ETH transfer
	BaseTransferCost = 21000
	// MaxTransactionSize - Ethereum transaction size limit in bytes (128KiB)
	MaxTransactionSize = 131072 // 128 * 1024
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

// EOAEntry represents a funded EOA account
type EOAEntry struct {
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
}

type Scenario struct {
	options    ScenarioOptions
	logger     *logrus.Entry
	walletPool *spamoor.WalletPool

	pendingChan   chan bool
	pendingWGroup sync.WaitGroup
	delegatorSeed []byte
	delegators    []*txbuilder.Wallet

	// FIFO queue for funded accounts
	eoaQueue      []EOAEntry
	eoaQueueMutex sync.Mutex

	// Semaphore for worker control
	workerSemaphore chan struct{}
	workerDone      chan struct{}
	workerWg        sync.WaitGroup
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

	if s.options.MaxBloating && !s.options.RandomCodeAddr && s.options.CodeAddr == "" {
		s.logger.Infof("no --code-addr specified, using ecrecover precompile as delegate: %s", common.HexToAddress("0x0000000000000000000000000000000000000001"))
	}

	if s.options.MaxWallets > 0 {
		s.walletPool.SetWalletCount(s.options.MaxWallets)
	} else if s.options.MaxBloating {
		// In max-bloating mode, if maxWallets is not set, use 1 wallet (the root wallet)
		s.walletPool.SetWalletCount(1)
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

	// Initialize FIFO queue and worker for EOA management
	s.eoaQueue = make([]EOAEntry, 0)
	s.workerSemaphore = make(chan struct{}, 1) // Buffered channel for semaphore
	s.workerDone = make(chan struct{})

	// Start the worker goroutine for writing EOAs to file
	if s.options.MaxBloating {
		s.workerWg.Add(1)
		go s.eoaWorker()
	}

	return nil
}

func (s *Scenario) Config() string {
	yamlBytes, _ := yaml.Marshal(&s.options)
	return string(yamlBytes)
}

// getNetworkBlockGasLimit retrieves the current block gas limit from the network
// It waits for a new block to be mined (with 30s timeout) to ensure fresh data
func (s *Scenario) getNetworkBlockGasLimit(ctx context.Context, client *txbuilder.Client) uint64 {
	// Create a timeout context for the entire operation (30 seconds)
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get the current block number first
	currentBlockNumber, err := client.GetEthClient().BlockNumber(timeoutCtx)
	if err != nil {
		s.logger.Warnf("failed to get current block number: %v, using fallback: %d", err, FallbackBlockGasLimit)
		return FallbackBlockGasLimit
	}

	s.logger.Debugf("waiting for new block to be mined (current: %d, timeout: 30s)", currentBlockNumber)

	// Wait for a new block to be mined (poll every 500ms)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var latestBlock *types.Block
	for {
		select {
		case <-timeoutCtx.Done():
			s.logger.Warnf("timeout waiting for new block to be mined, using fallback: %d", FallbackBlockGasLimit)
			return FallbackBlockGasLimit
		case <-ticker.C:
			// Check for a new block
			newBlockNumber, err := client.GetEthClient().BlockNumber(timeoutCtx)
			if err != nil {
				s.logger.Debugf("error checking block number: %v", err)
				continue
			}

			// If we have a new block, get its details
			if newBlockNumber > currentBlockNumber {
				latestBlock, err = client.GetEthClient().BlockByNumber(timeoutCtx, nil)
				if err != nil {
					s.logger.Debugf("error getting latest block details: %v", err)
					continue
				}
				s.logger.Debugf("new block mined: %d", newBlockNumber)
				goto blockFound
			}
		}
	}

blockFound:
	gasLimit := latestBlock.GasLimit()
	s.logger.Debugf("network block gas limit from fresh block #%d: %d", latestBlock.NumberU64(), gasLimit)
	return gasLimit
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

	// Use a fixed delegate contract address for maximum efficiency
	// In max bloating mode, we want all EOAs to delegate to the same existing contract
	// to benefit from reduced gas costs (PER_AUTH_BASE_COST vs PER_EMPTY_ACCOUNT_COST)
	// Precompiles are ideal as they're guaranteed to exist with code on all networks
	var codeAddr common.Address
	if s.options.CodeAddr != "" {
		codeAddr = common.HexToAddress(s.options.CodeAddr)
	} else {
		// Default to using the ecrecover precompile (0x1) as delegate target
		codeAddr = common.HexToAddress("0x0000000000000000000000000000000000000001")
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
	}

	return authorizations
}

func (s *Scenario) runMaxBloatingMode(ctx context.Context) error {
	s.logger.Infof("starting max bloating mode: self-adjusting to target block gas limit, continuous operation")

	// Get a client for network operations
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)

	// Get the actual network block gas limit
	networkGasLimit := s.getNetworkBlockGasLimit(ctx, client)
	targetGas := uint64(float64(networkGasLimit) * DefaultTargetGasRatio)

	// Calculate initial authorization count based on network gas limit and known gas cost per authorization
	initialAuthorizations := int(targetGas / GasPerAuthorization)

	// Dynamic authorization count - starts based on network parameters and adjusts based on actual performance
	currentAuthorizations := initialAuthorizations

	var blockCounter int

	for {
		select {
		case <-ctx.Done():
			s.logger.Errorf("max bloating mode stopping due to context cancellation")
			return ctx.Err()
		default:
		}

		blockCounter++

		// For the first iteration, we need to fund delegators before bloating
		// For subsequent iterations, funding happens after analysis
		if blockCounter == 1 {
			s.logger.Infof("════════════════ INITIAL FUNDING PHASE ════════════════")
			err := s.fundMaxBloatingDelegators(ctx, currentAuthorizations, blockCounter, networkGasLimit)
			if err != nil {
				s.logger.Errorf("failed to fund delegators for initial iteration: %v", err)
				time.Sleep(5 * time.Second) // Wait before retry
				blockCounter--              // Retry the same iteration
				continue
			}

			// Wait for funding transactions to be confirmed and included in blocks
			s.logger.Infof("Waiting for funding transactions to be confirmed...")
			time.Sleep(15 * time.Second) // Wait for at least one block to ensure funding txs are mined
		}

		// Send the max bloating transaction and wait for confirmation
		s.logger.Infof("════════════════ BLOATING PHASE #%d ════════════════", blockCounter)
		actualGasUsed, blockNumber, authCount, gasPerAuth, gasPerByte, gweiTotalFee, err := s.sendMaxBloatingTransaction(ctx, currentAuthorizations, targetGas, blockCounter)
		if err != nil {
			s.logger.Errorf("failed to send max bloating transaction for iteration %d: %v", blockCounter, err)
			time.Sleep(5 * time.Second) // Wait before retry
			continue
		}

		// Wait a moment to ensure no more transactions are being processed
		time.Sleep(2 * time.Second)

		// Open semaphore (green light) during analysis phase to allow worker to process EOA queue
		s.openWorkerSemaphore()

		s.logger.Infof("%%%%%%%%%%%%%%%%%%%% ANALYSIS PHASE #%d %%%%%%%%%%%%%%%%%%%%", blockCounter)
		s.logger.WithField("scenario", "setcodetx").Infof("MAX BLOATING TX MINED - Block #%s, Gas Used: %d, Authorizations: %d, Gas/Auth: %.1f, Gas/Byte: %.1f, Total Fee: %s gwei",
			blockNumber, actualGasUsed, authCount, gasPerAuth, gasPerByte, gweiTotalFee)

		// Self-adjust authorization count based on actual performance
		if actualGasUsed > 0 && authCount > 0 {
			gasPerAuth := float64(actualGasUsed) / float64(authCount)
			targetAuths := int(float64(targetGas) / gasPerAuth)

			// Calculate the adjustment needed
			authDifference := targetAuths - authCount

			if actualGasUsed < targetGas {
				// We're under target, increase authorization count with a slight safety margin
				newAuthorizations := currentAuthorizations + authDifference - 1

				if newAuthorizations > currentAuthorizations {
					s.logger.Infof("Adjusting authorizations: %d → %d (need %d more for target)",
						currentAuthorizations, newAuthorizations, authDifference)
					currentAuthorizations = newAuthorizations
				}
			} else if actualGasUsed > targetGas {
				// We're over target, reduce to reach max block utilization
				excess := actualGasUsed - targetGas
				newAuthorizations := currentAuthorizations - int(excess) + 1

				s.logger.Infof("Reducing authorizations: %d → %d (excess: %d gas)",
					currentAuthorizations, newAuthorizations, excess)
				currentAuthorizations = newAuthorizations

			} else {
				s.logger.Infof("Target achieved! Gas Used: %d / Target: %d", actualGasUsed, targetGas)
			}
		}

		// Transaction confirmed - small delay before next iteration
		time.Sleep(2 * time.Second)

		// Now fund delegators for the next iteration (except on the last iteration)
		// This ensures funding happens AFTER bloating transactions are confirmed
		s.logger.Infof("════════════════ FUNDING PHASE #%d (for next iteration) ════════════════", blockCounter)
		err = s.fundMaxBloatingDelegators(ctx, currentAuthorizations, blockCounter+1, networkGasLimit)
		if err != nil {
			s.logger.Errorf("failed to fund delegators for next iteration: %v", err)
			// Don't fail the entire loop, just log the error and continue
		}

		// Wait for funding transactions to be confirmed before next bloating phase
		s.logger.Infof("Waiting for funding transactions to be confirmed before next iteration...")
		time.Sleep(15 * time.Second) // Ensure funding txs are mined in separate blocks
	}
}

func (s *Scenario) fundMaxBloatingDelegators(ctx context.Context, targetCount int, iteration int, gasLimit uint64) error {
	// Close semaphore (red light) during funding phase
	s.closeWorkerSemaphore()

	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return fmt.Errorf("no client available for funding delegators")
	}

	// Use root wallet since we set child wallet count to 0 in max-bloating mode
	wallet := s.walletPool.GetRootWallet()

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

	var confirmedCount int64
	sentCount := uint64(0)
	delegatorIndex := uint64(iteration * 1000000) // Large offset per iteration to avoid conflicts

	// Calculate approximate transactions per block based on gas limit
	// Standard transfer = BaseTransferCost gas.
	var maxTxsPerBlock = gasLimit / uint64(BaseTransferCost)

	for {
		// Check if we have enough confirmed transactions
		confirmed := atomic.LoadInt64(&confirmedCount)
		if confirmed >= int64(targetCount) {
			// We have minimum required, but let's check if we should fill the current block
			// If we've sent transactions recently, wait a bit to see if block gets filled
			if sentCount > 0 && (sentCount%100) > 50 {
				// Continue to fill the block
			} else {
				break
			}
		}

		// Generate unique delegator address
		delegator, err := s.prepareDelegator(delegatorIndex)
		if err != nil {
			s.logger.Errorf("could not prepare delegator %v for funding: %v", delegatorIndex, err)
			delegatorIndex++
			continue
		}

		// Build funding transaction
		delegatorAddr := delegator.GetAddress()
		txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       BaseTransferCost, // Standard ETH transfer gas
			To:        &delegatorAddr,
			Value:     fundingAmount,
			Data:      []byte{},
		})
		if err != nil {
			s.logger.Errorf("failed to build funding tx for delegator %d: %v", delegatorIndex, err)
			delegatorIndex++
			continue
		}

		tx, err := wallet.BuildDynamicFeeTx(txData)
		if err != nil {
			s.logger.Errorf("failed to build funding transaction for delegator %d: %v", delegatorIndex, err)
			delegatorIndex++
			continue
		}

		// Send funding transaction with no retries to avoid duplicates
		err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
			Client:          client,
			MaxRebroadcasts: 0, // No retries to avoid duplicates
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					return // Don't log individual failures
				}
				if receipt != nil && receipt.Status == 1 {
					atomic.AddInt64(&confirmedCount, 1)

					// Add successfully funded delegator to EOA queue
					s.addEOAToQueue(delegator.GetAddress().Hex(), fmt.Sprintf("%x", delegator.GetPrivateKey().D))

					// No progress logging - only log when target is reached
				}
			},
			LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
				// Only log actual send failures, not confirmation failures
				if err != nil {
					s.logger.Debugf("funding tx send failed: %v", err)
				}
			},
		})

		if err != nil {
			delegatorIndex++
			continue
		}

		sentCount++
		delegatorIndex++

		// Check if we should continue filling the block
		confirmed = atomic.LoadInt64(&confirmedCount)
		if confirmed >= int64(targetCount) && sentCount%maxTxsPerBlock < 100 {
			// We have enough confirmed and we're at the end of a block cycle, stop for now
			break
		}

		// Small delay between transactions to ensure proper nonce ordering
		// Reduce delay as we get more efficient
		if sentCount < 100 {
			time.Sleep(10 * time.Millisecond)
		} else {
			time.Sleep(5 * time.Millisecond)
		}

		// Add context cancellation check
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	// Wait for any remaining transactions to be included
	time.Sleep(3 * time.Second)

	return nil
}

func (s *Scenario) sendMaxBloatingTransaction(ctx context.Context, targetAuthorizations int, targetGasLimit uint64, blockCounter int) (uint64, string, int, float64, float64, string, error) {
	client := s.walletPool.GetClient(spamoor.SelectClientByIndex, 0, s.options.ClientGroup)
	if client == nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("no client available for sending max bloating transaction")
	}

	// Use root wallet since we set child wallet count to 0 in max-bloating mode
	wallet := s.walletPool.GetRootWallet()

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
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to get suggested fees: %w", err)
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
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to parse call data: %w", err)
		}
		txCallData = dataBytes
	}

	// Build the authorizations for maximum state bloat
	authorizations := s.buildMaxBloatingAuthorizations(targetAuthorizations, blockCounter)

	// Check transaction size and split into batches if needed
	batches := s.splitAuthorizationsBatches(authorizations, len(txCallData))

	if len(batches) == 1 {
		// Single transaction - use existing logic
		return s.sendSingleMaxBloatingTransaction(ctx, batches[0], txCallData, feeCap, tipCap, amount, toAddr, targetGasLimit, wallet, client)
	} else {
		// Multiple transactions needed - send them as a batch
		return s.sendBatchedMaxBloatingTransactions(ctx, batches, txCallData, feeCap, tipCap, amount, toAddr, targetGasLimit, wallet, client)
	}
}

// sendSingleMaxBloatingTransaction sends a single transaction (original logic)
func (s *Scenario) sendSingleMaxBloatingTransaction(ctx context.Context, authorizations []types.SetCodeAuthorization, txCallData []byte, feeCap, tipCap *big.Int, amount *uint256.Int, toAddr common.Address, targetGasLimit uint64, wallet *txbuilder.Wallet, client *txbuilder.Client) (uint64, string, int, float64, float64, string, error) {
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
		return 0, "", 0, 0, 0, "", fmt.Errorf("failed to build transaction metadata: %w", err)
	}

	tx, err := wallet.BuildSetCodeTx(txData)
	if err != nil {
		return 0, "", 0, 0, 0, "", fmt.Errorf("failed to build transaction: %w", err)
	}

	// Log actual transaction size
	txSize := len(tx.Data())
	if encoded, err := tx.MarshalBinary(); err == nil {
		txSize = len(encoded)
	}
	sizeKiB := float64(txSize) / 1024.0
	exceedsLimit := txSize > MaxTransactionSize
	limitKiB := float64(MaxTransactionSize) / 1024.0

	s.logger.WithField("scenario", "setcodetx").Infof("MAX BLOATING TX SIZE: %d bytes (%.2f KiB) | Limit: %d bytes (%.1f KiB) | %d authorizations | Exceeds limit: %v",
		txSize, sizeKiB, MaxTransactionSize, limitKiB, len(authorizations), exceedsLimit)

	// Use channels to capture transaction results
	resultChan := make(chan struct {
		gasUsed      uint64
		blockNumber  string
		authCount    int
		gasPerAuth   float64
		gasPerByte   float64
		gweiTotalFee string
		err          error
	}, 1)

	// Send the transaction
	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
		Client:              client,
		MaxRebroadcasts:     10,
		RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
		OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			if err != nil {
				s.logger.WithField("rpc", client.GetName()).Errorf("max bloating tx failed: %v", err)
				resultChan <- struct {
					gasUsed      uint64
					blockNumber  string
					authCount    int
					gasPerAuth   float64
					gasPerByte   float64
					gweiTotalFee string
					err          error
				}{0, "", 0, 0, 0, "", err}
				return
			}
			if receipt == nil {
				resultChan <- struct {
					gasUsed      uint64
					blockNumber  string
					authCount    int
					gasPerAuth   float64
					gasPerByte   float64
					gweiTotalFee string
					err          error
				}{0, "", 0, 0, 0, "", fmt.Errorf("no receipt received")}
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
			gasPerByte := gasPerAuth / EstimatedBytesPerAuth

			resultChan <- struct {
				gasUsed      uint64
				blockNumber  string
				authCount    int
				gasPerAuth   float64
				gasPerByte   float64
				gweiTotalFee string
				err          error
			}{receipt.GasUsed, receipt.BlockNumber.String(),
				authCount,
				gasPerAuth,
				gasPerByte,
				gweiTotalFee.String(), nil}
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
		return 0, "", 0, 0, 0, "", fmt.Errorf("failed to send max bloating transaction: %w", err)
	}

	// Wait for transaction confirmation
	result := <-resultChan
	return result.gasUsed, result.blockNumber, result.authCount, result.gasPerAuth, result.gasPerByte, result.gweiTotalFee, result.err
}

// sendBatchedMaxBloatingTransactions sends multiple transactions when size limit is exceeded
func (s *Scenario) sendBatchedMaxBloatingTransactions(ctx context.Context, batches [][]types.SetCodeAuthorization, txCallData []byte, feeCap, tipCap *big.Int, amount *uint256.Int, toAddr common.Address, targetGasLimit uint64, wallet *txbuilder.Wallet, client *txbuilder.Client) (uint64, string, int, float64, float64, string, error) {
	s.logger.Infof("Transaction size limit exceeded, sending %d batched transactions with minimal delay", len(batches))

	// Aggregate results
	var totalGasUsed uint64
	var totalAuthCount int
	var totalFees *big.Int = big.NewInt(0)
	var lastBlockNumber string

	// Create result channels for all batches upfront
	resultChans := make([]chan struct {
		gasUsed      uint64
		blockNumber  string
		authCount    int
		gweiTotalFee string
		err          error
	}, len(batches))

	// Send all batches quickly with minimal delay to increase chance of same block inclusion
	for batchIndex, batch := range batches {
		// Create result channel for this batch
		resultChans[batchIndex] = make(chan struct {
			gasUsed      uint64
			blockNumber  string
			authCount    int
			gweiTotalFee string
			err          error
		}, 1)

		s.logger.Infof("Sending batch %d/%d with %d authorizations", batchIndex+1, len(batches), len(batch))

		// Calculate appropriate gas limit for this batch based on authorization count
		// Each authorization needs ~26000 gas, plus some overhead for the transaction itself
		batchGasLimit := uint64(len(batch))*GasPerAuthorization + BaseTransferCost + uint64(len(txCallData)*16)

		// Ensure we don't exceed the target limit per transaction
		maxGasPerTx := targetGasLimit
		if batchGasLimit > maxGasPerTx {
			batchGasLimit = maxGasPerTx
		}

		// Build the transaction for this batch
		txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Gas:       batchGasLimit,
			To:        &toAddr,
			Value:     amount,
			Data:      txCallData,
			AuthList:  batch,
		})
		if err != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to build batch %d transaction metadata: %w", batchIndex+1, err)
		}

		tx, err := wallet.BuildSetCodeTx(txData)
		if err != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to build batch %d transaction: %w", batchIndex+1, err)
		}

		// Log actual transaction size for this batch
		txSize := len(tx.Data())
		if encoded, err := tx.MarshalBinary(); err == nil {
			txSize = len(encoded)
		}
		sizeKiB := float64(txSize) / 1024.0
		s.logger.WithField("scenario", "setcodetx").Infof("BATCH %d/%d TX SIZE: %d bytes (%.2f KiB) | %d authorizations | Gas limit: %d",
			batchIndex+1, len(batches), txSize, sizeKiB, len(batch), batchGasLimit)

		// Send the transaction immediately without waiting for confirmation
		resultChan := resultChans[batchIndex]
		err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &txbuilder.SendTransactionOptions{
			Client:              client,
			MaxRebroadcasts:     10,
			RebroadcastInterval: time.Duration(s.options.Rebroadcast) * time.Second,
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					s.logger.WithField("rpc", client.GetName()).Errorf("batch %d tx failed: %v", batchIndex+1, err)
					resultChan <- struct {
						gasUsed      uint64
						blockNumber  string
						authCount    int
						gweiTotalFee string
						err          error
					}{0, "", 0, "", err}
					return
				}
				if receipt == nil {
					resultChan <- struct {
						gasUsed      uint64
						blockNumber  string
						authCount    int
						gweiTotalFee string
						err          error
					}{0, "", 0, "", fmt.Errorf("batch %d: no receipt received", batchIndex+1)}
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

				s.logger.WithField("rpc", client.GetName()).Infof("Batch %d/%d confirmed: %d gas, block %s",
					batchIndex+1, len(batches), receipt.GasUsed, receipt.BlockNumber.String())

				resultChan <- struct {
					gasUsed      uint64
					blockNumber  string
					authCount    int
					gweiTotalFee string
					err          error
				}{receipt.GasUsed, receipt.BlockNumber.String(), len(batch), gweiTotalFee.String(), nil}
			},
			LogFn: func(client *txbuilder.Client, retry int, rebroadcast int, err error) {
				logger := s.logger.WithField("rpc", client.GetName())
				if err != nil {
					logger.Errorf("failed sending batch %d tx: %v", batchIndex+1, err)
				} else if retry > 0 || rebroadcast > 0 {
					logger.Debugf("successfully sent batch %d tx (retry/rebroadcast)", batchIndex+1)
				}
			},
		})

		if err != nil {
			wallet.ResetPendingNonce(ctx, client)
			return 0, "", 0, 0, 0, "", fmt.Errorf("failed to send batch %d transaction: %w", batchIndex+1, err)
		}

		// Small delay to avoid overwhelming the network but keep transactions close together
		if batchIndex < len(batches)-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	// Now wait for all batch confirmations
	s.logger.Infof("All %d batches sent, waiting for confirmations...", len(batches))
	blockNumbers := make(map[string]int) // Track which blocks contain our transactions
	for batchIndex := range batches {
		result := <-resultChans[batchIndex]
		if result.err != nil {
			return 0, "", 0, 0, 0, "", fmt.Errorf("batch %d failed: %w", batchIndex+1, result.err)
		}

		// Aggregate successful results
		totalGasUsed += result.gasUsed
		totalAuthCount += result.authCount
		lastBlockNumber = result.blockNumber

		// Track block numbers
		blockNumbers[result.blockNumber]++

		// Parse and add fee
		if feeGwei, ok := new(big.Int).SetString(result.gweiTotalFee, 10); ok {
			totalFees.Add(totalFees, feeGwei)
		}
	}

	// Log block distribution
	s.logger.Infof("Bloating transactions were included in %d blocks:", len(blockNumbers))
	for blockNum, txCount := range blockNumbers {
		s.logger.Infof("  Block #%s: %d bloating transaction(s)", blockNum, txCount)
	}

	// Calculate aggregate metrics
	gasPerAuth := float64(totalGasUsed) / float64(totalAuthCount)
	gasPerByte := gasPerAuth / EstimatedBytesPerAuth

	// Log summary of batched transactions
	s.logger.WithField("scenario", "setcodetx").Infof("BATCHED MAX BLOATING SUMMARY: %d batches sent | Total: %d gas, %d auths, %s gwei fees",
		len(batches), totalGasUsed, totalAuthCount, totalFees.String())

	return totalGasUsed, lastBlockNumber, totalAuthCount, gasPerAuth, gasPerByte, totalFees.String(), nil
}

// eoaWorker runs in a separate goroutine and writes funded EOAs to EOAs.json
// when the semaphore is open (green). It sleeps when the semaphore is closed (red).
func (s *Scenario) eoaWorker() {
	defer s.workerWg.Done()

	for {
		select {
		case <-s.workerDone:
			// Shutdown signal received
			return
		case <-s.workerSemaphore:
			// Semaphore is green, process the queue
			s.processEOAQueue()
		}
	}
}

// processEOAQueue drains the EOA queue and writes entries to EOAs.json
func (s *Scenario) processEOAQueue() {
	for {
		// Check if there are items in the queue
		s.eoaQueueMutex.Lock()
		if len(s.eoaQueue) == 0 {
			s.eoaQueueMutex.Unlock()
			return // Queue is empty, exit processing
		}

		// Dequeue all items (FIFO)
		eoasToWrite := make([]EOAEntry, len(s.eoaQueue))
		copy(eoasToWrite, s.eoaQueue)
		s.eoaQueue = s.eoaQueue[:0] // Clear the queue
		s.eoaQueueMutex.Unlock()

		// Write to file
		err := s.writeEOAsToFile(eoasToWrite)
		if err != nil {
			s.logger.Errorf("failed to write EOAs to file: %v", err)
			// Re-queue the items if write failed
			s.eoaQueueMutex.Lock()
			s.eoaQueue = append(eoasToWrite, s.eoaQueue...)
			s.eoaQueueMutex.Unlock()
			return
		}

	}
}

// writeEOAsToFile appends EOA entries to EOAs.json file
func (s *Scenario) writeEOAsToFile(eoas []EOAEntry) error {
	if len(eoas) == 0 {
		return nil
	}

	fileName := "EOAs.json"

	// Read existing entries if file exists
	var existingEntries []EOAEntry
	if data, err := os.ReadFile(fileName); err == nil {
		json.Unmarshal(data, &existingEntries)
	}

	// Append new entries
	allEntries := append(existingEntries, eoas...)

	// Write back to file
	data, err := json.MarshalIndent(allEntries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal EOA entries: %w", err)
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write EOAs.json: %w", err)
	}

	return nil
}

// addEOAToQueue adds a funded EOA to the queue
func (s *Scenario) addEOAToQueue(address, privateKey string) {
	s.eoaQueueMutex.Lock()
	defer s.eoaQueueMutex.Unlock()

	entry := EOAEntry{
		Address:    address,
		PrivateKey: privateKey,
	}

	s.eoaQueue = append(s.eoaQueue, entry)
}

// openWorkerSemaphore opens the semaphore (green light) allowing the worker to process
func (s *Scenario) openWorkerSemaphore() {
	select {
	case s.workerSemaphore <- struct{}{}:
		// Semaphore opened successfully
	default:
		// Semaphore already open, do nothing
	}
}

// closeWorkerSemaphore closes the semaphore (red light) putting the worker to sleep
func (s *Scenario) closeWorkerSemaphore() {
	select {
	case <-s.workerSemaphore:
		// Semaphore closed successfully
	default:
		// Semaphore already closed, do nothing
	}
}

// shutdownWorker signals the worker to stop and waits for it to finish
func (s *Scenario) shutdownWorker() {
	if s.options.MaxBloating {
		close(s.workerDone)
		s.workerWg.Wait()
	}
}

// calculateTransactionSize estimates the serialized size of a transaction with given authorizations
func (s *Scenario) calculateTransactionSize(authCount int, callDataSize int) int {
	// Rough estimation based on RLP encoding structure:
	// - Base transaction overhead: ~200 bytes
	// - Each SetCodeAuthorization: ~135 bytes (ChainId(8) + Address(20) + Nonce(8) + YParity(1) + R(32) + S(32) + RLP overhead)
	// - Call data: variable size
	// - Additional RLP encoding overhead: ~50 bytes

	baseSize := 200 + callDataSize + 50
	authSize := authCount * 135
	return baseSize + authSize
}

// splitAuthorizationsBatches splits authorizations into batches that fit within transaction size limit
func (s *Scenario) splitAuthorizationsBatches(authorizations []types.SetCodeAuthorization, callDataSize int) [][]types.SetCodeAuthorization {
	if len(authorizations) == 0 {
		return [][]types.SetCodeAuthorization{authorizations}
	}

	// Calculate how many authorizations can fit in one transaction
	maxAuthsPerTx := (MaxTransactionSize - 200 - callDataSize - 50) / 135
	if maxAuthsPerTx <= 0 {
		s.logger.Warnf("Transaction call data too large, using minimal batch size of 1")
		maxAuthsPerTx = 1
	}

	// If all authorizations fit in one transaction, return as single batch
	if len(authorizations) <= maxAuthsPerTx {
		estimatedSize := s.calculateTransactionSize(len(authorizations), callDataSize)
		s.logger.Infof("All %d authorizations fit in single transaction (estimated size: %d bytes)", len(authorizations), estimatedSize)
		return [][]types.SetCodeAuthorization{authorizations}
	}

	// Split into multiple batches
	var batches [][]types.SetCodeAuthorization
	for i := 0; i < len(authorizations); i += maxAuthsPerTx {
		end := i + maxAuthsPerTx
		if end > len(authorizations) {
			end = len(authorizations)
		}
		batch := authorizations[i:end]
		estimatedSize := s.calculateTransactionSize(len(batch), callDataSize)
		s.logger.Infof("Created batch %d with %d authorizations (estimated size: %d bytes)", len(batches)+1, len(batch), estimatedSize)
		batches = append(batches, batch)
	}

	s.logger.Infof("Split %d authorizations into %d batches (max %d auths per batch)", len(authorizations), len(batches), maxAuthsPerTx)
	return batches
}
