package spamoor

import (
	"context"
	"fmt"
	"sync"

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
	// SubmitCount is the number of times to submit the transaction in the first attempt (default 3)
	SubmitCount int

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

	// Maximum number of pending transactions per wallet
	// If 0, no limit is enforced per wallet
	PendingLimit uint64

	// Maximum number of retries for failed submissions
	// If 0, no retries are attempted
	MaxRetries int

	// Pool of clients to assign to wallet groups
	// If set, assigns client 0 to first wallet, client 1 to second wallet, etc.
	// Cycles through the pool if there are more wallets than clients
	ClientPool  *ClientPool
	ClientGroup string // optional client group filter

	// Optional logging callback called after every LogInterval confirmed transactions
	LogFn func(confirmedCount int, totalCount int)

	// Interval for calling LogFn (number of confirmed transactions)
	// If 0, LogFn is never called
	LogInterval int
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

// batchState tracks the state of transaction submission for a single wallet
type batchState struct {
	txs          []*types.Transaction
	sem          chan struct{}
	completeChan chan int // Signals when a tx at index completes
	errorChan    chan error
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
	batchMap := make(map[*Wallet][]*types.Transaction)
	batchMap[wallet] = txs
	receipts, err := p.SendMultiTransactionBatch(ctx, batchMap, opts)
	if err != nil {
		return nil, err
	}
	return receipts[wallet], nil
}

// SendMultiTransactionBatch submits transactions for multiple wallets with sliding window submission.
// Returns receipts in the same order as input transactions for each wallet.
// Respects both per-wallet PendingLimit and GlobalPendingLimit to control concurrent submissions.
// Implements retry logic with MaxRetries for failed submissions.
//
// Example:
//
//	options := &BatchOptions{
//	    SendTransactionOptions: SendTransactionOptions{Rebroadcast: true},
//	    PendingLimit: 50, // 50 pending per wallet
//	    MaxRetries: 3,    // 3 retries per transaction
//	}
//	receipts, err := submitter.SendMultiTransactionBatch(ctx, walletTxs, options)
func (p *TxPool) SendMultiTransactionBatch(ctx context.Context, walletTxs map[*Wallet][]*types.Transaction, opts *BatchOptions) (map[*Wallet][]*types.Receipt, error) {
	if len(walletTxs) == 0 {
		return make(map[*Wallet][]*types.Receipt), nil
	}

	if opts == nil {
		opts = &BatchOptions{
			SendTransactionOptions: SendTransactionOptions{Rebroadcast: true},
			MaxRetries:             3,
		}
	}

	// Count total transactions and initialize result structures
	totalTxs := 0
	resultsMutex := sync.Mutex{}
	receipts := make(map[*Wallet][]*types.Receipt)
	errors := make(map[*Wallet][]error)

	// Global confirmed transaction counter for LogFn callback
	var confirmedCount int

	for wallet, txs := range walletTxs {
		totalTxs += len(txs)
		receipts[wallet] = make([]*types.Receipt, len(txs))
		errors[wallet] = make([]error, len(txs))
	}

	if totalTxs == 0 {
		return make(map[*Wallet][]*types.Receipt), nil
	}

	if opts.LogFn != nil {
		opts.LogFn(0, totalTxs)
	}

	// Assign clients to wallets if ClientPool is provided
	walletIndexMap := make(map[*Wallet]int)
	walletIndex := 0
	for wallet := range walletTxs {
		walletIndexMap[wallet] = walletIndex
		walletIndex++
	}

	// Set up limits
	walletLimit := opts.PendingLimit
	if walletLimit == 0 {
		walletLimit = 1000000 // Effectively no limit per wallet
	}

	// Per-wallet state tracking
	walletStates := make(map[*Wallet]*batchState)
	for wallet, txs := range walletTxs {
		walletStates[wallet] = &batchState{
			txs:          txs,
			sem:          make(chan struct{}, walletLimit),
			completeChan: make(chan int, len(txs)),
			errorChan:    make(chan error, 1),
		}
	}

	// Error handling
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	// Start a goroutine for each wallet to manage its sliding window
	for wallet, state := range walletStates {
		wg.Add(1)
		go func(wallet *Wallet, state *batchState) {
			defer wg.Done()

			// Process transactions in order with sliding window
			for txIndex := 0; txIndex < len(state.txs); txIndex++ {
				tx := state.txs[txIndex]

				// Wait for semaphore slot
				select {
				case state.sem <- struct{}{}:
					// Got wallet semaphore
				case <-ctx.Done():
					return
				}

				// Submit transaction with retry logic
				go func(txIndex int, tx *types.Transaction) {
					defer func() {
						<-state.sem // Release semaphore
						state.completeChan <- txIndex
					}()

					maxRetries := opts.MaxRetries
					if maxRetries <= 0 {
						maxRetries = 1
					}

					originalOnComplete := opts.SendTransactionOptions.OnComplete

					var lastErr error
					var receipt *types.Receipt

				attemptLoop:
					for attempt := 0; attempt < maxRetries; attempt++ {
						select {
						case <-ctx.Done():
							errors[wallet][txIndex] = ctx.Err()
							return
						default:
						}

						// Create completion callback
						sendOpts := opts.SendTransactionOptions

						// Override client if assigned
						if opts.ClientPool != nil {
							sendOpts.Client = opts.ClientPool.GetClient(SelectClientByIndex, walletIndexMap[wallet]+attempt, opts.ClientGroup)
						}

						completed := make(chan struct {
							receipt *types.Receipt
							err     error
						}, 1)

						sendOpts.OnComplete = func(tx *types.Transaction, receipt *types.Receipt, err error) {
							// Send completion signal
							completed <- struct {
								receipt *types.Receipt
								err     error
							}{receipt, err}
						}

						// Submit transaction
						err := p.SendTransaction(ctx, wallet, tx, &sendOpts)
						if err != nil {
							lastErr = err
							if attempt == maxRetries-1 {
								// Last attempt failed at submission level
								break attemptLoop
							}
							continue attemptLoop // Retry
						}

						// Wait for transaction completion
						select {
						case <-ctx.Done():
							lastErr = ctx.Err()
							break attemptLoop
						case result := <-completed:
							receipt = result.receipt
							lastErr = result.err

							// If transaction failed but we have retries left, retry
							if result.err != nil {
								if attempt == maxRetries-1 {
									// Last attempt failed at confirmation level
									break attemptLoop
								}
								continue attemptLoop
							} else if receipt != nil && sendOpts.OnConfirm != nil {
								sendOpts.OnConfirm(tx, receipt)
							}

							break attemptLoop // Success or final failure
						}
					}

					var finalErr error
					if lastErr != nil {
						finalErr = fmt.Errorf("failed to submit after %d attempts: %w", maxRetries, lastErr)
						state.errorChan <- lastErr // Signal hard error
					}

					resultsMutex.Lock()
					receipts[wallet][txIndex] = receipt
					errors[wallet][txIndex] = finalErr

					// Track confirmed transactions and call LogFn if needed
					callLogFn := false
					if opts.LogFn != nil && opts.LogInterval > 0 {
						confirmedCount++
						if confirmedCount%opts.LogInterval == 0 {
							callLogFn = true
						}
					}
					confirmedCountCopy := confirmedCount
					resultsMutex.Unlock()

					if callLogFn {
						opts.LogFn(confirmedCountCopy, totalTxs)
					}

					if originalOnComplete != nil {
						originalOnComplete(tx, receipt, lastErr)
					}
				}(txIndex, tx)
			}

			// Wait for all transactions to complete
			for completedCount := 0; completedCount < len(state.txs); {
				select {
				case <-ctx.Done():
					return
				case <-state.completeChan:
					completedCount++
				case err := <-state.errorChan:
					// Hard error occurred, cancel everything
					if err != nil {
						cancel()
						return
					}
				}
			}
		}(wallet, state)
	}

	// Wait for all wallets to complete
	wg.Wait()

	// Check for any hard errors and return the first one found
	var firstError error
	for wallet, walletErrors := range errors {
		for i, err := range walletErrors {
			if err != nil && firstError == nil {
				firstError = fmt.Errorf("wallet %v transaction %d failed: %w", wallet, i, err)
			}
		}
	}

	return receipts, firstError
}
