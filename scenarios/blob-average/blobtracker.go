package blobaverage

import (
	"sync"
	"time"
)

// BlockBlobInfo stores blob count information for a single block
type BlockBlobInfo struct {
	BlockNumber uint64
	Timestamp   uint64
	BlobCount   uint64
}

// BlobTracker tracks blob counts across blocks within a configurable time window
type BlobTracker struct {
	mutex           sync.RWMutex
	blocks          []*BlockBlobInfo
	trackingWindow  time.Duration
	totalBlobCount  uint64
	totalBlockCount uint64
}

// NewBlobTracker creates a new BlobTracker with the specified tracking window
func NewBlobTracker(trackingWindow time.Duration) *BlobTracker {
	return &BlobTracker{
		blocks:         make([]*BlockBlobInfo, 0),
		trackingWindow: trackingWindow,
	}
}

// AddBlock adds a new block's blob count to the tracker
func (bt *BlobTracker) AddBlock(blockNumber uint64, timestamp uint64, blobCount uint64) {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	// Add the new block
	bt.blocks = append(bt.blocks, &BlockBlobInfo{
		BlockNumber: blockNumber,
		Timestamp:   timestamp,
		BlobCount:   blobCount,
	})

	// Prune old blocks outside the tracking window
	bt.pruneOldBlocks(timestamp)
}

// pruneOldBlocks removes blocks that are outside the tracking window
// Must be called with mutex held
func (bt *BlobTracker) pruneOldBlocks(currentTimestamp uint64) {
	if len(bt.blocks) == 0 {
		return
	}

	cutoffTime := currentTimestamp - uint64(bt.trackingWindow.Seconds())

	// Find the first block that's within the window
	firstValidIdx := 0
	for i, block := range bt.blocks {
		if block.Timestamp >= cutoffTime {
			firstValidIdx = i
			break
		}
		firstValidIdx = i + 1
	}

	// Remove old blocks
	if firstValidIdx > 0 {
		bt.blocks = bt.blocks[firstValidIdx:]
	}

	// Recalculate totals
	bt.totalBlobCount = 0
	bt.totalBlockCount = uint64(len(bt.blocks))
	for _, block := range bt.blocks {
		bt.totalBlobCount += block.BlobCount
	}
}

// GetAverageBlobCount returns the average blob count per block within the tracking window
func (bt *BlobTracker) GetAverageBlobCount() float64 {
	bt.mutex.RLock()
	defer bt.mutex.RUnlock()

	if bt.totalBlockCount == 0 {
		return 0
	}

	return float64(bt.totalBlobCount) / float64(bt.totalBlockCount)
}

// GetBlobDeficit calculates how many blobs are needed to reach the target average
// Returns a positive number if below target, negative if above target
func (bt *BlobTracker) GetBlobDeficit(targetAverage float64) float64 {
	bt.mutex.RLock()
	defer bt.mutex.RUnlock()

	if bt.totalBlockCount == 0 {
		// No blocks tracked yet, assume we need some blobs
		return targetAverage
	}

	// Calculate what the total blob count should be to achieve the target average
	targetTotalBlobs := targetAverage * float64(bt.totalBlockCount)

	// Return the deficit (positive means we need more blobs)
	return targetTotalBlobs - float64(bt.totalBlobCount)
}

// GetStats returns the current tracking statistics
func (bt *BlobTracker) GetStats() (totalBlobs uint64, totalBlocks uint64, avgBlobs float64) {
	bt.mutex.RLock()
	defer bt.mutex.RUnlock()

	totalBlobs = bt.totalBlobCount
	totalBlocks = bt.totalBlockCount
	if totalBlocks > 0 {
		avgBlobs = float64(totalBlobs) / float64(totalBlocks)
	}
	return
}

// GetBlockCount returns the number of blocks currently being tracked
func (bt *BlobTracker) GetBlockCount() uint64 {
	bt.mutex.RLock()
	defer bt.mutex.RUnlock()
	return bt.totalBlockCount
}
