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
)

type WalletSelectionMode uint8

var (
	SelectWalletByIndex    WalletSelectionMode = 0
	SelectWalletRandom     WalletSelectionMode = 1
	SelectWalletRoundRobin WalletSelectionMode = 2
)

type WalletPoolConfig struct {
	WalletCount    uint64
	WalletPrefund  *uint256.Int
	WalletMinfund  *uint256.Int
	RefillInterval uint64
	WalletSeed     string
}

type WalletPool struct {
	ctx        context.Context
	config     WalletPoolConfig
	logger     logrus.FieldLogger
	rootWallet *txbuilder.Wallet
	clientPool *ClientPool
	txpool     *txbuilder.TxPool

	childWallets   []*txbuilder.Wallet
	selectionMutex sync.Mutex
	rrWalletIdx    int
}

func NewWalletPool(ctx context.Context, logger logrus.FieldLogger, rootWallet *txbuilder.Wallet, clientPool *ClientPool, txpool *txbuilder.TxPool) *WalletPool {
	return &WalletPool{
		ctx:          ctx,
		logger:       logger,
		rootWallet:   rootWallet,
		clientPool:   clientPool,
		txpool:       txpool,
		childWallets: make([]*txbuilder.Wallet, 0),
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

func (pool *WalletPool) SetWalletCount(count uint64) {
	pool.config.WalletCount = count
}

func (pool *WalletPool) SetWalletPrefund(prefund *uint256.Int) {
	pool.config.WalletPrefund = prefund
}

func (pool *WalletPool) SetWalletMinfund(minfund *uint256.Int) {
	pool.config.WalletMinfund = minfund
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

func (pool *WalletPool) GetWalletIndex(address common.Address) int {
	if pool.rootWallet.GetAddress() == address {
		return 0
	}

	for i, wallet := range pool.childWallets {
		if wallet.GetAddress() == address {
			return i + 1
		}
	}

	return -1
}

func (pool *WalletPool) GetWalletCount() uint64 {
	return uint64(len(pool.childWallets))
}

func (pool *WalletPool) PrepareWallets() error {
	if pool.childWallets != nil {
		return nil
	}

	seed := pool.config.WalletSeed

	if pool.config.WalletCount == 0 {
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

					childWallet, fundingTx, err := pool.prepareChildWallet(childIdx, client, seed)
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
				err := pool.sendTxRange(fundingTxList[txIdx:endIdx], client)
				if err != nil {
					return err
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
	go pool.watchWalletBalancesLoop()

	return nil
}

func (pool *WalletPool) prepareChildWallet(childIdx uint64, client *txbuilder.Client, seed string) (*txbuilder.Wallet, *types.Transaction, error) {
	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes, childIdx)
	if seed != "" {
		seedBytes := []byte(seed)
		idxBytes = append(idxBytes, seedBytes...)
	}
	parentKey := crypto.FromECDSA(pool.rootWallet.GetPrivateKey())
	childKey := sha256.Sum256(append(parentKey, idxBytes...))

	childWallet, err := txbuilder.NewWallet(fmt.Sprintf("%x", childKey))
	if err != nil {
		return nil, nil, err
	}
	err = client.UpdateWallet(pool.ctx, childWallet)
	if err != nil {
		return nil, nil, err
	}
	tx, err := pool.buildWalletFundingTx(childWallet, client)
	if err != nil {
		return nil, nil, err
	}
	if tx != nil {
		childWallet.AddBalance(tx.Value())
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
	fundingTxs := make([]*types.Transaction, pool.config.WalletCount)
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
			tx, err := pool.buildWalletFundingTx(childWallet, client)
			if err != nil {
				walletErr = err
				return
			}
			if tx != nil {
				childWallet.AddBalance(tx.Value())
			}

			fundingTxs[childIdx] = tx
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
			err := pool.sendTxRange(fundingTxList[txIdx:endIdx], client)
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
	tx, err := pool.buildWalletFundingTx(childWallet, client)
	if err != nil {
		return nil, err
	}

	if tx != nil {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		var confirmErr error

		err := pool.txpool.SendTransaction(context.Background(), childWallet, tx, &txbuilder.SendTransactionOptions{
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

func (pool *WalletPool) buildWalletFundingTx(childWallet *txbuilder.Wallet, client *txbuilder.Client) (*types.Transaction, error) {
	if childWallet.GetBalance().Cmp(pool.config.WalletMinfund.ToBig()) >= 0 {
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
		Value:     pool.config.WalletPrefund,
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

func (pool *WalletPool) sendTxRange(txList []*types.Transaction, client *txbuilder.Client) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	for idx := range txList {
		err := func(idx int) error {
			tx := txList[idx]

			return pool.txpool.SendTransaction(pool.ctx, pool.rootWallet, tx, &txbuilder.SendTransactionOptions{
				Client: client,
				OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					defer wg.Done()

					if err != nil {
						pool.logger.Warnf("could not send funding tx %v: %v", tx.Hash().String(), err)
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
					pool.rootWallet.SubBalance(totalAmount)
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
