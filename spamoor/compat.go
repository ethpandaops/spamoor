package spamoor

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethpandaops/spamoor/txtypes"
)

// Compatibility helpers for consumers that use spamoor as a library and hold
// go-ethereum transaction and receipt values.
//
// The engine works on txtypes values throughout, since a transaction type
// go-ethereum has not implemented cannot be represented as a types.Transaction. These
// shims convert at the boundary and fail rather than degrade when a value has no
// go-ethereum equivalent.
//
// This is the only file in package spamoor that may import go-ethereum's core/types;
// scripts/check-core-imports.sh enforces that.

// GethTxConfirmFn is the go-ethereum-typed form of TxConfirmFn.
type GethTxConfirmFn func(tx *types.Transaction, receipt *types.Receipt)

// GethTxCompleteFn is the go-ethereum-typed form of TxCompleteFn.
type GethTxCompleteFn func(tx *types.Transaction, receipt *types.Receipt, err error)

// GethTxEncodeFn is the go-ethereum-typed form of TxEncodeFn.
type GethTxEncodeFn func(tx *types.Transaction) ([]byte, error)

// invokeConfirm dispatches to whichever confirmation callback is set.
func (o *SendTransactionOptions) invokeConfirm(tx *txtypes.Transaction, receipt *txtypes.Receipt) {
	if o.OnConfirm != nil {
		o.OnConfirm(tx, receipt)
	}

	if o.OnGethConfirm == nil {
		return
	}

	gethTx, err := tx.ToGethTx()
	if err != nil {
		return
	}

	o.OnGethConfirm(gethTx, receipt.ToGethReceipt())
}

// invokeComplete dispatches to whichever completion callback is set.
func (o *SendTransactionOptions) invokeComplete(tx *txtypes.Transaction, receipt *txtypes.Receipt, txErr error) {
	if o.OnComplete != nil {
		o.OnComplete(tx, receipt, txErr)
	}

	if o.OnGethComplete == nil {
		return
	}

	gethTx, err := tx.ToGethTx()
	if err != nil {
		return
	}

	o.OnGethComplete(gethTx, receipt.ToGethReceipt(), txErr)
}

// invokeEncode dispatches to whichever encode callback is set. An empty result means
// the caller wants the default encoding.
func (o *SendTransactionOptions) invokeEncode(tx *txtypes.Transaction) ([]byte, error) {
	if o.OnEncode != nil {
		return o.OnEncode(tx)
	}

	if o.OnGethEncode == nil {
		return nil, nil
	}

	gethTx, err := tx.ToGethTx()
	if err != nil {
		return nil, err
	}

	return o.OnGethEncode(gethTx)
}

// hasEncodeFn reports whether any encode callback is set.
func (o *SendTransactionOptions) hasEncodeFn() bool {
	return o.OnEncode != nil || o.OnGethEncode != nil
}

// SendGethTransaction submits a go-ethereum transaction.
//
// Deprecated: build transactions with the wallet's Build*Tx methods and use
// SendTransaction. This shim exists for library consumers holding go-ethereum values.
func (p *TxPool) SendGethTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, opts *SendTransactionOptions) error {
	converted, err := txtypes.FromGethTx(tx)
	if err != nil {
		return fmt.Errorf("failed converting transaction: %w", err)
	}

	return p.SendTransaction(ctx, wallet, converted, opts)
}

// BuildGethDynamicFeeTx builds and signs an EIP-1559 transaction from go-ethereum
// transaction data.
//
// Deprecated: use BuildDynamicFeeTx with txtypes.DynamicFeeTx.
func (wallet *Wallet) BuildGethDynamicFeeTx(txData *types.DynamicFeeTx) (*txtypes.Transaction, error) {
	return wallet.BuildDynamicFeeTx(&txtypes.DynamicFeeTx{
		ChainID:    txData.ChainID,
		Nonce:      txData.Nonce,
		GasTipCap:  txData.GasTipCap,
		GasFeeCap:  txData.GasFeeCap,
		Gas:        txData.Gas,
		To:         txData.To,
		Value:      txData.Value,
		Data:       txData.Data,
		AccessList: txData.AccessList,
	})
}

// BuildGethLegacyTx builds and signs a legacy transaction from go-ethereum
// transaction data.
//
// Deprecated: use BuildLegacyTx with txtypes.LegacyTx.
func (wallet *Wallet) BuildGethLegacyTx(txData *types.LegacyTx) (*txtypes.Transaction, error) {
	return wallet.BuildLegacyTx(&txtypes.LegacyTx{
		Nonce:    txData.Nonce,
		GasPrice: txData.GasPrice,
		Gas:      txData.Gas,
		To:       txData.To,
		Value:    txData.Value,
		Data:     txData.Data,
	})
}

// BuildGethBlobTx builds and signs an EIP-4844 transaction from go-ethereum
// transaction data.
//
// Deprecated: use BuildBlobTx with txtypes.BlobTx.
func (wallet *Wallet) BuildGethBlobTx(txData *types.BlobTx) (*txtypes.Transaction, error) {
	return wallet.BuildBlobTx(&txtypes.BlobTx{
		ChainID:    txData.ChainID,
		Nonce:      txData.Nonce,
		GasTipCap:  txData.GasTipCap,
		GasFeeCap:  txData.GasFeeCap,
		Gas:        txData.Gas,
		To:         txData.To,
		Value:      txData.Value,
		Data:       txData.Data,
		AccessList: txData.AccessList,
		BlobFeeCap: txData.BlobFeeCap,
		BlobHashes: txData.BlobHashes,
		Sidecar:    txData.Sidecar,
	})
}

// BuildGethSetCodeTx builds and signs an EIP-7702 transaction from go-ethereum
// transaction data.
//
// Deprecated: use BuildSetCodeTx with txtypes.SetCodeTx.
func (wallet *Wallet) BuildGethSetCodeTx(txData *types.SetCodeTx) (*txtypes.Transaction, error) {
	return wallet.BuildSetCodeTx(&txtypes.SetCodeTx{
		ChainID:    txData.ChainID,
		Nonce:      txData.Nonce,
		GasTipCap:  txData.GasTipCap,
		GasFeeCap:  txData.GasFeeCap,
		Gas:        txData.Gas,
		To:         txData.To,
		Value:      txData.Value,
		Data:       txData.Data,
		AccessList: txData.AccessList,
		AuthList:   txData.AuthList,
	})
}
