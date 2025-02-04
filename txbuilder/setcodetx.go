package txbuilder

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

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
