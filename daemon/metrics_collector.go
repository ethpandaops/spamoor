package daemon

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/sirupsen/logrus"
)

// TxPoolMetricsCollector subscribes to TxPool block updates for advanced metrics collection
type TxPoolMetricsCollector struct {
	subscription uint64
	txPool       *spamoor.TxPool

	// Multiple time windows with different granularities
	shortWindow *MultiGranularityMetrics // 30min with per-block precision
	longWindow  *MultiGranularityMetrics // 6h with per-32-block precision

	// Real-time update subscribers
	subscribers      map[uint64]chan *MetricsUpdate
	subscribersMux   sync.RWMutex
	nextSubscriberID uint64

	logger *logrus.Entry
}

// MultiGranularityMetrics manages time-windowed metrics with configurable block granularity
type MultiGranularityMetrics struct {
	mutex            sync.RWMutex
	windowDuration   time.Duration
	blockGranularity uint64 // blocks to aggregate per data point (1 = per block, 32 = per 32 blocks)
	dataPoints       []BlockDataPoint
	currentSpammers  map[uint64]*SpammerSnapshot // current state for each spammer
	maxDataPoints    int                         // prevent memory growth
}

// BlockDataPoint represents metrics for one or more blocks (based on granularity)
type BlockDataPoint struct {
	Timestamp        time.Time                    `json:"timestamp"`
	StartBlockNumber uint64                       `json:"startBlockNumber"`
	EndBlockNumber   uint64                       `json:"endBlockNumber"`
	BlockCount       uint64                       `json:"blockCount"`
	TotalGasUsed     uint64                       `json:"totalGasUsed"`
	SpammerGasData   map[uint64]*SpammerBlockData `json:"spammerGasData"` // spammerID -> data
	OthersGasUsed    uint64                       `json:"othersGasUsed"`
}

// SpammerBlockData represents a spammer's metrics within a data point
type SpammerBlockData struct {
	GasUsed          uint64 `json:"gasUsed"`
	ConfirmedTxCount uint64 `json:"confirmedTxCount"`
	PendingTxCount   uint64 `json:"pendingTxCount"`   // snapshot at end of period
	SubmittedTxCount uint64 `json:"submittedTxCount"` // total transactions submitted (lifetime)
}

// SpammerSnapshot represents current state of a spammer (not aggregated)
type SpammerSnapshot struct {
	SpammerID        uint64
	PendingTxCount   uint64
	TotalConfirmedTx uint64 // lifetime total
	TotalSubmittedTx uint64 // lifetime total submitted
	LastUpdate       time.Time
}

// MetricsUpdate represents a real-time update to metrics
type MetricsUpdate struct {
	BlockNumber     uint64
	NewDataPoint    *BlockDataPoint
	UpdatedSpammers map[uint64]*SpammerSnapshot
	Timestamp       time.Time
}

// NewTxPoolMetricsCollector creates and initializes the metrics collector
func NewTxPoolMetricsCollector(txPool *spamoor.TxPool) *TxPoolMetricsCollector {
	collector := &TxPoolMetricsCollector{
		txPool:      txPool,
		shortWindow: NewMultiGranularityMetrics(30*time.Minute, 1), // 30min, per-block
		longWindow:  NewMultiGranularityMetrics(6*time.Hour, 32),   // 6h, per-32-blocks
		subscribers: make(map[uint64]chan *MetricsUpdate),
		logger:      logrus.WithField("component", "metrics_collector"),
	}

	// Subscribe to bulk block updates
	collector.subscription = txPool.SubscribeToBulkBlockUpdates(collector.handleBulkBlockUpdate)
	collector.logger.Infof("Metrics collector initialized with bulk subscription ID: %d", collector.subscription)

	return collector
}

// NewMultiGranularityMetrics creates a new multi-granularity metrics collector
func NewMultiGranularityMetrics(windowDuration time.Duration, blockGranularity uint64) *MultiGranularityMetrics {
	// Calculate reasonable max data points based on granularity
	maxPoints := int((windowDuration.Minutes() / 12) * 60 / float64(blockGranularity)) // assume 12s blocks
	if maxPoints < 100 {
		maxPoints = 100
	}
	if maxPoints > 10000 {
		maxPoints = 10000
	}

	return &MultiGranularityMetrics{
		windowDuration:   windowDuration,
		blockGranularity: blockGranularity,
		dataPoints:       make([]BlockDataPoint, 0),
		currentSpammers:  make(map[uint64]*SpammerSnapshot),
		maxDataPoints:    maxPoints,
	}
}

// handleBulkBlockUpdate processes each new block via TxPool bulk subscription
func (mc *TxPoolMetricsCollector) handleBulkBlockUpdate(blockNumber uint64, allWalletPoolStats map[*spamoor.WalletPool]*spamoor.WalletPoolBlockStats) {
	// Extract block data
	var block *types.Block
	var receipts []*types.Receipt

	for _, stats := range allWalletPoolStats {
		if stats.Block != nil {
			block = stats.Block
			receipts = stats.Receipts
			break
		}
	}

	// If no block data from spammer stats, we need to get block data another way
	// For now, we'll handle the case where we have spammer data
	if block == nil {
		mc.logger.Tracef("No block data available for block %d", blockNumber)
		return
	}

	// Process the block for both time windows
	newDataPoint := mc.processBlockForWindow(mc.shortWindow, blockNumber, block, allWalletPoolStats, receipts)
	mc.processBlockForWindow(mc.longWindow, blockNumber, block, allWalletPoolStats, receipts)

	// Send real-time update to subscribers (only for short window)
	if newDataPoint != nil {
		update := &MetricsUpdate{
			BlockNumber:     blockNumber,
			NewDataPoint:    newDataPoint,
			UpdatedSpammers: mc.shortWindow.GetSpammerSnapshots(),
			Timestamp:       time.Now(),
		}
		mc.notifySubscribers(update)
	}
}

// processBlockForWindow processes a block for a specific time window
func (mc *TxPoolMetricsCollector) processBlockForWindow(
	window *MultiGranularityMetrics,
	blockNumber uint64,
	block *types.Block,
	allWalletPoolStats map[*spamoor.WalletPool]*spamoor.WalletPoolBlockStats,
	receipts []*types.Receipt,
) *BlockDataPoint {
	window.mutex.Lock()
	defer window.mutex.Unlock()

	now := time.Now()
	totalGasUsed := block.GasUsed()

	// Calculate spammer metrics for this block
	spammerGasData := make(map[uint64]*SpammerBlockData)
	totalSpammerGas := uint64(0)

	for walletPool, stats := range allWalletPoolStats {
		spammerID := walletPool.GetSpammerID()
		gasUsed := mc.calculateGasUsedFromReceipts(stats.Receipts)
		pendingCount, submittedCount := mc.getWalletPoolTxCounts(walletPool)

		spammerGasData[spammerID] = &SpammerBlockData{
			GasUsed:          gasUsed,
			ConfirmedTxCount: stats.ConfirmedTxCount,
			PendingTxCount:   pendingCount,
			SubmittedTxCount: submittedCount,
		}

		totalSpammerGas += gasUsed

		// Update spammer snapshot
		if existing, exists := window.currentSpammers[spammerID]; exists {
			existing.PendingTxCount = pendingCount
			existing.TotalConfirmedTx += stats.ConfirmedTxCount
			existing.TotalSubmittedTx = submittedCount
			existing.LastUpdate = now
		} else {
			window.currentSpammers[spammerID] = &SpammerSnapshot{
				SpammerID:        spammerID,
				PendingTxCount:   pendingCount,
				TotalConfirmedTx: stats.ConfirmedTxCount,
				TotalSubmittedTx: submittedCount,
				LastUpdate:       now,
			}
		}
	}

	othersGasUsed := uint64(0)
	if totalGasUsed > totalSpammerGas {
		othersGasUsed = totalGasUsed - totalSpammerGas
	}

	// Add or aggregate data point based on granularity
	updatedDataPoint := window.addDataPoint(blockNumber, now, totalGasUsed, spammerGasData, othersGasUsed)

	// Clean up old data
	window.pruneOldData(now)

	mc.logger.Tracef("Processed block %d for window (granularity %d): total_gas=%d, spammer_gas=%d, others_gas=%d",
		blockNumber, window.blockGranularity, totalGasUsed, totalSpammerGas, othersGasUsed)

	// Return the data point only for short window (granularity 1)
	if window.blockGranularity == 1 {
		return updatedDataPoint
	}
	return nil
}

// addDataPoint adds a block's data to the window, aggregating based on granularity
// Returns the data point that was either updated or newly created
func (mgm *MultiGranularityMetrics) addDataPoint(
	blockNumber uint64,
	timestamp time.Time,
	totalGasUsed uint64,
	spammerGasData map[uint64]*SpammerBlockData,
	othersGasUsed uint64,
) *BlockDataPoint {
	// Check if we should aggregate with the last data point
	if len(mgm.dataPoints) > 0 {
		lastPoint := &mgm.dataPoints[len(mgm.dataPoints)-1]

		// If this block falls within the current granularity window, aggregate
		if blockNumber < lastPoint.EndBlockNumber+mgm.blockGranularity {
			// Aggregate into existing data point
			lastPoint.EndBlockNumber = blockNumber
			lastPoint.BlockCount++
			lastPoint.TotalGasUsed += totalGasUsed
			lastPoint.OthersGasUsed += othersGasUsed
			lastPoint.Timestamp = timestamp // Update to latest timestamp

			// Aggregate spammer data
			for spammerID, data := range spammerGasData {
				if existing, exists := lastPoint.SpammerGasData[spammerID]; exists {
					existing.GasUsed += data.GasUsed
					existing.ConfirmedTxCount += data.ConfirmedTxCount
					existing.PendingTxCount = data.PendingTxCount     // Keep latest pending count
					existing.SubmittedTxCount = data.SubmittedTxCount // Keep latest submitted count
				} else {
					lastPoint.SpammerGasData[spammerID] = &SpammerBlockData{
						GasUsed:          data.GasUsed,
						ConfirmedTxCount: data.ConfirmedTxCount,
						PendingTxCount:   data.PendingTxCount,
						SubmittedTxCount: data.SubmittedTxCount,
					}
				}
			}
			return lastPoint
		}
	}

	// Create new data point
	spammerDataCopy := make(map[uint64]*SpammerBlockData)
	for spammerID, data := range spammerGasData {
		spammerDataCopy[spammerID] = &SpammerBlockData{
			GasUsed:          data.GasUsed,
			ConfirmedTxCount: data.ConfirmedTxCount,
			PendingTxCount:   data.PendingTxCount,
			SubmittedTxCount: data.SubmittedTxCount,
		}
	}

	newPoint := BlockDataPoint{
		Timestamp:        timestamp,
		StartBlockNumber: blockNumber,
		EndBlockNumber:   blockNumber,
		BlockCount:       1,
		TotalGasUsed:     totalGasUsed,
		SpammerGasData:   spammerDataCopy,
		OthersGasUsed:    othersGasUsed,
	}

	mgm.dataPoints = append(mgm.dataPoints, newPoint)
	return &mgm.dataPoints[len(mgm.dataPoints)-1]
}

// pruneOldData removes data points outside the time window
func (mgm *MultiGranularityMetrics) pruneOldData(now time.Time) {
	cutoff := now.Add(-mgm.windowDuration)

	// Remove old data points
	validPoints := make([]BlockDataPoint, 0, len(mgm.dataPoints))
	for _, point := range mgm.dataPoints {
		if point.Timestamp.After(cutoff) {
			validPoints = append(validPoints, point)
		}
	}

	// Enforce max data points limit
	if len(validPoints) > mgm.maxDataPoints {
		start := len(validPoints) - mgm.maxDataPoints
		validPoints = validPoints[start:]
	}

	mgm.dataPoints = validPoints
}

// calculateGasUsedFromReceipts calculates total gas used from transaction receipts
func (mc *TxPoolMetricsCollector) calculateGasUsedFromReceipts(receipts []*types.Receipt) uint64 {
	totalGas := uint64(0)
	for _, receipt := range receipts {
		if receipt != nil {
			totalGas += receipt.GasUsed
		}
	}
	return totalGas
}

// getWalletPoolTxCounts efficiently gets both pending and submitted transaction counts for a spammer
// by iterating through wallets only once
func (mc *TxPoolMetricsCollector) getWalletPoolTxCounts(walletPool *spamoor.WalletPool) (pendingCount, submittedCount uint64) {
	if walletPool == nil {
		return 0, 0
	}

	allWallets := walletPool.GetAllWallets()
	totalPending := uint64(0)
	totalSubmitted := uint64(0)

	for _, wallet := range allWallets {
		// Get pending count
		pendingNonce := wallet.GetNonce()
		confirmedNonce := wallet.GetConfirmedNonce()
		if pendingNonce > confirmedNonce {
			totalPending += pendingNonce - confirmedNonce
		}

		// Get submitted count
		totalSubmitted += wallet.GetSubmittedTxCount()
	}

	return totalPending, totalSubmitted
}

// GetShortWindowMetrics returns the 30-minute per-block metrics
func (mc *TxPoolMetricsCollector) GetShortWindowMetrics() *MultiGranularityMetrics {
	return mc.shortWindow
}

// GetLongWindowMetrics returns the 6-hour per-32-block metrics
func (mc *TxPoolMetricsCollector) GetLongWindowMetrics() *MultiGranularityMetrics {
	return mc.longWindow
}

// GetDataPoints returns a copy of data points for a specific window
func (mgm *MultiGranularityMetrics) GetDataPoints() []BlockDataPoint {
	mgm.mutex.RLock()
	defer mgm.mutex.RUnlock()

	result := make([]BlockDataPoint, len(mgm.dataPoints))
	copy(result, mgm.dataPoints)
	return result
}

// GetSpammerSnapshots returns current spammer snapshots
func (mgm *MultiGranularityMetrics) GetSpammerSnapshots() map[uint64]*SpammerSnapshot {
	mgm.mutex.RLock()
	defer mgm.mutex.RUnlock()

	result := make(map[uint64]*SpammerSnapshot)
	for id, snapshot := range mgm.currentSpammers {
		result[id] = &SpammerSnapshot{
			SpammerID:        snapshot.SpammerID,
			PendingTxCount:   snapshot.PendingTxCount,
			TotalConfirmedTx: snapshot.TotalConfirmedTx,
			TotalSubmittedTx: snapshot.TotalSubmittedTx,
			LastUpdate:       snapshot.LastUpdate,
		}
	}
	return result
}

// GetTimeRange returns the time range of data points
func (mgm *MultiGranularityMetrics) GetTimeRange() (time.Time, time.Time) {
	mgm.mutex.RLock()
	defer mgm.mutex.RUnlock()

	if len(mgm.dataPoints) == 0 {
		now := time.Now()
		return now.Add(-mgm.windowDuration), now
	}

	return mgm.dataPoints[0].Timestamp, mgm.dataPoints[len(mgm.dataPoints)-1].Timestamp
}

// Shutdown stops the metrics collector
func (mc *TxPoolMetricsCollector) Shutdown() {
	if mc.subscription != 0 {
		mc.txPool.UnsubscribeFromBulkBlockUpdates(mc.subscription)
		mc.logger.Infof("Unsubscribed from bulk block updates (subscription ID: %d)", mc.subscription)
	}

	// Close all subscriber channels
	mc.subscribersMux.Lock()
	for id, ch := range mc.subscribers {
		close(ch)
		delete(mc.subscribers, id)
	}
	mc.subscribersMux.Unlock()
}

// Subscribe creates a new subscription for real-time metrics updates
func (mc *TxPoolMetricsCollector) Subscribe() (uint64, <-chan *MetricsUpdate) {
	mc.subscribersMux.Lock()
	defer mc.subscribersMux.Unlock()

	id := mc.nextSubscriberID
	mc.nextSubscriberID++

	// Create buffered channel to prevent blocking
	ch := make(chan *MetricsUpdate, 10)
	mc.subscribers[id] = ch

	mc.logger.Debugf("New metrics subscriber registered: %d", id)
	return id, ch
}

// Unsubscribe removes a subscription
func (mc *TxPoolMetricsCollector) Unsubscribe(id uint64) {
	mc.subscribersMux.Lock()
	defer mc.subscribersMux.Unlock()

	if ch, exists := mc.subscribers[id]; exists {
		close(ch)
		delete(mc.subscribers, id)
		mc.logger.Debugf("Metrics subscriber unregistered: %d", id)
	}
}

// notifySubscribers sends updates to all active subscribers
func (mc *TxPoolMetricsCollector) notifySubscribers(update *MetricsUpdate) {
	mc.subscribersMux.RLock()
	defer mc.subscribersMux.RUnlock()

	for id, ch := range mc.subscribers {
		select {
		case ch <- update:
			// Successfully sent update
		default:
			// Channel is full, skip this update
			mc.logger.Warnf("Metrics subscriber %d channel is full, skipping update", id)
		}
	}
}
