package spamoortypes

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// Wallet defines the interface for an Ethereum wallet with private key management,
// nonce tracking, and balance management. It provides thread-safe operations for
// transaction building, nonce management, and balance updates.
type Wallet interface {
	// Address and Identity Management
	GetAddress() common.Address
	SetAddress(address common.Address)
	GetPrivateKey() *ecdsa.PrivateKey

	// Chain Configuration
	GetChainId() *big.Int
	SetChainId(chainid *big.Int)

	// Nonce Management
	GetNonce() uint64
	GetConfirmedNonce() uint64
	SetConfirmedNonce(nonce uint64)
	GetNextNonce() uint64
	SetNonce(nonce uint64)
	ResetPendingNonce(ctx context.Context, client Client)

	// Balance Management
	GetBalance() *big.Int
	SetBalance(balance *big.Int)
	AddBalance(amount *big.Int)
	SubBalance(amount *big.Int)

	// Transaction Building
	BuildDynamicFeeTx(txData *types.DynamicFeeTx) (*types.Transaction, error)
	BuildBlobTx(txData *types.BlobTx) (*types.Transaction, error)
	BuildSetCodeTx(txData *types.SetCodeTx) (*types.Transaction, error)
	BuildBoundTx(ctx context.Context, txData *txbuilder.TxMetadata, buildFn func(transactOpts *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error)

	// Transaction Replacement
	ReplaceDynamicFeeTx(txData *types.DynamicFeeTx, nonce uint64) (*types.Transaction, error)
	ReplaceBlobTx(txData *types.BlobTx, nonce uint64) (*types.Transaction, error)

	// Transaction Inclusion
	ProcessTransactionInclusion(blockNumber uint64, tx *types.Transaction, receipt *types.Receipt)
	RevertTransactionReceival(tx *types.Transaction)
	ProcessTransactionReceival(tx *types.Transaction)
	ProcessStaleTransactions(blockNumber uint64, nonce uint64)
	GetTxNonceChan(targetNonce uint64) (*NonceStatus, bool)
	GetLastConfirmation() uint64
	SetLastConfirmation(nonce uint64)
	GetPendingTxCount() int
}

// NonceStatus tracks the confirmation status of a transaction with a specific nonce
type NonceStatus struct {
	Receipt *types.Receipt
	Channel chan bool
}
