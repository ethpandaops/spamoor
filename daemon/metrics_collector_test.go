package daemon

import (
	"sync"
	"testing"
	"time"
)

// newAggregatingMetrics returns a window with a very large granularity so every
// block aggregates into a single data point. That keeps the collector mutating
// one shared SpammerGasData map, which is the case the copy is meant to protect.
func newAggregatingMetrics() *MultiGranularityMetrics {
	m := &MultiGranularityMetrics{
		blockGranularity: 1_000_000,
		windowDuration:   time.Hour,
		maxDataPoints:    1 << 20,
		currentSpammers:  make(map[uint64]*SpammerSnapshot),
	}
	m.mutex.Lock()
	m.addDataPoint(&BlockDataPoint{Timestamp: time.Now(), StartBlockNumber: 1, EndBlockNumber: 1, BlockCount: 1, TotalGasUsed: 1, SpammerGasData: map[uint64]*SpammerBlockData{0: {}}})
	m.mutex.Unlock()
	return m
}

// GetDataPoints must return data that is independent from the live collector.
// A later mutation of the collector must not be visible through an earlier result.
func TestGetDataPointsReturnsIndependentCopy(t *testing.T) {
	m := newAggregatingMetrics()

	points := m.GetDataPoints()
	if len(points) != 1 {
		t.Fatalf("expected 1 data point, got %d", len(points))
	}
	if _, ok := points[0].SpammerGasData[7]; ok {
		t.Fatal("spammer 7 should not be present yet")
	}

	// Add a new spammer to the live data point after the copy was taken.
	m.mutex.Lock()
	m.addDataPoint(&BlockDataPoint{Timestamp: time.Now(), StartBlockNumber: 2, EndBlockNumber: 2, BlockCount: 1, TotalGasUsed: 10, SpammerGasData: map[uint64]*SpammerBlockData{7: {GasUsed: 10}}})
	m.mutex.Unlock()

	if _, ok := points[0].SpammerGasData[7]; ok {
		t.Fatal("earlier result changed after the collector was mutated; the map is still shared")
	}
}

// TestGetDataPointsConcurrentReaders exercises the reader pattern used by the
// dashboard and SSE handlers (GetDataPoints followed by ranging SpammerGasData)
// while the collector aggregates new blocks. Run with -race to confirm the
// returned maps are not shared with the collector.
func TestGetDataPointsConcurrentReaders(t *testing.T) {
	m := newAggregatingMetrics()

	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := uint64(2)
		for {
			select {
			case <-stop:
				return
			default:
			}
			m.mutex.Lock()
			m.addDataPoint(&BlockDataPoint{Timestamp: time.Now(), StartBlockNumber: i, EndBlockNumber: i, BlockCount: 1, TotalGasUsed: 1, SpammerGasData: map[uint64]*SpammerBlockData{i: {GasUsed: i}}})
			m.mutex.Unlock()
			i++
		}
	}()

	for r := 0; r < 4; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				for _, p := range m.GetDataPoints() {
					for id, data := range p.SpammerGasData {
						_, _ = id, data
					}
				}
			}
		}()
	}

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()
}

// A data point covers blockGranularity blocks, then the next one starts. Counters add
// up across its blocks, snapshots keep the latest value.
func TestAddDataPointAggregatesByGranularity(t *testing.T) {
	m := &MultiGranularityMetrics{blockGranularity: 4, windowDuration: time.Hour, maxDataPoints: 100}

	for block := uint64(10); block < 20; block++ {
		m.addDataPoint(&BlockDataPoint{
			Timestamp: time.Now(), StartBlockNumber: block, EndBlockNumber: block, BlockCount: 1,
			TotalGasUsed: 100, GasLimit: 1000, TxCount: 5, Reorgs: 1,
			SpammerGasData: map[uint64]*SpammerBlockData{
				7: {GasUsed: 60, ConfirmedTxCount: 2, PendingTxCount: block, SubmittedTxDelta: 3, FailedTxCount: 1, DelayMedianMs: 1000, DelayP95Ms: block * 100},
			},
		})
	}

	if len(m.dataPoints) != 3 {
		t.Fatalf("expected 3 data points, got %d", len(m.dataPoints))
	}

	first, last := m.dataPoints[0], m.dataPoints[2]
	if first.StartBlockNumber != 10 || first.EndBlockNumber != 13 || first.BlockCount != 4 || last.BlockCount != 2 {
		t.Fatalf("unexpected block ranges: %d-%d (%d), last %d blocks", first.StartBlockNumber, first.EndBlockNumber, first.BlockCount, last.BlockCount)
	}

	if first.TotalGasUsed != 400 || first.GasLimit != 4000 || first.TxCount != 20 || first.Reorgs != 4 {
		t.Fatalf("block values not summed: %+v", first)
	}

	data := first.SpammerGasData[7]
	if data.GasUsed != 240 || data.ConfirmedTxCount != 8 || data.SubmittedTxDelta != 12 || data.FailedTxCount != 4 {
		t.Fatalf("spammer counters not summed: %+v", data)
	}

	if data.PendingTxCount != 13 || data.DelayMedianMs != 1000 || data.DelayP95Ms != 1300 {
		t.Fatalf("unexpected snapshot values: %+v", data)
	}
}

// A block that is processed again after a reorg replaces its earlier data point in the
// per-block window instead of being added to it.
func TestAddDataPointReplacesReprocessedBlock(t *testing.T) {
	m := &MultiGranularityMetrics{blockGranularity: 1, windowDuration: time.Hour, maxDataPoints: 100}

	add := func(block, gas uint64) {
		m.addDataPoint(&BlockDataPoint{Timestamp: time.Now(), StartBlockNumber: block, EndBlockNumber: block, BlockCount: 1, TotalGasUsed: gas})
	}

	add(10, 100)
	add(11, 200)

	// block 12 saw recovery work and submissions before it was reorged out
	m.addDataPoint(&BlockDataPoint{
		Timestamp: time.Now(), StartBlockNumber: 12, EndBlockNumber: 12, BlockCount: 1, TotalGasUsed: 300, Rebroadcasts: 4, Reorgs: 1,
		SpammerGasData: map[uint64]*SpammerBlockData{7: {ConfirmedTxCount: 9, SubmittedTxDelta: 6, FailedTxCount: 2}},
	})

	add(11, 50)

	if point := m.dataPoints[1]; point.Rebroadcasts != 4 || point.Reorgs != 1 {
		t.Fatalf("recovery counters of the dropped point were lost: %+v", point)
	}

	if data := m.dataPoints[1].SpammerGasData[7]; data == nil || data.SubmittedTxDelta != 6 || data.FailedTxCount != 2 || data.ConfirmedTxCount != 0 {
		t.Fatalf("unexpected spammer values carried over: %+v", data)
	}

	if len(m.dataPoints) != 2 || m.dataPoints[1].StartBlockNumber != 11 || m.dataPoints[1].TotalGasUsed != 50 {
		t.Fatalf("unexpected data points after reprocessing: %+v", m.dataPoints)
	}

	// The transactions confirmed again count for the point that replaces the dropped
	// ones, but not a second time in an aggregated window.
	reprocessed := &BlockDataPoint{
		Timestamp: time.Now(), StartBlockNumber: 11, EndBlockNumber: 11, BlockCount: 1,
		SpammerGasData: map[uint64]*SpammerBlockData{7: {ConfirmedTxCount: 2, reincludedTxCount: 5}},
	}

	m.addDataPoint(reprocessed)

	if got := m.dataPoints[1].SpammerGasData[7].ConfirmedTxCount; got != 7 {
		t.Fatalf("per-block window counts %d confirmed transactions, want 7", got)
	}

	aggregated := &MultiGranularityMetrics{blockGranularity: 32, windowDuration: time.Hour, maxDataPoints: 100}
	aggregated.addDataPoint(reprocessed)

	if got := aggregated.dataPoints[0].SpammerGasData[7].ConfirmedTxCount; got != 2 {
		t.Fatalf("aggregated window counts %d confirmed transactions, want 2", got)
	}

	if reprocessed.SpammerGasData[7].ConfirmedTxCount != 2 {
		t.Fatal("the block's own data point was changed")
	}
}

func TestSpammerBlockDataMerge(t *testing.T) {
	data := &SpammerBlockData{GasUsed: 10, ConfirmedTxCount: 1, PendingTxCount: 2, SubmittedTxDelta: 3, FailedTxCount: 1, DelayMedianMs: 1000, DelayP95Ms: 2000}
	data.merge(&SpammerBlockData{GasUsed: 30, ConfirmedTxCount: 3, PendingTxCount: 5, SubmittedTxDelta: 4, FailedTxCount: 2, DelayMedianMs: 2000, DelayP95Ms: 1500})

	if data.GasUsed != 40 || data.ConfirmedTxCount != 4 || data.PendingTxCount != 7 || data.SubmittedTxDelta != 7 || data.FailedTxCount != 3 {
		t.Fatalf("counters not added: %+v", data)
	}

	if data.DelayMedianMs != 1750 || data.DelayP95Ms != 2000 {
		t.Fatalf("unexpected delays: %+v", data)
	}
}
