package spamoor

import (
	"sync"
	"testing"
	"time"
)

func newSelectionTestClient(groups ...string) *Client {
	return &Client{enabled: true, clientGroups: groups}
}

// The round-robin cursor is shared across calls that can see differently
// sized candidate sets (different group filters, or clients toggling
// enabled). Selecting from a smaller set than the cursor's current value
// must wrap instead of indexing out of range.
func TestGetClient_RoundRobinWrapsOnSmallerCandidateSet(t *testing.T) {
	pool := &ClientPool{
		goodClients: []*Client{
			newSelectionTestClient("default"),
			newSelectionTestClient("default"),
			newSelectionTestClient("default"),
			newSelectionTestClient("default"),
			newSelectionTestClient("default"),
			newSelectionTestClient("builders"),
		},
	}

	// Advance the shared cursor past the size of the "builders" group.
	pool.GetClient()
	pool.GetClient()

	for i := 0; i < 10; i++ {
		c := pool.GetClient(WithClientGroup("builders"))
		if c == nil {
			t.Fatalf("expected a client, got nil on iteration %d", i)
		}
		if !c.HasGroup("builders") {
			t.Fatalf("expected a builders client, got one in groups %v", c.clientGroups)
		}
	}
}

// A client set that shrinks between calls (a client goes unhealthy or
// disabled) must not leave the cursor pointing past the end of the next
// candidate list.
func TestGetClient_RoundRobinWrapsOnShrinkingPool(t *testing.T) {
	c1 := newSelectionTestClient("default")
	c2 := newSelectionTestClient("default")
	c3 := newSelectionTestClient("default")
	pool := &ClientPool{goodClients: []*Client{c1, c2, c3}}

	pool.GetClient()
	pool.GetClient()

	pool.goodClients = []*Client{c1, c2}

	for i := 0; i < 10; i++ {
		if c := pool.GetClient(); c == nil {
			t.Fatalf("expected a client, got nil on iteration %d", i)
		}
	}
}

// ByIndex selection was already bounds-safe; keep it that way.
func TestGetClient_ByIndexModeWrapsWithModulo(t *testing.T) {
	pool := &ClientPool{goodClients: []*Client{
		newSelectionTestClient("builders"),
	}}

	c := pool.GetClient(WithClientSelectionMode(SelectClientByIndex, 999), WithClientGroup("builders"))
	if c == nil {
		t.Fatal("expected a client, got nil")
	}
}

// The health-check loop replaces goodClients from a different goroutine than
// GetClient's readers. The write must be covered by the same lock the reader
// takes, or this fails under -race.
func TestClientPool_GoodClientsUpdateIsRaceFreeWithGetClient(t *testing.T) {
	mk := func() *Client { return newSelectionTestClient("default") }
	pool := &ClientPool{goodClients: []*Client{mk(), mk(), mk()}}

	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			pool.selectionMutex.Lock()
			pool.goodClients = []*Client{mk(), mk()}
			pool.selectionMutex.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			pool.GetClient()
		}
	}()

	time.Sleep(150 * time.Millisecond)
	close(stop)
	wg.Wait()
}
