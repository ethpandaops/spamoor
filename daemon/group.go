package daemon

import (
	"fmt"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v3"

	"github.com/ethpandaops/spamoor/daemon/configs"
	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/daemon/logscope"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/sirupsen/logrus"
)

// memberShare holds the resolved shared throughput/count for a single group member.
// Nil pointers mean the corresponding value is not managed by the group.
type memberShare struct {
	throughput *uint64
	count      *uint64
	maxPending *uint64
}

// NewGroup creates a new spammer group. A group is stored as a regular spammers row
// with the reserved sentinel scenario name; overlayConfig holds the sparse YAML overlay
// applied to members and groupConfig holds the mode/totals metadata.
func (d *Daemon) NewGroup(name string, description string, overlayConfig string, groupConfig *configs.GroupConfig, userEmail string) (*Spammer, error) {
	// Validate the overlay is parseable YAML so members never fail at resolve time.
	if overlayConfig != "" {
		var probe map[string]any
		if err := yaml.Unmarshal([]byte(overlayConfig), &probe); err != nil {
			return nil, fmt.Errorf("invalid overlay config: %w", err)
		}
	}

	if groupConfig == nil {
		groupConfig = &configs.GroupConfig{ThroughputMode: configs.GroupModeIndependent}
	}
	if groupConfig.ThroughputMode == "" {
		groupConfig.ThroughputMode = configs.GroupModeIndependent
	}

	groupConfigJSON, err := groupConfig.Marshal()
	if err != nil {
		return nil, err
	}

	// Reuse the shared scenario id counter so group ids never collide with spammer ids.
	d.spammerIdMtx.Lock()
	scenarioCounter := 0
	d.db.GetSpamoorState("scenario_counter", &scenarioCounter)
	if scenarioCounter < 100 {
		scenarioCounter = 100
	} else {
		scenarioCounter++
	}
	d.db.SetSpamoorState(nil, "scenario_counter", scenarioCounter)
	d.spammerIdMtx.Unlock()

	dbEntity := &db.Spammer{
		ID:          int64(scenarioCounter),
		Scenario:    scenario.GroupScenarioName,
		Name:        name,
		Description: description,
		Config:      overlayConfig,
		Status:      int(SpammerStatusPaused),
		CreatedAt:   time.Now().Unix(),
		GroupConfig: groupConfigJSON,
	}

	logger := logscope.NewLogger(&logscope.ScopeOptions{
		Parent:     d.logger.WithField("group_id", dbEntity.ID),
		BufferSize: 1000,
	})
	logger.GetLogger().SetLevel(logrus.GetLevel())

	group := &Spammer{
		daemon:   d,
		dbEntity: dbEntity,
		logscope: logger,
		logger:   logger.GetLogger(),
	}

	err = d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		if err := d.db.InsertSpammer(tx, dbEntity); err != nil {
			return err
		}
		if d.auditLogger != nil && userEmail != "" {
			return d.auditLogger.LogSpammerCreate(tx, userEmail, dbEntity.ID, name, scenario.GroupScenarioName, overlayConfig)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	d.spammerMapMtx.Lock()
	d.spammerMap[group.dbEntity.ID] = group
	d.spammerMapMtx.Unlock()

	d.emitSpammerSnapshot(group, SpammerEventCreated)

	return group, nil
}

// UpdateGroup updates a group's name, description, overlay config and group config.
func (d *Daemon) UpdateGroup(id int64, name string, description string, overlayConfig string, groupConfig *configs.GroupConfig, userEmail string) error {
	group := d.GetSpammer(id)
	if group == nil || !group.IsGroup() {
		return fmt.Errorf("group not found")
	}

	if overlayConfig != "" {
		var probe map[string]any
		if err := yaml.Unmarshal([]byte(overlayConfig), &probe); err != nil {
			return fmt.Errorf("invalid overlay config: %w", err)
		}
	}

	if groupConfig == nil {
		groupConfig = &configs.GroupConfig{ThroughputMode: configs.GroupModeIndependent}
	}
	if groupConfig.ThroughputMode == "" {
		groupConfig.ThroughputMode = configs.GroupModeIndependent
	}

	groupConfigJSON, err := groupConfig.Marshal()
	if err != nil {
		return err
	}

	// When switching to shared mode, reject members whose scenario cannot take the
	// shared field(s) so the group never produces an invalid member config.
	if groupConfig.ThroughputMode == configs.GroupModeShared {
		if err := d.validateMembersForSharedMode(id, groupConfig); err != nil {
			return err
		}
	}

	d.spammerMapMtx.Lock()
	oldName := group.dbEntity.Name
	oldConfig := group.dbEntity.Config
	group.dbEntity.Name = name
	group.dbEntity.Description = description
	group.dbEntity.Config = overlayConfig
	group.dbEntity.GroupConfig = groupConfigJSON
	d.spammerMapMtx.Unlock()

	err = d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		if err := d.db.UpdateSpammer(tx, group.dbEntity); err != nil {
			return err
		}
		if d.auditLogger != nil && userEmail != "" {
			return d.auditLogger.LogSpammerUpdate(tx, userEmail, id, oldName, name, "", description, oldConfig, overlayConfig)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	d.emitSpammerSnapshot(group, SpammerEventUpdated)

	return nil
}

// scenarioSharedCompatible reports whether a scenario can participate in a shared-mode
// group: it must support at least one of the dimensions (throughput/total_count) the
// group actually manages. A scenario that supports neither managed dimension (e.g. a
// throughput-less scenario in a shared-throughput group) is rejected at write time.
func scenarioSharedCompatible(descriptor *scenario.Descriptor, groupConfig *configs.GroupConfig) bool {
	if descriptor == nil {
		return true
	}
	managesThroughput := groupConfig.TotalThroughput > 0
	managesCount := groupConfig.TotalCount > 0
	if !managesThroughput && !managesCount {
		// Shared mode but no totals set yet: nothing to split, so any scenario is fine.
		return true
	}
	if managesThroughput && configs.ScenarioSupportsSharedThroughput(descriptor) {
		return true
	}
	if managesCount && configs.ScenarioSupportsSharedCount(descriptor) {
		return true
	}
	return false
}

// validateMembersForSharedMode ensures every member of the group can participate in the
// shared split. Returns a descriptive error for the first incompatible member found.
func (d *Daemon) validateMembersForSharedMode(groupID int64, groupConfig *configs.GroupConfig) error {
	members := d.getGroupMembersFromMap(groupID)
	for _, m := range members {
		descriptor := scenarios.GetScenario(m.GetScenario())
		if !scenarioSharedCompatible(descriptor, groupConfig) {
			return fmt.Errorf("member %q (%s) has no throughput or total_count field and cannot join a shared-throughput group", m.GetName(), m.GetScenario())
		}
	}
	return nil
}

// startGroup starts every enabled member of the group. Members are snapshotted under
// the read lock and started outside it to avoid holding the map mutex across Start
// (which would block other operations). Starting individual members is independent of
// this and never started here.
func (s *Spammer) startGroup() error {
	members := s.daemon.getStartableGroupMembers(s.dbEntity.ID)
	var firstErr error
	for _, m := range members {
		if err := m.Start(); err != nil {
			s.logger.Errorf("failed to start group member %d (%s): %v", m.GetID(), m.GetName(), err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	// Reflect the group's aggregate (derived) status on the dashboard (deduped — a member
	// start above may have already emitted it).
	s.daemon.emitGroupStatusIfChanged(s.dbEntity.ID)
	return firstErr
}

// pauseGroup pauses every running member of the group using the same lock-free
// snapshot pattern as startGroup.
func (s *Spammer) pauseGroup() error {
	members := s.daemon.getGroupMembersFromMap(s.dbEntity.ID)
	var firstErr error
	for _, m := range members {
		if m.scenarioCancel == nil {
			continue
		}
		if err := m.Pause(); err != nil {
			s.logger.Errorf("failed to pause group member %d (%s): %v", m.GetID(), m.GetName(), err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	// Reflect the group's aggregate (derived) status on the dashboard (deduped).
	s.daemon.emitGroupStatusIfChanged(s.dbEntity.ID)
	return firstErr
}

// GetEffectiveStatus returns the display status for this spammer. For group rows the
// status is derived from members; for normal spammers it is the stored status.
func (s *Spammer) GetEffectiveStatus() int {
	if s.IsGroup() {
		return s.daemon.computeGroupStatus(s.dbEntity.ID)
	}
	return s.dbEntity.Status
}

// resolveEffectiveConfig returns the config a member should actually run with. For
// standalone spammers the stored config is returned verbatim (byte-for-byte). For group
// members it applies the group overlay and, in shared mode, the weight-based split.
func (s *Spammer) resolveEffectiveConfig() (string, error) {
	if s.dbEntity.GroupID == 0 {
		return s.dbEntity.Config, nil
	}

	group := s.daemon.GetSpammer(s.dbEntity.GroupID)
	if group == nil || !group.IsGroup() {
		// Orphaned member (parent missing): behave exactly like a standalone spammer.
		return s.dbEntity.Config, nil
	}

	var overlay map[string]any
	if group.dbEntity.Config != "" {
		if err := yaml.Unmarshal([]byte(group.dbEntity.Config), &overlay); err != nil {
			return "", fmt.Errorf("failed to parse group overlay: %w", err)
		}
	}

	groupCfg, err := configs.ParseGroupConfig(group.dbEntity.GroupConfig)
	if err != nil {
		return "", err
	}

	descriptor := scenarios.GetScenario(s.dbEntity.Scenario)
	if descriptor == nil {
		return "", fmt.Errorf("scenario %s not found", s.dbEntity.Scenario)
	}

	var sharedThroughput, sharedCount, sharedMaxPending *uint64
	if groupCfg.ThroughputMode == configs.GroupModeShared {
		shares := s.daemon.computeMemberShares(group.GetID(), groupCfg)
		if share, ok := shares[s.dbEntity.ID]; ok {
			sharedThroughput = share.throughput
			sharedCount = share.count
			sharedMaxPending = share.maxPending
		}
	}

	return configs.ResolveMemberConfig(descriptor, s.dbEntity.Config, overlay, sharedThroughput, sharedCount, sharedMaxPending)
}

// GetEffectiveConfig returns the resolved config a spammer would run with: the stored
// config for standalone spammers, or the group-resolved config for members.
func (d *Daemon) GetEffectiveConfig(id int64) (string, error) {
	spammer := d.GetSpammer(id)
	if spammer == nil {
		return "", fmt.Errorf("spammer not found")
	}
	if spammer.IsGroup() {
		return "", fmt.Errorf("groups have no runnable config")
	}
	return spammer.resolveEffectiveConfig()
}

// computeMemberShares computes the per-member shared throughput/count for a shared-mode
// group using the largest-remainder apportionment over enabled, weight>0 members.
func (d *Daemon) computeMemberShares(groupID int64, groupCfg *configs.GroupConfig) map[int64]memberShare {
	members := d.getGroupMembersFromMap(groupID)

	type active struct {
		id     int64
		weight uint64
		order  int
	}
	act := make([]active, 0, len(members))
	for _, m := range members {
		mc, err := configs.ParseMemberConfig(m.dbEntity.GroupConfig)
		if err != nil {
			continue
		}
		// Enabled members participate in the split; weight only sets the proportion (a
		// weight-0 member still runs, at the resolver's min-1 throughput).
		if !mc.Enabled {
			continue
		}
		act = append(act, active{id: m.GetID(), weight: mc.Weight, order: mc.SortOrder})
	}

	sort.SliceStable(act, func(a, b int) bool {
		if act[a].order != act[b].order {
			return act[a].order < act[b].order
		}
		return act[a].id < act[b].id
	})

	weights := make([]uint64, len(act))
	for i := range act {
		weights[i] = act[i].weight
	}

	var tpShares, cntShares, mpShares []uint64
	if groupCfg.TotalThroughput > 0 {
		tpShares = configs.Apportion(groupCfg.TotalThroughput, weights)
	}
	if groupCfg.TotalCount > 0 {
		cntShares = configs.Apportion(groupCfg.TotalCount, weights)
	}
	if groupCfg.TotalMaxPending > 0 {
		mpShares = configs.Apportion(groupCfg.TotalMaxPending, weights)
	}

	res := make(map[int64]memberShare, len(act))
	for i, a := range act {
		ms := memberShare{}
		if tpShares != nil {
			v := tpShares[i]
			ms.throughput = &v
		}
		if cntShares != nil {
			v := cntShares[i]
			ms.count = &v
		}
		if mpShares != nil {
			v := mpShares[i]
			ms.maxPending = &v
		}
		res[a.id] = ms
	}
	return res
}

// getGroupMembersFromMap returns a snapshot of all member spammers belonging to the
// given group id, taken under the read lock.
func (d *Daemon) getGroupMembersFromMap(groupID int64) []*Spammer {
	d.spammerMapMtx.RLock()
	defer d.spammerMapMtx.RUnlock()

	members := make([]*Spammer, 0, 8)
	for _, s := range d.spammerMap {
		if s.dbEntity.GroupID == groupID {
			members = append(members, s)
		}
	}
	return members
}

// GetGroupMembers returns the members of a group ordered by their configured sort_order
// (then by id for stability). Safe to call concurrently.
func (d *Daemon) GetGroupMembers(groupID int64) []*Spammer {
	members := d.getGroupMembersFromMap(groupID)
	sort.SliceStable(members, func(a, b int) bool {
		oa := memberSortOrder(members[a])
		ob := memberSortOrder(members[b])
		if oa != ob {
			return oa < ob
		}
		return members[a].GetID() < members[b].GetID()
	})
	return members
}

// memberSortOrder extracts a member's sort_order, defaulting to 0 on parse failure.
func memberSortOrder(s *Spammer) int {
	mc, err := configs.ParseMemberConfig(s.dbEntity.GroupConfig)
	if err != nil {
		return 0
	}
	return mc.SortOrder
}

// getStartableGroupMembers returns the members that should run when the group starts:
// those whose enabled flag is set. Weight does not gate participation.
func (d *Daemon) getStartableGroupMembers(groupID int64) []*Spammer {
	if group := d.GetSpammer(groupID); group == nil {
		return nil
	}

	members := d.GetGroupMembers(groupID)
	startable := make([]*Spammer, 0, len(members))
	for _, m := range members {
		mc, err := configs.ParseMemberConfig(m.dbEntity.GroupConfig)
		if err != nil {
			continue
		}
		// Group start only starts enabled members; weight does not gate participation.
		if mc.Enabled {
			startable = append(startable, m)
		}
	}
	return startable
}

// computeGroupStatus derives the aggregate status of a group from its members:
// running if any member runs, else failed if any failed, else finished if all
// finished, else paused.
func (d *Daemon) computeGroupStatus(groupID int64) int {
	members := d.getGroupMembersFromMap(groupID)
	if len(members) == 0 {
		return int(SpammerStatusPaused)
	}

	anyRunning := false
	anyFailed := false
	allFinished := true
	for _, m := range members {
		switch SpammerStatus(m.GetStatus()) {
		case SpammerStatusRunning:
			anyRunning = true
			allFinished = false
		case SpammerStatusFailed:
			anyFailed = true
			allFinished = false
		case SpammerStatusFinished:
			// keeps allFinished true
		default:
			allFinished = false
		}
	}

	switch {
	case anyRunning:
		return int(SpammerStatusRunning)
	case anyFailed:
		return int(SpammerStatusFailed)
	case allFinished:
		return int(SpammerStatusFinished)
	default:
		return int(SpammerStatusPaused)
	}
}

// AddSpammerToGroup assigns a spammer to a group with the given member metadata. It
// validates that the target is a real spammer (not a group), the group exists, and the
// scenario is compatible with a shared-mode group.
func (d *Daemon) AddSpammerToGroup(spammerID int64, groupID int64, member *configs.MemberConfig, userEmail string) error {
	spammer := d.GetSpammer(spammerID)
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}
	if spammer.IsGroup() {
		return fmt.Errorf("cannot add a group to another group")
	}

	group := d.GetSpammer(groupID)
	if group == nil || !group.IsGroup() {
		return fmt.Errorf("group not found")
	}

	groupCfg, err := configs.ParseGroupConfig(group.dbEntity.GroupConfig)
	if err != nil {
		return err
	}

	if groupCfg.ThroughputMode == configs.GroupModeShared {
		descriptor := scenarios.GetScenario(spammer.GetScenario())
		if !scenarioSharedCompatible(descriptor, groupCfg) {
			return fmt.Errorf("scenario %s has no throughput or total_count field and cannot join a shared-throughput group", spammer.GetScenario())
		}
	}

	if member == nil {
		member = &configs.MemberConfig{Weight: 1, Enabled: true}
	}

	return d.writeMembership(spammer, groupID, member, userEmail)
}

// UpdateGroupMember updates the member metadata (weight/enabled/sort_order) for a
// spammer already in a group.
func (d *Daemon) UpdateGroupMember(spammerID int64, member *configs.MemberConfig, userEmail string) error {
	spammer := d.GetSpammer(spammerID)
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}
	if spammer.GetGroupID() == 0 {
		return fmt.Errorf("spammer is not a group member")
	}
	if member == nil {
		return fmt.Errorf("member config required")
	}
	return d.writeMembership(spammer, spammer.GetGroupID(), member, userEmail)
}

// writeMembership persists a member's group_id and group_config JSON.
func (d *Daemon) writeMembership(spammer *Spammer, groupID int64, member *configs.MemberConfig, userEmail string) error {
	memberJSON, err := member.Marshal()
	if err != nil {
		return err
	}

	d.spammerMapMtx.Lock()
	spammer.dbEntity.GroupID = groupID
	spammer.dbEntity.GroupConfig = memberJSON
	d.spammerMapMtx.Unlock()

	err = d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return d.db.UpdateSpammer(tx, spammer.dbEntity)
	})
	if err != nil {
		return fmt.Errorf("failed to update group membership: %w", err)
	}

	if d.auditLogger != nil && userEmail != "" {
		_ = d.auditLogger.LogSpammerAction(userEmail, db.AuditActionSpammerUpdate, spammer.GetID(), spammer.GetName(), map[string]interface{}{
			"group_id":     groupID,
			"group_config": memberJSON,
		})
	}

	d.emitSpammerSnapshot(spammer, SpammerEventMembership)

	return nil
}

// RemoveSpammerFromGroup detaches a spammer from its group, leaving it a fully working
// standalone spammer (its stored config is never touched).
func (d *Daemon) RemoveSpammerFromGroup(spammerID int64, userEmail string) error {
	spammer := d.GetSpammer(spammerID)
	if spammer == nil {
		return fmt.Errorf("spammer not found")
	}
	if spammer.GetGroupID() == 0 {
		return nil
	}

	d.spammerMapMtx.Lock()
	spammer.dbEntity.GroupID = 0
	spammer.dbEntity.GroupConfig = ""
	d.spammerMapMtx.Unlock()

	err := d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return d.db.UpdateSpammer(tx, spammer.dbEntity)
	})
	if err != nil {
		return fmt.Errorf("failed to detach spammer from group: %w", err)
	}

	if d.auditLogger != nil && userEmail != "" {
		_ = d.auditLogger.LogSpammerAction(userEmail, db.AuditActionSpammerUpdate, spammerID, spammer.GetName(), map[string]interface{}{
			"group_id": 0,
		})
	}

	d.emitSpammerSnapshot(spammer, SpammerEventMembership)

	return nil
}

// ReorderGroupMembers assigns sort_order values to members of a group according to the
// position of their ids in orderedIDs. Ids not in the group are ignored.
func (d *Daemon) ReorderGroupMembers(groupID int64, orderedIDs []int64, userEmail string) error {
	group := d.GetSpammer(groupID)
	if group == nil || !group.IsGroup() {
		return fmt.Errorf("group not found")
	}

	for order, id := range orderedIDs {
		member := d.GetSpammer(id)
		if member == nil || member.GetGroupID() != groupID {
			continue
		}
		mc, err := configs.ParseMemberConfig(member.dbEntity.GroupConfig)
		if err != nil {
			mc = &configs.MemberConfig{Weight: 1, Enabled: true}
		}
		mc.SortOrder = order
		if err := d.writeMembership(member, groupID, mc, ""); err != nil {
			return err
		}
	}

	if d.auditLogger != nil && userEmail != "" {
		_ = d.auditLogger.LogSpammerAction(userEmail, db.AuditActionSpammerUpdate, groupID, group.GetName(), map[string]interface{}{
			"reorder": orderedIDs,
		})
	}

	d.emitSpammerEvent(&SpammerEvent{Type: SpammerEventReorder, ID: groupID, Order: orderedIDs})

	return nil
}

// DeleteGroup deletes a group row. When cascade is true the members are paused and
// deleted as well; otherwise members are detached to standalone spammers first.
func (d *Daemon) DeleteGroup(id int64, cascade bool, userEmail string) error {
	group := d.GetSpammer(id)
	if group == nil || !group.IsGroup() {
		return fmt.Errorf("group not found")
	}
	groupName := group.GetName()

	// Snapshot member ids before mutating anything (each op locks internally).
	members := d.getGroupMembersFromMap(id)

	for _, m := range members {
		if cascade {
			if err := d.DeleteSpammer(m.GetID(), userEmail); err != nil {
				d.logger.Errorf("failed to delete group member %d during cascade: %v", m.GetID(), err)
			}
		} else {
			if err := d.RemoveSpammerFromGroup(m.GetID(), userEmail); err != nil {
				d.logger.Errorf("failed to detach group member %d: %v", m.GetID(), err)
			}
		}
	}

	// Delete the group row itself.
	d.spammerMapMtx.Lock()
	err := d.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		if err := d.db.DeleteSpammer(tx, id); err != nil {
			return err
		}
		if d.auditLogger != nil {
			return d.auditLogger.LogSpammerDelete(tx, userEmail, id, groupName)
		}
		return nil
	})
	if err != nil {
		d.spammerMapMtx.Unlock()
		return fmt.Errorf("failed to delete group: %w", err)
	}
	delete(d.spammerMap, id)
	d.spammerMapMtx.Unlock()

	d.emitSpammerDeleted(id)

	return nil
}
