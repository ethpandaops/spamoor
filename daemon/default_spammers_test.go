package daemon

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/daemon/configs"
	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios"
)

// The embedded default spammer definitions must always be loadable and reference only
// known scenarios and groups defined within the same set.
func TestLoadDefaultSpammerConfigs(t *testing.T) {
	defaultConfigs, err := LoadDefaultSpammerConfigs()
	if err != nil {
		t.Fatalf("failed to load default spammer configs: %v", err)
	}
	if len(defaultConfigs) == 0 {
		t.Fatal("no default spammer configs found")
	}

	groups := make(map[string]bool, len(defaultConfigs))
	names := make(map[string]bool, len(defaultConfigs))
	keys := make(map[string]bool, len(defaultConfigs))
	for _, cfg := range defaultConfigs {
		if cfg.Name == "" {
			t.Errorf("default spammer with scenario %q has no name", cfg.Scenario)
		}
		if names[cfg.Name] {
			t.Errorf("duplicate default spammer name %q", cfg.Name)
		}
		names[cfg.Name] = true

		if cfg.Key == "" {
			t.Errorf("default spammer %q has no technical key", cfg.Name)
		}
		if keys[cfg.Key] {
			t.Errorf("duplicate default spammer key %q", cfg.Key)
		}
		keys[cfg.Key] = true

		if cfg.Scenario == scenario.GroupScenarioName {
			groups[cfg.Name] = true
		}
	}

	for _, cfg := range defaultConfigs {
		if cfg.Scenario == scenario.GroupScenarioName {
			continue
		}

		descriptor := scenarios.GetScenario(cfg.Scenario)
		if descriptor == nil {
			t.Errorf("default spammer %q references unknown scenario %q", cfg.Name, cfg.Scenario)
			continue
		}

		cfgNode := cfg.Config
		if _, err := configs.MergeScenarioConfiguration(descriptor, &cfgNode); err != nil {
			t.Errorf("default spammer %q config does not merge with scenario defaults: %v", cfg.Name, err)
		}

		if cfg.Group != "" && !groups[cfg.Group] {
			t.Errorf("default spammer %q references unknown group %q", cfg.Name, cfg.Group)
		}
	}
}

// Importing the defaults into a fresh database must create every definition exactly
// once (paused), link group members, and skip everything on a repeated import.
func TestImportDefaultSpammers(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	database := db.NewDatabase(&db.SqliteDatabaseConfig{
		File: filepath.Join(t.TempDir(), "test.db"),
	}, logger)
	if err := database.Init(); err != nil {
		t.Fatalf("failed to init database: %v", err)
	}
	defer database.Close()

	if err := database.ApplyEmbeddedDbSchema(-2); err != nil {
		t.Fatalf("failed to apply db schema: %v", err)
	}

	d := NewDaemon(context.Background(), logger, nil, nil, database)
	d.SetAuditLogger(NewAuditLogger(d, "", "user"))

	if err := d.ImportDefaultSpammers(nil, logger); err != nil {
		t.Fatalf("failed to import default spammers: %v", err)
	}

	// The defaults insertion is a system bootstrap and must not appear in the audit log.
	var auditCount int
	if err := database.ReaderDb.Get(&auditCount, "SELECT COUNT(*) FROM audit_logs"); err != nil {
		t.Fatalf("failed to count audit logs: %v", err)
	}
	if auditCount != 0 {
		t.Errorf("default spammer import wrote %d audit log entries, expected none", auditCount)
	}

	defaultConfigs, err := LoadDefaultSpammerConfigs()
	if err != nil {
		t.Fatalf("failed to load default spammer configs: %v", err)
	}

	spammers := d.GetAllSpammers()
	if len(spammers) != len(defaultConfigs) {
		t.Fatalf("expected %d spammers after import, got %d", len(defaultConfigs), len(spammers))
	}

	byName := make(map[string]*Spammer, len(spammers))
	for _, s := range spammers {
		if s.GetStatus() != int(SpammerStatusPaused) {
			t.Errorf("spammer %q was not created paused (status %d)", s.GetName(), s.GetStatus())
		}
		byName[s.GetName()] = s
	}

	// Defaults occupy the reserved id range (< 100) in definition order: sequential ids
	// starting at 1, groups snapped to the next free multiple of 10.
	prevID := int64(0)
	for _, cfg := range defaultConfigs {
		s := byName[cfg.Name]
		if s == nil {
			continue
		}
		if s.GetID() >= 100 {
			t.Errorf("default spammer %q got id %d outside the reserved range", cfg.Name, s.GetID())
		}
		if s.GetID() <= prevID {
			t.Errorf("default spammer %q id %d is not increasing in definition order (previous %d)", cfg.Name, s.GetID(), prevID)
		}
		if cfg.Scenario == scenario.GroupScenarioName && s.GetID()%10 != 0 {
			t.Errorf("default group %q id %d is not a multiple of 10", cfg.Name, s.GetID())
		}
		prevID = s.GetID()
	}

	expectedIDs := map[string]int64{
		"eoatx-heavy":        1,
		"evm-fuzz-heavy":     7,
		"regular-chain-load": 10,
		"regular-eoatx":      11,
		"regular-blobs":      20,
		"fuzzing":            30,
		"fuzzing-tx-invalid": 33,
	}
	for _, cfg := range defaultConfigs {
		expected, ok := expectedIDs[cfg.Key]
		if !ok {
			continue
		}
		if s := byName[cfg.Name]; s != nil && s.GetID() != expected {
			t.Errorf("default %q: expected id %d, got %d", cfg.Key, expected, s.GetID())
		}
	}

	for _, cfg := range defaultConfigs {
		s := byName[cfg.Name]
		if s == nil {
			t.Errorf("default spammer %q was not imported", cfg.Name)
			continue
		}

		if cfg.Scenario == scenario.GroupScenarioName {
			if !s.IsGroup() {
				t.Errorf("default group %q was not imported as a group", cfg.Name)
			} else if len(d.GetGroupMembers(s.GetID())) == 0 {
				t.Errorf("default group %q has no members", cfg.Name)
			}
			continue
		}

		if cfg.Group != "" {
			parent := byName[cfg.Group]
			if parent == nil || s.GetGroupID() != parent.GetID() {
				t.Errorf("default spammer %q is not linked to group %q", cfg.Name, cfg.Group)
			}
		} else if s.GetGroupID() != 0 {
			t.Errorf("default spammer %q is unexpectedly linked to group %d", cfg.Name, s.GetGroupID())
		}
	}

	// A second import (e.g. first_launch state lost) must not duplicate anything.
	if err := d.ImportDefaultSpammers(nil, logger); err != nil {
		t.Fatalf("failed to re-import default spammers: %v", err)
	}
	if got := len(d.GetAllSpammers()); got != len(defaultConfigs) {
		t.Fatalf("re-import duplicated spammers: expected %d, got %d", len(defaultConfigs), got)
	}
}

// --startup-defaults values must resolve technical keys to display names and pass
// unknown values through for direct name matching.
func TestResolveDefaultStartNames(t *testing.T) {
	defaultConfigs, err := LoadDefaultSpammerConfigs()
	if err != nil {
		t.Fatalf("failed to load default spammer configs: %v", err)
	}

	nameByKey := make(map[string]string, len(defaultConfigs))
	for _, cfg := range defaultConfigs {
		nameByKey[cfg.Key] = cfg.Name
	}
	if _, ok := nameByKey["regular-chain-load"]; !ok {
		t.Fatal("expected a default with key 'regular-chain-load'")
	}
	if _, ok := nameByKey["fuzzing"]; !ok {
		t.Fatal("expected a default with key 'fuzzing'")
	}

	resolved := resolveDefaultStartNames(defaultConfigs, []string{
		"regular-chain-load",
		"fuzzing",
		"Some Custom Name",
	})

	expected := []string{nameByKey["regular-chain-load"], nameByKey["fuzzing"], "Some Custom Name"}
	if len(resolved) != len(expected) {
		t.Fatalf("expected %d resolved names, got %d", len(expected), len(resolved))
	}
	for i := range expected {
		if resolved[i] != expected[i] {
			t.Errorf("resolved[%d]: expected %q, got %q", i, expected[i], resolved[i])
		}
	}
}
