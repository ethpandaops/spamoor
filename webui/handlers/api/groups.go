package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/ethpandaops/spamoor/daemon/configs"
)

// SpammerGroupRequest is the body for creating or updating a spammer group.
type SpammerGroupRequest struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	Config              string `json:"config"`                // sparse overlay YAML
	ThroughputMode      string `json:"throughput_mode"`       // "independent" or "shared"
	TotalThroughput     uint64 `json:"total_throughput"`      // shared mode: total tx/slot
	TotalCount          uint64 `json:"total_count"`           // shared mode: total tx count
	TotalMaxPending     uint64 `json:"total_max_pending"`     // shared mode: total pending budget (0 = 2x throughput)
	AutoRestartFailed   bool   `json:"auto_restart_failed"`   // restart failed members after a cooldown
	AutoRestartCooldown uint64 `json:"auto_restart_cooldown"` // cooldown in seconds (0 = default 300)
}

func (r *SpammerGroupRequest) groupConfig() *configs.GroupConfig {
	mode := r.ThroughputMode
	if mode == "" {
		mode = configs.GroupModeIndependent
	}
	return &configs.GroupConfig{
		ThroughputMode:      mode,
		TotalThroughput:     r.TotalThroughput,
		TotalCount:          r.TotalCount,
		TotalMaxPending:     r.TotalMaxPending,
		AutoRestartFailed:   r.AutoRestartFailed,
		AutoRestartCooldown: r.AutoRestartCooldown,
	}
}

// SetMemberRequest is the body for assigning a spammer to a group or updating its
// membership metadata.
type SetMemberRequest struct {
	GroupID   int64  `json:"group_id"`
	Weight    uint64 `json:"weight"`
	Enabled   bool   `json:"enabled"`
	SortOrder int    `json:"sort_order"`
}

// ReorderMembersRequest is the body for reordering a group's members.
type ReorderMembersRequest struct {
	Order []int64 `json:"order"`
}

// CreateSpammerGroup godoc
// @Id createSpammerGroup
// @Summary Create a spammer group
// @Tags SpammerGroup
// @Description Creates a new spammer group with a shared overlay and throughput mode
// @Accept json
// @Produce json
// @Param request body SpammerGroupRequest true "Group configuration"
// @Success 200 {object} int64 "Group ID"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer-group [post]
func (ah *APIHandler) CreateSpammerGroup(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	var req SpammerGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userEmail := ah.getUserEmail(r)

	group, err := ah.daemon.NewGroup(req.Name, req.Description, req.Config, req.groupConfig(), userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group.GetID())
}

// UpdateSpammerGroup godoc
// @Id updateSpammerGroup
// @Summary Update a spammer group
// @Tags SpammerGroup
// @Description Updates a group's name, description, overlay and throughput mode/totals
// @Accept json
// @Param id path int true "Group ID"
// @Param request body SpammerGroupRequest true "Group configuration"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Group not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer-group/{id} [put]
func (ah *APIHandler) UpdateSpammerGroup(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	var req SpammerGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userEmail := ah.getUserEmail(r)

	if err := ah.daemon.UpdateGroup(id, req.Name, req.Description, req.Config, req.groupConfig(), userEmail); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SetSpammerGroup godoc
// @Id setSpammerGroup
// @Summary Add or update a spammer's group membership
// @Tags SpammerGroup
// @Description Assigns a spammer to a group (or updates its weight/enabled/sort_order)
// @Accept json
// @Param id path int true "Spammer ID"
// @Param request body SetMemberRequest true "Membership configuration"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id}/group [put]
func (ah *APIHandler) SetSpammerGroup(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	var req SetMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userEmail := ah.getUserEmail(r)
	member := &configs.MemberConfig{
		Weight:    req.Weight,
		Enabled:   req.Enabled,
		SortOrder: req.SortOrder,
	}

	if err := ah.daemon.AddSpammerToGroup(id, req.GroupID, member, userEmail); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveSpammerGroup godoc
// @Id removeSpammerGroup
// @Summary Detach a spammer from its group
// @Tags SpammerGroup
// @Description Removes a spammer from its group, leaving it a standalone spammer
// @Param id path int true "Spammer ID"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid spammer ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id}/group [delete]
func (ah *APIHandler) RemoveSpammerGroup(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	userEmail := ah.getUserEmail(r)

	if err := ah.daemon.RemoveSpammerFromGroup(id, userEmail); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ReorderGroupMembers godoc
// @Id reorderGroupMembers
// @Summary Reorder a group's members
// @Tags SpammerGroup
// @Description Sets the display/sort order of a group's members by id
// @Accept json
// @Param id path int true "Group ID"
// @Param request body ReorderMembersRequest true "Ordered member ids"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Group not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer-group/{id}/members/order [put]
func (ah *APIHandler) ReorderGroupMembers(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	var req ReorderMembersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userEmail := ah.getUserEmail(r)

	if err := ah.daemon.ReorderGroupMembers(id, req.Order, userEmail); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetEffectiveConfig godoc
// @Id getEffectiveConfig
// @Summary Get a spammer's effective (resolved) config
// @Tags SpammerGroup
// @Description Returns the config a spammer would actually run with. For group members
// @Description this includes the group overlay and (in shared mode) the resolved
// @Description throughput/count split and derived max_wallets.
// @Produce text/plain
// @Param id path int true "Spammer ID"
// @Success 200 {string} string "Resolved YAML configuration"
// @Failure 400 {string} string "Invalid spammer ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id}/effective-config [get]
func (ah *APIHandler) GetEffectiveConfig(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	resolved, err := ah.daemon.GetEffectiveConfig(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write([]byte(resolved))
}
