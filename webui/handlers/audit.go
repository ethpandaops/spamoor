package handlers

import (
	"net/http"
	"time"

	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/webui/server"
)

type AuditPage struct {
	Stats   *AuditPageStats   `json:"stats"`
	Logs    []*AuditPageEntry `json:"logs"`
	Filters *AuditPageFilters `json:"filters"`
}

type AuditPageStats struct {
	TotalLogs        int64            `json:"total_logs"`
	LogsByActionType map[string]int64 `json:"logs_by_action_type"`
	LogsByEntityType map[string]int64 `json:"logs_by_entity_type"`
	LogsByUser       map[string]int64 `json:"logs_by_user"`
}

type AuditPageEntry struct {
	ID         int64     `json:"id"`
	UserEmail  string    `json:"user_email"`
	ActionType string    `json:"action_type"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	EntityName string    `json:"entity_name"`
	Diff       string    `json:"diff"`
	Metadata   string    `json:"metadata"`
	Timestamp  time.Time `json:"timestamp"`
}

type AuditPageFilters struct {
	UserEmail  string `json:"user_email"`
	ActionType string `json:"action_type"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

// Audit will return the "audit" page using a go template
func (fh *FrontendHandler) Audit(w http.ResponseWriter, r *http.Request) {
	var templateFiles = append(server.LayoutTemplateFiles,
		"audit/audit.html",
	)

	var pageTemplate = server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "audit", "/audit", "Audit Logs", templateFiles)

	var pageError error
	data.Data, pageError = fh.getAuditPageData(r)
	if pageError != nil {
		server.HandlePageError(w, r, pageError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "audit.go", "Audit", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return
	}
}

func (fh *FrontendHandler) getAuditPageData(r *http.Request) (*AuditPage, error) {
	// Get query parameters for filtering
	query := r.URL.Query()

	filters := AuditPageFilters{
		UserEmail:  query.Get("user_email"),
		ActionType: query.Get("action_type"),
		EntityType: query.Get("entity_type"),
		EntityID:   query.Get("entity_id"),
	}

	// Get recent logs (last 100)
	logFilters := db.AuditLogFilters{
		UserEmail:  filters.UserEmail,
		ActionType: filters.ActionType,
		EntityType: filters.EntityType,
		EntityID:   filters.EntityID,
		Limit:      100,
	}

	dbLogs, err := fh.daemon.GetDatabase().GetAuditLogs(logFilters)
	if err != nil {
		return nil, err
	}

	// Convert to page format
	logs := make([]*AuditPageEntry, len(dbLogs))
	for i, log := range dbLogs {
		logs[i] = &AuditPageEntry{
			ID:         log.ID,
			UserEmail:  log.UserEmail,
			ActionType: log.ActionType,
			EntityType: log.EntityType,
			EntityID:   log.EntityID,
			EntityName: log.EntityName,
			Diff:       log.Diff,
			Metadata:   log.Metadata,
			Timestamp:  time.Unix(log.Timestamp, 0),
		}
	}

	// Get statistics for the last 30 days
	startTime := time.Now().AddDate(0, 0, -30).Unix()
	statsFilters := db.AuditLogFilters{
		StartTime: startTime,
		Limit:     10000, // Large limit for stats
	}

	allLogs, err := fh.daemon.GetDatabase().GetAuditLogs(statsFilters)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	logsByActionType := make(map[string]int64)
	logsByEntityType := make(map[string]int64)
	logsByUser := make(map[string]int64)

	for _, log := range allLogs {
		logsByActionType[log.ActionType]++
		logsByEntityType[log.EntityType]++
		logsByUser[log.UserEmail]++
	}

	stats := &AuditPageStats{
		TotalLogs:        int64(len(allLogs)),
		LogsByActionType: logsByActionType,
		LogsByEntityType: logsByEntityType,
		LogsByUser:       logsByUser,
	}

	return &AuditPage{
		Stats:   stats,
		Logs:    logs,
		Filters: &filters,
	}, nil
}
