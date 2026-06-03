package configs

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenario"
)

// optsWithThroughput mimics a scenario that accepts throughput/total_count/max_wallets.
type optsWithThroughput struct {
	TotalCount uint64 `yaml:"total_count"`
	Throughput uint64 `yaml:"throughput"`
	MaxWallets uint64 `yaml:"max_wallets"`
	BaseFee    uint64 `yaml:"base_fee"`
}

// optsNoThroughput mimics a scenario that has neither throughput nor max_wallets.
type optsNoThroughput struct {
	TargetAverage float64 `yaml:"target_average"`
	BaseFee       uint64  `yaml:"base_fee"`
}

func descThroughput() *scenario.Descriptor {
	return &scenario.Descriptor{Name: "test-tp", DefaultOptions: optsWithThroughput{}}
}

func descNoThroughput() *scenario.Descriptor {
	return &scenario.Descriptor{Name: "test-notp", DefaultOptions: optsNoThroughput{}}
}

// optsWithMaxPending mimics a scenario that also accepts a max_pending field.
type optsWithMaxPending struct {
	Throughput uint64 `yaml:"throughput"`
	MaxPending uint64 `yaml:"max_pending"`
	MaxWallets uint64 `yaml:"max_wallets"`
}

func descMaxPending() *scenario.Descriptor {
	return &scenario.Descriptor{Name: "test-mp", DefaultOptions: optsWithMaxPending{}}
}

func parseYAML(t *testing.T, s string) map[string]any {
	t.Helper()
	m := map[string]any{}
	if err := yaml.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("failed to parse yaml: %v", err)
	}
	return m
}

// asUint64 coerces a decoded YAML numeric value to uint64 for comparison.
func asUint64(t *testing.T, v any) uint64 {
	t.Helper()
	switch n := v.(type) {
	case int:
		return uint64(n)
	case int64:
		return uint64(n)
	case uint64:
		return n
	case float64:
		return uint64(n)
	default:
		t.Fatalf("value %v (%T) is not numeric", v, v)
		return 0
	}
}

func TestApportion(t *testing.T) {
	tests := []struct {
		name    string
		total   uint64
		weights []uint64
		want    []uint64
	}{
		{"empty", 100, nil, []uint64{}},
		{"single", 100, []uint64{5}, []uint64{100}},
		{"even split", 100, []uint64{1, 1, 1, 1}, []uint64{25, 25, 25, 25}},
		{"weighted 20/50/30", 100, []uint64{20, 50, 30}, []uint64{20, 50, 30}},
		{"largest remainder", 10, []uint64{1, 1, 1}, []uint64{4, 3, 3}},
		{"all zero equal", 9, []uint64{0, 0, 0}, []uint64{3, 3, 3}},
		{"all zero remainder", 10, []uint64{0, 0, 0}, []uint64{4, 3, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Apportion(tt.total, tt.weights)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Apportion(%d, %v) = %v, want %v", tt.total, tt.weights, got, tt.want)
			}

			if len(tt.weights) > 0 {
				var sum uint64
				for _, s := range got {
					sum += s
				}
				if sum != tt.total {
					t.Fatalf("shares %v sum to %d, want %d", got, sum, tt.total)
				}
			}
		})
	}
}

func TestResolveMemberConfig_OverlayWinsAndGated(t *testing.T) {
	member := "base_fee: 10\nthroughput: 5\n"
	overlay := map[string]any{
		"base_fee":    99,  // accepted -> overlay wins
		"nonexistent": 123, // not a scenario field -> skipped
		"target_avg":  1,   // not a scenario field -> skipped
	}

	out, err := ResolveMemberConfig(descThroughput(), member, overlay, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	got := parseYAML(t, out)
	if asUint64(t, got["base_fee"]) != 99 {
		t.Fatalf("overlay value must win, got %v", got["base_fee"])
	}
	if _, ok := got["nonexistent"]; ok {
		t.Fatal("fields the scenario rejects must be skipped")
	}
	if _, ok := got["target_avg"]; ok {
		t.Fatal("unknown overlay fields must be skipped")
	}
}

func TestResolveMemberConfig_SharedThroughputAndMaxWallets(t *testing.T) {
	tp := uint64(100)
	out, err := ResolveMemberConfig(descThroughput(), "throughput: 5\n", nil, &tp, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	got := parseYAML(t, out)
	if asUint64(t, got["throughput"]) != 100 {
		t.Fatalf("shared throughput must be injected, got %v", got["throughput"])
	}
	// max_wallets = clamp(round(100/4)=25, 20, 1000) = 25
	if asUint64(t, got["max_wallets"]) != 25 {
		t.Fatalf("max_wallets = %v, want 25", got["max_wallets"])
	}
}

func TestResolveMemberConfig_MaxWalletsClamp(t *testing.T) {
	tests := []struct {
		throughput uint64
		want       uint64
	}{
		{40, 20},     // round(10) -> clamp up to 20
		{100, 25},    // round(25)
		{8000, 1000}, // round(2000) -> clamp down to 1000
	}
	for _, tt := range tests {
		tp := tt.throughput
		out, err := ResolveMemberConfig(descThroughput(), "", nil, &tp, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		got := parseYAML(t, out)
		if asUint64(t, got["max_wallets"]) != tt.want {
			t.Fatalf("throughput=%d: max_wallets=%v, want %d", tt.throughput, got["max_wallets"], tt.want)
		}
	}
}

func TestResolveMemberConfig_NeverInjectZeroThroughput(t *testing.T) {
	zero := uint64(0)
	out, err := ResolveMemberConfig(descThroughput(), "", nil, &zero, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := parseYAML(t, out)
	if asUint64(t, got["throughput"]) != 1 {
		t.Fatalf("a 0 share must be bumped to 1, got %v", got["throughput"])
	}
}

func TestResolveMemberConfig_NotInjectedForUnsupportedScenario(t *testing.T) {
	tp := uint64(100)
	out, err := ResolveMemberConfig(descNoThroughput(), "base_fee: 7\n", nil, &tp, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := parseYAML(t, out)
	if _, ok := got["throughput"]; ok {
		t.Fatal("throughput must not be injected when scenario lacks it")
	}
	if _, ok := got["max_wallets"]; ok {
		t.Fatal("max_wallets must not be injected when scenario lacks it")
	}
	if asUint64(t, got["base_fee"]) != 7 {
		t.Fatalf("base_fee = %v, want 7", got["base_fee"])
	}
}

func TestResolveMemberConfig_SharedCountOnly(t *testing.T) {
	count := uint64(1000)
	out, err := ResolveMemberConfig(descThroughput(), "", nil, nil, &count, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := parseYAML(t, out)
	if asUint64(t, got["total_count"]) != 1000 {
		t.Fatalf("total_count = %v, want 1000", got["total_count"])
	}
	// no throughput -> max_wallets from count: round(1000/50)=20
	if asUint64(t, got["max_wallets"]) != 20 {
		t.Fatalf("max_wallets = %v, want 20", got["max_wallets"])
	}
}

func TestResolveMemberConfig_EmptyOverlayPreservesMemberConfig(t *testing.T) {
	member := "base_fee: 12\nthroughput: 40\ntotal_count: 0\n"
	out, err := ResolveMemberConfig(descThroughput(), member, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	got := parseYAML(t, out)
	if asUint64(t, got["base_fee"]) != 12 {
		t.Fatalf("base_fee = %v, want 12", got["base_fee"])
	}
	if asUint64(t, got["throughput"]) != 40 {
		t.Fatalf("throughput = %v, want 40", got["throughput"])
	}
	// max_wallets derived from throughput 40 -> clamp(10,20,1000) = 20
	if asUint64(t, got["max_wallets"]) != 20 {
		t.Fatalf("max_wallets = %v, want 20", got["max_wallets"])
	}
}

func TestResolveMemberConfig_MaxPendingDerivedAndOverride(t *testing.T) {
	tp := uint64(50)

	// default: max_pending = 2 x throughput
	out, err := ResolveMemberConfig(descMaxPending(), "", nil, &tp, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := parseYAML(t, out); asUint64(t, got["max_pending"]) != 100 {
		t.Fatalf("derived max_pending = %v, want 100", got["max_pending"])
	}

	// explicit group share overrides the derivation
	mp := uint64(30)
	out, err = ResolveMemberConfig(descMaxPending(), "", nil, &tp, nil, &mp)
	if err != nil {
		t.Fatal(err)
	}
	if got := parseYAML(t, out); asUint64(t, got["max_pending"]) != 30 {
		t.Fatalf("override max_pending = %v, want 30", got["max_pending"])
	}

	// a 0 explicit share is bumped to 1 (never injected as 0 = unlimited)
	zero := uint64(0)
	out, err = ResolveMemberConfig(descMaxPending(), "", nil, &tp, nil, &zero)
	if err != nil {
		t.Fatal(err)
	}
	if got := parseYAML(t, out); asUint64(t, got["max_pending"]) != 1 {
		t.Fatalf("zero override max_pending = %v, want 1", got["max_pending"])
	}

	// scenarios without a max_pending field never get it injected
	out, err = ResolveMemberConfig(descThroughput(), "", nil, &tp, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := parseYAML(t, out)["max_pending"]; ok {
		t.Fatal("max_pending must not be injected when scenario lacks the field")
	}
}

func TestParseGroupConfigDefaults(t *testing.T) {
	cfg, err := ParseGroupConfig("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ThroughputMode != GroupModeIndependent {
		t.Fatalf("default mode = %q, want %q", cfg.ThroughputMode, GroupModeIndependent)
	}

	cfg, err = ParseGroupConfig(`{"throughput_mode":"shared","total_throughput":100}`)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ThroughputMode != GroupModeShared || cfg.TotalThroughput != 100 {
		t.Fatalf("parsed group config = %+v", cfg)
	}
}

func TestParseMemberConfigDefaults(t *testing.T) {
	cfg, err := ParseMemberConfig("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Weight != 1 || !cfg.Enabled {
		t.Fatalf("default member config = %+v", cfg)
	}

	cfg, err = ParseMemberConfig(`{"weight":3,"enabled":false,"sort_order":2}`)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Weight != 3 || cfg.Enabled || cfg.SortOrder != 2 {
		t.Fatalf("parsed member config = %+v", cfg)
	}
}

func TestResolveRunConfigs_WeightZeroMemberStillRuns(t *testing.T) {
	lookup := func(name string) *scenario.Descriptor {
		if name == "test-tp" {
			return descThroughput()
		}
		return nil
	}

	all := []SpammerConfig{
		{
			Scenario:    scenario.GroupScenarioName,
			Name:        "grp",
			GroupConfig: map[string]any{"throughput_mode": "shared", "total_throughput": 100},
		},
		{Scenario: "test-tp", Name: "a", Group: "grp", GroupConfig: map[string]any{"weight": 100}},
		// enabled (default) but weight 0: still runs, at the resolver's min-1 throughput.
		{Scenario: "test-tp", Name: "z", Group: "grp", GroupConfig: map[string]any{"weight": 0}},
	}

	resolved, err := ResolveRunConfigs(all, lookup)
	if err != nil {
		t.Fatal(err)
	}
	if len(resolved) != 2 {
		t.Fatalf("expected 2 runnable members (weight-0 still runs), got %d", len(resolved))
	}
	byName := map[string]uint64{}
	for _, r := range resolved {
		byName[r.Config.Name] = asUint64(t, parseYAML(t, r.ConfigYAML)["throughput"])
	}
	if byName["a"] != 100 {
		t.Fatalf("member a throughput = %d, want 100", byName["a"])
	}
	if byName["z"] != 1 {
		t.Fatalf("weight-0 member z throughput = %d, want 1 (min-1 guard)", byName["z"])
	}
}

func TestResolveRunConfigs_SharedSplit(t *testing.T) {
	lookup := func(name string) *scenario.Descriptor {
		if name == "test-tp" {
			return descThroughput()
		}
		return nil
	}

	all := []SpammerConfig{
		{
			Scenario: scenario.GroupScenarioName,
			Name:     "grp",
			Config:   map[string]any{"base_fee": 50},
			GroupConfig: map[string]any{
				"throughput_mode":  "shared",
				"total_throughput": 100,
			},
		},
		{Scenario: "test-tp", Name: "a", Group: "grp", GroupConfig: map[string]any{"weight": 20}},
		{Scenario: "test-tp", Name: "b", Group: "grp", GroupConfig: map[string]any{"weight": 50}},
		{Scenario: "test-tp", Name: "c", Group: "grp", GroupConfig: map[string]any{"weight": 30}},
	}

	resolved, err := ResolveRunConfigs(all, lookup)
	if err != nil {
		t.Fatal(err)
	}
	if len(resolved) != 3 {
		t.Fatalf("expected 3 runnable members, got %d", len(resolved))
	}

	wantTP := map[string]uint64{"a": 20, "b": 50, "c": 30}
	for _, r := range resolved {
		got := parseYAML(t, r.ConfigYAML)
		if asUint64(t, got["throughput"]) != wantTP[r.Config.Name] {
			t.Fatalf("member %s throughput = %v, want %d", r.Config.Name, got["throughput"], wantTP[r.Config.Name])
		}
		if asUint64(t, got["base_fee"]) != 50 {
			t.Fatalf("member %s overlay base_fee = %v, want 50", r.Config.Name, got["base_fee"])
		}
	}
}
