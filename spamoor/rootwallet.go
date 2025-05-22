package spamoor

import (
	"context"
	"sync"

	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/utils"
)

type RootWallet struct {
	wallet     *Wallet
	walletLock sync.Mutex
	txbatcher  *TxBatcher
}

func InitRootWallet(ctx context.Context, privkey string, client *Client, logger logrus.FieldLogger) (*RootWallet, error) {
	rootWallet, err := NewWallet(privkey)
	if err != nil {
		return nil, err
	}

	err = client.UpdateWallet(ctx, rootWallet)
	if err != nil {
		return nil, err
	}

	if logger != nil {
		logger.Infof(
			"initialized root wallet (addr: %v balance: %v ETH, nonce: %v)",
			rootWallet.GetAddress().String(),
			utils.WeiToEther(uint256.MustFromBig(rootWallet.GetBalance())).Uint64(),
			rootWallet.GetNonce(),
		)
	}

	return &RootWallet{
		wallet: rootWallet,
	}, nil
}

func (wallet *RootWallet) GetWallet() *Wallet {
	return wallet.wallet
}

func (wallet *RootWallet) WithWalletLock(txCount int, lockedLogFn func(), lockedFn func() error) error {
	if !wallet.walletLock.TryLock() {
		if lockedLogFn != nil {
			lockedLogFn()
		}
		wallet.walletLock.Lock()
	}

	defer wallet.walletLock.Unlock()

	return lockedFn()
}

func (wallet *RootWallet) GetTxBatcher() *TxBatcher {
	return wallet.txbatcher
}

func (wallet *RootWallet) InitTxBatcher(ctx context.Context, txpool *TxPool) {
	wallet.txbatcher = NewTxBatcher(txpool)
}
