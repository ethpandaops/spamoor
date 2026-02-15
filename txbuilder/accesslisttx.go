package txbuilder

import (
	"github.com/ethereum/go-ethereum/core/types"
)

// LegacyTx creates a legacy transaction from the provided transaction metadata.
// It constructs a LegacyTx with gas price, gas limit, recipient address, value, and data.
func AccessListTx(txData *TxMetadata) (*types.AccessListTx, error) {
	tx := types.AccessListTx{
		GasPrice:   txData.GasFeeCap.ToBig(),
		Gas:        txData.Gas,
		To:         txData.To,
		Value:      txData.Value.ToBig(),
		Data:       txData.Data,
		AccessList: txData.AccessList,
	}
	return &tx, nil
}
