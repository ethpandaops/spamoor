package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/daemon/configs"
	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/scenario"
)

func autoRestartTestDaemon(t *testing.T, groupCfg *configs.GroupConfig, memberCfg *configs.MemberConfig, memberStatus SpammerStatus) (*Daemon, *Spammer) {
	t.Helper()

	d := NewDaemon(context.Background(), logrus.New(), nil, nil, nil)

	gcJSON, err := groupCfg.Marshal()
	if err != nil {
		t.Fatalf("failed to marshal group config: %v", err)
	}
	mcJSON, err := memberCfg.Marshal()
	if err != nil {
		t.Fatalf("failed to marshal member config: %v", err)
	}

	group := &Spammer{daemon: d, dbEntity: &db.Spammer{
		ID: 1, Scenario: scenario.GroupScenarioName, GroupConfig: gcJSON,
	}}
	member := &Spammer{daemon: d, dbEntity: &db.Spammer{
		ID: 2, Scenario: "eoatx", GroupID: 1, GroupConfig: mcJSON, Status: int(memberStatus),
	}}
	d.spammerMap[1] = group
	d.spammerMap[2] = member

	return d, member
}

func TestAutoRestartParams(t *testing.T) {
	tests := []struct {
		name         string
		groupCfg     configs.GroupConfig
		memberCfg    configs.MemberConfig
		memberStatus SpammerStatus
		wantEligible bool
		wantCooldown time.Duration
	}{
		{
			name:         "failed member of opted-in group",
			groupCfg:     configs.GroupConfig{AutoRestartFailed: true, AutoRestartCooldown: 60},
			memberCfg:    configs.MemberConfig{Weight: 1, Enabled: true},
			memberStatus: SpammerStatusFailed,
			wantEligible: true,
			wantCooldown: 60 * time.Second,
		},
		{
			name:         "cooldown defaults to 300s when unset",
			groupCfg:     configs.GroupConfig{AutoRestartFailed: true},
			memberCfg:    configs.MemberConfig{Weight: 1, Enabled: true},
			memberStatus: SpammerStatusFailed,
			wantEligible: true,
			wantCooldown: configs.DefaultAutoRestartCooldownSecs * time.Second,
		},
		{
			name:         "group did not opt in",
			groupCfg:     configs.GroupConfig{},
			memberCfg:    configs.MemberConfig{Weight: 1, Enabled: true},
			memberStatus: SpammerStatusFailed,
			wantEligible: false,
		},
		{
			name:         "normally paused member is not restarted",
			groupCfg:     configs.GroupConfig{AutoRestartFailed: true},
			memberCfg:    configs.MemberConfig{Weight: 1, Enabled: true},
			memberStatus: SpammerStatusPaused,
			wantEligible: false,
		},
		{
			name:         "finished member is not restarted",
			groupCfg:     configs.GroupConfig{AutoRestartFailed: true},
			memberCfg:    configs.MemberConfig{Weight: 1, Enabled: true},
			memberStatus: SpammerStatusFinished,
			wantEligible: false,
		},
		{
			name:         "disabled member is not restarted",
			groupCfg:     configs.GroupConfig{AutoRestartFailed: true},
			memberCfg:    configs.MemberConfig{Weight: 1, Enabled: false},
			memberStatus: SpammerStatusFailed,
			wantEligible: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, member := autoRestartTestDaemon(t, &tt.groupCfg, &tt.memberCfg, tt.memberStatus)

			cooldown, ok := d.autoRestartParams(member)
			if ok != tt.wantEligible {
				t.Fatalf("eligible = %v, want %v", ok, tt.wantEligible)
			}
			if ok && cooldown != tt.wantCooldown {
				t.Fatalf("cooldown = %v, want %v", cooldown, tt.wantCooldown)
			}
		})
	}
}

func TestAutoRestartParamsStandaloneAndOrphan(t *testing.T) {
	d, member := autoRestartTestDaemon(t,
		&configs.GroupConfig{AutoRestartFailed: true},
		&configs.MemberConfig{Weight: 1, Enabled: true},
		SpammerStatusFailed)

	// Standalone spammer (no group) is never auto-restarted.
	standalone := &Spammer{daemon: d, dbEntity: &db.Spammer{
		ID: 3, Scenario: "eoatx", Status: int(SpammerStatusFailed),
	}}
	d.spammerMap[3] = standalone
	if _, ok := d.autoRestartParams(standalone); ok {
		t.Fatal("standalone spammer must not be eligible for auto-restart")
	}

	// Orphaned member (parent group deleted) is not eligible either.
	delete(d.spammerMap, 1)
	if _, ok := d.autoRestartParams(member); ok {
		t.Fatal("orphaned member must not be eligible for auto-restart")
	}
}
