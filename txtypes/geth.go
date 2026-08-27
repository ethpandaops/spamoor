package txtypes

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// Conversions between spamoor and go-ethereum transactions, used by abigen contract
// bindings and by the compatibility shims in package spamoor.
//
// Both directions round-trip through the wire encoding rather than mapping fields, so
// they need no maintenance as types gain fields and fail loudly on any disagreement.

// FromGethTx converts a go-ethereum transaction into a spamoor transaction. Blob
// sidecars are preserved.
func FromGethTx(tx *types.Transaction) (*Transaction, error) {
	if tx == nil {
		return nil, fmt.Errorf("cannot convert nil transaction")
	}

	encoded, err := tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed encoding go-ethereum transaction: %w", err)
	}

	converted, err := DecodeTx(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed decoding go-ethereum transaction: %w", err)
	}

	if converted.Hash() != tx.Hash() {
		return nil, fmt.Errorf("transaction hash mismatch after conversion: %s != %s", converted.Hash(), tx.Hash())
	}

	return converted, nil
}

// ToGethTx converts a spamoor transaction into a go-ethereum transaction. Transaction
// types go-ethereum cannot represent return an error rather than a degraded value.
func (tx *Transaction) ToGethTx() (*types.Transaction, error) {
	encoded, err := tx.MarshalNetwork()
	if err != nil {
		return nil, fmt.Errorf("failed encoding transaction: %w", err)
	}

	converted := new(types.Transaction)
	if err := converted.UnmarshalBinary(encoded); err != nil {
		return nil, fmt.Errorf("transaction type 0x%02x is not representable in go-ethereum: %w", tx.Type(), err)
	}

	return converted, nil
}

// FromGethReceipt converts a go-ethereum receipt into a spamoor receipt.
func FromGethReceipt(receipt *types.Receipt) *Receipt {
	if receipt == nil {
		return nil
	}

	return &Receipt{
		Type:              receipt.Type,
		PostState:         receipt.PostState,
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		Bloom:             receipt.Bloom,
		Logs:              receipt.Logs,
		TxHash:            receipt.TxHash,
		ContractAddress:   receipt.ContractAddress,
		GasUsed:           receipt.GasUsed,
		EffectiveGasPrice: receipt.EffectiveGasPrice,
		BlobGasUsed:       receipt.BlobGasUsed,
		BlobGasPrice:      receipt.BlobGasPrice,
		BlockHash:         receipt.BlockHash,
		BlockNumber:       receipt.BlockNumber,
		TransactionIndex:  receipt.TransactionIndex,
	}
}

// ToGethReceipt converts a spamoor receipt into a go-ethereum receipt. Type-specific
// extensions have no representation there and are dropped.
func (r *Receipt) ToGethReceipt() *types.Receipt {
	if r == nil {
		return nil
	}

	blockNumber := r.BlockNumber
	if blockNumber == nil {
		blockNumber = new(big.Int)
	}

	return &types.Receipt{
		Type:              r.Type,
		PostState:         r.PostState,
		Status:            r.Status,
		CumulativeGasUsed: r.CumulativeGasUsed,
		Bloom:             r.Bloom,
		Logs:              r.Logs,
		TxHash:            r.TxHash,
		ContractAddress:   r.ContractAddress,
		GasUsed:           r.GasUsed,
		EffectiveGasPrice: r.EffectiveGasPrice,
		BlobGasUsed:       r.BlobGasUsed,
		BlobGasPrice:      r.BlobGasPrice,
		BlockHash:         r.BlockHash,
		BlockNumber:       blockNumber,
		TransactionIndex:  r.TransactionIndex,
	}
}
