package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/ethpandaops/spamoor/daemon/db"
)

// newBlockedPauseSpammer returns a running spammer whose Pause blocks until
// runningChan is closed, mimicking a scenario that takes a while to wind down.
func newBlockedPauseSpammer(d *Daemon, id int64, runningChan chan struct{}) *Spammer {
	s := &Spammer{
		daemon:         d,
		dbEntity:       &db.Spammer{ID: id, Scenario: "x", Name: "s"},
		logger:         silentSpammerLogger(),
		scenarioCancel: func() {},
		runningChan:    runningChan,
	}
	s.running.Store(true)

	return s
}

// Deleting a running spammer must not hold the spammer map lock across the
// blocking Pause call; other readers would otherwise freeze for up to 10 seconds.
// The entry is removed while Pause is blocked, so the delete finishes through the
// reacquire recheck without needing a database.
func TestDeleteSpammer_DoesNotBlockReadersWhilePausing(t *testing.T) {
	dctx, dcancel := context.WithCancel(context.Background())
	defer dcancel()
	d := &Daemon{ctx: dctx, cancel: dcancel, spammerMap: make(map[int64]*Spammer, 1)}

	runningChan := make(chan struct{})
	d.spammerMap[1] = newBlockedPauseSpammer(d, 1, runningChan)

	result := make(chan error, 1)
	go func() {
		result <- d.DeleteSpammer(1, "")
	}()

	// Wait until the delete has entered Pause (it releases the lock before).
	time.Sleep(50 * time.Millisecond)

	read := make(chan struct{})
	go func() {
		_ = d.GetSpammer(1)
		close(read)
	}()

	select {
	case <-read:
	case <-time.After(2 * time.Second):
		t.Fatal("GetSpammer was blocked while DeleteSpammer waited in Pause")
	}

	// Simulate a concurrent delete of the same id finishing meanwhile.
	d.spammerMapMtx.Lock()
	delete(d.spammerMap, 1)
	d.spammerMapMtx.Unlock()

	close(runningChan)

	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("expected DeleteSpammer to return cleanly for an already removed entry, got: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("DeleteSpammer did not return after Pause was released")
	}
}
