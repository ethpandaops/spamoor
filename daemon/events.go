package daemon

import (
	"github.com/ethpandaops/spamoor/daemon/configs"
)

// Spammer event types emitted on the lifecycle event stream.
const (
	// SpammerEventCreated is emitted when a spammer or group is created.
	SpammerEventCreated = "created"
	// SpammerEventUpdated is emitted when a spammer's or group's name/description/
	// scenario or group settings change.
	SpammerEventUpdated = "updated"
	// SpammerEventStatus is emitted when a spammer's run state changes (and for the
	// parent group's derived status).
	SpammerEventStatus = "status"
	// SpammerEventMembership is emitted when a spammer's group membership changes
	// (added, detached, weight/enabled/order updated).
	SpammerEventMembership = "membership"
	// SpammerEventReorder is emitted when a group's members are reordered.
	SpammerEventReorder = "reorder"
	// SpammerEventDeleted is emitted when a spammer or group is deleted.
	SpammerEventDeleted = "deleted"
)

// SpammerEventInfo is the safe, public snapshot of a spammer carried on the event
// stream. It deliberately contains NO config, seed or logs so the stream can be served
// unauthenticated. Its shape mirrors the spammer list API entry.
type SpammerEventInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scenario    string `json:"scenario"`
	Status      int    `json:"status"`
	IsGroup     bool   `json:"is_group"`
	GroupID     int64  `json:"group_id"`

	// Group (parent) summary, set when IsGroup is true.
	ThroughputMode  string `json:"throughput_mode,omitempty"`
	TotalThroughput uint64 `json:"total_throughput,omitempty"`
	TotalCount      uint64 `json:"total_count,omitempty"`
	TotalMaxPending uint64 `json:"total_max_pending,omitempty"`

	// Member fields, set when GroupID != 0. Enabled is always serialized so the client
	// can distinguish a disabled member (false) from a non-member (field absent).
	Weight    uint64 `json:"weight,omitempty"`
	Enabled   bool   `json:"enabled"`
	SortOrder int    `json:"sort_order,omitempty"`
}

// SpammerEvent is a single lifecycle event published to subscribers of the spammer
// event stream.
type SpammerEvent struct {
	Type    string            `json:"type"`
	ID      int64             `json:"id"`
	Spammer *SpammerEventInfo `json:"spammer,omitempty"`
	Order   []int64           `json:"order,omitempty"` // reorder events only
}

// SubscribeEvents registers a new subscriber for spammer lifecycle events and returns a
// subscription id and a buffered receive channel. Mirrors the metrics collector pub/sub.
func (d *Daemon) SubscribeEvents() (uint64, <-chan *SpammerEvent) {
	d.eventSubMtx.Lock()
	defer d.eventSubMtx.Unlock()

	id := d.eventSubID
	d.eventSubID++

	ch := make(chan *SpammerEvent, 64)
	d.eventSubs[id] = ch
	return id, ch
}

// UnsubscribeEvents removes a subscription and closes its channel.
func (d *Daemon) UnsubscribeEvents(id uint64) {
	d.eventSubMtx.Lock()
	defer d.eventSubMtx.Unlock()

	if ch, ok := d.eventSubs[id]; ok {
		close(ch)
		delete(d.eventSubs, id)
	}
}

// emitSpammerEvent delivers an event to all subscribers without blocking. Slow
// subscribers that have filled their buffer simply miss the event (the dashboard
// re-syncs from the API on the next event or reconnect).
//
// Callers MUST NOT hold spammerMapMtx: building a group snapshot reads the map under
// RLock, which would deadlock against the write lock.
func (d *Daemon) emitSpammerEvent(event *SpammerEvent) {
	d.eventSubMtx.RLock()
	defer d.eventSubMtx.RUnlock()

	for _, ch := range d.eventSubs {
		select {
		case ch <- event:
		default:
		}
	}
}

// buildSpammerEventInfo builds the safe public snapshot for a spammer or group.
func (d *Daemon) buildSpammerEventInfo(s *Spammer) *SpammerEventInfo {
	info := &SpammerEventInfo{
		ID:          s.GetID(),
		Name:        s.GetName(),
		Description: s.GetDescription(),
		Scenario:    s.GetScenario(),
		Status:      s.GetEffectiveStatus(),
		IsGroup:     s.IsGroup(),
		GroupID:     s.GetGroupID(),
	}

	switch {
	case s.IsGroup():
		if gc, err := configs.ParseGroupConfig(s.GetGroupConfig()); err == nil {
			info.ThroughputMode = gc.ThroughputMode
			info.TotalThroughput = gc.TotalThroughput
			info.TotalCount = gc.TotalCount
			info.TotalMaxPending = gc.TotalMaxPending
		}
	case s.GetGroupID() != 0:
		if mc, err := configs.ParseMemberConfig(s.GetGroupConfig()); err == nil {
			info.Weight = mc.Weight
			info.Enabled = mc.Enabled
			info.SortOrder = mc.SortOrder
		}
	}

	return info
}

// emitSpammerSnapshot emits an event carrying the spammer's current safe snapshot.
func (d *Daemon) emitSpammerSnapshot(s *Spammer, eventType string) {
	if s == nil {
		return
	}
	d.emitSpammerEvent(&SpammerEvent{
		Type:    eventType,
		ID:      s.GetID(),
		Spammer: d.buildSpammerEventInfo(s),
	})
}

// emitSpammerDeleted emits a deletion event for the given id.
func (d *Daemon) emitSpammerDeleted(id int64) {
	d.emitSpammerEvent(&SpammerEvent{Type: SpammerEventDeleted, ID: id})
}

// emitStatusChange emits a status event for the spammer and, when it is a group member,
// a status event for the parent group — but only if the group's derived status actually
// changed. Must be called without holding spammerMapMtx.
func (d *Daemon) emitStatusChange(s *Spammer) {
	if s == nil {
		return
	}
	d.emitSpammerSnapshot(s, SpammerEventStatus)

	if s.GetGroupID() != 0 {
		d.emitGroupStatusIfChanged(s.GetGroupID())
	}
}

// emitGroupStatusIfChanged emits a status event for the group only when its derived
// status differs from the last one observed. This collapses the burst that would
// otherwise fire once per member when a group's members all start/stop together (the
// group's cached status field is reused as the last-observed marker; the authoritative
// status is always derived). Must be called without holding spammerMapMtx.
func (d *Daemon) emitGroupStatusIfChanged(groupID int64) {
	group := d.GetSpammer(groupID)
	if group == nil || !group.IsGroup() {
		return
	}

	// Compute the derived status before taking the write lock (computeGroupStatus takes
	// the read lock internally).
	derived := d.computeGroupStatus(groupID)

	d.spammerMapMtx.Lock()
	changed := group.dbEntity.Status != derived
	if changed {
		group.dbEntity.Status = derived
	}
	d.spammerMapMtx.Unlock()

	if changed {
		d.emitSpammerSnapshot(group, SpammerEventStatus)
	}
}
