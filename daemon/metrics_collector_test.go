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
	m.addDataPoint(1, time.Now(), 1, map[uint64]*SpammerBlockData{0: {}}, 0)
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
	m.addDataPoint(2, time.Now(), 10, map[uint64]*SpammerBlockData{7: {GasUsed: 10}}, 0)
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
			m.addDataPoint(i, time.Now(), 1, map[uint64]*SpammerBlockData{i: {GasUsed: i}}, 0)
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
