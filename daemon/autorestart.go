package daemon

import (
	"time"

	"github.com/ethpandaops/spamoor/daemon/configs"
)

// maybeScheduleAutoRestart schedules a delayed restart for a group member that stopped
// in the failed state, if its parent group has auto-restart enabled. Normal stops
// (paused/finished) never reach this path. The restart is re-validated when the
// cooldown expires so manual intervention (restart, pause, membership or group config
// changes, deletion) in the meantime always wins.
func (d *Daemon) maybeScheduleAutoRestart(s *Spammer) {
	cooldown, ok := d.autoRestartParams(s)
	if !ok {
		return
	}

	d.autoRestartMtx.Lock()
	if _, pending := d.autoRestartPending[s.GetID()]; pending {
		d.autoRestartMtx.Unlock()
		return
	}
	d.autoRestartPending[s.GetID()] = struct{}{}
	d.autoRestartMtx.Unlock()

	s.logger.Infof("scheduling auto-restart of failed spammer in %s", cooldown)

	go func() {
		timer := time.NewTimer(cooldown)
		defer timer.Stop()

		fired := false
		select {
		case <-d.ctx.Done():
		case <-timer.C:
			fired = true
		}

		// Clear the pending marker before restarting so a quick re-failure of the
		// restarted member can schedule its own timer without racing this one.
		d.autoRestartMtx.Lock()
		delete(d.autoRestartPending, s.GetID())
		d.autoRestartMtx.Unlock()

		if fired {
			d.runAutoRestart(s.GetID())
		}
	}()
}

// autoRestartParams reports whether the spammer is currently eligible for an
// auto-restart and returns the configured cooldown. Eligible means: a failed,
// enabled member of a group that has auto-restart turned on.
func (d *Daemon) autoRestartParams(s *Spammer) (time.Duration, bool) {
	if s == nil || s.IsGroup() || s.GetGroupID() == 0 {
		return 0, false
	}
	if s.GetStatus() != int(SpammerStatusFailed) {
		return 0, false
	}

	group := d.GetSpammer(s.GetGroupID())
	if group == nil || !group.IsGroup() {
		return 0, false
	}

	groupCfg, err := configs.ParseGroupConfig(group.GetGroupConfig())
	if err != nil || !groupCfg.AutoRestartFailed {
		return 0, false
	}

	memberCfg, err := configs.ParseMemberConfig(s.GetGroupConfig())
	if err != nil || !memberCfg.Enabled {
		return 0, false
	}

	return time.Duration(groupCfg.RestartCooldownSecs()) * time.Second, true
}

// runAutoRestart performs the actual restart after the cooldown, re-validating the
// spammer's state first so stale timers never override manual actions.
func (d *Daemon) runAutoRestart(id int64) {
	if d.ctx.Err() != nil {
		return
	}

	spammer := d.GetSpammer(id)
	if spammer == nil {
		return
	}
	if _, ok := d.autoRestartParams(spammer); !ok {
		return
	}

	spammer.logger.Info("auto-restarting failed spammer")
	if err := spammer.Start(); err != nil {
		spammer.logger.Errorf("failed to auto-restart spammer: %v", err)
	}
}
