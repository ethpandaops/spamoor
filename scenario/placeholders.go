package scenario

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
)

// Create2FactoryWalletName is the very well known wallet that deploys the shared CREATE2 factory.
const Create2FactoryWalletName = "create2-factory-deployer"

// Create2FactoryAddress returns the shared CREATE2 factory deployed by factorydeploytx.
func Create2FactoryAddress(walletPool *spamoor.WalletPool) common.Address {
	return crypto.CreateAddress(walletPool.GetVeryWellKnownWalletAddress(Create2FactoryWalletName), 0)
}

// FactoryPlaceholders returns the placeholders that are fixed for a wallet pool: {factory_address}.
func FactoryPlaceholders(walletPool *spamoor.WalletPool) utils.Placeholders {
	return utils.Placeholders{
		"factory_address": func(string) (string, error) {
			return Create2FactoryAddress(walletPool).Hex(), nil
		},
	}
}

// TxPlaceholders returns the per-transaction placeholders (see utils.TxPlaceholders) plus {factory_address}.
func TxPlaceholders(walletPool *spamoor.WalletPool, txIdx uint64) utils.Placeholders {
	placeholders := utils.TxPlaceholders(txIdx)
	for name, fn := range FactoryPlaceholders(walletPool) {
		placeholders[name] = fn
	}
	return placeholders
}
