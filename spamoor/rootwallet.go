package spamoor

import (
	"context"

	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/txbuilder"
	"github.com/ethpandaops/spamoor/utils"
)

func InitRootWallet(ctx context.Context, privkey string, client *txbuilder.Client, logger logrus.FieldLogger) (*txbuilder.Wallet, error) {
	rootWallet, err := txbuilder.NewWallet(privkey)
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

	return rootWallet, nil
}
