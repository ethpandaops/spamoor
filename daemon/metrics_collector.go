package daemon

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txtypes"
)

// TxPoolMetricsCollector subscribes to TxPool block updates for advanced metrics collection
type TxPoolMetricsCollector struct {
	subscription uint64
	txPool       *spamoor.TxPool

	// Multiple time windows with different granularities
	shortWindow *MultiGranularityMetrics // 30min with per-block precision
	longWindow  *MultiGranularityMetrics // 6h with per-32-block precision

	// Cumulative counters seen at the previous block, to derive per-block values
	lastCounters map[*spamoor.WalletPool]poolCounters
	lastRecovery spamoor.RecoveryStats

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

	GasLimit     uint64 `json:"gasLimit"`     // summed over the blocks of the point
	BaseFee      uint64 `json:"baseFee"`      // wei, of the last block
	TxCount      uint64 `json:"txCount"`      // all transactions in the blocks
	Reorgs       uint64 `json:"reorgs"`       // chain reorganizations handled
	ReorgedTxs   uint64 `json:"reorgedTxs"`   // transactions reorged out
	Rebroadcasts uint64 `json:"rebroadcasts"` // pending transactions sent again
	GapFills     uint64 `json:"gapFills"`     // nonce gaps filled
}

// SpammerBlockData represents a spammer's metrics within a data point
type SpammerBlockData struct {
	GasUsed          uint64 `json:"gasUsed"`
	ConfirmedTxCount uint64 `json:"confirmedTxCount"`
	PendingTxCount   uint64 `json:"pendingTxCount"`   // snapshot at end of period
	SubmittedTxCount uint64 `json:"submittedTxCount"` // total transactions submitted (lifetime)

	SubmittedTxDelta uint64 `json:"submittedTxDelta"` // transactions submitted in the period
	FailedTxCount    uint64 `json:"failedTxCount"`    // transactions that failed in the period
	DelayMedianMs    uint64 `json:"delayMedianMs"`    // submission to inclusion, median
	DelayP95Ms       uint64 `json:"delayP95Ms"`       // submission to inclusion, 95th percentile

	// confirmed again after a reorg, not part of ConfirmedTxCount
	reincludedTxCount uint64
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

		lastCounters: make(map[*spamoor.WalletPool]poolCounters, 16),
		lastRecovery: txPool.GetRecoveryStats(),
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
func (mc *TxPoolMetricsCollector) handleBulkBlockUpdate(blockNumber uint64, globalBlockStats *spamoor.GlobalBlockStats) {
	// Extract block data
	block := globalBlockStats.Block

	// If no block data from spammer stats, we need to get block data another way
	// For now, we'll handle the case where we have spammer data
	if block == nil {
		mc.logger.Tracef("No block data available for block %d", blockNumber)
		return
	}

	// Process the block for both time windows
	point := mc.buildBlockPoint(blockNumber, block, globalBlockStats.WalletPoolStats)
	newDataPoint := mc.shortWindow.addBlock(point)
	mc.longWindow.addBlock(point)

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

// poolCounters are a wallet pool's cumulative counters as seen at a block.
type poolCounters struct {
	submitted uint64
	failed    uint64
}

// buildBlockPoint derives the data point of a single block.
func (mc *TxPoolMetricsCollector) buildBlockPoint(
	blockNumber uint64,
	block *txtypes.Block,
	allWalletPoolStats map[*spamoor.WalletPool]*spamoor.WalletPoolBlockStats,
) *BlockDataPoint {
	point := &BlockDataPoint{
		Timestamp:        time.Now(),
		StartBlockNumber: blockNumber,
		EndBlockNumber:   blockNumber,
		BlockCount:       1,
		TotalGasUsed:     block.GasUsed,
		SpammerGasData:   make(map[uint64]*SpammerBlockData, len(allWalletPoolStats)),
		GasLimit:         block.GasLimit,
		TxCount:          uint64(len(block.Transactions)),
	}

	if block.BaseFee != nil && block.BaseFee.IsUint64() {
		point.BaseFee = block.BaseFee.Uint64()
	}

	recovery := mc.txPool.GetRecoveryStats()
	point.Reorgs = recovery.Reorgs - mc.lastRecovery.Reorgs
	point.ReorgedTxs = recovery.ReorgedTxs - mc.lastRecovery.ReorgedTxs
	point.Rebroadcasts = recovery.Rebroadcasts - mc.lastRecovery.Rebroadcasts
	point.GapFills = recovery.GapFills - mc.lastRecovery.GapFills
	mc.lastRecovery = recovery

	totalSpammerGas := uint64(0)
	counters := make(map[*spamoor.WalletPool]poolCounters, len(allWalletPoolStats))

	for walletPool, stats := range allWalletPoolStats {
		spammerID := walletPool.GetSpammerID()
		gasUsed := mc.calculateGasUsedFromReceipts(stats.Receipts)

		// The counters are per pool: a restarted spammer has a new pool that starts
		// from zero, which may briefly live next to the old one.
		current := poolCounters{submitted: walletPool.GetSubmittedTxCount(), failed: walletPool.GetFailedTxCount()}
		previous := mc.lastCounters[walletPool]
		counters[walletPool] = current

		data := &SpammerBlockData{
			GasUsed:          gasUsed,
			ConfirmedTxCount: stats.ConfirmedTxCount,
			PendingTxCount:   walletPool.GetPendingTxCount(),
			SubmittedTxCount: current.submitted,
			SubmittedTxDelta: current.submitted - previous.submitted,
			FailedTxCount:    current.failed - previous.failed,
			DelayMedianMs:    uint64(stats.MedianInclusionDelay.Milliseconds()),
			DelayP95Ms:       uint64(stats.P95InclusionDelay.Milliseconds()),

			reincludedTxCount: stats.ReincludedTxCount,
		}

		if existing := point.SpammerGasData[spammerID]; existing != nil {
			existing.merge(data)
		} else {
			point.SpammerGasData[spammerID] = data
		}

		totalSpammerGas += gasUsed
	}

	// only the pools that are still there are remembered
	mc.lastCounters = counters

	if point.TotalGasUsed > totalSpammerGas {
		point.OthersGasUsed = point.TotalGasUsed - totalSpammerGas
	}

	return point
}

// addBlock adds a block's data point to the window and returns a copy of the window's
// current point for a per-block window.
func (mgm *MultiGranularityMetrics) addBlock(block *BlockDataPoint) *BlockDataPoint {
	mgm.mutex.Lock()
	defer mgm.mutex.Unlock()

	for spammerID, data := range block.SpammerGasData {
		snapshot := mgm.currentSpammers[spammerID]
		if snapshot == nil {
			snapshot = &SpammerSnapshot{SpammerID: spammerID}
			mgm.currentSpammers[spammerID] = snapshot
		}

		snapshot.PendingTxCount = data.PendingTxCount
		snapshot.TotalConfirmedTx += data.ConfirmedTxCount
		snapshot.TotalSubmittedTx = data.SubmittedTxCount
		snapshot.LastUpdate = block.Timestamp
	}

	updated := mgm.addDataPoint(block)
	mgm.pruneOldData(block.Timestamp)

	if mgm.blockGranularity == 1 {
		dataPointCopy := updated.copy()
		return &dataPointCopy
	}

	return nil
}

// merge adds another set of values for the same period: counters add up, the median
// delay is weighted by confirmed transactions.
func (data *SpammerBlockData) merge(other *SpammerBlockData) {
	if confirmed := data.ConfirmedTxCount + other.ConfirmedTxCount; confirmed > 0 {
		data.DelayMedianMs = (data.DelayMedianMs*data.ConfirmedTxCount + other.DelayMedianMs*other.ConfirmedTxCount) / confirmed
	}

	data.DelayP95Ms = max(data.DelayP95Ms, other.DelayP95Ms)
	data.GasUsed += other.GasUsed
	data.ConfirmedTxCount += other.ConfirmedTxCount
	data.PendingTxCount += other.PendingTxCount
	data.SubmittedTxCount += other.SubmittedTxCount
	data.SubmittedTxDelta += other.SubmittedTxDelta
	data.FailedTxCount += other.FailedTxCount
	data.reincludedTxCount += other.reincludedTxCount
}

// addDataPoint adds a block to the current data point, or starts a new one once the
// current one covers blockGranularity blocks.
func (mgm *MultiGranularityMetrics) addDataPoint(block *BlockDataPoint) *BlockDataPoint {
	// A block processed again after a reorg replaces what the per-block window holds
	// for its height and above, and with it the transactions counted there: the ones
	// confirmed again count for the new points. An aggregated point cannot give
	// anything back, so there they stay counted once, where they were first.
	if mgm.blockGranularity == 1 {
		for len(mgm.dataPoints) > 0 && mgm.dataPoints[len(mgm.dataPoints)-1].StartBlockNumber >= block.StartBlockNumber {
			mgm.dataPoints = mgm.dataPoints[:len(mgm.dataPoints)-1]
		}

		counted := block.copy()
		for _, data := range counted.SpammerGasData {
			data.ConfirmedTxCount += data.reincludedTxCount
		}

		block = &counted
	}

	if len(mgm.dataPoints) > 0 {
		lastPoint := &mgm.dataPoints[len(mgm.dataPoints)-1]

		if block.StartBlockNumber < lastPoint.StartBlockNumber+mgm.blockGranularity {
			lastPoint.EndBlockNumber = block.EndBlockNumber
			lastPoint.BlockCount++
			lastPoint.TotalGasUsed += block.TotalGasUsed
			lastPoint.OthersGasUsed += block.OthersGasUsed
			lastPoint.Timestamp = block.Timestamp
			lastPoint.GasLimit += block.GasLimit
			lastPoint.BaseFee = block.BaseFee
			lastPoint.TxCount += block.TxCount
			lastPoint.Reorgs += block.Reorgs
			lastPoint.ReorgedTxs += block.ReorgedTxs
			lastPoint.Rebroadcasts += block.Rebroadcasts
			lastPoint.GapFills += block.GapFills

			for spammerID, data := range block.SpammerGasData {
				existing := lastPoint.SpammerGasData[spammerID]
				if existing == nil {
					dataCopy := *data
					lastPoint.SpammerGasData[spammerID] = &dataCopy

					continue
				}

				existing.merge(data)
				existing.PendingTxCount = data.PendingTxCount
				existing.SubmittedTxCount = data.SubmittedTxCount
			}

			return lastPoint
		}
	}

	mgm.dataPoints = append(mgm.dataPoints, block.copy())

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
func (mc *TxPoolMetricsCollector) calculateGasUsedFromReceipts(receipts []*txtypes.Receipt) uint64 {
	totalGas := uint64(0)
	for _, receipt := range receipts {
		if receipt != nil {
			totalGas += receipt.GasUsed
		}
	}
	return totalGas
}

// GetShortWindowMetrics returns the 30-minute per-block metrics
func (mc *TxPoolMetricsCollector) GetShortWindowMetrics() *MultiGranularityMetrics {
	return mc.shortWindow
}

// GetLongWindowMetrics returns the 6-hour per-32-block metrics
func (mc *TxPoolMetricsCollector) GetLongWindowMetrics() *MultiGranularityMetrics {
	return mc.longWindow
}

// copy returns a deep copy of the data point. The SpammerGasData map is copied
// so callers can keep reading the result after releasing the window lock without
// racing the collector, which keeps mutating the live map as new blocks are
// aggregated into the current data point.
func (p *BlockDataPoint) copy() BlockDataPoint {
	result := *p
	if p.SpammerGasData != nil {
		result.SpammerGasData = make(map[uint64]*SpammerBlockData, len(p.SpammerGasData))
		for id, data := range p.SpammerGasData {
			if data == nil {
				result.SpammerGasData[id] = nil
				continue
			}
			dataCopy := *data
			result.SpammerGasData[id] = &dataCopy
		}
	}
	return result
}

// GetDataPoints returns a deep copy of data points for a specific window.
// The nested maps are copied so the returned data is safe to read without
// holding the window lock.
func (mgm *MultiGranularityMetrics) GetDataPoints() []BlockDataPoint {
	mgm.mutex.RLock()
	defer mgm.mutex.RUnlock()

	result := make([]BlockDataPoint, len(mgm.dataPoints))
	for i := range mgm.dataPoints {
		result[i] = mgm.dataPoints[i].copy()
	}
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
