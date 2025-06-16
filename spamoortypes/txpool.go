package spamoortypes

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

// TxConfirmFn is a callback function called when a transaction is confirmed or fails.
// It receives the transaction, receipt (if successful), and any error that occurred.
type TxConfirmFn func(tx *types.Transaction, receipt *types.Receipt, err error)

// TxLogFn is a callback function for logging transaction submission attempts.
// It receives the client used, retry count, rebroadcast count, and any error.
type TxLogFn func(client Client, retry int, rebroadcast int, err error)

// TxRebroadcastFn is a callback function called before each transaction rebroadcast.
// It receives the transaction, send options, and the client being used for rebroadcast.
type TxRebroadcastFn func(tx *types.Transaction, options *SendTransactionOptions, client Client)

// SendTransactionOptions contains options for transaction submission including
// client selection, confirmation callbacks, rebroadcast settings, and logging.
type SendTransactionOptions struct {
	Client             Client
	ClientGroup        string
	ClientsStartOffset int

	OnConfirm     TxConfirmFn
	LogFn         TxLogFn
	OnRebroadcast TxRebroadcastFn

	Rebroadcast      bool
	TransactionBytes []byte
}

// TxPool defines the interface for managing transaction submission, confirmation tracking,
// and chain reorganization handling. It monitors blockchain blocks, tracks transaction
// confirmations, handles reorgs by re-submitting affected transactions, and provides
// transaction awaiting functionality with automatic rebroadcasting.
type TxPool interface {

	// Transaction Management
	SendTransaction(ctx context.Context, wallet Wallet, tx *types.Transaction, options *SendTransactionOptions) error
	AwaitTransaction(ctx context.Context, wallet Wallet, tx *types.Transaction) (*types.Receipt, error)
	SendAndAwaitTxRange(ctx context.Context, wallet Wallet, txs []*types.Transaction, options *SendTransactionOptions) error

	// Gas Limit Management
	GetCurrentGasLimit() uint64
	GetCurrentGasLimitWithInit() (uint64, error)
	InitializeGasLimit() error
}
