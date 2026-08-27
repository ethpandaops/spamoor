package txbuilder

import (
	"github.com/ethpandaops/spamoor/txtypes"
)

// LegacyTx creates a legacy transaction from the provided transaction metadata.
// It constructs a LegacyTx with gas price, gas limit, recipient address, value, and data.
func LegacyTx(txData *TxMetadata) (*txtypes.LegacyTx, error) {
	tx := txtypes.LegacyTx{
		GasPrice: txData.GasFeeCap.ToBig(),
		Gas:      txData.Gas,
		To:       txData.To,
		Value:    txData.Value.ToBig(),
		Data:     txData.Data,
	}
	return &tx, nil
}
