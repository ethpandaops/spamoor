package txbuilder

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

// SetCodeTx creates a set code transaction (EIP-7702) from the provided transaction metadata.
// It constructs a SetCodeTx that can authorize code changes for externally owned accounts.
// The transaction must have a valid 'To' address as it cannot be used for contract deployment.
// Includes authorization list for account code delegation as specified in EIP-7702.
func SetCodeTx(txData *TxMetadata) (*types.SetCodeTx, error) {
	if txData.To == nil {
		return nil, fmt.Errorf("to cannot be nil for setcode transactions")
	}
	tx := types.SetCodeTx{
		GasTipCap:  txData.GasTipCap,
		GasFeeCap:  txData.GasFeeCap,
		Gas:        txData.Gas,
		To:         *txData.To,
		Value:      txData.Value,
		Data:       txData.Data,
		AccessList: txData.AccessList,
		AuthList:   txData.AuthList,
	}
	return &tx, nil
}
