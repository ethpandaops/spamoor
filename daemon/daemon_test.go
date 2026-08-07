package daemon

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/daemon/db"
)

func quietLogger() *logrus.Entry {
	lg := logrus.New()
	lg.SetOutput(io.Discard)
	return lg.WithField("test", "daemon")
}

// TestDeleteSpammer_DoesNotFreezeOtherReadersWhilePausing verifies that
// deleting a running spammer no longer holds spammerMapMtx across the
// blocking Pause call. Before the fix, GetSpammer (and every other reader
// that needs the same lock) would be frozen daemon-wide for however long
// Pause took, up to its 10 second timeout.
func TestDeleteSpammer_DoesNotFreezeOtherReadersWhilePausing(t *testing.T) {
	dctx, dcancel := context.WithCancel(context.Background())
	defer dcancel()
	d := &Daemon{ctx: dctx, cancel: dcancel, logger: quietLogger(), spammerMap: make(map[int64]*Spammer)}

	runningChan := make(chan struct{}) // never closed here -> Pause blocks until released
	s := &Spammer{
		daemon:         d,
		dbEntity:       &db.Spammer{ID: 1, Scenario: "x", Name: "s1"},
		logger:         quietLogger(),
		scenarioCancel: func() {}, // non-nil so DeleteSpammer calls Pause
		runningChan:    runningChan,
	}
	s.running.Store(true)
	d.spammerMap[1] = s

	go func() {
		// DeleteSpammer will eventually hit a nil d.db once Pause releases;
		// recover that expected panic so it can't crash the test binary.
		defer func() { _ = recover() }()
		_ = d.DeleteSpammer(1, "")
	}()

	time.Sleep(50 * time.Millisecond) // let it acquire the lock, unlock, and enter Pause

	done := make(chan struct{})
	go func() {
		_ = d.GetSpammer(1)
		close(done)
	}()

	select {
	case <-done:
		// GetSpammer returned promptly - not frozen behind DeleteSpammer's Pause.
	case <-time.After(2 * time.Second):
		t.Fatal("GetSpammer was frozen while DeleteSpammer was blocked in Pause")
	}

	close(runningChan) // release Pause so the delete goroutine can unwind
	time.Sleep(50 * time.Millisecond)
}

// TestDeleteSpammer_SkipsAlreadyRemovedEntryAfterPauseRace verifies the
// recheck added when the lock is reacquired after Pause: releasing the lock
// around a blocking call means a concurrent DeleteSpammer for the same id
// can now finish while this call is still waiting. Without the recheck,
// this call would carry on to a database delete for an entry that is
// already gone.
func TestDeleteSpammer_SkipsAlreadyRemovedEntryAfterPauseRace(t *testing.T) {
	dctx, dcancel := context.WithCancel(context.Background())
	defer dcancel()
	d := &Daemon{ctx: dctx, cancel: dcancel, logger: quietLogger(), spammerMap: make(map[int64]*Spammer)}

	runningChan := make(chan struct{})
	s := &Spammer{
		daemon:         d,
		dbEntity:       &db.Spammer{ID: 1, Scenario: "x", Name: "s1"},
		logger:         quietLogger(),
		scenarioCancel: func() {},
		runningChan:    runningChan,
	}
	s.running.Store(true)
	d.spammerMap[1] = s

	resultChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultChan <- fmt.Errorf("panic: %v", r)
			}
		}()
		resultChan <- d.DeleteSpammer(1, "")
	}()

	time.Sleep(50 * time.Millisecond) // let it acquire the lock, unlock, and enter Pause

	// Simulate a second, concurrent delete of the same id finishing while
	// the lock is released for Pause.
	d.spammerMapMtx.Lock()
	delete(d.spammerMap, 1)
	d.spammerMapMtx.Unlock()

	close(runningChan) // release Pause

	select {
	case err := <-resultChan:
		if err != nil {
			t.Fatalf("expected DeleteSpammer to find nothing left to do, got: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("DeleteSpammer did not return")
	}
}
