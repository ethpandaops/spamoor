package daemon

import (
	"context"
	"io"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/daemon/db"
)

func silentSpammerLogger() *logrus.Entry {
	lg := logrus.New()
	lg.SetOutput(io.Discard)
	return lg.WithField("test", "spammer")
}

// newFastFailSpammer builds a Spammer whose runScenario returns almost
// immediately with no database and no scenario lookup: the daemon context is
// pre-cancelled and the scenario name does not exist.
func newFastFailSpammer(d *Daemon, id int64) *Spammer {
	return &Spammer{
		daemon:   d,
		logger:   silentSpammerLogger(),
		dbEntity: &db.Spammer{ID: id, Scenario: "no-such-scenario"},
	}
}

// TestRunScenario_ConcurrentStartDoesNotPanic drives runScenario from several
// goroutines released at the same time, on the same spammer, over many
// rounds. Only one call per round may pass the running guard; every other
// call must return immediately without touching runningChan, so no round
// should ever produce a panic or a data race.
func TestRunScenario_ConcurrentStartDoesNotPanic(t *testing.T) {
	dctx, dcancel := context.WithCancel(context.Background())
	dcancel()
	d := &Daemon{ctx: dctx}

	const rounds = 500
	const goroutines = 8

	for r := 0; r < rounds; r++ {
		s := newFastFailSpammer(d, int64(r+1))

		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(goroutines)
		for g := 0; g < goroutines; g++ {
			go func() {
				defer wg.Done()
				defer func() {
					if rec := recover(); rec != nil {
						t.Errorf("runScenario panicked: %v", rec)
					}
				}()
				<-start
				s.runScenario()
			}()
		}
		close(start)
		wg.Wait()
	}
}

// TestRunScenario_RunningFlagClearedBeforeChannelCloses verifies that by the
// time runScenario returns, the running flag has already been cleared.
// Closing runningChan is what unblocks a caller waiting in Pause, so if
// running were still true at that point, a Start issued right after Pause
// returns could see stale state and silently do nothing.
func TestRunScenario_RunningFlagClearedBeforeChannelCloses(t *testing.T) {
	dctx, dcancel := context.WithCancel(context.Background())
	dcancel()
	d := &Daemon{ctx: dctx}
	s := newFastFailSpammer(d, 1)

	s.runScenario()

	select {
	case <-s.runningChan:
	default:
		t.Fatal("expected runScenario to close runningChan")
	}
	if s.running.Load() {
		t.Fatal("expected running to already be false once runScenario has returned")
	}
}
