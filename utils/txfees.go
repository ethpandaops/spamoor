package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// TxFees represents the fees associated with a transaction including
// the fee amount and blob fee amount.
type TxFees struct {
	FeeAmount     big.Int
	BlobFeeAmount big.Int
	TxBaseFee     big.Int
	BlobBaseFee   big.Int
}

func GetTransactionFees(tx *types.Transaction, receipt *types.Receipt) *TxFees {
	effectiveGasPrice := receipt.EffectiveGasPrice
	if effectiveGasPrice == nil {
		effectiveGasPrice = big.NewInt(0)
	}
	blobGasPrice := receipt.BlobGasPrice
	if blobGasPrice == nil {
		blobGasPrice = big.NewInt(0)
	}

	txFees := &TxFees{
		TxBaseFee:   *effectiveGasPrice,
		BlobBaseFee: *blobGasPrice,
	}

	txFees.FeeAmount.Mul(effectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
	txFees.BlobFeeAmount.Mul(blobGasPrice, big.NewInt(int64(receipt.BlobGasUsed)))

	return txFees
}

func (txFees *TxFees) TotalFeeGwei() (res big.Int) {
	res.Add(&txFees.FeeAmount, &txFees.BlobFeeAmount)
	res.Div(&res, big.NewInt(1000000000))
	return
}

func (txFees *TxFees) TxFeeGwei() (res big.Int) {
	res.Div(&txFees.FeeAmount, big.NewInt(1000000000))
	return
}

func (txFees *TxFees) TxBaseFeeGwei() (res big.Int) {
	res.Div(&txFees.TxBaseFee, big.NewInt(1000000000))
	return
}

func (txFees *TxFees) BlobFeeGwei() (res big.Int) {
	res.Div(&txFees.BlobFeeAmount, big.NewInt(1000000000))
	return
}

func (txFees *TxFees) BlobBaseFeeGwei() (res big.Int) {
	res.Div(&txFees.BlobBaseFee, big.NewInt(1000000000))
	return
}
