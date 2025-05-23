package spamoor

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
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
}

type FundingRequest struct {
	Wallet *Wallet
	Amount *uint256.Int
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

func (pool *WalletPool) GetContext() context.Context {
	return pool.ctx
}

func (pool *WalletPool) GetTxPool() *TxPool {
	return pool.txpool
}

func (pool *WalletPool) GetClientPool() *ClientPool {
	return pool.clientPool
}

func (pool *WalletPool) GetRootWallet() *RootWallet {
	return pool.rootWallet
}

func (pool *WalletPool) GetChainId() *big.Int {
	return pool.rootWallet.wallet.GetChainId()
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

func (pool *WalletPool) SetRunFundings(runFundings bool) {
	pool.runFundings = runFundings
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

func (pool *WalletPool) GetClient(mode ClientSelectionMode, input int, group string) *Client {
	return pool.clientPool.GetClient(mode, input, group)
}

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
	}
	return pool.childWallets[input]
}

func (pool *WalletPool) GetWellKnownWallet(name string) *Wallet {
	return pool.wellKnownWallets[name]
}

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

func (pool *WalletPool) GetAllWallets() []*Wallet {
	wallets := make([]*Wallet, len(pool.childWallets)+len(pool.wellKnownWallets))
	for i, config := range pool.wellKnownNames {
		wallets[i] = pool.wellKnownWallets[config.Name]
	}
	copy(wallets[len(pool.wellKnownWallets):], pool.childWallets)
	return wallets
}

func (pool *WalletPool) GetConfiguredWalletCount() uint64 {
	return pool.config.WalletCount
}

func (pool *WalletPool) GetWalletCount() uint64 {
	return uint64(len(pool.childWallets))
}

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

func (pool *WalletPool) prepareWellKnownWallet(config *WellKnownWalletConfig, client *Client, seed string) (*Wallet, *FundingRequest, error) {
	idxBytes := make([]byte, len(config.Name))
	copy(idxBytes, config.Name)
	if seed != "" {
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
		if batcher != nil {
			for txIdx := 0; txIdx < reqTxCount; txIdx += BatcherTxLimit {
				batch := fundingReqs[txIdx:min(txIdx+BatcherTxLimit, reqTxCount)]
				tx, err := pool.buildWalletFundingBatchTx(batch, client, batcher)
				if err != nil {
					return err
				}
				txList = append(txList, tx)
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
			err := pool.txpool.SendAndAwaitTxRange(pool.ctx, pool.rootWallet.wallet, txList[txIdx:endIdx], &SendTransactionOptions{
				Client: client,
				OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					if err != nil {
						pool.logger.Warnf("could not send funding tx %v: %v", tx.Hash().String(), err)
					}
				},
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

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

type reclaimTx struct {
	tx     *types.Transaction
	wallet *Wallet
}

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
				err := pool.txpool.SendTransaction(ctx, tx.wallet, tx.tx, &SendTransactionOptions{
					Client: client,
					OnConfirm: func(_ *types.Transaction, receipt *types.Receipt, err error) {
						defer wg.Done()
						if err != nil {
							pool.logger.Warnf("reclaim tx %v failed: %v", tx.tx.Hash().String(), err)
							return
						}

						effectiveGasPrice := receipt.EffectiveGasPrice
						if effectiveGasPrice == nil {
							effectiveGasPrice = big.NewInt(0)
						}
						feeAmount := new(big.Int).Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))

						tx.wallet.SubBalance(big.NewInt(0).Add(tx.tx.Value(), feeAmount))
					},
				})
				if err != nil {
					pool.logger.Warnf("could not send reclaim tx %v: %v", tx.tx.Hash().String(), err)
				}
			}(tx)
		}
		wg.Wait()
	}

	return nil
}

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

func (pool *WalletPool) collectPoolWallets(walletMap map[common.Address]*Wallet) {
	walletMap[pool.rootWallet.wallet.GetAddress()] = pool.rootWallet.wallet
	for _, wallet := range pool.childWallets {
		walletMap[wallet.GetAddress()] = wallet
	}
	for _, wallet := range pool.wellKnownWallets {
		walletMap[wallet.GetAddress()] = wallet
	}
}
