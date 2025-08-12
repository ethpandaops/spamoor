package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/gorilla/mux"
)

// AuditLogEntry represents an audit log entry in the API response
type AuditLogEntry struct {
	ID         int64  `json:"id"`
	UserEmail  string `json:"user_email"`
	ActionType string `json:"action_type"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	EntityName string `json:"entity_name"`
	Diff       string `json:"diff"`
	Metadata   string `json:"metadata"`
	Timestamp  int64  `json:"timestamp"`
}

// AuditLogResponse represents the paginated audit log response
type AuditLogResponse struct {
	Logs       []AuditLogEntry `json:"logs"`
	TotalCount int64           `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// GetAuditLogs godoc
// @Id getAuditLogs
// @Summary Get audit logs
// @Tags Audit
// @Description Returns paginated audit logs with optional filtering
// @Produce json
// @Param user_email query string false "Filter by user email"
// @Param action_type query string false "Filter by action type"
// @Param entity_type query string false "Filter by entity type (spammer, client, system)"
// @Param entity_id query string false "Filter by entity ID"
// @Param start_time query int false "Start time (unix timestamp)"
// @Param end_time query int false "End time (unix timestamp)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 50, max: 1000)"
// @Success 200 {object} AuditLogResponse "Success"
// @Failure 400 {string} string "Invalid parameters"
// @Failure 500 {string} string "Server Error"
// @Router /api/audit-logs [get]
func (ah *APIHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Pagination parameters
	page := 1
	if pageStr := query.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 50
	if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 1000 {
			pageSize = ps
		}
	}

	// Filter parameters
	filters := db.AuditLogFilters{
		UserEmail:  query.Get("user_email"),
		ActionType: query.Get("action_type"),
		EntityType: query.Get("entity_type"),
		EntityID:   query.Get("entity_id"),
		Limit:      pageSize,
		Offset:     (page - 1) * pageSize,
	}

	// Parse time filters
	if startTimeStr := query.Get("start_time"); startTimeStr != "" {
		if startTime, err := strconv.ParseInt(startTimeStr, 10, 64); err == nil {
			filters.StartTime = startTime
		}
	}

	if endTimeStr := query.Get("end_time"); endTimeStr != "" {
		if endTime, err := strconv.ParseInt(endTimeStr, 10, 64); err == nil {
			filters.EndTime = endTime
		}
	}

	// Get total count for pagination
	totalCount, err := ah.daemon.GetDatabase().CountAuditLogs(filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to count audit logs: %v", err), http.StatusInternalServerError)
		return
	}

	// Get audit logs
	dbLogs, err := ah.daemon.GetDatabase().GetAuditLogs(filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get audit logs: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert to API format
	logs := make([]AuditLogEntry, len(dbLogs))
	for i, log := range dbLogs {
		logs[i] = AuditLogEntry{
			ID:         log.ID,
			UserEmail:  log.UserEmail,
			ActionType: log.ActionType,
			EntityType: log.EntityType,
			EntityID:   log.EntityID,
			EntityName: log.EntityName,
			Diff:       log.Diff,
			Metadata:   log.Metadata,
			Timestamp:  log.Timestamp,
		}
	}

	// Calculate total pages
	totalPages := int(totalCount)/pageSize + 1
	if int(totalCount)%pageSize == 0 && totalCount > 0 {
		totalPages--
	}

	response := AuditLogResponse{
		Logs:       logs,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAuditLog godoc
// @Id getAuditLog
// @Summary Get a specific audit log entry
// @Tags Audit
// @Description Returns a specific audit log entry by ID
// @Produce json
// @Param id path int true "Audit log ID"
// @Success 200 {object} AuditLogEntry "Success"
// @Failure 400 {string} string "Invalid audit log ID"
// @Failure 404 {string} string "Audit log not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/audit-logs/{id} [get]
func (ah *APIHandler) GetAuditLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid audit log ID", http.StatusBadRequest)
		return
	}

	log, err := ah.daemon.GetDatabase().GetAuditLogByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "Audit log not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get audit log: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := AuditLogEntry{
		ID:         log.ID,
		UserEmail:  log.UserEmail,
		ActionType: log.ActionType,
		EntityType: log.EntityType,
		EntityID:   log.EntityID,
		EntityName: log.EntityName,
		Diff:       log.Diff,
		Metadata:   log.Metadata,
		Timestamp:  log.Timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AuditLogStatsResponse represents audit log statistics
type AuditLogStatsResponse struct {
	TotalLogs        int64                `json:"total_logs"`
	LogsByActionType map[string]int64     `json:"logs_by_action_type"`
	LogsByEntityType map[string]int64     `json:"logs_by_entity_type"`
	LogsByUser       map[string]int64     `json:"logs_by_user"`
	RecentActivity   []AuditLogEntry      `json:"recent_activity"`
	TimeRange        map[string]time.Time `json:"time_range"`
}

// GetAuditLogStats godoc
// @Id getAuditLogStats
// @Summary Get audit log statistics
// @Tags Audit
// @Description Returns statistics about audit logs including counts by type and recent activity
// @Produce json
// @Param days query int false "Number of days to include in stats (default: 30)"
// @Success 200 {object} AuditLogStatsResponse "Success"
// @Failure 500 {string} string "Server Error"
// @Router /api/audit-logs/stats [get]
func (ah *APIHandler) GetAuditLogStats(w http.ResponseWriter, r *http.Request) {
	// Parse days parameter
	days := 30
	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	startTime := time.Now().AddDate(0, 0, -days).Unix()
	endTime := time.Now().Unix()

	// Get recent activity (last 20 logs)
	recentFilters := db.AuditLogFilters{
		StartTime: startTime,
		EndTime:   endTime,
		Limit:     20,
	}

	recentLogs, err := ah.daemon.GetDatabase().GetAuditLogs(recentFilters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get recent audit logs: %v", err), http.StatusInternalServerError)
		return
	}

	// Get all logs in time range for statistics
	allFilters := db.AuditLogFilters{
		StartTime: startTime,
		EndTime:   endTime,
		Limit:     10000, // Large limit to get all logs for stats
	}

	allLogs, err := ah.daemon.GetDatabase().GetAuditLogs(allFilters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get audit logs for stats: %v", err), http.StatusInternalServerError)
		return
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

	// Convert recent logs to API format
	recentActivity := make([]AuditLogEntry, len(recentLogs))
	for i, log := range recentLogs {
		recentActivity[i] = AuditLogEntry{
			ID:         log.ID,
			UserEmail:  log.UserEmail,
			ActionType: log.ActionType,
			EntityType: log.EntityType,
			EntityID:   log.EntityID,
			EntityName: log.EntityName,
			Diff:       log.Diff,
			Metadata:   log.Metadata,
			Timestamp:  log.Timestamp,
		}
	}

	response := AuditLogStatsResponse{
		TotalLogs:        int64(len(allLogs)),
		LogsByActionType: logsByActionType,
		LogsByEntityType: logsByEntityType,
		LogsByUser:       logsByUser,
		RecentActivity:   recentActivity,
		TimeRange: map[string]time.Time{
			"start": time.Unix(startTime, 0),
			"end":   time.Unix(endTime, 0),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
