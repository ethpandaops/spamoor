package framework

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/spamoortypes"
	"github.com/ethpandaops/spamoor/testing/load/validators"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// LoadTestEngine orchestrates load test execution
type LoadTestEngine struct {
	config      *LoadTestConfig
	clients     []spamoortypes.Client
	wallets     []spamoortypes.Wallet
	validator   *validators.TransactionValidator
	metrics     *MetricsCollector
	logger      *logrus.Entry
	
	// Runtime state
	running      bool
	startTime    time.Time
	endTime      time.Time
	cancelFunc   context.CancelFunc
	
	// Synchronization
	mutex           sync.RWMutex
	waitGroup       sync.WaitGroup
	lastPendingWarn time.Time // Rate limit pending warnings
}

// NewLoadTestEngine creates a new load test engine
func NewLoadTestEngine(config *LoadTestConfig, clients []spamoortypes.Client, wallets []spamoortypes.Wallet, logger *logrus.Entry) *LoadTestEngine {
	return &LoadTestEngine{
		config:    config,
		clients:   clients,
		wallets:   wallets,
		validator: validators.NewTransactionValidator(logger),
		metrics:   NewMetricsCollector(),
		logger:    logger.WithField("component", "load-test-engine"),
	}
}

// RunLoadTest executes a load test with the configured parameters
func (e *LoadTestEngine) RunLoadTest(ctx context.Context) (*LoadTestResult, error) {
	e.mutex.Lock()
	if e.running {
		e.mutex.Unlock()
		return nil, fmt.Errorf("load test is already running")
	}
	e.running = true
	e.startTime = time.Now()
	e.mutex.Unlock()
	
	defer func() {
		e.mutex.Lock()
		e.running = false
		e.endTime = time.Now()
		e.mutex.Unlock()
	}()
	
	// Create cancellable context
	testCtx, cancel := context.WithTimeout(ctx, e.config.Timeout)
	e.cancelFunc = cancel
	defer cancel()
	
	e.logger.WithFields(logrus.Fields{
		"duration":         e.config.Duration,
		"target_tps":       e.config.TargetTPS,
		"client_count":     e.config.ClientCount,
		"wallet_count":     e.config.WalletCount,
		"transaction_type": e.config.TransactionType,
	}).Info("Starting load test")
	
	// Start transaction validation
	e.validator.Start()
	defer e.validator.Stop()
	
	// Start metrics collection - do this after validator to ensure proper timing
	e.metrics.Start()
	defer e.metrics.Stop()
	
	// Warmup phase
	if e.config.WarmupDuration > 0 {
		e.logger.WithField("duration", e.config.WarmupDuration).Info("Starting warmup phase")
		
		// Start progress indicator for warmup
		warmupCtx, warmupCancel := context.WithTimeout(testCtx, e.config.WarmupDuration)
		defer warmupCancel()
		go e.logProgress(warmupCtx, "Warmup", e.config.WarmupDuration)
		
		if err := e.runWarmup(testCtx); err != nil {
			return nil, fmt.Errorf("warmup failed: %w", err)
		}
	}
	
	// Main load test phase
	e.logger.Info("Starting main load test phase")
	
	// Start progress indicator for main phase
	mainCtx, mainCancel := context.WithTimeout(testCtx, e.config.Duration)
	defer mainCancel()
	go e.logProgress(mainCtx, "Main Load Test", e.config.Duration)
	
	if err := e.runMainPhase(testCtx); err != nil {
		return nil, fmt.Errorf("main phase failed: %w", err)
	}
	
	// Wait for pending transactions to complete
	e.logger.Info("Waiting for pending transactions to complete")
	e.waitGroup.Wait()
	
	// Ensure metrics are stopped before generating report
	e.metrics.Stop()
	
	// Generate final report
	result := e.metrics.GenerateReport(e.config)
	
	e.logger.WithFields(logrus.Fields{
		"total_transactions": result.TotalTransactions,
		"successful_txs":     result.SuccessfulTxs,
		"failed_txs":        result.FailedTxs,
		"average_tps":       result.AverageTPS,
		"average_latency":   result.AverageLatency,
		"error_rate":        result.ErrorRate,
		"test_successful":   result.IsSuccessful(),
	}).Info("Load test completed")
	
	return result, nil
}

// runWarmup runs the warmup phase to prepare the system
func (e *LoadTestEngine) runWarmup(ctx context.Context) error {
	warmupCtx, cancel := context.WithTimeout(ctx, e.config.WarmupDuration)
	defer cancel()
	
	// Run a reduced load during warmup (10% of target TPS)
	warmupTPS := e.config.TargetTPS / 10
	if warmupTPS < 1 {
		warmupTPS = 1
	}
	
	return e.runTransactionLoad(warmupCtx, warmupTPS, false)
}

// runMainPhase runs the main load test phase
func (e *LoadTestEngine) runMainPhase(ctx context.Context) error {
	mainCtx, cancel := context.WithTimeout(ctx, e.config.Duration)
	defer cancel()
	
	return e.runTransactionLoad(mainCtx, e.config.TargetTPS, true)
}

// runTransactionLoad runs transaction load generation
func (e *LoadTestEngine) runTransactionLoad(ctx context.Context, targetTPS int, recordMetrics bool) error {
	interval := time.Second / time.Duration(targetTPS)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	txCount := uint64(0)
	pendingCount := int64(0)
	
	// Track transaction counts for progress reporting
	var lastTxCount uint64
	progressTicker := time.NewTicker(15 * time.Second) // Report every 15 seconds
	defer progressTicker.Stop()
	
	// Start transaction progress reporting
	if recordMetrics {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-progressTicker.C:
					currentCount := txCount
					txRate := float64(currentCount-lastTxCount) / 15.0 // TPS over last 15 seconds
					e.logger.WithFields(logrus.Fields{
						"total_sent":    currentCount,
						"pending":       pendingCount,
						"recent_tps":    fmt.Sprintf("%.1f", txRate),
					}).Info("Transaction progress")
					lastTxCount = currentCount
				}
			}
		}()
	}
	
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			// Check if we have too many pending transactions
			if pendingCount >= int64(e.config.MaxPending) {
				// Rate limit the warning to avoid spam (max once per 5 seconds)
				now := time.Now()
				if now.Sub(e.lastPendingWarn) > 5*time.Second {
					e.logger.WithField("pending_count", pendingCount).Warn("Skipping transactions due to max pending limit")
					e.lastPendingWarn = now
				}
				continue
			}
			
			// Select wallet and client for this transaction
			wallet := e.selectWallet(txCount)
			client := e.selectClient(txCount)
			
			// Increment pending count and wait group
			e.waitGroup.Add(1)
			pendingCount++
			
			// Generate and send transaction asynchronously
			go func(txIdx uint64, w spamoortypes.Wallet, c spamoortypes.Client) {
				defer e.waitGroup.Done()
				defer func() {
					pendingCount--
				}()
				
				if err := e.generateAndSendTransaction(ctx, txIdx, w, c, recordMetrics); err != nil {
					e.logger.WithError(err).WithField("tx_index", txIdx).Debug("Failed to send transaction")
					if recordMetrics {
						e.metrics.RecordCriticalError(fmt.Sprintf("Transaction %d failed: %v", txIdx, err))
					}
				}
			}(txCount, wallet, client)
			
			txCount++
		}
	}
}

// generateAndSendTransaction generates and sends a single transaction
func (e *LoadTestEngine) generateAndSendTransaction(ctx context.Context, txIndex uint64, wallet spamoortypes.Wallet, client spamoortypes.Client, recordMetrics bool) error {
	// Build transaction based on type
	tx, err := e.buildTransaction(wallet)
	if err != nil {
		return fmt.Errorf("failed to build transaction: %w", err)
	}
	
	txHash := tx.Hash()
	
	// Record transaction start time
	if recordMetrics {
		e.metrics.RecordTransactionSent(txHash)
	}
	
	// Capture transaction for validation
	e.validator.CaptureTransaction(tx, wallet.GetAddress())
	
	// Send transaction
	startTime := time.Now()
	err = client.SendTransaction(ctx, tx)
	
	if err != nil {
		if recordMetrics {
			e.metrics.RecordTransactionError(txHash, err)
		}
		return fmt.Errorf("failed to send transaction: %w", err)
	}
	
	// Record successful completion
	if recordMetrics {
		e.metrics.RecordTransactionCompleted(txHash)
	}
	
	e.logger.WithFields(logrus.Fields{
		"tx_hash":    txHash.Hex(),
		"tx_index":   txIndex,
		"wallet":     wallet.GetAddress().Hex(),
		"latency":    time.Since(startTime),
	}).Debug("Transaction sent successfully")
	
	return nil
}

// buildTransaction builds a transaction based on the configured type
func (e *LoadTestEngine) buildTransaction(wallet spamoortypes.Wallet) (*types.Transaction, error) {
	switch e.config.TransactionType {
	case "dynamicfee":
		return e.buildDynamicFeeTransaction(wallet)
	case "blob":
		return e.buildBlobTransaction(wallet)
	case "setcode":
		return e.buildSetCodeTransaction(wallet)
	case "mixed":
		// Randomly select transaction type for mixed testing
		txTypes := []string{"dynamicfee", "blob", "setcode"}
		selectedType := txTypes[int(time.Now().UnixNano())%len(txTypes)]
		e.config.TransactionType = selectedType
		return e.buildTransaction(wallet)
	default:
		return e.buildDynamicFeeTransaction(wallet)
	}
}

// buildDynamicFeeTransaction builds a dynamic fee transaction
func (e *LoadTestEngine) buildDynamicFeeTransaction(wallet spamoortypes.Wallet) (*types.Transaction, error) {
	// Generate a random recipient address
	recipient := common.HexToAddress(fmt.Sprintf("0x%040x", time.Now().UnixNano()))
	
	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.NewInt(e.config.BaseFee + e.config.TipFee),
		GasTipCap: uint256.NewInt(e.config.TipFee),
		Gas:       e.config.GasLimit,
		To:        &recipient,
		Value:     uint256.NewInt(e.config.Amount),
		Data:      nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic fee tx data: %w", err)
	}
	
	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to build dynamic fee transaction: %w", err)
	}
	
	return tx, nil
}

// buildBlobTransaction builds a blob transaction
func (e *LoadTestEngine) buildBlobTransaction(wallet spamoortypes.Wallet) (*types.Transaction, error) {
	// Generate a random recipient address
	recipient := common.HexToAddress(fmt.Sprintf("0x%040x", time.Now().UnixNano()))
	
	// Create simple blob data
	blobData := make([]byte, 1024) // Smaller blob for testing
	for i := range blobData {
		blobData[i] = byte(i % 256)
	}
	
	// Convert blob data to hex string format
	blobHex := "0x" + common.Bytes2Hex(blobData)
	
	txData, err := txbuilder.BuildBlobTx(&txbuilder.TxMetadata{
		GasFeeCap:  uint256.NewInt(e.config.BaseFee + e.config.TipFee),
		GasTipCap:  uint256.NewInt(e.config.TipFee),
		BlobFeeCap: uint256.NewInt(e.config.BaseFee),
		Gas:        e.config.GasLimit,
		To:         &recipient,
		Value:      uint256.NewInt(e.config.Amount),
		Data:       nil,
	}, [][]string{{blobHex}})
	if err != nil {
		return nil, fmt.Errorf("failed to create blob tx data: %w", err)
	}
	
	tx, err := wallet.BuildBlobTx(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to build blob transaction: %w", err)
	}
	
	return tx, nil
}

// buildSetCodeTransaction builds a set code transaction
func (e *LoadTestEngine) buildSetCodeTransaction(wallet spamoortypes.Wallet) (*types.Transaction, error) {
	// Generate a random recipient address
	recipient := common.HexToAddress(fmt.Sprintf("0x%040x", time.Now().UnixNano()))
	
	// Simple contract bytecode (just returns)
	contractCode := []byte{0x60, 0x00, 0x60, 0x00, 0xf3} // PUSH1 0, PUSH1 0, RETURN
	
	txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.NewInt(e.config.BaseFee + e.config.TipFee),
		GasTipCap: uint256.NewInt(e.config.TipFee),
		Gas:       e.config.GasLimit,
		To:        &recipient,
		Value:     uint256.NewInt(e.config.Amount),
		Data:      contractCode,
		AuthList:  []types.SetCodeAuthorization{}, // Empty auth list for testing
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create set code tx data: %w", err)
	}
	
	tx, err := wallet.BuildSetCodeTx(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to build set code transaction: %w", err)
	}
	
	return tx, nil
}

// selectWallet selects a wallet for the given transaction index
func (e *LoadTestEngine) selectWallet(txIndex uint64) spamoortypes.Wallet {
	if len(e.wallets) == 0 {
		return nil
	}
	
	return e.wallets[txIndex%uint64(len(e.wallets))]
}

// selectClient selects a client for the given transaction index
func (e *LoadTestEngine) selectClient(txIndex uint64) spamoortypes.Client {
	if len(e.clients) == 0 {
		return nil
	}
	
	return e.clients[txIndex%uint64(len(e.clients))]
}

// Stop stops the running load test
func (e *LoadTestEngine) Stop() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	
	if e.cancelFunc != nil {
		e.cancelFunc()
	}
}

// IsRunning returns true if the load test is currently running
func (e *LoadTestEngine) IsRunning() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	return e.running
}

// logProgress logs periodic progress updates during test phases
func (e *LoadTestEngine) logProgress(ctx context.Context, phase string, duration time.Duration) {
	ticker := time.NewTicker(10 * time.Second) // Progress every 10 seconds
	defer ticker.Stop()
	
	start := time.Now()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			elapsed := time.Since(start)
			remaining := duration - elapsed
			if remaining < 0 {
				remaining = 0
			}
			
			progress := float64(elapsed) / float64(duration) * 100
			if progress > 100 {
				progress = 100
			}
			
			e.logger.WithFields(logrus.Fields{
				"phase":     phase,
				"elapsed":   elapsed.Round(time.Second),
				"remaining": remaining.Round(time.Second),
				"progress":  fmt.Sprintf("%.1f%%", progress),
			}).Info("Test progress")
		}
	}
}