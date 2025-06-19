package spamoor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

// TxConfirmFn is a callback function called when a transaction is confirmed or fails.
// It receives the transaction, receipt (if successful), and any error that occurred.
type TxConfirmFn func(tx *types.Transaction, receipt *types.Receipt)

// TxCompleteFn is a callback function called when transaction processing is complete (confirmed or failed).
// Always called regardless of success/failure.
type TxCompleteFn func(tx *types.Transaction, receipt *types.Receipt, err error)

// TxLogFn is a callback function for logging transaction submission attempts.
// It receives the client used, retry count, rebroadcast count, and any error.
type TxLogFn func(client *Client, retry int, rebroadcast int, err error)

// TxEncodeFn is a callback function called to encode a transaction to bytes.
// It receives the transaction and should return the encoded bytes.
type TxEncodeFn func(tx *types.Transaction) ([]byte, error)

// SendTransactionOptions contains options for transaction submission including
// client selection, confirmation callbacks, rebroadcast settings, and logging.
type SendTransactionOptions struct {
	// Client to use for sending (optional, uses pool selection if nil)
	Client *Client
	// ClientGroup to prefer when selecting clients
	ClientGroup string
	// ClientsStartOffset for client selection
	ClientsStartOffset int

	// Enable reliable rebroadcasting
	Rebroadcast bool

	// Callbacks
	OnConfirm  TxConfirmFn  // Called only if tx was sent successfully and confirmed
	OnComplete TxCompleteFn // Always called when processing completes
	OnEncode   TxEncodeFn   // Called to encode tx to bytes on-demand
	LogFn      TxLogFn      // Custom logging function (uses default if nil)
}

// BatchOptions contains options for batch transaction submission.
type BatchOptions struct {
	SendTransactionOptions

	// Maximum number of pending transactions to allow concurrently
	// If 0, no limit is enforced
	PendingLimit uint64
}

func GetDefaultLogFn(logger logrus.FieldLogger, txTypeName string, txIdx string, tx *types.Transaction) TxLogFn {
	if txTypeName != "" {
		txTypeName = txTypeName + " "
	}

	return func(client *Client, retry int, rebroadcast int, err error) {
		logger = logger.WithField("rpc", client.GetName())
		if tx != nil {
			logger = logger.WithField("nonce", tx.Nonce())
		}

		if retry == 0 && rebroadcast > 0 {
			logger.Debugf("rebroadcasting %stx %s", txTypeName, txIdx)
		}
		if retry > 0 {
			logger = logger.WithField("retry", retry)
		}
		if rebroadcast > 0 {
			logger = logger.WithField("rebroadcast", rebroadcast)
		}
		if err != nil {
			logger.Debugf("failed sending %stx %s: %v", txTypeName, txIdx, err)
		} else if retry > 0 || rebroadcast > 0 {
			logger.Debugf("successfully sent %stx %s", txTypeName, txIdx)
		}
	}
}

// SendTransaction submits a single transaction with the given options.
// This is a lower-level interface that provides access to all callback options.
func (p *TxPool) SendTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, opts *SendTransactionOptions) error {
	return p.submitTransaction(ctx, wallet, tx, opts, true)
}

// Await waits for a transaction to be confirmed and returns its receipt.
// It monitors the blockchain for the transaction and handles reorgs by continuing
// to wait if the transaction gets reorged out of the chain.
func (p *TxPool) AwaitTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction) (*types.Receipt, error) {
	return p.awaitTransaction(ctx, wallet, tx, nil)
}

// SendAndAwaitTransaction submits a transaction with custom options and waits for confirmation.
// Allows specifying client preferences, rebroadcast settings, etc.
//
// Example:
//
//	options := &SendOptions{
//	    Client: specificClient,
//	    Rebroadcast: true,
//	}
//	receipt, err := submitter.SendAndAwaitTransaction(ctx, wallet, tx, options)
func (p *TxPool) SendAndAwaitTransaction(ctx context.Context, wallet *Wallet, tx *types.Transaction, opts *SendTransactionOptions) (*types.Receipt, error) {
	if opts == nil {
		opts = &SendTransactionOptions{Rebroadcast: true}
	}

	resultChan := make(chan struct {
		receipt *types.Receipt
		err     error
	}, 1)

	// Create a copy of options to avoid modifying the original
	sendOpts := *opts

	// Override OnComplete to capture the result
	originalOnComplete := sendOpts.OnComplete
	sendOpts.OnComplete = func(tx *types.Transaction, receipt *types.Receipt, err error) {
		// Call the original callback if it exists
		if originalOnComplete != nil {
			originalOnComplete(tx, receipt, err)
		}

		// Send the result to our channel
		resultChan <- struct {
			receipt *types.Receipt
			err     error
		}{receipt, err}
	}

	// Submit the transaction
	err := p.SendTransaction(ctx, wallet, tx, &sendOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction: %w", err)
	}

	// Wait for confirmation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-resultChan:
		return result.receipt, result.err
	}
}

// SendBatchWithOptions submits multiple transactions with custom options and waits for all confirmations.
// Returns receipts in the same order as input transactions.
// Respects PendingLimit to control concurrent submissions.
//
// Example:
//
//	options := &BatchOptions{
//	    SendOptions: SendOptions{Rebroadcast: true},
//	    PendingLimit: 100,
//	}
//	receipts, err := submitter.SendBatchWithOptions(ctx, wallet, txs, options)
func (p *TxPool) SendTransactionBatch(ctx context.Context, wallet *Wallet, txs []*types.Transaction, opts *BatchOptions) ([]*types.Receipt, error) {
	if len(txs) == 0 {
		return nil, nil
	}

	if opts == nil {
		opts = &BatchOptions{
			SendTransactionOptions: SendTransactionOptions{Rebroadcast: true},
		}
	}

	receipts := make([]*types.Receipt, len(txs))
	errors := make([]error, len(txs))

	// Use semaphore for pending limit if specified
	var pendingSem chan struct{}
	if opts.PendingLimit > 0 {
		pendingSem = make(chan struct{}, opts.PendingLimit)
	}

	var pendingCount int64
	wg := sync.WaitGroup{}
	wg.Add(len(txs))

	// Submit all transactions with potential pending limit
	for i, tx := range txs {
		// Acquire semaphore if pending limit is set
		if pendingSem != nil {
			select {
			case pendingSem <- struct{}{}:
				// Got semaphore, proceed
			case <-ctx.Done():
				// Context cancelled while waiting for semaphore
				errors[i] = ctx.Err()
				wg.Done()
				continue
			}
		}

		atomic.AddInt64(&pendingCount, 1)

		go func(index int, transaction *types.Transaction) {
			sendOpts := opts.SendTransactionOptions // Copy the options

			// Set up completion callback
			originalOnComplete := sendOpts.OnComplete
			sendOpts.OnComplete = func(tx *types.Transaction, receipt *types.Receipt, err error) {
				defer func() {
					atomic.AddInt64(&pendingCount, -1)
					wg.Done()

					// Release semaphore if we're using one
					if pendingSem != nil {
						<-pendingSem
					}
				}()

				receipts[index] = receipt
				errors[index] = err

				// Call the original callback if it exists
				if originalOnComplete != nil {
					originalOnComplete(tx, receipt, err)
				}
			}

			// Submit transaction
			p.SendTransaction(ctx, wallet, transaction, &sendOpts)
		}(i, tx)
	}

	// Wait for all transactions to complete
	wg.Wait()

	// Check for any errors
	var firstError error
	for i, err := range errors {
		if err != nil && firstError == nil {
			firstError = fmt.Errorf("transaction %d failed: %w", i, err)
		}
	}

	return receipts, firstError
}
