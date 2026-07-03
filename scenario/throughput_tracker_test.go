package scenario

import (
	"testing"
)

func TestThroughputTrackerAverage(t *testing.T) {
	tracker := newThroughputTracker()

	// no data recorded yet
	if got := tracker.getAverageThroughput(5, 10); got != 0 {
		t.Fatalf("expected 0 for empty tracker, got %f", got)
	}

	for block := uint64(1); block <= 5; block++ {
		tracker.recordCompletion(block, 10)
	}

	if got := tracker.getAverageThroughput(5, 5); got != 10 {
		t.Fatalf("expected average of 10 over 5 blocks, got %f", got)
	}

	if got := tracker.getAverageThroughput(3, 5); got != 10 {
		t.Fatalf("expected average of 10 over 3 blocks, got %f", got)
	}

	// requesting more blocks than tracked falls back to the tracked range
	if got := tracker.getAverageThroughput(50, 5); got != 10 {
		t.Fatalf("expected average of 10 when requesting more blocks than tracked, got %f", got)
	}

	// zero parameters yield no average
	if got := tracker.getAverageThroughput(0, 5); got != 0 {
		t.Fatalf("expected 0 for zero block count, got %f", got)
	}
	if got := tracker.getAverageThroughput(5, 0); got != 0 {
		t.Fatalf("expected 0 for zero current block, got %f", got)
	}

	// a window far past the recorded blocks contains no data
	if got := tracker.getAverageThroughput(3, 1000); got != 0 {
		t.Fatalf("expected 0 for a window without recorded blocks, got %f", got)
	}
}

func TestThroughputTrackerUpdatesExistingBlock(t *testing.T) {
	tracker := newThroughputTracker()

	tracker.recordCompletion(1, 4)
	tracker.recordCompletion(2, 4)
	// recording the same block again adds to the existing entry
	tracker.recordCompletion(1, 6)

	if got := tracker.getAverageThroughput(2, 2); got != 7 {
		t.Fatalf("expected average of 7 ((4+6)+4)/2, got %f", got)
	}
}

func TestThroughputTrackerRingWraparound(t *testing.T) {
	tracker := newThroughputTracker()

	// record more blocks than the ring buffer holds
	for block := uint64(1); block <= 150; block++ {
		tracker.recordCompletion(block, block)
	}

	// the most recent blocks are still tracked: blocks 141..150 sum to 1455
	if got := tracker.getAverageThroughput(10, 150); got != 145.5 {
		t.Fatalf("expected average of 145.5 over the last 10 blocks, got %f", got)
	}

	// blocks that have been evicted from the ring no longer contribute
	if got := tracker.getAverageThroughput(10, 50); got != 0 {
		t.Fatalf("expected 0 for evicted blocks, got %f", got)
	}
}
