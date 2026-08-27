package txbuilder

import (
	"github.com/ethpandaops/spamoor/txtypes"
)

// DynFeeTx creates a dynamic fee transaction (EIP-1559) from the provided transaction metadata.
// It constructs a DynamicFeeTx with gas tip cap, gas fee cap, gas limit, recipient address,
// value, data, and access list. This transaction type supports the EIP-1559 fee market
// with separate base fee and priority fee components.
func DynFeeTx(txData *TxMetadata) (*txtypes.DynamicFeeTx, error) {
	tx := txtypes.DynamicFeeTx{
		GasTipCap:  txData.GasTipCap.ToBig(),
		GasFeeCap:  txData.GasFeeCap.ToBig(),
		Gas:        txData.Gas,
		To:         txData.To,
		Value:      txData.Value.ToBig(),
		Data:       txData.Data,
		AccessList: txData.AccessList,
	}
	return &tx, nil
}
