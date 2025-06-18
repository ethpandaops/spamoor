package scenario

import (
	"sync"
)

// throughputTracker uses a ring buffer for efficient throughput calculation
type throughputTracker struct {
	mutex        sync.RWMutex
	buffer       []uint64 // ring buffer of transaction counts per block
	blockNumbers []uint64 // corresponding block numbers
	size         int      // buffer size
	head         int      // current write position
	count        int      // number of valid entries
}

func newThroughputTracker() *throughputTracker {
	size := 100 // Track last 100 blocks
	return &throughputTracker{
		buffer:       make([]uint64, size),
		blockNumbers: make([]uint64, size),
		size:         size,
		head:         0,
		count:        0,
	}
}

func (t *throughputTracker) recordCompletion(blockNumber uint64, count uint64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Check if this block is already recorded (update existing entry)
	for i := 0; i < t.count; i++ {
		idx := (t.head - 1 - i + t.size) % t.size
		if t.blockNumbers[idx] == blockNumber {
			t.buffer[idx] += count
			return
		}
	}

	// Add new entry
	t.buffer[t.head] = count
	t.blockNumbers[t.head] = blockNumber
	t.head = (t.head + 1) % t.size

	if t.count < t.size {
		t.count++
	}
}

func (t *throughputTracker) getAverageThroughput(blockCount uint64, currentBlock uint64) float64 {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if blockCount == 0 || currentBlock == 0 || t.count == 0 {
		return 0
	}

	// Limit block count to what we have tracked
	if blockCount > uint64(t.count) {
		blockCount = uint64(t.count)
	}

	var total uint64
	var validBlocks uint64

	// Walk backwards from head to get the most recent blocks
	for i := uint64(0); i < blockCount && i < uint64(t.count); i++ {
		idx := (t.head - 1 - int(i) + t.size) % t.size
		blockNum := t.blockNumbers[idx]

		// Only include blocks in our range
		if blockNum > currentBlock-blockCount && blockNum <= currentBlock {
			total += t.buffer[idx]
			validBlocks++
		}
	}

	if validBlocks == 0 {
		return 0
	}

	// Calculate average over the actual block range, not just blocks with transactions
	return float64(total) / float64(blockCount)
}
