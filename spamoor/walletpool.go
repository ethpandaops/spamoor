package spamoor

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
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

type WalletSelectionMode uint8

var (
	SelectWalletByIndex    WalletSelectionMode = 0
	SelectWalletRandom     WalletSelectionMode = 1
	SelectWalletRoundRobin WalletSelectionMode = 2
)

type WalletPoolConfig struct {
	WalletCount    uint64       `yaml:"wallet_count,omitempty"`
	RefillAmount   *uint256.Int `yaml:"refill_amount"`
	RefillBalance  *uint256.Int `yaml:"refill_balance"`
	RefillInterval uint64       `yaml:"refill_interval"`
	WalletSeed     string       `yaml:"seed"`
}

type WellKnownWalletConfig struct {
	Name          string
	RefillAmount  *uint256.Int
	RefillBalance *uint256.Int
}

type WalletPool struct {
	ctx        context.Context
	config     WalletPoolConfig
	logger     logrus.FieldLogger
	rootWallet *txbuilder.Wallet
	clientPool *ClientPool
	txpool     *txbuilder.TxPool

	childWallets     []*txbuilder.Wallet
	wellKnownNames   []*WellKnownWalletConfig
	wellKnownWallets map[string]*txbuilder.Wallet
	selectionMutex   sync.Mutex
	rrWalletIdx      int
}

func GetDefaultWalletConfig(scenarioName string) *WalletPoolConfig {
	return &WalletPoolConfig{
		WalletSeed:     fmt.Sprintf("%v-%v", scenarioName, rand.Intn(1000000)),
		WalletCount:    0,
		RefillAmount:   uint256.NewInt(5000000000000000000),
		RefillBalance:  uint256.NewInt(1000000000000000000),
		RefillInterval: 600,
	}
}

func NewWalletPool(ctx context.Context, logger logrus.FieldLogger, rootWallet *txbuilder.Wallet, clientPool *ClientPool, txpool *txbuilder.TxPool) *WalletPool {
	return &WalletPool{
		ctx:              ctx,
		logger:           logger,
		rootWallet:       rootWallet,
		clientPool:       clientPool,
		txpool:           txpool,
		childWallets:     make([]*txbuilder.Wallet, 0),
		wellKnownWallets: make(map[string]*txbuilder.Wallet),
	}
}

func (pool *WalletPool) GetContext() context.Context {
	return pool.ctx
}

func (pool *WalletPool) GetTxPool() *txbuilder.TxPool {
	return pool.txpool
}

func (pool *WalletPool) GetClientPool() *ClientPool {
	return pool.clientPool
}

func (pool *WalletPool) GetRootWallet() *txbuilder.Wallet {
	return pool.rootWallet
}

func (pool *WalletPool) LoadConfig(configYaml string) error {
	err := yaml.Unmarshal([]byte(configYaml), &pool.config)
	if err != nil {
		return err
	}

	return nil
}

func (pool *WalletPool) MarshalConfig() (string, error) {
	yamlBytes, err := yaml.Marshal(&pool.config)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}

func (pool *WalletPool) SetWalletCount(count uint64) {
	pool.config.WalletCount = count
}

func (pool *WalletPool) AddWellKnownWallet(config *WellKnownWalletConfig) {
	pool.wellKnownNames = append(pool.wellKnownNames, config)
}

func (pool *WalletPool) SetRefillAmount(amount *uint256.Int) {
	pool.config.RefillAmount = amount
}

func (pool *WalletPool) SetRefillBalance(balance *uint256.Int) {
	pool.config.RefillBalance = balance
}

func (pool *WalletPool) SetWalletSeed(seed string) {
	pool.config.WalletSeed = seed
}

func (pool *WalletPool) SetRefillInterval(interval uint64) {
	pool.config.RefillInterval = interval
}

func (pool *WalletPool) GetClient(mode ClientSelectionMode, input int) *txbuilder.Client {
	return pool.clientPool.GetClient(mode, input)
}

func (pool *WalletPool) GetWallet(mode WalletSelectionMode, input int) *txbuilder.Wallet {
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
	}
	return pool.childWallets[input]
}

func (pool *WalletPool) GetWellKnownWallet(name string) *txbuilder.Wallet {
	return pool.wellKnownWallets[name]
}

func (pool *WalletPool) GetWalletName(address common.Address) string {
	if pool.rootWallet.GetAddress() == address {
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

func (pool *WalletPool) GetAllWallets() []*txbuilder.Wallet {
	wallets := make([]*txbuilder.Wallet, len(pool.childWallets)+len(pool.wellKnownWallets))
	for i, config := range pool.wellKnownNames {
		wallets[i] = pool.wellKnownWallets[config.Name]
	}
	copy(wallets[len(pool.wellKnownWallets):], pool.childWallets)
	return wallets
}

func (pool *WalletPool) GetWalletCount() uint64 {
	return uint64(len(pool.childWallets))
}

func (pool *WalletPool) PrepareWallets(runFundings bool) error {
	if len(pool.childWallets) > 0 {
		return nil
	}

	seed := pool.config.WalletSeed

	if pool.config.WalletCount == 0 && len(pool.wellKnownWallets) == 0 {
		pool.childWallets = make([]*txbuilder.Wallet, 0)
	} else {
		var client *txbuilder.Client
		var fundingTxs []*types.Transaction

		for i := 0; i < 3; i++ {
			client = pool.clientPool.GetClient(SelectClientRandom, 0) // send all preparation transactions via this client to avoid rejections due to nonces
			pool.childWallets = make([]*txbuilder.Wallet, 0, pool.config.WalletCount)
			fundingTxs = make([]*types.Transaction, 0, pool.config.WalletCount)

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

					childWallet, fundingTx, err := pool.prepareWellKnownWallet(config, client, seed, runFundings)
					if err != nil {
						pool.logger.Errorf("could not prepare well known wallet %v: %v", config.Name, err)
						walletErr = err
						return
					}

					walletsMutex.Lock()
					fundingTxs = append(fundingTxs, fundingTx)
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

					childWallet, fundingTx, err := pool.prepareChildWallet(childIdx, client, seed, runFundings)
					if err != nil {
						pool.logger.Errorf("could not prepare child wallet %v: %v", childIdx, err)
						walletErr = err
						return
					}

					walletsMutex.Lock()
					pool.childWallets = append(pool.childWallets, childWallet)
					fundingTxs = append(fundingTxs, fundingTx)
					walletsMutex.Unlock()
				}(childIdx)
			}
			wg.Wait()

			if len(pool.childWallets) > 0 {
				break
			}
		}

		if runFundings {
			fundingTxList := make([]*types.Transaction, 0, len(fundingTxs))
			for _, tx := range fundingTxs {
				if tx != nil {
					fundingTxList = append(fundingTxList, tx)
				}
			}

			if len(fundingTxList) > 0 {
				sort.Slice(fundingTxList, func(a int, b int) bool {
					return fundingTxList[a].Nonce() < fundingTxList[b].Nonce()
				})

				pool.logger.Infof("funding child wallets... (0/%v)", len(fundingTxList))
				for txIdx := 0; txIdx < len(fundingTxList); txIdx += 200 {
					endIdx := txIdx + 200
					if txIdx > 0 {
						pool.logger.Infof("funding child wallets... (%v/%v)", txIdx, len(fundingTxList))
					}
					if endIdx > len(fundingTxList) {
						endIdx = len(fundingTxList)
					}
					err := pool.SendTxRange(fundingTxList[txIdx:endIdx], client, pool.rootWallet, func(tx *types.Transaction, receipt *types.Receipt, err error) {
						if err != nil {
							pool.logger.Warnf("could not send funding tx %v: %v", tx.Hash().String(), err)
						}
					})
					if err != nil {
						return err
					}
				}
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
	if runFundings {
		go pool.watchWalletBalancesLoop()
	}

	return nil
}

func (pool *WalletPool) prepareChildWallet(childIdx uint64, client *txbuilder.Client, seed string, runFunding bool) (*txbuilder.Wallet, *types.Transaction, error) {
	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes, childIdx)
	if seed != "" {
		seedBytes := []byte(seed)
		idxBytes = append(idxBytes, seedBytes...)
	}
	parentKey := crypto.FromECDSA(pool.rootWallet.GetPrivateKey())
	childKey := sha256.Sum256(append(parentKey, idxBytes...))

	return pool.prepareWallet(fmt.Sprintf("%x", childKey), client, runFunding, pool.config.RefillAmount, pool.config.RefillBalance)
}

func (pool *WalletPool) prepareWellKnownWallet(config *WellKnownWalletConfig, client *txbuilder.Client, seed string, runFunding bool) (*txbuilder.Wallet, *types.Transaction, error) {
	idxBytes := make([]byte, len(config.Name))
	copy(idxBytes, config.Name)
	if seed != "" {
		seedBytes := []byte(seed)
		idxBytes = append(idxBytes, seedBytes...)
	}
	parentKey := crypto.FromECDSA(pool.rootWallet.GetPrivateKey())
	childKey := sha256.Sum256(append(parentKey, idxBytes...))

	refillAmount := pool.config.RefillAmount
	refillBalance := pool.config.RefillBalance

	if config.RefillAmount != nil {
		refillAmount = config.RefillAmount
	}
	if config.RefillBalance != nil {
		refillBalance = config.RefillBalance
	}

	return pool.prepareWallet(fmt.Sprintf("%x", childKey), client, runFunding, refillAmount, refillBalance)
}

func (pool *WalletPool) prepareWallet(privkey string, client *txbuilder.Client, runFunding bool, refillAmount *uint256.Int, refillBalance *uint256.Int) (*txbuilder.Wallet, *types.Transaction, error) {
	childWallet, err := txbuilder.NewWallet(privkey)
	if err != nil {
		return nil, nil, err
	}
	err = client.UpdateWallet(pool.ctx, childWallet)
	if err != nil {
		return nil, nil, err
	}

	var tx *types.Transaction
	if runFunding {
		tx, err = pool.buildWalletFundingTx(childWallet, client, refillAmount, refillBalance)
		if err != nil {
			return nil, nil, err
		}
		if tx != nil {
			childWallet.AddBalance(tx.Value())
		}
	}
	return childWallet, tx, nil
}

func (pool *WalletPool) watchWalletBalancesLoop() {
	sleepTime := time.Duration(pool.config.RefillInterval) * time.Second
	for {
		select {
		case <-pool.ctx.Done():
			return
		case <-time.After(sleepTime):
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

func (pool *WalletPool) resupplyChildWallets() error {
	client := pool.clientPool.GetClient(SelectClientRandom, 0)

	err := client.UpdateWallet(pool.ctx, pool.rootWallet)
	if err != nil {
		return err
	}

	var walletErr error
	wg := &sync.WaitGroup{}
	wl := make(chan bool, 50)

	wellKnownCount := uint64(len(pool.wellKnownWallets))
	fundingTxs := make([]*types.Transaction, pool.config.WalletCount+wellKnownCount)

	for idx, config := range pool.wellKnownNames {
		wellKnownWallet := pool.wellKnownWallets[config.Name]
		if wellKnownWallet == nil {
			continue
		}

		wg.Add(1)
		wl <- true
		go func(idx int, childWallet *txbuilder.Wallet, config *WellKnownWalletConfig) {
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
			tx, err := pool.buildWalletFundingTx(childWallet, client, refillAmount, refillBalance)
			if err != nil {
				walletErr = err
				return
			}
			if tx != nil {
				childWallet.AddBalance(tx.Value())
			}

			fundingTxs[idx] = tx
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
			tx, err := pool.buildWalletFundingTx(childWallet, client, pool.config.RefillAmount, pool.config.RefillBalance)
			if err != nil {
				walletErr = err
				return
			}
			if tx != nil {
				childWallet.AddBalance(tx.Value())
			}

			fundingTxs[wellKnownCount+childIdx] = tx
		}(childIdx)
	}
	wg.Wait()
	if walletErr != nil {
		return walletErr
	}

	fundingTxList := []*types.Transaction{}
	for _, tx := range fundingTxs {
		if tx != nil {
			fundingTxList = append(fundingTxList, tx)
		}
	}

	if len(fundingTxList) > 0 {
		sort.Slice(fundingTxList, func(a int, b int) bool {
			return fundingTxList[a].Nonce() < fundingTxList[b].Nonce()
		})

		lastNonce := uint64(0)
		for idx, tx := range fundingTxList {
			if idx == 0 {
				lastNonce = tx.Nonce()
				continue
			}

			if tx.Nonce() != lastNonce+1 {
				panic(fmt.Sprintf("Error: nonce mismatch: %v != %v + 1\n", tx.Nonce(), lastNonce))
			}
			lastNonce = tx.Nonce()
		}

		pool.logger.Infof("funding child wallets... (0/%v)", len(fundingTxList))
		for txIdx := 0; txIdx < len(fundingTxList); txIdx += 200 {
			endIdx := txIdx + 200
			if txIdx > 0 {
				pool.logger.Infof("funding child wallets... (%v/%v)", txIdx, len(fundingTxList))
			}
			if endIdx > len(fundingTxList) {
				endIdx = len(fundingTxList)
			}
			err := pool.SendTxRange(fundingTxList[txIdx:endIdx], client, pool.rootWallet, func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					pool.logger.Warnf("could not send funding tx %v: %v", tx.Hash().String(), err)
				}
			})
			if err != nil {
				return err
			}
		}
		pool.logger.Infof("funded child wallets... (%v/%v)", len(fundingTxList), len(fundingTxList))
	} else {
		pool.logger.Infof("checked child wallets (no funding needed)")
	}

	return nil
}

func (pool *WalletPool) CheckChildWalletBalance(childWallet *txbuilder.Wallet) (*types.Transaction, error) {
	client := pool.clientPool.GetClient(SelectClientRandom, 0)
	balance, err := client.GetBalanceAt(pool.ctx, childWallet.GetAddress())
	if err != nil {
		return nil, err
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

	tx, err := pool.buildWalletFundingTx(childWallet, client, refillAmount, refillBalance)
	if err != nil {
		return nil, err
	}

	if tx != nil {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		var confirmErr error

		err := pool.txpool.SendTransaction(pool.ctx, childWallet, tx, &txbuilder.SendTransactionOptions{
			OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
				if err != nil {
					confirmErr = err
				}
				wg.Done()
			},
		})
		if err != nil {
			return tx, err
		}

		wg.Wait()
		if confirmErr != nil {
			return tx, confirmErr
		}
	}

	return tx, nil
}

func (pool *WalletPool) buildWalletFundingTx(childWallet *txbuilder.Wallet, client *txbuilder.Client, refillAmount *uint256.Int, refillBalance *uint256.Int) (*types.Transaction, error) {
	if childWallet.GetBalance().Cmp(refillBalance.ToBig()) >= 0 {
		// no refill needed
		return nil, nil
	}

	if client == nil {
		client = pool.clientPool.GetClient(SelectClientByIndex, 0)
	}
	feeCap, tipCap, err := client.GetSuggestedFee(pool.ctx)
	if err != nil {
		return nil, err
	}
	if feeCap.Cmp(big.NewInt(400000000000)) < 0 {
		feeCap = big.NewInt(400000000000)
	}
	if tipCap.Cmp(big.NewInt(3000000000)) < 0 {
		tipCap = big.NewInt(3000000000)
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
	tx, err := pool.rootWallet.BuildDynamicFeeTx(refillTx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (pool *WalletPool) SendTxRange(txList []*types.Transaction, client *txbuilder.Client, wallet *txbuilder.Wallet, confirmCb func(tx *types.Transaction, receipt *types.Receipt, err error)) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	for idx := range txList {
		err := func(idx int) error {
			tx := txList[idx]

			return pool.txpool.SendTransaction(pool.ctx, wallet, tx, &txbuilder.SendTransactionOptions{
				Client: client,
				OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					defer wg.Done()

					if err != nil {
						if confirmCb != nil {
							confirmCb(tx, receipt, err)
						}
						return
					}

					feeAmount := big.NewInt(0)
					if receipt == nil {
						pool.logger.Warnf("no receipt for funding tx %v", tx.Hash().String())
					} else {
						effectiveGasPrice := receipt.EffectiveGasPrice
						if effectiveGasPrice == nil {
							effectiveGasPrice = big.NewInt(0)
						}
						feeAmount = feeAmount.Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
					}

					totalAmount := big.NewInt(0).Add(tx.Value(), feeAmount)
					wallet.SubBalance(totalAmount)

					if confirmCb != nil {
						confirmCb(tx, receipt, nil)
					}
				},

				MaxRebroadcasts:     10,
				RebroadcastInterval: 30 * time.Second,
			})
		}(idx)

		if err != nil {
			return err
		}
		wg.Add(1)
	}

	wg.Done()
	wg.Wait()
	return nil
}
