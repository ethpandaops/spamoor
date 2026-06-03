package handlers

import (
	"net/http"
	"sort"
	"time"

	"github.com/ethpandaops/spamoor/daemon/configs"
	"github.com/ethpandaops/spamoor/webui/server"
)

type IndexPage struct {
	Spammers              []*IndexPageSpammer
	StartupDelayActive    bool
	StartupDelayRemaining int64
}

type IndexPageSpammer struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Scenario    string    `json:"scenario"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`

	// Group fields
	IsGroup bool  `json:"is_group"`
	GroupID int64 `json:"group_id"`

	// Group (parent) fields, populated when IsGroup is true.
	ThroughputMode  string              `json:"throughput_mode,omitempty"`
	TotalThroughput uint64              `json:"total_throughput,omitempty"`
	TotalCount      uint64              `json:"total_count,omitempty"`
	TotalMaxPending uint64              `json:"total_max_pending,omitempty"`
	Members         []*IndexPageSpammer `json:"members,omitempty"`

	// Member fields, populated when GroupID != 0.
	Weight    uint64 `json:"weight,omitempty"`
	Enabled   bool   `json:"enabled,omitempty"`
	SortOrder int    `json:"sort_order,omitempty"`
}

// Index will return the "index" page using a go template
func (fh *FrontendHandler) Index(w http.ResponseWriter, r *http.Request) {
	var templateFiles = append(server.LayoutTemplateFiles,
		"index/index.html",
	)

	var pageTemplate = server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "index", "/", "Spamoor Dashboard", templateFiles)

	var pageError error
	data.Data, pageError = fh.getIndexPageData()
	if pageError != nil {
		server.HandlePageError(w, r, pageError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "index.go", "Index", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return
	}
}

func (fh *FrontendHandler) getIndexPageData() (*IndexPage, error) {
	spammers := fh.daemon.GetAllSpammers()

	// Build a lookup from spammer id to its model, and bucket members by group.
	groups := make(map[int64]*IndexPageSpammer)
	membersByGroup := make(map[int64][]*IndexPageSpammer)

	// Top-level entries preserve the daemon's id-descending order (standalone spammers
	// and group rows interleaved). Members are nested under their group instead.
	topLevel := make([]*IndexPageSpammer, 0, len(spammers))

	for _, s := range spammers {
		model := &IndexPageSpammer{
			ID:          s.GetID(),
			Name:        s.GetName(),
			Description: s.GetDescription(),
			Scenario:    s.GetScenario(),
			Status:      s.GetEffectiveStatus(),
			CreatedAt:   time.Unix(s.GetCreatedAt(), 0),
			IsGroup:     s.IsGroup(),
			GroupID:     s.GetGroupID(),
		}

		switch {
		case s.IsGroup():
			if gc, err := configs.ParseGroupConfig(s.GetGroupConfig()); err == nil {
				model.ThroughputMode = gc.ThroughputMode
				model.TotalThroughput = gc.TotalThroughput
				model.TotalCount = gc.TotalCount
				model.TotalMaxPending = gc.TotalMaxPending
			}
			groups[s.GetID()] = model
			topLevel = append(topLevel, model)
		case s.GetGroupID() != 0:
			if mc, err := configs.ParseMemberConfig(s.GetGroupConfig()); err == nil {
				model.Weight = mc.Weight
				model.Enabled = mc.Enabled
				model.SortOrder = mc.SortOrder
			}
			membersByGroup[s.GetGroupID()] = append(membersByGroup[s.GetGroupID()], model)
		default:
			topLevel = append(topLevel, model)
		}
	}

	// Attach members to their group, ordered by sort_order then id. Orphaned members
	// (group missing) are promoted to top-level so they remain visible and controllable.
	for groupID, members := range membersByGroup {
		sort.SliceStable(members, func(a, b int) bool {
			if members[a].SortOrder != members[b].SortOrder {
				return members[a].SortOrder < members[b].SortOrder
			}
			return members[a].ID < members[b].ID
		})
		if group, ok := groups[groupID]; ok {
			group.Members = members
		} else {
			topLevel = append(topLevel, members...)
		}
	}

	// Check if startup delay is active
	startupDelayActive := fh.daemon.IsInStartupDelay()
	startupDelayRemaining := int64(0)
	if startupDelayActive {
		startupDelayRemaining = int64(fh.daemon.GetStartupDelayRemaining().Seconds())
	}

	return &IndexPage{
		Spammers:              topLevel,
		StartupDelayActive:    startupDelayActive,
		StartupDelayRemaining: startupDelayRemaining,
	}, nil
}
