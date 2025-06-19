package spamoor

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// WalletSelectionMode defines how wallets are selected from the pool.
type WalletSelectionMode uint8

var (
	// SelectWalletByIndex selects a wallet by index (modulo pool size).
	SelectWalletByIndex WalletSelectionMode = 0
	// SelectWalletRandom selects a random wallet from the pool.
	SelectWalletRandom WalletSelectionMode = 1
	// SelectWalletRoundRobin selects wallets in round-robin fashion.
	SelectWalletRoundRobin WalletSelectionMode = 2
	// SelectWalletByPendingTxCount selects a wallet by pending tx count (lowest pending tx count first).
	SelectWalletByPendingTxCount WalletSelectionMode = 3
)

// WalletPoolConfig contains configuration settings for the wallet pool including
// wallet count, funding amounts, and automatic refill behavior.
type WalletPoolConfig struct {
	WalletCount    uint64       `yaml:"wallet_count,omitempty"`
	RefillAmount   *uint256.Int `yaml:"refill_amount"`
	RefillBalance  *uint256.Int `yaml:"refill_balance"`
	RefillInterval uint64       `yaml:"refill_interval"`
	WalletSeed     string       `yaml:"seed"`
}

// WellKnownWalletConfig defines configuration for a named wallet with custom funding settings.
// Well-known wallets have specific names and can have different refill amounts than regular wallets.
type WellKnownWalletConfig struct {
	Name          string
	RefillAmount  *uint256.Int
	RefillBalance *uint256.Int
	VeryWellKnown bool
}

// WalletPool manages a pool of child wallets derived from a root wallet with automatic funding
// and balance monitoring. It provides wallet selection strategies, automatic refills when balances
// drop below thresholds, and batch funding operations for efficiency.
type WalletPool struct {
	ctx         context.Context
	config      WalletPoolConfig
	logger      logrus.FieldLogger
	rootWallet  *RootWallet
	clientPool  *ClientPool
	txpool      *TxPool
	runFundings bool

	childWallets     []*Wallet
	wellKnownNames   []*WellKnownWalletConfig
	wellKnownWallets map[string]*Wallet
	selectionMutex   sync.Mutex
	rrWalletIdx      int
	reclaimedFunds   bool

	// Optional callback to track transaction results for metrics
	transactionTracker func(err error)
}

// FundingRequest represents a request to fund a wallet with a specific amount.
// Used internally for batch funding operations.
type FundingRequest struct {
	Wallet *Wallet
	Amount *uint256.Int
}

// GetDefaultWalletConfig returns default wallet pool configuration for a given scenario.
// It generates a random seed and sets reasonable defaults for refill amounts and intervals.
func GetDefaultWalletConfig(scenarioName string) *WalletPoolConfig {
	return &WalletPoolConfig{
		WalletSeed:     fmt.Sprintf("%v-%v", scenarioName, rand.Intn(1000000)),
		WalletCount:    0,
		RefillAmount:   uint256.NewInt(5000000000000000000),
		RefillBalance:  uint256.NewInt(1000000000000000000),
		RefillInterval: 600,
	}
}

// NewWalletPool creates a new wallet pool with the specified dependencies.
// The pool must be configured and prepared with PrepareWallets() before use.
func NewWalletPool(ctx context.Context, logger logrus.FieldLogger, rootWallet *RootWallet, clientPool *ClientPool, txpool *TxPool) *WalletPool {
	return &WalletPool{
		ctx:              ctx,
		logger:           logger,
		rootWallet:       rootWallet,
		clientPool:       clientPool,
		txpool:           txpool,
		childWallets:     make([]*Wallet, 0),
		wellKnownWallets: make(map[string]*Wallet),
		runFundings:      true,
	}
}

// GetContext returns the context associated with this wallet pool.
func (pool *WalletPool) GetContext() context.Context {
	return pool.ctx
}

// GetTxPool returns the transaction pool used by this wallet pool.
func (pool *WalletPool) GetTxPool() *TxPool {
	return pool.txpool
}

// GetSubmitter returns the transaction submitter used by this wallet pool.
func (pool *WalletPool) GetSubmitter() *TxSubmitter {
	return pool.txpool.submitter
}

// GetClientPool returns the client pool used for blockchain interactions.
func (pool *WalletPool) GetClientPool() *ClientPool {
	return pool.clientPool
}

// GetRootWallet returns the root wallet that funds all child wallets.
func (pool *WalletPool) GetRootWallet() *RootWallet {
	return pool.rootWallet
}

// GetChainId returns the chain ID from the root wallet.
func (pool *WalletPool) GetChainId() *big.Int {
	return pool.rootWallet.wallet.GetChainId()
}

// LoadConfig loads wallet pool configuration from YAML string.
func (pool *WalletPool) LoadConfig(configYaml string) error {
	err := yaml.Unmarshal([]byte(configYaml), &pool.config)
	if err != nil {
		return err
	}

	return nil
}

// MarshalConfig returns the current configuration as a YAML string.
func (pool *WalletPool) MarshalConfig() (string, error) {
	yamlBytes, err := yaml.Marshal(&pool.config)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}

// SetWalletCount sets the number of child wallets to create.
func (pool *WalletPool) SetWalletCount(count uint64) {
	pool.config.WalletCount = count
}

// SetRunFundings enables or disables automatic wallet funding.
// When disabled, wallets will not be automatically refilled when their balance drops.
func (pool *WalletPool) SetRunFundings(runFundings bool) {
	pool.runFundings = runFundings
}

// AddWellKnownWallet adds a named wallet with custom funding configuration.
// Well-known wallets are created alongside regular numbered wallets.
func (pool *WalletPool) AddWellKnownWallet(config *WellKnownWalletConfig) {
	pool.wellKnownNames = append(pool.wellKnownNames, config)
}

// SetRefillAmount sets the amount sent to wallets when they need funding.
func (pool *WalletPool) SetRefillAmount(amount *uint256.Int) {
	pool.config.RefillAmount = amount
}

// SetRefillBalance sets the balance threshold below which wallets are automatically refilled.
func (pool *WalletPool) SetRefillBalance(balance *uint256.Int) {
	pool.config.RefillBalance = balance
}

// SetWalletSeed sets the seed used for deterministic wallet generation.
// The same seed will always generate the same set of wallets.
func (pool *WalletPool) SetWalletSeed(seed string) {
	pool.config.WalletSeed = seed
}

// SetRefillInterval sets the interval in seconds between automatic balance checks.
func (pool *WalletPool) SetRefillInterval(interval uint64) {
	pool.config.RefillInterval = interval
}

// SetTransactionTracker sets the optional callback to track transaction results for metrics.
func (pool *WalletPool) SetTransactionTracker(tracker func(err error)) {
	pool.transactionTracker = tracker
}

// GetTransactionTracker returns the transaction tracking callback if set.
func (pool *WalletPool) GetTransactionTracker() func(err error) {
	return pool.transactionTracker
}

// GetClient returns a client from the client pool using the specified selection strategy.
func (pool *WalletPool) GetClient(mode ClientSelectionMode, input int, group string) *Client {
	return pool.clientPool.GetClient(mode, input, group)
}

// GetWallet returns a wallet from the pool using the specified selection strategy.
// Returns nil if no wallets are available.
func (pool *WalletPool) GetWallet(mode WalletSelectionMode, input int) *Wallet {
	pool.selectionMutex.Lock()
	defer pool.selectionMutex.Unlock()

	if len(pool.childWallets) == 0 {
		return nil
	}

	switch mode {
	case SelectWalletByIndex:
		input = input % len(pool.childWallets)
	case SelectWalletRandom:
		input = rand.Intn(len(pool.childWallets))
	case SelectWalletRoundRobin:
		input = pool.rrWalletIdx
		pool.rrWalletIdx++
		if pool.rrWalletIdx >= len(pool.childWallets) {
			pool.rrWalletIdx = 0
		}
	case SelectWalletByPendingTxCount:
		minPendingCount := uint64(math.MaxUint64)
		minPendingIndexes := []int{}
		for i, wallet := range pool.childWallets {
			pendingCount := wallet.GetNonce() - wallet.GetConfirmedNonce()
			if pendingCount < minPendingCount {
				minPendingCount = pendingCount
				minPendingIndexes = []int{i}
			} else if pendingCount == minPendingCount {
				minPendingIndexes = append(minPendingIndexes, i)
			}
		}
		input = input % len(minPendingIndexes)
		input = minPendingIndexes[input]
	}
	return pool.childWallets[input]
}

// GetWellKnownWallet returns a well-known wallet by name.
// Returns nil if the wallet doesn't exist.
func (pool *WalletPool) GetWellKnownWallet(name string) *Wallet {
	return pool.wellKnownWallets[name]
}

// GetVeryWellKnownWalletAddress derives the address of a "very well known" wallet
// without registering it. Very well known wallets are derived only from the root
// wallet's private key and the wallet name, without any scenario seed.
// This makes them consistent across different scenario runs.
func (pool *WalletPool) GetVeryWellKnownWalletAddress(name string) common.Address {
	idxBytes := make([]byte, len(name))
	copy(idxBytes, name)
	// VeryWellKnown wallets don't use the seed, so we skip adding it

	parentKey := crypto.FromECDSA(pool.rootWallet.wallet.GetPrivateKey())
	childKey := sha256.Sum256(append(parentKey, idxBytes...))

	// Derive private key and then address
	privateKey, err := crypto.HexToECDSA(fmt.Sprintf("%x", childKey))
	if err != nil {
		return common.Address{}
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

// GetWalletName returns a human-readable name for the given wallet address.
// Returns "root" for the root wallet, numbered names for child wallets,
// custom names for well-known wallets, or "unknown" if not found.
func (pool *WalletPool) GetWalletName(address common.Address) string {
	if pool.rootWallet.wallet.GetAddress() == address {
		return "root"
	}

	for i, wallet := range pool.childWallets {
		if wallet.GetAddress() == address {
			return fmt.Sprintf("%d", i+1)
		}
	}

	for _, config := range pool.wellKnownNames {
		wallet := pool.wellKnownWallets[config.Name]
		if wallet != nil && wallet.GetAddress() == address {
			return config.Name
		}
	}

	return "unknown"
}

// GetAllWallets returns a slice containing all wallets (well-known and child wallets).
// The root wallet is not included in this list.
func (pool *WalletPool) GetAllWallets() []*Wallet {
	wallets := make([]*Wallet, len(pool.childWallets)+len(pool.wellKnownWallets))
	for i, config := range pool.wellKnownNames {
		wallets[i] = pool.wellKnownWallets[config.Name]
	}
	copy(wallets[len(pool.wellKnownWallets):], pool.childWallets)
	return wallets
}

// GetConfiguredWalletCount returns the configured number of child wallets.
func (pool *WalletPool) GetConfiguredWalletCount() uint64 {
	return pool.config.WalletCount
}

// GetWalletCount returns the actual number of child wallets created.
func (pool *WalletPool) GetWalletCount() uint64 {
	return uint64(len(pool.childWallets))
}

// PrepareWallets creates all configured wallets and funds them if needed.
// It generates deterministic wallets based on the root wallet and seed,
// then funds any wallets below the refill threshold. Also starts the
// automatic balance monitoring if funding is enabled.
func (pool *WalletPool) PrepareWallets() error {
	if len(pool.childWallets) > 0 {
		return nil
	}

	seed := pool.config.WalletSeed

	if pool.config.WalletCount == 0 && len(pool.wellKnownWallets) == 0 {
		pool.childWallets = make([]*Wallet, 0)
	} else {
		var client *Client
		var fundingReqs []*FundingRequest

		for i := 0; i < 3; i++ {
			client = pool.clientPool.GetClient(SelectClientRandom, 0, "") // send all preparation transactions via this client to avoid rejections due to nonces
			if client == nil {
				return fmt.Errorf("no client available")
			}

			pool.childWallets = make([]*Wallet, 0, pool.config.WalletCount)
			fundingReqs = make([]*FundingRequest, 0, pool.config.WalletCount)

			var walletErr error
			wg := &sync.WaitGroup{}
			wl := make(chan bool, 50)
			walletsMutex := &sync.Mutex{}

			for _, config := range pool.wellKnownNames {
				wg.Add(1)
				wl <- true
				go func(config *WellKnownWalletConfig) {
					defer func() {
						<-wl
						wg.Done()
					}()
					if walletErr != nil {
						fmt.Printf("Error: %v\n", walletErr)
						return
					}

					childWallet, fundingReq, err := pool.prepareWellKnownWallet(config, client, seed)
					if err != nil {
						pool.logger.Errorf("could not prepare well known wallet %v: %v", config.Name, err)
						walletErr = err
						return
					}

					walletsMutex.Lock()
					if fundingReq != nil {
						fundingReqs = append(fundingReqs, fundingReq)
					}
					pool.wellKnownWallets[config.Name] = childWallet
					walletsMutex.Unlock()
				}(config)
			}

			for childIdx := uint64(0); childIdx < pool.config.WalletCount; childIdx++ {
				wg.Add(1)
				wl <- true
				go func(childIdx uint64) {
					defer func() {
						<-wl
						wg.Done()
					}()
					if walletErr != nil {
						fmt.Printf("Error: %v\n", walletErr)
						return
					}

					childWallet, fundingReq, err := pool.prepareChildWallet(childIdx, client, seed)
					if err != nil {
						pool.logger.Errorf("could not prepare child wallet %v: %v", childIdx, err)
						walletErr = err
						return
					}

					walletsMutex.Lock()
					pool.childWallets = append(pool.childWallets, childWallet)
					if fundingReq != nil {
						fundingReqs = append(fundingReqs, fundingReq)
					}
					walletsMutex.Unlock()
				}(childIdx)
			}
			wg.Wait()

			if len(pool.childWallets) > 0 {
				break
			}
		}

		if pool.runFundings && len(fundingReqs) > 0 {
			err := pool.processFundingRequests(fundingReqs)
			if err != nil {
				return err
			}
		}

		for childIdx, childWallet := range pool.childWallets {
			pool.logger.Debugf(
				"initialized child wallet %4d (addr: %v, balance: %v ETH, nonce: %v)",
				childIdx,
				childWallet.GetAddress().String(),
				utils.WeiToEther(uint256.MustFromBig(childWallet.GetBalance())).Uint64(),
				childWallet.GetNonce(),
			)
		}

		pool.logger.Infof("initialized %v child wallets", pool.config.WalletCount)
	}

	// watch wallet balances
	if pool.runFundings {
		go pool.watchWalletBalancesLoop()
	}

	return nil
}

// prepareChildWallet creates a child wallet derived from the root wallet using deterministic key generation.
// It generates a private key by hashing the root wallet's private key with the child index and seed.
// Returns the wallet, funding request (if needed), and any error.
func (pool *WalletPool) prepareChildWallet(childIdx uint64, client *Client, seed string) (*Wallet, *FundingRequest, error) {
	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes, childIdx)
	if seed != "" {
		seedBytes := []byte(seed)
		idxBytes = append(idxBytes, seedBytes...)
	}
	parentKey := crypto.FromECDSA(pool.rootWallet.wallet.GetPrivateKey())
	childKey := sha256.Sum256(append(parentKey, idxBytes...))

	return pool.prepareWallet(fmt.Sprintf("%x", childKey), client, pool.config.RefillAmount, pool.config.RefillBalance)
}

// prepareWellKnownWallet creates a named wallet derived from the root wallet using deterministic key generation.
// It generates a private key by hashing the root wallet's private key with the wallet name and seed.
// Uses custom refill amounts from the config if specified, otherwise falls back to pool defaults.
func (pool *WalletPool) prepareWellKnownWallet(config *WellKnownWalletConfig, client *Client, seed string) (*Wallet, *FundingRequest, error) {
	idxBytes := make([]byte, len(config.Name))
	copy(idxBytes, config.Name)
	if seed != "" && !config.VeryWellKnown {
		seedBytes := []byte(seed)
		idxBytes = append(idxBytes, seedBytes...)
	}
	parentKey := crypto.FromECDSA(pool.rootWallet.wallet.GetPrivateKey())
	childKey := sha256.Sum256(append(parentKey, idxBytes...))

	refillAmount := pool.config.RefillAmount
	refillBalance := pool.config.RefillBalance

	if config.RefillAmount != nil {
		refillAmount = config.RefillAmount
	}
	if config.RefillBalance != nil {
		refillBalance = config.RefillBalance
	}

	return pool.prepareWallet(fmt.Sprintf("%x", childKey), client, refillAmount, refillBalance)
}

// prepareWallet creates a wallet from a private key and checks if it needs funding.
// Updates the wallet's state from the blockchain and creates a funding request if
// the wallet's balance is below the specified refill threshold.
func (pool *WalletPool) prepareWallet(privkey string, client *Client, refillAmount *uint256.Int, refillBalance *uint256.Int) (*Wallet, *FundingRequest, error) {
	childWallet, err := NewWallet(privkey)
	if err != nil {
		return nil, nil, err
	}
	err = client.UpdateWallet(pool.ctx, childWallet)
	if err != nil {
		return nil, nil, err
	}

	var fundingReq *FundingRequest
	if pool.runFundings && childWallet.GetBalance().Cmp(refillBalance.ToBig()) < 0 {
		fundingReq = &FundingRequest{
			Wallet: childWallet,
			Amount: refillAmount,
		}
	}
	return childWallet, fundingReq, nil
}

// watchWalletBalancesLoop runs continuously to monitor and refill wallet balances.
// It periodically checks all wallets and funds those below the refill threshold.
// Exits when the context is cancelled or funds have been reclaimed.
func (pool *WalletPool) watchWalletBalancesLoop() {
	sleepTime := time.Duration(pool.config.RefillInterval) * time.Second
	for {
		select {
		case <-pool.ctx.Done():
			return
		case <-time.After(sleepTime):
		}

		if pool.reclaimedFunds {
			return
		}

		err := pool.resupplyChildWallets()
		if err != nil {
			pool.logger.Warnf("could not check & resupply chile wallets: %v", err)
			sleepTime = 1 * time.Minute
		} else {
			sleepTime = time.Duration(pool.config.RefillInterval) * time.Second
		}
	}
}

// resupplyChildWallets checks all wallets and creates funding requests for those below the refill threshold.
// It updates wallet states from the blockchain and processes any needed funding requests in batch.
func (pool *WalletPool) resupplyChildWallets() error {
	client := pool.clientPool.GetClient(SelectClientRandom, 0, "")
	if client == nil {
		return fmt.Errorf("no client available")
	}

	err := client.UpdateWallet(pool.ctx, pool.rootWallet.wallet)
	if err != nil {
		return err
	}

	var walletErr error
	wg := &sync.WaitGroup{}
	wl := make(chan bool, 50)

	wellKnownCount := uint64(len(pool.wellKnownWallets))
	fundingReqs := make([]*FundingRequest, 0, pool.config.WalletCount+wellKnownCount)
	reqsMutex := &sync.Mutex{}

	for idx, config := range pool.wellKnownNames {
		wellKnownWallet := pool.wellKnownWallets[config.Name]
		if wellKnownWallet == nil {
			continue
		}

		wg.Add(1)
		wl <- true
		go func(idx int, childWallet *Wallet, config *WellKnownWalletConfig) {
			defer func() {
				<-wl
				wg.Done()
			}()
			if walletErr != nil {
				return
			}

			refillAmount := pool.config.RefillAmount
			refillBalance := pool.config.RefillBalance

			if config.RefillAmount != nil {
				refillAmount = config.RefillAmount
			}
			if config.RefillBalance != nil {
				refillBalance = config.RefillBalance
			}

			err := client.UpdateWallet(pool.ctx, childWallet)
			if err != nil {
				walletErr = err
				return
			}

			if childWallet.GetBalance().Cmp(refillBalance.ToBig()) < 0 {
				reqsMutex.Lock()
				fundingReqs = append(fundingReqs, &FundingRequest{
					Wallet: childWallet,
					Amount: refillAmount,
				})
				reqsMutex.Unlock()
			}
		}(idx, wellKnownWallet, config)
	}

	for childIdx := uint64(0); childIdx < pool.config.WalletCount; childIdx++ {
		wg.Add(1)
		wl <- true
		go func(childIdx uint64) {
			defer func() {
				<-wl
				wg.Done()
			}()
			if walletErr != nil {
				return
			}

			childWallet := pool.childWallets[childIdx]
			err := client.UpdateWallet(pool.ctx, childWallet)
			if err != nil {
				walletErr = err
				return
			}
			if childWallet.GetBalance().Cmp(pool.config.RefillBalance.ToBig()) < 0 {
				reqsMutex.Lock()
				fundingReqs = append(fundingReqs, &FundingRequest{
					Wallet: childWallet,
					Amount: pool.config.RefillAmount,
				})
				reqsMutex.Unlock()
			}
		}(childIdx)
	}
	wg.Wait()
	if walletErr != nil {
		return walletErr
	}

	if len(fundingReqs) > 0 {
		err := pool.processFundingRequests(fundingReqs)
		if err != nil {
			return err
		}
	} else {
		pool.logger.Infof("checked child wallets (no funding needed)")
	}

	return nil
}

// processFundingRequests handles a batch of funding requests by creating and sending transactions.
// It can use either individual transactions or batch transactions via the batcher contract for efficiency.
// Processes transactions in chunks to avoid overwhelming the network.
func (pool *WalletPool) processFundingRequests(fundingReqs []*FundingRequest) error {
	client := pool.clientPool.GetClient(SelectClientRandom, 0, "")
	if client == nil {
		return fmt.Errorf("no client available")
	}

	reqTxCount := len(fundingReqs)
	batchTxCount := reqTxCount
	batcher := pool.rootWallet.GetTxBatcher()
	if batcher != nil {
		err := batcher.Deploy(pool.ctx, pool.rootWallet.wallet, client)
		if err != nil {
			return fmt.Errorf("failed to deploy batcher: %v", err)
		}

		batchTxCount = len(fundingReqs) / BatcherTxLimit
		if len(fundingReqs)%BatcherTxLimit != 0 {
			batchTxCount++
		}
	}

	return pool.rootWallet.WithWalletLock(pool.ctx, batchTxCount, func() {
		pool.logger.Infof("root wallet is locked, waiting for other funding txs to finish...")
	}, func() error {
		txList := make([]*types.Transaction, 0, batchTxCount)
		batchTxMap := map[common.Hash][]*FundingRequest{}
		if batcher != nil {
			for txIdx := 0; txIdx < reqTxCount; txIdx += BatcherTxLimit {
				batch := fundingReqs[txIdx:min(txIdx+BatcherTxLimit, reqTxCount)]
				tx, err := pool.buildWalletFundingBatchTx(batch, client, batcher)
				if err != nil {
					return err
				}
				txList = append(txList, tx)
				batchTxMap[tx.Hash()] = batch
			}
		} else {
			for _, req := range fundingReqs {
				tx, err := pool.buildWalletFundingTx(req.Wallet, client, req.Amount)
				if err != nil {
					return err
				}
				txList = append(txList, tx)
			}
		}

		pool.logger.Infof("funding child wallets... (0/%v)", len(txList))
		for txIdx := 0; txIdx < len(txList); txIdx += 200 {
			endIdx := txIdx + 200
			if txIdx > 0 {
				pool.logger.Infof("funding child wallets... (%v/%v)", txIdx, len(txList))
			}
			if endIdx > len(txList) {
				endIdx = len(txList)
			}
			_, err := pool.txpool.submitter.SendBatch(pool.ctx, pool.rootWallet.wallet, txList[txIdx:endIdx], &BatchOptions{
				SendTransactionOptions: SendTransactionOptions{
					Client: client,
					OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
						if err != nil {
							pool.logger.Warnf("could not send funding tx %v: %v", tx.Hash().String(), err)
							return
						}

						batch, ok := batchTxMap[tx.Hash()]
						if ok {
							for _, req := range batch {
								req.Wallet.AddBalance(req.Amount.ToBig())
							}
						}
					},
				},
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// buildWalletFundingTx creates a transaction to fund a single wallet with the specified amount.
// It gets suggested fees from the client and builds a dynamic fee transaction.
func (pool *WalletPool) buildWalletFundingTx(childWallet *Wallet, client *Client, refillAmount *uint256.Int) (*types.Transaction, error) {
	if client == nil {
		client = pool.clientPool.GetClient(SelectClientByIndex, 0, "")
		if client == nil {
			return nil, fmt.Errorf("no client available")
		}
	}
	feeCap, tipCap, err := client.GetSuggestedFee(pool.ctx)
	if err != nil {
		return nil, err
	}
	if feeCap.Cmp(big.NewInt(400000000000)) < 0 {
		feeCap = big.NewInt(400000000000)
	}
	if tipCap.Cmp(big.NewInt(200000000000)) < 0 {
		tipCap = big.NewInt(200000000000)
	}

	toAddr := childWallet.GetAddress()
	refillTx, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       21000,
		To:        &toAddr,
		Value:     refillAmount,
	})
	if err != nil {
		return nil, err
	}
	tx, err := pool.rootWallet.wallet.BuildDynamicFeeTx(refillTx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// buildWalletFundingBatchTx creates a transaction to fund multiple wallets using the batcher contract.
// It calculates the total amount needed, encodes the funding requests as calldata,
// and builds a transaction to the batcher contract with appropriate gas limits.
func (pool *WalletPool) buildWalletFundingBatchTx(requests []*FundingRequest, client *Client, batcher *TxBatcher) (*types.Transaction, error) {
	if client == nil {
		client = pool.clientPool.GetClient(SelectClientByIndex, 0, "")
		if client == nil {
			return nil, fmt.Errorf("no client available")
		}
	}
	feeCap, tipCap, err := client.GetSuggestedFee(pool.ctx)
	if err != nil {
		return nil, err
	}
	if feeCap.Cmp(big.NewInt(200000000000)) < 0 {
		feeCap = big.NewInt(200000000000)
	}
	if tipCap.Cmp(big.NewInt(100000000000)) < 0 {
		tipCap = big.NewInt(100000000000)
	}

	totalAmount := uint256.NewInt(0)
	for _, req := range requests {
		totalAmount = totalAmount.Add(totalAmount, req.Amount)
	}

	batchData, err := batcher.GetRequestCalldata(requests)
	if err != nil {
		return nil, err
	}

	toAddr := batcher.GetAddress()
	refillTx, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       BatcherBaseGas + BatcherGasPerTx*uint64(len(requests)),
		To:        &toAddr,
		Value:     totalAmount,
		Data:      batchData,
	})
	if err != nil {
		return nil, err
	}
	tx, err := pool.rootWallet.wallet.BuildDynamicFeeTx(refillTx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// reclaimTx holds a reclaim transaction and its associated wallet for fund recovery operations.
type reclaimTx struct {
	tx     *types.Transaction
	wallet *Wallet
}

// buildWalletReclaimTx creates a transaction to reclaim funds from a child wallet back to the target address.
// It calculates the maximum amount that can be reclaimed after accounting for transaction fees.
// Returns nil if the wallet doesn't have enough balance to cover fees.
func (pool *WalletPool) buildWalletReclaimTx(ctx context.Context, childWallet *Wallet, client *Client, target common.Address) (*types.Transaction, error) {
	if client == nil {
		client = pool.clientPool.GetClient(SelectClientByIndex, 0, "")
		if client == nil {
			return nil, fmt.Errorf("no client available")
		}
	}
	feeCap, tipCap, err := client.GetSuggestedFee(ctx)
	if err != nil {
		return nil, err
	}
	if feeCap.Cmp(big.NewInt(200000000000)) < 0 {
		feeCap = big.NewInt(200000000000)
	}
	if tipCap.Cmp(feeCap) < 0 {
		tipCap = feeCap
	}

	feeAmount := big.NewInt(0).Mul(tipCap, big.NewInt(21000))
	reclaimAmount := big.NewInt(0).Sub(childWallet.GetBalance(), feeAmount)

	if reclaimAmount.Cmp(big.NewInt(0)) <= 0 {
		return nil, nil
	}

	reclaimTx, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       21000,
		To:        &target,
		Value:     uint256.MustFromBig(reclaimAmount),
	})
	if err != nil {
		return nil, err
	}

	tx, err := childWallet.BuildDynamicFeeTx(reclaimTx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// collectPoolWallets adds all wallets (root, child, and well-known) to the provided map.
// This is used by the transaction pool to track which addresses belong to this wallet pool.
func (pool *WalletPool) collectPoolWallets(walletMap map[common.Address]*Wallet) {
	walletMap[pool.rootWallet.wallet.GetAddress()] = pool.rootWallet.wallet
	for _, wallet := range pool.childWallets {
		walletMap[wallet.GetAddress()] = wallet
	}
	for _, wallet := range pool.wellKnownWallets {
		walletMap[wallet.GetAddress()] = wallet
	}
}

// CheckChildWalletBalance checks and refills a specific wallet if needed.
// This can be used to manually trigger funding for a single wallet.
func (pool *WalletPool) CheckChildWalletBalance(childWallet *Wallet) error {
	client := pool.clientPool.GetClient(SelectClientRandom, 0, "")
	if client == nil {
		return fmt.Errorf("no client available")
	}

	balance, err := client.GetBalanceAt(pool.ctx, childWallet.GetAddress())
	if err != nil {
		return err
	}
	childWallet.SetBalance(balance)

	refillAmount := pool.config.RefillAmount
	refillBalance := pool.config.RefillBalance

	for _, config := range pool.wellKnownNames {
		wellKnownWallet := pool.wellKnownWallets[config.Name]
		if wellKnownWallet != nil && wellKnownWallet.GetAddress().String() == childWallet.GetAddress().String() {
			if config.RefillAmount != nil {
				refillAmount = config.RefillAmount
			}
			if config.RefillBalance != nil {
				refillBalance = config.RefillBalance
			}
		}
	}

	if childWallet.GetBalance().Cmp(refillBalance.ToBig()) >= 0 {
		return nil
	}

	return pool.processFundingRequests([]*FundingRequest{
		{
			Wallet: childWallet,
			Amount: refillAmount,
		},
	})
}

// ReclaimFunds reclaims all funds from child wallets back to the root wallet.
// This is typically called when shutting down to consolidate remaining funds.
// After calling this, automatic funding is disabled.
func (pool *WalletPool) ReclaimFunds(ctx context.Context, client *Client) error {
	pool.reclaimedFunds = true

	if client == nil {
		client = pool.clientPool.GetClient(SelectClientRandom, 0, "")
	}
	if client == nil {
		return fmt.Errorf("no client available")
	}

	reclaimMtx := sync.Mutex{}
	reclaimTxs := []*reclaimTx{}
	reclaimWg := sync.WaitGroup{}
	reclaimChan := make(chan struct{}, 100)

	reclaimWallet := func(wallet *Wallet) {
		reclaimWg.Add(1)
		reclaimChan <- struct{}{}

		go func() {
			defer func() {
				<-reclaimChan
				reclaimWg.Done()
			}()

			err := client.UpdateWallet(ctx, wallet)
			if err != nil {
				return
			}

			balance := wallet.GetBalance()
			if balance.Cmp(big.NewInt(0)) == 0 {
				return
			}

			tx, err := pool.buildWalletReclaimTx(ctx, wallet, client, pool.rootWallet.wallet.GetAddress())
			if err != nil {
				return
			}

			reclaimMtx.Lock()
			reclaimTxs = append(reclaimTxs, &reclaimTx{
				tx:     tx,
				wallet: wallet,
			})
			reclaimMtx.Unlock()
		}()
	}

	for _, wallet := range pool.childWallets {
		reclaimWallet(wallet)
	}
	for _, wallet := range pool.wellKnownWallets {
		reclaimWallet(wallet)
	}
	reclaimWg.Wait()

	if len(reclaimTxs) > 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(reclaimTxs))
		for _, tx := range reclaimTxs {
			go func(tx *reclaimTx) {
				pool.logger.Infof("sending reclaim tx %v (%v)", tx.tx.Hash().String(), utils.ReadableAmount(uint256.MustFromBig(tx.tx.Value())))
				pool.txpool.submitter.Send(ctx, tx.wallet, tx.tx, &SendTransactionOptions{
					Client: client,
					OnComplete: func(_ *types.Transaction, receipt *types.Receipt, err error) {
						wg.Done()
						if err != nil {
							pool.logger.Warnf("reclaim tx %v failed: %v", tx.tx.Hash().String(), err)
						}
					},
				})
			}(tx)
		}
		wg.Wait()
	}

	return nil
}
