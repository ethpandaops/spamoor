package txbuilder

import (
	"github.com/ethereum/go-ethereum/core/types"
)

// LegacyTx creates a legacy transaction from the provided transaction metadata.
// It constructs a LegacyTx with gas price, gas limit, recipient address, value, and data.
func LegacyTx(txData *TxMetadata) (*types.LegacyTx, error) {
	tx := types.LegacyTx{
		GasPrice: txData.GasFeeCap.ToBig(),
		Gas:      txData.Gas,
		To:       txData.To,
		Value:    txData.Value.ToBig(),
		Data:     txData.Data,
	}
	return &tx, nil
}
