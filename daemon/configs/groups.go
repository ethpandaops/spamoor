package configs

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/scenario"
)

// Throughput modes for a spammer group.
const (
	// GroupModeIndependent applies only the shared overlay; each member keeps its
	// own throughput/count.
	GroupModeIndependent = "independent"
	// GroupModeShared splits the group's total throughput/count across enabled
	// members by weight.
	GroupModeShared = "shared"
)

// max_wallets derivation bounds and divisors for group members.
const (
	groupMaxWalletsMin               = 20
	groupMaxWalletsMax               = 1000
	groupMaxWalletsThroughputDivisor = 4
	groupMaxWalletsCountDivisor      = 50
)

// maxPendingThroughputMultiplier scales a member's resolved throughput into its derived
// max_pending when the group does not set an explicit total.
const maxPendingThroughputMultiplier = 2

// DefaultAutoRestartCooldownSecs is the cooldown applied before auto-restarting a
// failed group member when the group does not configure an explicit value.
const DefaultAutoRestartCooldownSecs = 300

// Config field names that receive computed (non-overlay) handling.
const (
	fieldThroughput = "throughput"
	fieldTotalCount = "total_count"
	fieldMaxWallets = "max_wallets"
	fieldMaxPending = "max_pending"
)

// GroupConfig holds the JSON metadata stored in a group row's group_config column.
type GroupConfig struct {
	ThroughputMode  string `json:"throughput_mode"`
	TotalThroughput uint64 `json:"total_throughput"`
	TotalCount      uint64 `json:"total_count"`
	// TotalMaxPending, when > 0, is a group-wide concurrent-pending budget split across
	// enabled members by weight. When 0, each member's max_pending defaults to
	// maxPendingThroughputMultiplier * its resolved throughput.
	TotalMaxPending uint64 `json:"total_max_pending"`
	// AutoRestartFailed restarts members that stopped in the failed state after
	// AutoRestartCooldown seconds. Members stopped normally are never restarted.
	AutoRestartFailed bool `json:"auto_restart_failed"`
	// AutoRestartCooldown is the delay in seconds before a failed member is restarted.
	// 0 falls back to DefaultAutoRestartCooldownSecs.
	AutoRestartCooldown uint64 `json:"auto_restart_cooldown"`
}

// RestartCooldownSecs returns the configured auto-restart cooldown in seconds,
// substituting the default when unset.
func (c *GroupConfig) RestartCooldownSecs() uint64 {
	if c.AutoRestartCooldown == 0 {
		return DefaultAutoRestartCooldownSecs
	}
	return c.AutoRestartCooldown
}

// MemberConfig holds the JSON metadata stored in a member row's group_config column.
type MemberConfig struct {
	Weight    uint64 `json:"weight"`
	Enabled   bool   `json:"enabled"`
	SortOrder int    `json:"sort_order"`
}

// ParseGroupConfig parses a group row's group_config JSON, applying defaults.
// An empty string yields a default GroupConfig in independent mode.
func ParseGroupConfig(s string) (*GroupConfig, error) {
	cfg := &GroupConfig{ThroughputMode: GroupModeIndependent}
	if s == "" {
		return cfg, nil
	}
	if err := json.Unmarshal([]byte(s), cfg); err != nil {
		return nil, fmt.Errorf("failed to parse group config: %w", err)
	}
	if cfg.ThroughputMode == "" {
		cfg.ThroughputMode = GroupModeIndependent
	}
	return cfg, nil
}

// Marshal serializes the group config to its JSON storage representation.
func (c *GroupConfig) Marshal() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal group config: %w", err)
	}
	return string(data), nil
}

// ParseMemberConfig parses a member row's group_config JSON, applying defaults.
// An empty string yields a default MemberConfig (weight 1, enabled).
func ParseMemberConfig(s string) (*MemberConfig, error) {
	cfg := &MemberConfig{Weight: 1, Enabled: true}
	if s == "" {
		return cfg, nil
	}
	if err := json.Unmarshal([]byte(s), cfg); err != nil {
		return nil, fmt.Errorf("failed to parse member config: %w", err)
	}
	return cfg, nil
}

// Marshal serializes the member config to its JSON storage representation.
func (c *MemberConfig) Marshal() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal member config: %w", err)
	}
	return string(data), nil
}

// ScenarioSupportsSharedThroughput reports whether the scenario accepts a
// throughput field and can therefore participate in a shared-throughput group.
func ScenarioSupportsSharedThroughput(descriptor *scenario.Descriptor) bool {
	_, ok := scenario.GetScenarioValidFields(descriptor)[fieldThroughput]
	return ok
}

// ScenarioSupportsSharedCount reports whether the scenario accepts a total_count
// field and can therefore participate in a shared-count group.
func ScenarioSupportsSharedCount(descriptor *scenario.Descriptor) bool {
	_, ok := scenario.GetScenarioValidFields(descriptor)[fieldTotalCount]
	return ok
}

// Apportion distributes total across the given weights using the largest-remainder
// (Hamilton) method so the integer shares sum exactly to total. A weight of 0 yields
// a 0 share. If all weights are 0, total is split as evenly as possible. A 0 share for
// an enabled member is harmless for throughput because ResolveMemberConfig bumps an
// injected 0 to 1 (never "unlimited").
func Apportion(total uint64, weights []uint64) []uint64 {
	n := len(weights)
	shares := make([]uint64, n)
	if n == 0 {
		return shares
	}

	var sumW uint64
	for _, w := range weights {
		sumW += w
	}

	if sumW == 0 {
		// No weights: split as evenly as possible.
		base := total / uint64(n)
		rem := total % uint64(n)
		for i := range shares {
			shares[i] = base
			if uint64(i) < rem {
				shares[i]++
			}
		}
		return shares
	}

	type remainder struct {
		idx int
		rem uint64
	}
	remainders := make([]remainder, n)
	var allocated uint64
	for i, w := range weights {
		num := total * w
		shares[i] = num / sumW
		allocated += shares[i]
		remainders[i] = remainder{idx: i, rem: num % sumW}
	}

	// Distribute the leftover units to the largest fractional remainders. Ties are
	// broken in favor of the lower index by the stable sort, keeping results
	// deterministic.
	remaining := total - allocated
	sort.SliceStable(remainders, func(a, b int) bool {
		return remainders[a].rem > remainders[b].rem
	})
	for i := uint64(0); i < remaining && int(i) < n; i++ {
		shares[remainders[i].idx]++
	}

	return shares
}

// ResolveMemberConfig computes the effective YAML config for a group member. It never
// mutates the stored config and applies, in order:
//   - the group overlay (group wins) for fields the member's scenario accepts,
//   - a shared throughput value (gated, never 0) when sharedThroughput is set,
//   - a shared total_count value (gated) when sharedCount is set,
//   - a max_pending value: the explicit sharedMaxPending when set (gated, never 0),
//     otherwise derived as maxPendingThroughputMultiplier * the effective throughput,
//   - a derived max_wallets from the effective throughput (preferred) or total_count.
//
// Fields the member's scenario does not define are silently skipped, so the result is
// always valid against the scenario's config validator.
func ResolveMemberConfig(
	descriptor *scenario.Descriptor,
	memberConfig string,
	groupOverlay map[string]any,
	sharedThroughput *uint64,
	sharedCount *uint64,
	sharedMaxPending *uint64,
) (string, error) {
	base := map[string]any{}
	if memberConfig != "" {
		if err := yaml.Unmarshal([]byte(memberConfig), &base); err != nil {
			return "", fmt.Errorf("failed to parse member config: %w", err)
		}
	}
	if base == nil {
		base = map[string]any{}
	}

	valid := scenario.GetScenarioValidFields(descriptor)

	// 1. Apply the overlay; the group's value wins for each accepted field.
	for k, v := range groupOverlay {
		if _, ok := valid[k]; ok {
			base[k] = v
		}
	}

	// 2. Inject the shared throughput, never as 0 (which would mean "unlimited").
	if sharedThroughput != nil {
		if _, ok := valid[fieldThroughput]; ok {
			tp := *sharedThroughput
			if tp == 0 {
				tp = 1
			}
			base[fieldThroughput] = tp
		}
	}

	// 3. Inject the shared total_count.
	if sharedCount != nil {
		if _, ok := valid[fieldTotalCount]; ok {
			base[fieldTotalCount] = *sharedCount
		}
	}

	// 4. Resolve max_pending: an explicit group budget share (never 0) when provided,
	// otherwise scaled from the effective throughput.
	if _, ok := valid[fieldMaxPending]; ok {
		if sharedMaxPending != nil {
			mp := *sharedMaxPending
			if mp == 0 {
				mp = 1
			}
			base[fieldMaxPending] = mp
		} else if eff := mapUint64(base, fieldThroughput); eff > 0 {
			base[fieldMaxPending] = eff * maxPendingThroughputMultiplier
		}
	}

	// 5. Derive max_wallets from the effective throughput (preferred) or total_count.
	if _, ok := valid[fieldMaxWallets]; ok {
		effThroughput := mapUint64(base, fieldThroughput)
		effCount := mapUint64(base, fieldTotalCount)
		if mw := deriveMaxWallets(effThroughput, effCount); mw > 0 {
			base[fieldMaxWallets] = mw
		}
	}

	out, err := yaml.Marshal(base)
	if err != nil {
		return "", fmt.Errorf("failed to marshal resolved config: %w", err)
	}
	return string(out), nil
}

// deriveMaxWallets returns a clamped wallet count derived from the effective
// throughput (preferred) or total_count. It returns 0 when neither is set so the
// caller can leave max_wallets untouched.
func deriveMaxWallets(throughput, count uint64) uint64 {
	var base uint64
	switch {
	case throughput > 0:
		base = uint64(math.Round(float64(throughput) / groupMaxWalletsThroughputDivisor))
	case count > 0:
		base = uint64(math.Round(float64(count) / groupMaxWalletsCountDivisor))
	default:
		return 0
	}
	return clampUint64(base, groupMaxWalletsMin, groupMaxWalletsMax)
}

func clampUint64(v, min, max uint64) uint64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// mapUint64 reads a numeric value from a decoded YAML/JSON map and coerces it to
// uint64. Missing keys, negative values and non-numeric values yield 0.
func mapUint64(m map[string]any, key string) uint64 {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case int:
		if n < 0 {
			return 0
		}
		return uint64(n)
	case int64:
		if n < 0 {
			return 0
		}
		return uint64(n)
	case uint64:
		return n
	case uint:
		return uint64(n)
	case float64:
		if n < 0 {
			return 0
		}
		return uint64(n)
	case json.Number:
		i, err := n.Int64()
		if err != nil || i < 0 {
			return 0
		}
		return uint64(i)
	default:
		return 0
	}
}

// GroupConfigFromMap builds a GroupConfig from a decoded YAML map (export/import path).
func GroupConfigFromMap(m map[string]any) *GroupConfig {
	cfg := &GroupConfig{ThroughputMode: GroupModeIndependent}
	if m == nil {
		return cfg
	}
	if v, ok := m["throughput_mode"].(string); ok && v != "" {
		cfg.ThroughputMode = v
	}
	cfg.TotalThroughput = mapUint64(m, "total_throughput")
	cfg.TotalCount = mapUint64(m, "total_count")
	cfg.TotalMaxPending = mapUint64(m, "total_max_pending")
	if v, ok := m["auto_restart_failed"].(bool); ok {
		cfg.AutoRestartFailed = v
	}
	cfg.AutoRestartCooldown = mapUint64(m, "auto_restart_cooldown")
	return cfg
}

// MemberConfigFromMap builds a MemberConfig from a decoded YAML map (export/import path).
func MemberConfigFromMap(m map[string]any) *MemberConfig {
	cfg := &MemberConfig{Weight: 1, Enabled: true}
	if m == nil {
		return cfg
	}
	if _, ok := m["weight"]; ok {
		cfg.Weight = mapUint64(m, "weight")
	}
	if v, ok := m["enabled"].(bool); ok {
		cfg.Enabled = v
	}
	cfg.SortOrder = int(mapUint64(m, "sort_order"))
	return cfg
}

// ResolvedRunConfig pairs a runnable spammer config with its fully-resolved YAML.
// The Config field carries the original (sparse) import config for naming/metadata,
// while ConfigYAML is the complete config string to hand to the scenario, with the
// group overlay, shared throughput/count split and derived max_wallets already
// applied for group members.
type ResolvedRunConfig struct {
	Config     SpammerConfig
	ConfigYAML string
}

// ResolveRunConfigs expands a flat list of import configs into runnable configs for
// the CLI single-run path, giving group members the same effective config they would
// receive in the daemon. Group entries (scenario == "group") are removed from the
// result; every remaining entry gets its config merged with scenario defaults, and
// members additionally get the group overlay and (in shared mode) the weight-based
// throughput/count split applied. The lookup callback resolves a scenario name to its
// descriptor; it must handle every non-group scenario referenced.
func ResolveRunConfigs(all []SpammerConfig, lookup func(string) *scenario.Descriptor) ([]ResolvedRunConfig, error) {
	// Index group entries by name.
	groups := make(map[string]*SpammerConfig, len(all))
	for i := range all {
		if all[i].Scenario == scenario.GroupScenarioName {
			if _, dup := groups[all[i].Name]; dup {
				return nil, fmt.Errorf("duplicate group name %q", all[i].Name)
			}
			groups[all[i].Name] = &all[i]
		}
	}

	// Collect member indices per group, preserving input order.
	membersByGroup := make(map[string][]int)
	for i := range all {
		c := &all[i]
		if c.Scenario == scenario.GroupScenarioName {
			continue
		}
		if c.Group != "" {
			membersByGroup[c.Group] = append(membersByGroup[c.Group], i)
		}
	}

	// Pre-compute shared throughput/count shares per member (index into `all`).
	type shareInfo struct {
		throughput *uint64
		count      *uint64
		maxPending *uint64
	}
	memberShares := make(map[int]shareInfo)
	for groupName, idxs := range membersByGroup {
		group, ok := groups[groupName]
		if !ok {
			return nil, fmt.Errorf("spammer references unknown group %q", groupName)
		}
		gc := GroupConfigFromMap(group.GroupConfig)
		if gc.ThroughputMode != GroupModeShared {
			continue
		}

		type member struct {
			idx    int
			weight uint64
			order  int
		}
		active := make([]member, 0, len(idxs))
		for _, idx := range idxs {
			mc := MemberConfigFromMap(all[idx].GroupConfig)
			if !mc.Enabled {
				continue
			}
			active = append(active, member{idx: idx, weight: mc.Weight, order: mc.SortOrder})
		}
		sort.SliceStable(active, func(a, b int) bool {
			return active[a].order < active[b].order
		})

		weights := make([]uint64, len(active))
		for i := range active {
			weights[i] = active[i].weight
		}

		var tpShares, cntShares, mpShares []uint64
		if gc.TotalThroughput > 0 {
			tpShares = Apportion(gc.TotalThroughput, weights)
		}
		if gc.TotalCount > 0 {
			cntShares = Apportion(gc.TotalCount, weights)
		}
		if gc.TotalMaxPending > 0 {
			mpShares = Apportion(gc.TotalMaxPending, weights)
		}
		for i, m := range active {
			si := shareInfo{}
			if tpShares != nil {
				v := tpShares[i]
				si.throughput = &v
			}
			if cntShares != nil {
				v := cntShares[i]
				si.count = &v
			}
			if mpShares != nil {
				v := mpShares[i]
				si.maxPending = &v
			}
			memberShares[m.idx] = si
		}
	}

	// Produce the runnable list with resolved YAML.
	result := make([]ResolvedRunConfig, 0, len(all))
	for i := range all {
		c := all[i]
		if c.Scenario == scenario.GroupScenarioName {
			continue
		}

		descriptor := lookup(c.Scenario)
		if descriptor == nil {
			return nil, fmt.Errorf("unknown scenario: %s", c.Scenario)
		}

		merged, err := MergeScenarioConfiguration(descriptor, &c.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to merge config for %q: %w", c.Name, err)
		}

		finalYAML := merged
		if c.Group != "" {
			group := groups[c.Group]
			if group == nil {
				return nil, fmt.Errorf("spammer %q references unknown group %q", c.Name, c.Group)
			}
			si := memberShares[i]
			finalYAML, err = ResolveMemberConfig(descriptor, merged, NodeToMap(&group.Config), si.throughput, si.count, si.maxPending)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve config for member %q: %w", c.Name, err)
			}
		}

		result = append(result, ResolvedRunConfig{Config: c, ConfigYAML: finalYAML})
	}

	return result, nil
}
