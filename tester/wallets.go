package tester

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
	"github.com/holiman/uint256"
)

func (tester *Tester) PrepareWallets(seed string) error {
	rootWallet, err := txbuilder.NewWallet(tester.config.WalletPrivkey)
	if err != nil {
		return err
	}
	tester.rootWallet = rootWallet

	err = tester.GetClient(SelectRandom, 0).UpdateWallet(tester.ctx, tester.rootWallet)
	if err != nil {
		return err
	}

	tester.logger.Infof(
		"initialized root wallet (addr: %v balance: %v ETH, nonce: %v)",
		rootWallet.GetAddress().String(),
		utils.WeiToEther(uint256.MustFromBig(rootWallet.GetBalance())).Uint64(),
		rootWallet.GetNonce(),
	)

	if tester.config.WalletCount == 0 {
		tester.childWallets = make([]*txbuilder.Wallet, 0)
	} else {
		var client *txbuilder.Client
		var fundingTxs []*types.Transaction

		for i := 0; i < 3; i++ {
			client = tester.GetClient(SelectRandom, 0) // send all preparation transactions via this client to avoid rejections due to nonces
			tester.childWallets = make([]*txbuilder.Wallet, 0, tester.config.WalletCount)
			fundingTxs = make([]*types.Transaction, 0, tester.config.WalletCount)

			var walletErr error
			wg := &sync.WaitGroup{}
			wl := make(chan bool, 50)
			walletsMutex := &sync.Mutex{}
			for childIdx := uint64(0); childIdx < tester.config.WalletCount; childIdx++ {
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

					childWallet, fundingTx, err := tester.prepareChildWallet(childIdx, client, seed)
					if err != nil {
						tester.logger.Errorf("could not prepare child wallet %v: %v", childIdx, err)
						walletErr = err
						return
					}

					walletsMutex.Lock()
					tester.childWallets = append(tester.childWallets, childWallet)
					fundingTxs = append(fundingTxs, fundingTx)
					walletsMutex.Unlock()
				}(childIdx)
			}
			wg.Wait()

			if len(tester.childWallets) > 0 {
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

			tester.logger.Infof("funding child wallets... (0/%v)", len(fundingTxList))
			for txIdx := 0; txIdx < len(fundingTxList); txIdx += 200 {
				endIdx := txIdx + 200
				if txIdx > 0 {
					tester.logger.Infof("funding child wallets... (%v/%v)", txIdx, len(fundingTxList))
				}
				if endIdx > len(fundingTxList) {
					endIdx = len(fundingTxList)
				}
				err := tester.sendTxRange(fundingTxList[txIdx:endIdx], client)
				if err != nil {
					return err
				}
			}
		}

		for childIdx, childWallet := range tester.childWallets {
			tester.logger.Debugf(
				"initialized child wallet %4d (addr: %v, balance: %v ETH, nonce: %v)",
				childIdx,
				childWallet.GetAddress().String(),
				utils.WeiToEther(uint256.MustFromBig(childWallet.GetBalance())).Uint64(),
				childWallet.GetNonce(),
			)
		}

		tester.logger.Infof("initialized %v child wallets", tester.config.WalletCount)
	}

	return nil
}

func (tester *Tester) prepareChildWallet(childIdx uint64, client *txbuilder.Client, seed string) (*txbuilder.Wallet, *types.Transaction, error) {
	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes, childIdx)
	if seed != "" {
		seedBytes := []byte(seed)
		idxBytes = append(idxBytes, seedBytes...)
	}
	childKey := sha256.Sum256(append(common.FromHex(tester.config.WalletPrivkey), idxBytes...))

	childWallet, err := txbuilder.NewWallet(fmt.Sprintf("%x", childKey))
	if err != nil {
		return nil, nil, err
	}
	err = client.UpdateWallet(tester.ctx, childWallet)
	if err != nil {
		return nil, nil, err
	}
	tx, err := tester.buildWalletFundingTx(childWallet, client)
	if err != nil {
		return nil, nil, err
	}
	if tx != nil {
		childWallet.AddBalance(tx.Value())
	}
	return childWallet, tx, nil
}

func (tester *Tester) resupplyChildWallets() error {
	client := tester.GetClient(SelectRandom, 0)

	err := client.UpdateWallet(tester.ctx, tester.rootWallet)
	if err != nil {
		return err
	}

	var walletErr error
	wg := &sync.WaitGroup{}
	wl := make(chan bool, 50)
	fundingTxs := make([]*types.Transaction, tester.config.WalletCount)
	for childIdx := uint64(0); childIdx < tester.config.WalletCount; childIdx++ {
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

			childWallet := tester.childWallets[childIdx]
			err := client.UpdateWallet(tester.ctx, childWallet)
			if err != nil {
				walletErr = err
				return
			}
			tx, err := tester.buildWalletFundingTx(childWallet, client)
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

		tester.logger.Infof("funding child wallets... (0/%v)", len(fundingTxList))
		for txIdx := 0; txIdx < len(fundingTxList); txIdx += 200 {
			endIdx := txIdx + 200
			if txIdx > 0 {
				tester.logger.Infof("funding child wallets... (%v/%v)", txIdx, len(fundingTxList))
			}
			if endIdx > len(fundingTxList) {
				endIdx = len(fundingTxList)
			}
			err := tester.sendTxRange(fundingTxList[txIdx:endIdx], client)
			if err != nil {
				return err
			}
		}
		tester.logger.Infof("funded child wallets... (%v/%v)", len(fundingTxList), len(fundingTxList))
	} else {
		tester.logger.Infof("checked child wallets (no funding needed)")
	}

	return nil
}

func (tester *Tester) CheckChildWalletBalance(childWallet *txbuilder.Wallet) (*types.Transaction, error) {
	client := tester.GetClient(SelectRandom, 0)
	balance, err := client.GetBalanceAt(tester.ctx, childWallet.GetAddress())
	if err != nil {
		return nil, err
	}
	childWallet.SetBalance(balance)
	tx, err := tester.buildWalletFundingTx(childWallet, client)
	if err != nil {
		return nil, err
	}

	if tx != nil {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		var confirmErr error

		err := tester.GetTxPool().SendTransaction(context.Background(), childWallet, tx, &txbuilder.SendTransactionOptions{
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

func (tester *Tester) buildWalletFundingTx(childWallet *txbuilder.Wallet, client *txbuilder.Client) (*types.Transaction, error) {
	if childWallet.GetBalance().Cmp(tester.config.WalletMinfund.ToBig()) >= 0 {
		// no refill needed
		return nil, nil
	}

	if client == nil {
		client = tester.GetClient(SelectByIndex, 0)
	}
	feeCap, tipCap, err := client.GetSuggestedFee(tester.ctx)
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
		Value:     tester.config.WalletPrefund,
	})
	if err != nil {
		return nil, err
	}
	tx, err := tester.rootWallet.BuildDynamicFeeTx(refillTx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (tester *Tester) sendTxRange(txList []*types.Transaction, client *txbuilder.Client) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	for idx := range txList {
		err := func(idx int) error {
			tx := txList[idx]

			return tester.GetTxPool().SendTransaction(context.Background(), tester.rootWallet, tx, &txbuilder.SendTransactionOptions{
				Client: client,
				OnConfirm: func(tx *types.Transaction, receipt *types.Receipt, err error) {
					defer wg.Done()

					if err != nil {
						tester.logger.Warnf("could not send funding tx %v: %v", tx.Hash().String(), err)
						return
					}

					feeAmount := big.NewInt(0)
					if receipt == nil {
						tester.logger.Warnf("no receipt for funding tx %v", tx.Hash().String())
					} else {
						effectiveGasPrice := receipt.EffectiveGasPrice
						if effectiveGasPrice == nil {
							effectiveGasPrice = big.NewInt(0)
						}
						feeAmount = feeAmount.Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
					}

					totalAmount := big.NewInt(0).Add(tx.Value(), feeAmount)
					tester.rootWallet.SubBalance(totalAmount)
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
