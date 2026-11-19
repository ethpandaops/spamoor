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

func TestCounterDelta(t *testing.T) {
	last := map[uint64]uint64{}

	for _, step := range []struct{ current, want uint64 }{{5, 5}, {12, 7}, {12, 0}, {3, 3}} {
		if got := counterDelta(last, 1, step.current); got != step.want {
			t.Fatalf("counter at %d: got delta %d, want %d", step.current, got, step.want)
		}
	}
}
