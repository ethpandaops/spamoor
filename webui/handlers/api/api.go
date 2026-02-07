package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethpandaops/spamoor/daemon"
	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// CreateSpammerRequest represents the request body for creating a new spammer
type CreateSpammerRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	Scenario         string `json:"scenario"`
	Config           string `json:"config"`
	StartImmediately bool   `json:"startImmediately"`
}

// UpdateSpammerRequest represents the request body for updating a spammer
type UpdateSpammerRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Config      string `json:"config"`
}

// SpammerDetails represents detailed information about a spammer
type SpammerDetails struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scenario    string `json:"scenario"`
	Config      string `json:"config"`
	Status      int    `json:"status"`
}

// VersionResponse represents version information
type VersionResponse struct {
	Version string `json:"version"`
	Release string `json:"release"`
}

// ScenarioEntry represents a single scenario in the API response.
type ScenarioEntry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Plugin      string `json:"plugin,omitempty"` // Plugin name if from a plugin, empty for native scenarios
}

// ScenarioCategory represents a category of scenarios in the API response.
type ScenarioCategory struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Scenarios   []*ScenarioEntry    `json:"scenarios,omitempty"`
	Children    []*ScenarioCategory `json:"children,omitempty"`
}

// SpammerLibraryEntry represents a spammer config from the library
type SpammerLibraryEntry struct {
	File         string   `json:"file"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	SpammerCount int      `json:"spammer_count"`
	Scenarios    []string `json:"scenarios"`
	MinVersion   string   `json:"min_version,omitempty"`
}

// SpammerLibraryIndex represents the index of all available configs
type SpammerLibraryIndex struct {
	Generated time.Time             `json:"generated"`
	Configs   []SpammerLibraryEntry `json:"configs"`
	CachedAt  time.Time             `json:"cached_at"`
	BaseURL   string                `json:"base_url"`
}

// GitHubFile represents a file from GitHub API
type GitHubFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	GitURL      string `json:"git_url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
}

// GraphsDashboardResponse represents the main dashboard graphs data
type GraphsDashboardResponse struct {
	TimeRange  TimeRange            `json:"range"`
	Spammers   []SpammerMetricsData `json:"spammers"`
	Totals     TotalMetricsData     `json:"totals"`
	Others     OthersMetricsData    `json:"others"`
	DataPoints []GraphsDataPoint    `json:"data"`
}

// TimeRange represents the time range of collected metrics
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// SpammerMetricsData represents metrics for a single spammer
type SpammerMetricsData struct {
	ID               uint64 `json:"id"`
	Name             string `json:"name"`
	PendingTxCount   uint64 `json:"pending"`
	ConfirmedTxCount uint64 `json:"confirmed"`
	SubmittedTxCount uint64 `json:"submitted"`
	GasUsedInWindow  uint64 `json:"gasUsed"`
	LastUpdate       string `json:"updated"`
	Status           int    `json:"status"` // Running status from daemon
}

// TotalMetricsData represents aggregated metrics across all spammers
type TotalMetricsData struct {
	PendingTxCount   uint64 `json:"pending"`
	ConfirmedTxCount uint64 `json:"confirmed"`
	SubmittedTxCount uint64 `json:"submitted"`
	GasUsedInWindow  uint64 `json:"gasUsed"`
}

// OthersMetricsData represents metrics for non-spammer transactions
type OthersMetricsData struct {
	GasUsedInWindow uint64 `json:"gasUsed"`
}

// GraphsDataPoint represents a single time-series data point for the graphs
type GraphsDataPoint struct {
	Timestamp        time.Time                    `json:"ts"`
	StartBlockNumber uint64                       `json:"startBlock"`
	EndBlockNumber   uint64                       `json:"endBlock"`
	BlockCount       uint64                       `json:"blocks"`
	TotalGasUsed     uint64                       `json:"totalGas"`
	OthersGasUsed    uint64                       `json:"othersGas"`
	SpammerData      map[string]*SpammerBlockData `json:"spammers"` // spammerID -> detailed data
}

// SpammerBlockData represents a spammer's data within a time period
type SpammerBlockData struct {
	GasUsed          uint64 `json:"gas"`
	ConfirmedTxCount uint64 `json:"confirmed"`
	PendingTxCount   uint64 `json:"pending"`
	SubmittedTxCount uint64 `json:"submitted"`
}

// SpammerTimeSeriesResponse represents time-series data for a specific spammer
type SpammerTimeSeriesResponse struct {
	SpammerID   uint64                   `json:"spammerId"`
	SpammerName string                   `json:"spammerName"`
	TimeRange   TimeRange                `json:"timeRange"`
	DataPoints  []SpammerTimeSeriesPoint `json:"dataPoints"`
}

// SpammerTimeSeriesPoint represents a single time-series point for a spammer
type SpammerTimeSeriesPoint struct {
	Timestamp        time.Time `json:"timestamp"`
	BlockNumber      uint64    `json:"blockNumber"`
	GasUsed          uint64    `json:"gasUsed"`
	ConfirmedTxCount uint64    `json:"confirmedTxCount"`
	PendingTxCount   uint64    `json:"pendingTxCount"`
}

// Library cache structure
var (
	libraryCache      *SpammerLibraryIndex
	libraryCacheMutex sync.RWMutex
	cacheExpiry       = 10 * time.Minute
)

// GetScenarios godoc
// @Id getScenarios
// @Summary Get all scenarios
// @Tags Scenario
// @Description Returns a list of all scenarios organized by category
// @Produce json
// @Success 200 {array} ScenarioCategory "Success"
// @Failure 400 {string} string "Failure"
// @Failure 500 {string} string "Server Error"
// @Router /api/scenarios [get]
func (ah *APIHandler) GetScenarios(w http.ResponseWriter, r *http.Request) {
	categories := scenarios.GetScenarioCategories()
	result := make([]*ScenarioCategory, len(categories))

	for i, category := range categories {
		result[i] = convertCategory(category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// convertCategory recursively converts a scenario.Category to a ScenarioCategory.
// It looks up each scenario in the registry to get the plugin name and sorts scenarios alphabetically.
func convertCategory(category *scenario.Category) *ScenarioCategory {
	scenarioEntries := make([]*ScenarioEntry, len(category.Descriptors))
	for j, descriptor := range category.Descriptors {
		entry := &ScenarioEntry{
			Name:        descriptor.Name,
			Description: descriptor.Description,
		}

		// Look up the scenario in the registry to get the plugin name
		registryEntry := scenarios.GetScenarioEntry(descriptor.Name)
		if registryEntry != nil && registryEntry.Plugin != nil {
			entry.Plugin = registryEntry.Plugin.GetName()
		}

		scenarioEntries[j] = entry
	}

	// Sort scenarios alphabetically by name
	slices.SortFunc(scenarioEntries, func(a, b *ScenarioEntry) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})

	var children []*ScenarioCategory
	if len(category.Children) > 0 {
		children = make([]*ScenarioCategory, len(category.Children))
		for k, child := range category.Children {
			children[k] = convertCategory(child)
		}
	}

	return &ScenarioCategory{
		Name:        category.Name,
		Description: category.Description,
		Scenarios:   scenarioEntries,
		Children:    children,
	}
}

// GetVersion godoc
// @Id getVersion
// @Summary Get spamoor version
// @Tags Version
// @Description Returns the current spamoor version information
// @Produce json
// @Success 200 {object} VersionResponse "Success"
// @Router /api/version [get]
func (ah *APIHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	versionInfo := VersionResponse{
		Version: utils.BuildVersion,
		Release: utils.BuildRelease,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versionInfo)
}

// GetScenarioConfig godoc
// @Id getScenarioConfig
// @Summary Get scenario configuration
// @Tags Scenario
// @Description Returns the default configuration for a specific scenario
// @Produce text/plain
// @Param name path string true "Scenario name"
// @Success 200 {string} string "YAML configuration"
// @Failure 404 {string} string "Scenario not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/scenarios/{name}/config [get]
func (ah *APIHandler) GetScenarioConfig(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	scenarioName := vars["name"]

	scenario := scenarios.GetScenario(scenarioName)
	if scenario == nil {
		http.Error(w, "Scenario not found", http.StatusNotFound)
		return
	}

	configYaml, err := yaml.Marshal(scenario.DefaultOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write header comment first
	w.Header().Set("Content-Type", "application/x-yaml")
	fmt.Fprintf(w, "# wallet settings\n")
	fmt.Fprintf(w, "seed: %s-%v # seed for the wallet\n", scenarioName, rand.Intn(1000000))
	fmt.Fprintf(w, "refill_amount: 5000000000000000000 # refill 5 ETH when\n")
	fmt.Fprintf(w, "refill_balance: 1000000000000000000 # balance drops below 1 ETH\n")
	fmt.Fprintf(w, "refill_interval: 600 # check every 10 minutes\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "# scenario: %s\n", scenarioName)
	fmt.Fprintf(w, "%s\n", string(configYaml))
}

type SpammerListEntry struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Scenario    string `json:"scenario"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"` // RFC3339Nano formatted timestamp
}

// GetSpammerList godoc
// @Id getSpammerList
// @Summary Get all spammers
// @Tags Spammer
// @Description Returns a list of all configured spammers
// @Produce json
// @Success 200 {array} SpammerListEntry "Success"
// @Router /api/spammers [get]
func (ah *APIHandler) GetSpammerList(w http.ResponseWriter, r *http.Request) {
	spammers := ah.daemon.GetAllSpammers()
	response := make([]SpammerListEntry, len(spammers))

	for i, s := range spammers {
		response[i] = SpammerListEntry{
			ID:          s.GetID(),
			Name:        s.GetName(),
			Description: s.GetDescription(),
			Scenario:    s.GetScenario(),
			Status:      s.GetStatus(),
			CreatedAt:   time.Unix(s.GetCreatedAt(), 0).Format(time.RFC3339Nano),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateSpammer godoc
// @Id createSpammer
// @Summary Create a new spammer
// @Tags Spammer
// @Description Creates a new spammer with the given configuration
// @Accept json
// @Produce json
// @Param request body CreateSpammerRequest true "Spammer configuration"
// @Success 200 {object} int64 "Spammer ID"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer [post]
func (ah *APIHandler) CreateSpammer(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	var req struct {
		Name             string `json:"name"`
		Description      string `json:"description"`
		Scenario         string `json:"scenario"`
		Config           string `json:"config"`
		StartImmediately bool   `json:"startImmediately"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	spammer, err := ah.daemon.NewSpammer(req.Scenario, req.Config, req.Name, req.Description, req.StartImmediately, userEmail, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spammer.GetID())
}

// StartSpammer godoc
// @Id startSpammer
// @Summary Start a spammer
// @Tags Spammer
// @Description Starts a specific spammer
// @Param id path int true "Spammer ID"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid spammer ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id}/start [post]
func (ah *APIHandler) StartSpammer(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	err = ah.daemon.StartSpammer(id, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateSpammer godoc
// @Id updateSpammer
// @Summary Update a spammer
// @Tags Spammer
// @Description Updates an existing spammer's configuration
// @Accept json
// @Param id path int true "Spammer ID"
// @Param request body UpdateSpammerRequest true "Updated configuration"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id} [put]
func (ah *APIHandler) UpdateSpammer(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Config      string `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	spammer := ah.daemon.GetSpammer(id)
	if spammer == nil {
		http.Error(w, "Spammer not found", http.StatusNotFound)
		return
	}

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	err = ah.daemon.UpdateSpammer(id, req.Name, req.Description, req.Config, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PauseSpammer godoc
// @Id pauseSpammer
// @Summary Pause a spammer
// @Tags Spammer
// @Description Pauses a running spammer
// @Param id path int true "Spammer ID"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid spammer ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id}/pause [post]
func (ah *APIHandler) PauseSpammer(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	err = ah.daemon.PauseSpammer(id, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteSpammer godoc
// @Id deleteSpammer
// @Summary Delete a spammer
// @Tags Spammer
// @Description Deletes a spammer and stops it if running
// @Param id path int true "Spammer ID"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid spammer ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id} [delete]
func (ah *APIHandler) DeleteSpammer(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	err = ah.daemon.DeleteSpammer(id, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ReclaimFunds godoc
// @Id reclaimFunds
// @Summary Reclaim funds from a spammer
// @Tags Spammer
// @Description Reclaims funds from a spammer's wallet pool back to the root wallet
// @Param id path int true "Spammer ID"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid spammer ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id}/reclaim [post]
func (ah *APIHandler) ReclaimFunds(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	err = ah.daemon.ReclaimSpammer(id, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetSpammerDetails godoc
// @Id getSpammerDetails
// @Summary Get spammer details
// @Tags Spammer
// @Description Returns detailed information about a specific spammer
// @Produce json
// @Param id path int true "Spammer ID"
// @Success 200 {object} SpammerDetails "Success"
// @Failure 400 {string} string "Invalid spammer ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Router /api/spammer/{id} [get]
func (ah *APIHandler) GetSpammerDetails(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	spammer := ah.daemon.GetSpammer(id)
	if spammer == nil {
		http.Error(w, "Spammer not found", http.StatusNotFound)
		return
	}

	response := struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Scenario    string `json:"scenario"`
		Config      string `json:"config"`
		Status      int    `json:"status"`
	}{
		ID:          spammer.GetID(),
		Name:        spammer.GetName(),
		Description: spammer.GetDescription(),
		Scenario:    spammer.GetScenario(),
		Config:      spammer.GetConfig(),
		Status:      spammer.GetStatus(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type LogEntry struct {
	Time    string            `json:"time"`
	Level   string            `json:"level"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

// GetSpammerLogs godoc
// @Id getSpammerLogs
// @Summary Get spammer logs
// @Tags Spammer
// @Description Returns the most recent logs for a specific spammer
// @Produce json
// @Param id path int true "Spammer ID"
// @Success 200 {array} LogEntry "Success"
// @Failure 400 {string} string "Invalid spammer ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Router /api/spammer/{id}/logs [get]
func (ah *APIHandler) GetSpammerLogs(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	spammer := ah.daemon.GetSpammer(id)
	if spammer == nil {
		http.Error(w, "Spammer not found", http.StatusNotFound)
		return
	}

	// Get last 1000 log entries
	logScope := spammer.GetLogScope()
	entries := logScope.GetLogEntries(time.Time{}, 1000)

	logs := make([]LogEntry, len(entries))
	for i, entry := range entries {
		fields := make(map[string]string)
		for k, v := range entry.Data {
			fields[k] = fmt.Sprint(v)
		}

		logs[i] = LogEntry{
			Time:    entry.Time.Format(time.RFC3339Nano),
			Level:   entry.Level.String(),
			Message: entry.Message,
			Fields:  fields,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// StreamSpammerLogs godoc
// @Id streamSpammerLogs
// @Summary Stream spammer logs
// @Tags Spammer
// @Description Streams logs for a specific spammer using Server-Sent Events
// @Produce text/event-stream
// @Param id path int true "Spammer ID"
// @Param since query string false "Timestamp to start from (RFC3339Nano)"
// @Success 200 {string} string "SSE stream of log entries"
// @Failure 400 {string} string "Invalid spammer ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Streaming unsupported"
// @Router /api/spammer/{id}/logs/stream [get]
func (ah *APIHandler) StreamSpammerLogs(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	spammer := ah.daemon.GetSpammer(id)
	if spammer == nil {
		http.Error(w, "Spammer not found", http.StatusNotFound)
		return
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")

	// Get initial timestamp
	var lastTime time.Time
	if timeStr := r.URL.Query().Get("since"); timeStr != "" {
		if t, err := time.Parse(time.RFC3339Nano, timeStr); err == nil {
			lastTime = t
		} else {
			logrus.Warnf("Failed to parse timestamp %s: %v", timeStr, err)
		}
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Stream logs until client disconnects
	for {
		select {
		case <-r.Context().Done():
			return
		default:
			entries := spammer.GetLogScope().GetLogEntries(lastTime, 1000)
			if len(entries) > 0 {
				// Convert to JSON format
				logs := make([]LogEntry, len(entries))
				for i, entry := range entries {
					fields := make(map[string]string)
					for k, v := range entry.Data {
						fields[k] = fmt.Sprint(v)
					}

					logs[i] = LogEntry{
						Time:    entry.Time.Format(time.RFC3339Nano),
						Level:   entry.Level.String(),
						Message: entry.Message,
						Fields:  fields,
					}
					lastTime = entry.Time
				}

				// Send as SSE event
				data, _ := json.Marshal(logs)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// PendingTransactionEntry represents a pending transaction in the API response
type PendingTransactionEntry struct {
	Hash             string `json:"hash"`
	WalletAddress    string `json:"wallet_address"`
	WalletName       string `json:"wallet_name"`   // Human-readable wallet name
	ScenarioName     string `json:"scenario_name"` // Scenario name (if applicable)
	SpammerID        int64  `json:"spammer_id"`    // Spammer ID (if applicable)
	Nonce            uint64 `json:"nonce"`
	Value            string `json:"value"`              // In wei as string
	ValueFormatted   string `json:"value_formatted"`    // Formatted for display
	Fee              string `json:"fee"`                // In wei as string
	FeeFormatted     string `json:"fee_formatted"`      // Formatted for display
	BaseFee          string `json:"base_fee"`           // In wei as string
	BaseFeeFormatted string `json:"base_fee_formatted"` // Formatted for display
	SubmittedAt      string `json:"submitted_at"`       // RFC3339 timestamp
	RebroadcastCount uint64 `json:"rebroadcast_count"`
	LastRebroadcast  string `json:"last_rebroadcast"` // RFC3339 timestamp
	RLPEncoded       string `json:"rlp_encoded"`      // RLP encoded transaction as hex string
}

// GetPendingTransactions godoc
// @Id getPendingTransactions
// @Summary Get pending transactions (global or wallet-specific)
// @Tags Wallet
// @Description Returns a list of pending transactions, optionally filtered by wallet address
// @Produce json
// @Param wallet query string false "Wallet address to filter by (optional)"
// @Success 200 {array} PendingTransactionEntry "Success"
// @Failure 400 {string} string "Invalid wallet address"
// @Failure 500 {string} string "Server Error"
// @Router /api/pending-transactions [get]
func (ah *APIHandler) GetPendingTransactions(w http.ResponseWriter, r *http.Request) {
	// Check if authenticated (don't return error, just flag for filtering)
	isAuthenticated := ah.isAuthenticated(r)

	walletFilter := r.URL.Query().Get("wallet")

	allPendingTxs := make([]PendingTransactionEntry, 0)

	// Helper function to format wei values with appropriate units
	formatWeiSmart := func(wei string) string {
		if wei == "" || wei == "0" {
			return "0 wei"
		}

		// Convert to float64 for comparison
		weiFloat, _ := strconv.ParseFloat(wei, 64)

		// Thresholds
		const (
			gweiThreshold = 1000       // 1000 wei = switch to gwei
			ethThreshold  = 1000 * 1e9 // 1000 gwei = switch to eth
		)

		if weiFloat < gweiThreshold {
			// Show in wei for very small amounts
			return fmt.Sprintf("%s wei", wei)
		} else if weiFloat < ethThreshold {
			// Show in gwei for medium amounts
			gwei := weiFloat / 1e9
			return fmt.Sprintf("%.2f gwei", gwei)
		} else {
			// Show in ETH for large amounts
			eth := weiFloat / 1e18
			return fmt.Sprintf("%.6f ETH", eth)
		}
	}

	// Helper function to format wei values to Gwei (for gas fees)
	formatWeiToGwei := func(wei string) string {
		if wei == "" || wei == "0" {
			return "0 gwei"
		}
		// Convert wei to Gwei for display
		weiFloat, _ := strconv.ParseFloat(wei, 64)
		gwei := weiFloat / 1e9
		return fmt.Sprintf("%.2f gwei", gwei)
	}

	// Helper function to process pending transactions from a wallet
	processPendingTxs := func(wallet *spamoor.Wallet, walletName, scenarioName string, spammerID int64) {
		walletAddr := wallet.GetAddress().Hex()

		// Skip if wallet filter is specified and doesn't match
		if walletFilter != "" && walletAddr != walletFilter {
			return
		}

		// Get pending transactions for this wallet
		pendingTxs := wallet.GetPendingTxs()

		for _, pendingTx := range pendingTxs {
			tx := pendingTx.Tx

			// Calculate fee (gas price * gas limit)
			gasPrice := tx.GasPrice()
			gasLimit := tx.Gas()
			fee := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasLimit))

			// Get base fee (assume it's the gas price for now)
			baseFee := gasPrice

			// Format timestamps
			submittedAt := pendingTx.Submitted.Format(time.RFC3339)
			lastRebroadcast := pendingTx.LastRebroadcast.Format(time.RFC3339)

			// Only include RLP encoded transaction for authenticated users
			rlpHex := ""
			if isAuthenticated {
				rlpBytes, err := tx.MarshalBinary()
				if err == nil {
					rlpHex = "0x" + hex.EncodeToString(rlpBytes)
				}
			}

			entry := PendingTransactionEntry{
				Hash:             tx.Hash().Hex(),
				WalletAddress:    walletAddr,
				WalletName:       walletName,
				ScenarioName:     scenarioName,
				SpammerID:        spammerID,
				Nonce:            tx.Nonce(),
				Value:            tx.Value().String(),
				ValueFormatted:   formatWeiSmart(tx.Value().String()),
				Fee:              fee.String(),
				FeeFormatted:     formatWeiSmart(fee.String()),
				BaseFee:          baseFee.String(),
				BaseFeeFormatted: formatWeiToGwei(baseFee.String()),
				SubmittedAt:      submittedAt,
				RebroadcastCount: pendingTx.RebroadcastCount,
				LastRebroadcast:  lastRebroadcast,
				RLPEncoded:       rlpHex,
			}

			allPendingTxs = append(allPendingTxs, entry)
		}
	}

	// Process root wallet's pending transactions
	rootWallet := ah.daemon.GetRootWallet()
	if rootWallet != nil && rootWallet.GetWallet() != nil {
		processPendingTxs(rootWallet.GetWallet(), "Root Wallet", "", -1)
	}

	// Process each spammer's wallets
	spammers := ah.daemon.GetAllSpammers()
	for _, spammer := range spammers {
		walletPool := spammer.GetWalletPool()
		if walletPool == nil {
			continue
		}

		// Get spammer information
		spammerID := spammer.GetID()
		spammerName := spammer.GetName()
		scenarioName := spammer.GetScenario()

		// Get all wallets from the pool
		wallets := walletPool.GetAllWallets()

		for _, wallet := range wallets {
			// Create a descriptive wallet name
			walletName := walletPool.GetWalletName(wallet.GetAddress())
			if spammerName != "" {
				walletName = fmt.Sprintf("%s (%s)", walletName, spammerName)
			}

			processPendingTxs(wallet, walletName, scenarioName, spammerID)
		}
	}

	// Sort by submission time (newest first)
	slices.SortFunc(allPendingTxs, func(a, b PendingTransactionEntry) int {
		timeA, _ := time.Parse(time.RFC3339, a.SubmittedAt)
		timeB, _ := time.Parse(time.RFC3339, b.SubmittedAt)
		return timeB.Compare(timeA)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allPendingTxs)
}

// ClientEntry represents a client in the API response
type ClientEntry struct {
	Index        int      `json:"index"`
	Name         string   `json:"name"`
	Group        string   `json:"group"`  // First group for backward compatibility
	Groups       []string `json:"groups"` // All groups
	Type         string   `json:"type"`   // Client type (client, builder)
	Version      string   `json:"version"`
	BlockHeight  uint64   `json:"block_height"`
	IsReady      bool     `json:"ready"`
	RpcHost      string   `json:"rpc_host"`
	Enabled      bool     `json:"enabled"`
	NameOverride string   `json:"name_override,omitempty"`
}

// UpdateClientGroupRequest represents the request body for updating a client group
type UpdateClientGroupRequest struct {
	Group  string   `json:"group,omitempty"`  // Single group for backward compatibility
	Groups []string `json:"groups,omitempty"` // Multiple groups
}

// UpdateClientEnabledRequest represents the request body for updating a client's enabled state
type UpdateClientEnabledRequest struct {
	Enabled bool `json:"enabled"`
}

// UpdateClientNameRequest represents the request body for updating a client's name override
type UpdateClientNameRequest struct {
	NameOverride string `json:"name_override"`
}

// UpdateClientTypeRequest represents the request body for updating a client's type
type UpdateClientTypeRequest struct {
	ClientType string `json:"client_type"`
}

// GetClients godoc
// @Id getClients
// @Summary Get all clients
// @Tags Client
// @Description Returns a list of all clients with their details
// @Produce json
// @Success 200 {array} ClientEntry "Success"
// @Router /api/clients [get]
func (ah *APIHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	// Check if authenticated (don't return error, just flag for filtering)
	isAuthenticated := ah.isAuthenticated(r)

	ctx := r.Context()
	goodClients := ah.daemon.GetClientPool().GetAllGoodClients()
	allClients := ah.daemon.GetClientPool().GetAllClients()

	response := make([]ClientEntry, len(allClients))

	for i, client := range allClients {
		blockHeight, _ := client.GetLastBlockHeight()

		version, err := client.GetClientVersion(ctx)
		if err != nil {
			version = "Unknown"
		}

		rpcHost := ""
		if isAuthenticated {
			rpcHost = client.GetRPCHost()
		}

		response[i] = ClientEntry{
			Index:        i,
			Name:         client.GetName(),
			Group:        client.GetClientGroup(),
			Groups:       client.GetClientGroups(),
			Type:         client.GetClientType().String(),
			Version:      version,
			BlockHeight:  blockHeight,
			IsReady:      slices.Contains(goodClients, client),
			RpcHost:      rpcHost,
			Enabled:      client.IsEnabled(),
			NameOverride: client.GetNameOverride(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateClientGroup godoc
// @Id updateClientGroup
// @Summary Update client group
// @Tags Client
// @Description Updates the group(s) for a specific client. Supports both single group (backward compatibility) and multiple groups.
// @Accept json
// @Param index path int true "Client index"
// @Param request body UpdateClientGroupRequest true "New group name(s)"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid client index"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Client not found"
// @Router /api/client/{index}/group [put]
func (ah *APIHandler) UpdateClientGroup(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	index, err := strconv.Atoi(vars["index"])
	if err != nil {
		http.Error(w, "Invalid client index", http.StatusBadRequest)
		return
	}

	var req UpdateClientGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allClients := ah.daemon.GetClientPool().GetAllClients()
	if index < 0 || index >= len(allClients) {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	client := allClients[index]

	// Get existing config for audit logging
	existingConfig, err := ah.daemon.GetClientConfig(client.GetRPCHost())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Handle both single group (backward compatibility) and multiple groups
	var groups []string
	if len(req.Groups) > 0 {
		groups = req.Groups
		client.SetClientGroups(req.Groups)
	} else if req.Group != "" {
		groups = []string{req.Group}
		client.SetClientGroups([]string{req.Group})
	} else {
		http.Error(w, "Either 'group' or 'groups' must be provided", http.StatusBadRequest)
		return
	}

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	// Update client config with new groups as tags
	tagsStr := strings.Join(groups, ",")
	err = ah.daemon.UpdateClientConfig(client.GetRPCHost(), existingConfig.Name, tagsStr, existingConfig.ClientType, existingConfig.Enabled, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateClientEnabled godoc
// @Id updateClientEnabled
// @Summary Update client enabled state
// @Tags Client
// @Description Updates the enabled state for a specific client
// @Accept json
// @Param index path int true "Client index"
// @Param request body UpdateClientEnabledRequest true "New enabled state"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid client index"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Client not found"
// @Router /api/client/{index}/enabled [put]
func (ah *APIHandler) UpdateClientEnabled(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	index, err := strconv.Atoi(vars["index"])
	if err != nil {
		http.Error(w, "Invalid client index", http.StatusBadRequest)
		return
	}

	var req UpdateClientEnabledRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allClients := ah.daemon.GetClientPool().GetAllClients()
	if index < 0 || index >= len(allClients) {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	client := allClients[index]

	// Get existing config for audit logging
	existingConfig, err := ah.daemon.GetClientConfig(client.GetRPCHost())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client.SetEnabled(req.Enabled)

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	// Update client config with new enabled state
	err = ah.daemon.UpdateClientConfig(client.GetRPCHost(), existingConfig.Name, existingConfig.Tags, existingConfig.ClientType, req.Enabled, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateClientName godoc
// @Id updateClientName
// @Summary Update client name override
// @Tags Client
// @Description Updates the name override for a specific client
// @Accept json
// @Param index path int true "Client index"
// @Param request body UpdateClientNameRequest true "New name override"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid client index"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Client not found"
// @Router /api/client/{index}/name [put]
func (ah *APIHandler) UpdateClientName(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	index, err := strconv.Atoi(vars["index"])
	if err != nil {
		http.Error(w, "Invalid client index", http.StatusBadRequest)
		return
	}

	var req UpdateClientNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allClients := ah.daemon.GetClientPool().GetAllClients()
	if index < 0 || index >= len(allClients) {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	client := allClients[index]

	// Get existing config for audit logging
	existingConfig, err := ah.daemon.GetClientConfig(client.GetRPCHost())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client.SetNameOverride(req.NameOverride)

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	// Update client config with new name override
	err = ah.daemon.UpdateClientConfig(client.GetRPCHost(), req.NameOverride, existingConfig.Tags, existingConfig.ClientType, existingConfig.Enabled, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateClientType godoc
// @Id updateClientType
// @Summary Update client type
// @Tags Client
// @Description Updates the type for a specific client (e.g., 'client' or 'builder')
// @Accept json
// @Param index path int true "Client index"
// @Param request body UpdateClientTypeRequest true "New client type"
// @Success 200 "Success"
// @Failure 400 {string} string "Invalid client index or type"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Client not found"
// @Router /api/client/{index}/type [put]
func (ah *APIHandler) UpdateClientType(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	index, err := strconv.Atoi(vars["index"])
	if err != nil {
		http.Error(w, "Invalid client index", http.StatusBadRequest)
		return
	}

	var req UpdateClientTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate client type
	if req.ClientType != "" && req.ClientType != "client" && req.ClientType != "builder" {
		http.Error(w, "Invalid client type. Must be 'client' or 'builder'", http.StatusBadRequest)
		return
	}

	allClients := ah.daemon.GetClientPool().GetAllClients()
	if index < 0 || index >= len(allClients) {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	client := allClients[index]

	// Get existing config for audit logging
	existingConfig, err := ah.daemon.GetClientConfig(client.GetRPCHost())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client.SetClientTypeOverride(req.ClientType)

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	// Update client config with new type
	err = ah.daemon.UpdateClientConfig(client.GetRPCHost(), existingConfig.Name, existingConfig.Tags, req.ClientType, existingConfig.Enabled, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ExportSpammersRequest represents the request body for exporting spammers
type ExportSpammersRequest struct {
	SpammerIDs []int64 `json:"spammer_ids,omitempty"` // If empty, exports all spammers
}

// ImportSpammersRequest represents the request body for importing spammers
type ImportSpammersRequest struct {
	Input string `json:"input"` // Can be YAML data or a URL
}

// SendTransactionRequest represents the request body for sending a transaction from root wallet
type SendTransactionRequest struct {
	To       string `json:"to"`                 // Target address (required)
	Value    string `json:"value"`              // Amount in specified unit (required)
	Unit     string `json:"unit"`               // Unit: "eth", "gwei", or "wei" (required)
	Data     string `json:"data,omitempty"`     // Hex encoded calldata (optional, default: "0x")
	GasLimit uint64 `json:"gasLimit,omitempty"` // Gas limit (optional, default: 21000 for simple transfers)
	MaxFee   string `json:"maxFee,omitempty"`   // Max fee per gas in gwei (optional)
	MaxTip   string `json:"maxTip,omitempty"`   // Max priority fee per gas in gwei (optional)
}

// SendTransactionResponse represents the response after sending a transaction
type SendTransactionResponse struct {
	TxHash string `json:"txHash"`
	Nonce  uint64 `json:"nonce"`
}

// ExportSpammers godoc
// @Id exportSpammers
// @Summary Export spammers to YAML
// @Tags Spammer
// @Description Exports specified spammers or all spammers to YAML format
// @Accept json
// @Produce text/plain
// @Param request body ExportSpammersRequest false "Spammer IDs to export (optional)"
// @Success 200 {string} string "YAML configuration"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammers/export [post]
func (ah *APIHandler) ExportSpammers(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	var req ExportSpammersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	yamlData, err := ah.daemon.ExportSpammers(req.SpammerIDs...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Content-Disposition", "attachment; filename=spammers-export.yaml")
	w.Write([]byte(yamlData))
}

// ImportSpammers godoc
// @Id importSpammers
// @Summary Import spammers from YAML data or URL
// @Tags Spammer
// @Description Imports spammers from YAML data or URL with validation and deduplication
// @Accept json
// @Produce json
// @Param request body ImportSpammersRequest true "Import configuration"
// @Success 200 {object} daemon.ImportResult "Success"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammers/import [post]
func (ah *APIHandler) ImportSpammers(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	var req ImportSpammersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Input == "" {
		http.Error(w, "input is required", http.StatusBadRequest)
		return
	}

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	result, err := ah.daemon.ImportSpammers(req.Input, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// fetchFileContent fetches file content from a URL
func fetchFileContent(downloadURL string) (string, error) {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch file, status: %d", resp.StatusCode)
	}

	var content []byte
	if resp.ContentLength > 0 {
		content = make([]byte, resp.ContentLength)
		_, err = resp.Body.Read(content)
	} else {
		// Read all content if ContentLength is unknown
		content, err = io.ReadAll(resp.Body)
	}

	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	return string(content), nil
}

// refreshLibraryCache fetches and caches the spammer library index
func refreshLibraryCache() error {
	// Fetch the _index.yaml file directly
	indexURL := "https://raw.githubusercontent.com/ethpandaops/spamoor/master/spammer-configs/_index.yaml"
	content, err := fetchFileContent(indexURL)
	if err != nil {
		return fmt.Errorf("failed to fetch index file: %w", err)
	}

	// Parse the index YAML
	var index SpammerLibraryIndex
	if err := yaml.Unmarshal([]byte(content), &index); err != nil {
		return fmt.Errorf("failed to parse index file: %w", err)
	}

	// Set the cached at time and base URL
	index.CachedAt = time.Now()
	index.BaseURL = "https://raw.githubusercontent.com/ethpandaops/spamoor/master/spammer-configs/"

	libraryCacheMutex.Lock()
	libraryCache = &index
	libraryCacheMutex.Unlock()

	return nil
}

// getLibraryIndex returns the cached library index, refreshing if needed
func getLibraryIndex() (*SpammerLibraryIndex, error) {
	libraryCacheMutex.RLock()
	if libraryCache != nil && time.Since(libraryCache.CachedAt) < cacheExpiry {
		defer libraryCacheMutex.RUnlock()
		return libraryCache, nil
	}
	libraryCacheMutex.RUnlock()

	// Cache is expired or doesn't exist, refresh it
	if err := refreshLibraryCache(); err != nil {
		// If refresh fails, return cached data if available
		libraryCacheMutex.RLock()
		defer libraryCacheMutex.RUnlock()
		if libraryCache != nil {
			logrus.Warnf("Failed to refresh library cache, using cached data: %v", err)
			return libraryCache, nil
		}
		return nil, err
	}

	libraryCacheMutex.RLock()
	defer libraryCacheMutex.RUnlock()
	return libraryCache, nil
}

// GetSpammerLibraryIndex godoc
// @Id getSpammerLibraryIndex
// @Summary Get spammer library index
// @Tags SpammerLibrary
// @Description Returns the index of available spammer configurations from GitHub
// @Produce json
// @Success 200 {object} SpammerLibraryIndex "Success"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer-library/index [get]
func (ah *APIHandler) GetSpammerLibraryIndex(w http.ResponseWriter, r *http.Request) {
	index, err := getLibraryIndex()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch library index: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(index)
}

// @Summary Get graphs dashboard data
// @Tags Graphs
// @Description Returns comprehensive graphs data for the dashboard including all spammers, totals, and time-series data
// @Produce json
// @Success 200 {object} GraphsDashboardResponse "Success"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server Error"
// @Router /api/graphs/dashboard [get]
func (ah *APIHandler) GetGraphsDashboard(w http.ResponseWriter, r *http.Request) {
	shortWindow := ah.daemon.GetShortWindowMetrics()
	if shortWindow == nil {
		http.Error(w, "Transaction metrics collection is disabled", http.StatusServiceUnavailable)
		return
	}

	// Get time range from short window (30 minutes)
	startTime, endTime := shortWindow.GetTimeRange()

	// Get current spammer snapshots
	spammerSnapshots := shortWindow.GetSpammerSnapshots()
	spammers := make([]SpammerMetricsData, 0, len(spammerSnapshots))

	totalPending := uint64(0)
	totalConfirmed := uint64(0)
	totalSubmitted := uint64(0)

	// Calculate gas used in window from data points
	dataPoints := shortWindow.GetDataPoints()
	spammerGasInWindow := make(map[uint64]uint64)

	for _, point := range dataPoints {
		for spammerID, spammerData := range point.SpammerGasData {
			spammerGasInWindow[spammerID] += spammerData.GasUsed
		}
	}

	totalGasUsed := uint64(0)
	for spammerID, snapshot := range spammerSnapshots {
		spammerName := ah.daemon.GetSpammerName(spammerID)
		gasInWindow := spammerGasInWindow[spammerID]

		// Get spammer status
		status := 0 // Default to stopped
		if spammer := ah.daemon.GetSpammer(int64(spammerID)); spammer != nil {
			status = spammer.GetStatus()
		}

		spammers = append(spammers, SpammerMetricsData{
			ID:               spammerID,
			Name:             spammerName,
			PendingTxCount:   snapshot.PendingTxCount,
			ConfirmedTxCount: snapshot.TotalConfirmedTx,
			SubmittedTxCount: snapshot.TotalSubmittedTx,
			GasUsedInWindow:  gasInWindow,
			LastUpdate:       snapshot.LastUpdate.Format(time.RFC3339),
			Status:           status,
		})

		totalPending += snapshot.PendingTxCount
		totalConfirmed += snapshot.TotalConfirmedTx
		totalSubmitted += snapshot.TotalSubmittedTx
		totalGasUsed += gasInWindow
	}

	// Convert data points to API format
	chartDataPoints := make([]GraphsDataPoint, len(dataPoints))
	totalOthersGas := uint64(0)

	for i, point := range dataPoints {
		// Convert spammer data
		spammerData := make(map[string]*SpammerBlockData)
		for spammerID, data := range point.SpammerGasData {
			spammerData[fmt.Sprintf("%d", spammerID)] = &SpammerBlockData{
				GasUsed:          data.GasUsed,
				ConfirmedTxCount: data.ConfirmedTxCount,
				PendingTxCount:   data.PendingTxCount,
				SubmittedTxCount: data.SubmittedTxCount,
			}
		}

		chartDataPoints[i] = GraphsDataPoint{
			Timestamp:        point.Timestamp,
			StartBlockNumber: point.StartBlockNumber,
			EndBlockNumber:   point.EndBlockNumber,
			BlockCount:       point.BlockCount,
			TotalGasUsed:     point.TotalGasUsed,
			OthersGasUsed:    point.OthersGasUsed,
			SpammerData:      spammerData,
		}

		totalOthersGas += point.OthersGasUsed
	}

	response := GraphsDashboardResponse{
		TimeRange: TimeRange{
			Start: startTime,
			End:   endTime,
		},
		Spammers: spammers,
		Totals: TotalMetricsData{
			PendingTxCount:   totalPending,
			ConfirmedTxCount: totalConfirmed,
			SubmittedTxCount: totalSubmitted,
			GasUsedInWindow:  totalGasUsed,
		},
		Others: OthersMetricsData{
			GasUsedInWindow: totalOthersGas,
		},
		DataPoints: chartDataPoints,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Get time-series data for a specific spammer
// @Tags Graphs
// @Description Returns detailed time-series graphs data for a specific spammer
// @Produce json
// @Param id path int true "Spammer ID"
// @Success 200 {object} SpammerTimeSeriesResponse "Success"
// @Failure 400 {string} string "Invalid spammer ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/graphs/spammer/{id}/timeseries [get]
func (ah *APIHandler) GetSpammerTimeSeries(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	vars := mux.Vars(r)
	spammerIDStr := vars["id"]

	spammerID, err := strconv.ParseUint(spammerIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	metricsData := ah.daemon.GetShortWindowMetrics()
	if metricsData == nil {
		http.Error(w, "Transaction metrics collection is disabled", http.StatusServiceUnavailable)
		return
	}

	// Check if spammer exists
	spammerSnapshots := metricsData.GetSpammerSnapshots()
	if _, exists := spammerSnapshots[spammerID]; !exists {
		http.Error(w, "Spammer not found", http.StatusNotFound)
		return
	}

	spammerName := ah.daemon.GetSpammerName(spammerID)
	startTime, endTime := metricsData.GetTimeRange()

	// Build time series data from data points
	dataPoints := metricsData.GetDataPoints()
	spammerDataPoints := make([]SpammerTimeSeriesPoint, 0, len(dataPoints))

	for _, point := range dataPoints {
		if spammerData, exists := point.SpammerGasData[spammerID]; exists {
			spammerDataPoints = append(spammerDataPoints, SpammerTimeSeriesPoint{
				Timestamp:        point.Timestamp,
				GasUsed:          spammerData.GasUsed,
				ConfirmedTxCount: spammerData.ConfirmedTxCount,
				PendingTxCount:   spammerData.PendingTxCount,
			})
		}
	}

	response := SpammerTimeSeriesResponse{
		SpammerID:   spammerID,
		SpammerName: spammerName,
		TimeRange: TimeRange{
			Start: startTime,
			End:   endTime,
		},
		DataPoints: spammerDataPoints,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Stream real-time graphs updates
// @Tags Graphs
// @Description Provides real-time graphs updates via Server-Sent Events (SSE)
// @Produce text/event-stream
// @Success 200 {string} string "SSE stream"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/graphs/stream [get]
func (ah *APIHandler) StreamGraphs(w http.ResponseWriter, r *http.Request) {
	// Check if tx metrics are enabled
	metricsCollector := ah.daemon.GetMetricsCollector()
	if metricsCollector == nil {
		http.Error(w, "Transaction metrics collection is disabled", http.StatusServiceUnavailable)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Create per-connection SSE state
	connectionState := &SSEState{
		lastSpammerHashes: make(map[uint64]string),
	}

	// Send initial data
	shortWindow := ah.daemon.GetShortWindowMetrics()
	if shortWindow != nil {
		ah.sendCurrentSpammerDataWithState(w, flusher, shortWindow, connectionState)
	}

	subscriptionID, updateChan := metricsCollector.Subscribe()
	defer metricsCollector.Unsubscribe(subscriptionID)

	// Set up a fallback ticker for updates every 30 seconds (reduced from 5)
	fallbackTicker := time.NewTicker(30 * time.Second)
	defer fallbackTicker.Stop()

	// Context for cleanup
	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updateChan:
			// Send real-time update based on metrics update
			ah.sendMetricsUpdateWithState(w, flusher, update, connectionState)
		case <-fallbackTicker.C:
			// Fallback update in case we miss real-time updates
			shortWindow := ah.daemon.GetShortWindowMetrics()
			if shortWindow == nil {
				continue
			}

			ah.sendCurrentSpammerDataWithState(w, flusher, shortWindow, connectionState)
		}
	}
}

// SSE state tracking
type SSEState struct {
	lastDataPointCount int
	lastSpammerHashes  map[uint64]string
}

// Removed global sseState - now using per-connection state

// sendCurrentSpammerDataWithState sends current spammer data via SSE only if there are changes
func (ah *APIHandler) sendCurrentSpammerDataWithState(w http.ResponseWriter, flusher http.Flusher, shortWindow *daemon.MultiGranularityMetrics, state *SSEState) {
	spammerSnapshots := shortWindow.GetSpammerSnapshots()
	dataPoints := shortWindow.GetDataPoints()

	// Check if there are new data points
	hasNewDataPoints := len(dataPoints) > state.lastDataPointCount

	// Calculate gas usage from data points
	spammerGasInWindow := make(map[uint64]uint64)
	for _, point := range dataPoints {
		for spammerID, spammerData := range point.SpammerGasData {
			spammerGasInWindow[spammerID] += spammerData.GasUsed
		}
	}

	// Check if spammer data has changed
	hasSpammerChanges := false
	newHashes := make(map[uint64]string)

	for spammerID, snapshot := range spammerSnapshots {
		// Create a simple hash of the spammer state
		gasInWindow := spammerGasInWindow[spammerID]

		// Get spammer status for hash
		status := 0
		if spammer := ah.daemon.GetSpammer(int64(spammerID)); spammer != nil {
			status = spammer.GetStatus()
		}

		currentHash := fmt.Sprintf("%d-%d-%d-%d-%d-%s",
			snapshot.PendingTxCount,
			snapshot.TotalConfirmedTx,
			snapshot.TotalSubmittedTx,
			gasInWindow,
			status,
			snapshot.LastUpdate.Format(time.RFC3339))

		newHashes[spammerID] = currentHash

		if oldHash, exists := state.lastSpammerHashes[spammerID]; !exists || oldHash != currentHash {
			hasSpammerChanges = true
		}
	}

	// Only send data if there are changes
	if !hasNewDataPoints && !hasSpammerChanges {
		return
	}

	// Build response data
	data := make(map[string]interface{})

	// Add new data points if any
	if hasNewDataPoints {
		newDataPoints := dataPoints[state.lastDataPointCount:]
		convertedDataPoints := make([]GraphsDataPoint, len(newDataPoints))

		for i, point := range newDataPoints {
			convertedDataPoints[i] = GraphsDataPoint{
				Timestamp:        point.Timestamp,
				StartBlockNumber: point.StartBlockNumber,
				EndBlockNumber:   point.EndBlockNumber,
				BlockCount:       point.BlockCount,
				TotalGasUsed:     point.TotalGasUsed,
				OthersGasUsed:    point.OthersGasUsed,
				SpammerData:      make(map[string]*SpammerBlockData),
			}

			// Convert spammer data with string keys
			for spammerID, spammerData := range point.SpammerGasData {
				convertedDataPoints[i].SpammerData[fmt.Sprintf("%d", spammerID)] = &SpammerBlockData{
					GasUsed:          spammerData.GasUsed,
					ConfirmedTxCount: spammerData.ConfirmedTxCount,
					PendingTxCount:   spammerData.PendingTxCount,
					SubmittedTxCount: spammerData.SubmittedTxCount,
				}
			}
		}

		data["newDataPoints"] = convertedDataPoints
		state.lastDataPointCount = len(dataPoints)
	}

	// Add spammer updates and detect new spammers
	if hasSpammerChanges {
		newSpammers := make([]map[string]interface{}, 0)

		for spammerID, snapshot := range spammerSnapshots {
			spammerName := ah.daemon.GetSpammerName(spammerID)
			gasInWindow := spammerGasInWindow[spammerID]

			// Get spammer status
			status := 0 // Default to stopped
			if spammer := ah.daemon.GetSpammer(int64(spammerID)); spammer != nil {
				status = spammer.GetStatus()
			}

			spammerData := map[string]interface{}{
				"id":        spammerID,
				"name":      spammerName,
				"pending":   snapshot.PendingTxCount,
				"confirmed": snapshot.TotalConfirmedTx,
				"submitted": snapshot.TotalSubmittedTx,
				"gasUsed":   gasInWindow,
				"updated":   snapshot.LastUpdate.Format(time.RFC3339),
				"status":    status,
			}

			// Check if this is a new spammer
			if _, exists := state.lastSpammerHashes[spammerID]; !exists {
				newSpammers = append(newSpammers, spammerData)
			} else {
				data[fmt.Sprintf("spammer_%d", spammerID)] = spammerData
			}
		}

		// Include new spammers in the response
		if len(newSpammers) > 0 {
			data["newSpammers"] = newSpammers
		}

		state.lastSpammerHashes = newHashes
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Failed to marshal graphs data: %v", err)
		return
	}

	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	flusher.Flush()
}

// sendMetricsUpdateWithState sends real-time metrics updates via SSE
func (ah *APIHandler) sendMetricsUpdateWithState(w http.ResponseWriter, flusher http.Flusher, update *daemon.MetricsUpdate, state *SSEState) {
	if update == nil {
		return
	}

	data := make(map[string]interface{})

	// Add new data point if available
	if update.NewDataPoint != nil {
		convertedDataPoint := GraphsDataPoint{
			Timestamp:        update.NewDataPoint.Timestamp,
			StartBlockNumber: update.NewDataPoint.StartBlockNumber,
			EndBlockNumber:   update.NewDataPoint.EndBlockNumber,
			BlockCount:       update.NewDataPoint.BlockCount,
			TotalGasUsed:     update.NewDataPoint.TotalGasUsed,
			OthersGasUsed:    update.NewDataPoint.OthersGasUsed,
			SpammerData:      make(map[string]*SpammerBlockData),
		}

		// Convert spammer data with string keys
		for spammerID, spammerData := range update.NewDataPoint.SpammerGasData {
			convertedDataPoint.SpammerData[fmt.Sprintf("%d", spammerID)] = &SpammerBlockData{
				GasUsed:          spammerData.GasUsed,
				ConfirmedTxCount: spammerData.ConfirmedTxCount,
				PendingTxCount:   spammerData.PendingTxCount,
				SubmittedTxCount: spammerData.SubmittedTxCount,
			}
		}

		data["newDataPoints"] = []GraphsDataPoint{convertedDataPoint}
	}

	// Add updated spammer snapshots
	if len(update.UpdatedSpammers) > 0 {
		// Get the complete window data to calculate cumulative gas usage
		shortWindow := ah.daemon.GetShortWindowMetrics()
		spammerGasInWindow := make(map[uint64]uint64)
		if shortWindow != nil {
			dataPoints := shortWindow.GetDataPoints()
			for _, point := range dataPoints {
				for spammerID, spammerData := range point.SpammerGasData {
					spammerGasInWindow[spammerID] += spammerData.GasUsed
				}
			}
		}

		newSpammers := make([]map[string]interface{}, 0)

		for spammerID, snapshot := range update.UpdatedSpammers {
			spammerName := ah.daemon.GetSpammerName(spammerID)

			// Get spammer status
			status := 0 // Default to stopped
			if spammer := ah.daemon.GetSpammer(int64(spammerID)); spammer != nil {
				status = spammer.GetStatus()
			}

			// Use cumulative gas usage from the entire window
			gasInWindow := spammerGasInWindow[spammerID]

			spammerData := map[string]interface{}{
				"id":        spammerID,
				"name":      spammerName,
				"pending":   snapshot.PendingTxCount,
				"confirmed": snapshot.TotalConfirmedTx,
				"submitted": snapshot.TotalSubmittedTx,
				"gasUsed":   gasInWindow,
				"updated":   snapshot.LastUpdate.Format(time.RFC3339),
				"status":    status,
			}

			// Check if this is a new spammer
			if _, exists := state.lastSpammerHashes[spammerID]; !exists {
				newSpammers = append(newSpammers, spammerData)
			} else {
				data[fmt.Sprintf("spammer_%d", spammerID)] = spammerData
			}
		}

		// Include new spammers in the response
		if len(newSpammers) > 0 {
			data["newSpammers"] = newSpammers
		}
	}

	// Only send if we have data
	if len(data) == 0 {
		return
	}

	// Update state to track what we've sent
	if update.NewDataPoint != nil {
		// Get current window data to update our state
		shortWindow := ah.daemon.GetShortWindowMetrics()
		if shortWindow != nil {
			dataPoints := shortWindow.GetDataPoints()
			state.lastDataPointCount = len(dataPoints)
		}
	}

	// Update spammer hashes for real-time updates
	if len(update.UpdatedSpammers) > 0 {
		for spammerID, snapshot := range update.UpdatedSpammers {
			// Get spammer status for hash
			status := 0
			if spammer := ah.daemon.GetSpammer(int64(spammerID)); spammer != nil {
				status = spammer.GetStatus()
			}

			// Get gas usage
			gasInWindow := uint64(0)
			if shortWindow := ah.daemon.GetShortWindowMetrics(); shortWindow != nil {
				dataPoints := shortWindow.GetDataPoints()
				for _, point := range dataPoints {
					if spammerData, exists := point.SpammerGasData[spammerID]; exists {
						gasInWindow += spammerData.GasUsed
					}
				}
			}

			currentHash := fmt.Sprintf("%d-%d-%d-%d-%d-%s",
				snapshot.PendingTxCount,
				snapshot.TotalConfirmedTx,
				snapshot.TotalSubmittedTx,
				gasInWindow,
				status,
				snapshot.LastUpdate.Format(time.RFC3339))

			state.lastSpammerHashes[spammerID] = currentHash
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Failed to marshal real-time metrics update: %v", err)
		return
	}

	fmt.Fprintf(w, "data: %s\n\n", jsonData)
	flusher.Flush()
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// sendError sends a JSON error response
func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// SendTransaction godoc
// @Id sendTransaction
// @Summary Send a transaction from the root wallet
// @Tags Wallet
// @Description Sends a transaction from the root wallet with specified parameters
// @Accept json
// @Produce json
// @Param request body SendTransactionRequest true "Transaction parameters"
// @Success 200 {object} SendTransactionResponse "Success"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {object} ErrorResponse "Server Error"
// @Router /api/root-wallet/send-transaction [post]
func (ah *APIHandler) SendTransaction(w http.ResponseWriter, r *http.Request) {
	if !ah.checkAuth(w, r) {
		return
	}

	var req SendTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.To == "" {
		sendError(w, "target address is required", http.StatusBadRequest)
		return
	}
	if req.Value == "" {
		sendError(w, "value is required", http.StatusBadRequest)
		return
	}
	if req.Unit == "" {
		sendError(w, "unit is required", http.StatusBadRequest)
		return
	}

	// Validate unit
	if req.Unit != "eth" && req.Unit != "gwei" && req.Unit != "wei" {
		sendError(w, "unit must be 'eth', 'gwei', or 'wei'", http.StatusBadRequest)
		return
	}

	// Get root wallet
	rootWallet := ah.daemon.GetRootWallet()
	if rootWallet == nil || rootWallet.GetWallet() == nil {
		sendError(w, "root wallet not available", http.StatusInternalServerError)
		return
	}

	// Convert value to wei
	valueFloat, err := strconv.ParseFloat(req.Value, 64)
	if err != nil {
		sendError(w, "invalid value format", http.StatusBadRequest)
		return
	}

	var valueWei *big.Int
	switch req.Unit {
	case "eth":
		// Convert ETH to wei (1 ETH = 10^18 wei)
		ethInWei := new(big.Float).Mul(big.NewFloat(valueFloat), big.NewFloat(1e18))
		valueWei = new(big.Int)
		ethInWei.Int(valueWei)
	case "gwei":
		// Convert gwei to wei (1 gwei = 10^9 wei)
		gweiInWei := new(big.Float).Mul(big.NewFloat(valueFloat), big.NewFloat(1e9))
		valueWei = new(big.Int)
		gweiInWei.Int(valueWei)
	case "wei":
		// Already in wei
		valueWei = new(big.Int)
		valueWei.SetString(req.Value, 10)
	}

	// Parse calldata
	var calldata []byte
	if req.Data != "" {
		// Remove 0x prefix if present
		dataStr := strings.TrimPrefix(req.Data, "0x")
		calldata, err = hex.DecodeString(dataStr)
		if err != nil {
			sendError(w, "invalid calldata format", http.StatusBadRequest)
			return
		}
	}

	// Set default gas limit if not provided
	gasLimit := req.GasLimit
	if gasLimit == 0 {
		if len(calldata) > 0 {
			gasLimit = 100000 // Higher limit for contract calls
		} else {
			gasLimit = 21000 // Standard transfer
		}
	}

	// Get user email for audit logging
	userEmail := "api"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	// Build and send transaction
	wallet := rootWallet.GetWallet()
	toAddr := common.HexToAddress(req.To)

	// Get current gas prices if not specified
	ctx := r.Context()
	client := ah.daemon.GetClientPool().GetClient()
	if client == nil {
		sendError(w, "no available client", http.StatusInternalServerError)
		return
	}

	// Get suggested gas prices
	var maxFeePerGas *big.Int
	var maxPriorityFeePerGas *big.Int

	if req.MaxFee != "" || req.MaxTip != "" {
		// Use provided values - if only one is provided, get the other from network
		baseFee, err := client.GetEthClient().SuggestGasPrice(ctx)
		if err != nil {
			sendError(w, fmt.Sprintf("failed to get gas price: %v", err), http.StatusInternalServerError)
			return
		}

		if req.MaxFee != "" {
			maxFeeFloat, err := strconv.ParseFloat(req.MaxFee, 64)
			if err != nil {
				sendError(w, "invalid maxFee format", http.StatusBadRequest)
				return
			}
			maxFeePerGas = new(big.Int).SetUint64(uint64(maxFeeFloat * 1e9)) // Convert gwei to wei
		} else {
			// Default max fee if not provided
			defaultTip := big.NewInt(2e9)
			maxFeePerGas = new(big.Int).Add(baseFee, defaultTip)
		}

		if req.MaxTip != "" {
			maxTipFloat, err := strconv.ParseFloat(req.MaxTip, 64)
			if err != nil {
				sendError(w, "invalid maxTip format", http.StatusBadRequest)
				return
			}
			maxPriorityFeePerGas = new(big.Int).SetUint64(uint64(maxTipFloat * 1e9)) // Convert gwei to wei
		} else {
			// Default tip if not provided
			maxPriorityFeePerGas = big.NewInt(2e9)
		}

		// Ensure priority fee doesn't exceed max fee
		if maxPriorityFeePerGas.Cmp(maxFeePerGas) > 0 {
			sendError(w, "max priority fee per gas cannot be higher than max fee per gas", http.StatusBadRequest)
			return
		}
	} else {
		// Get suggested values from network
		baseFee, err := client.GetEthClient().SuggestGasPrice(ctx)
		if err != nil {
			sendError(w, fmt.Sprintf("failed to get gas price: %v", err), http.StatusInternalServerError)
			return
		}

		// Set priority fee to 2 gwei by default
		maxPriorityFeePerGas = big.NewInt(2e9)

		// Add 10% buffer to base fee + priority fee for maxFeePerGas
		maxFeePerGas = new(big.Int).Add(baseFee, maxPriorityFeePerGas)
		maxFeePerGas = new(big.Int).Mul(maxFeePerGas, big.NewInt(110))
		maxFeePerGas.Div(maxFeePerGas, big.NewInt(100))
	}

	// Build transaction
	txData := &types.DynamicFeeTx{
		To:        &toAddr,
		Value:     valueWei,
		Gas:       gasLimit,
		GasFeeCap: maxFeePerGas,
		GasTipCap: maxPriorityFeePerGas,
		Data:      calldata,
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		sendError(w, fmt.Sprintf("failed to build transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Send transaction
	err = ah.daemon.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		sendError(w, fmt.Sprintf("failed to send transaction: %v", err), http.StatusInternalServerError)
		return
	}

	// Log to audit log
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		maxFeeStr := fmt.Sprintf("%.2f gwei", float64(maxFeePerGas.Uint64())/1e9)
		maxTipStr := fmt.Sprintf("%.2f gwei", float64(maxPriorityFeePerGas.Uint64())/1e9)
		dataStr := req.Data
		if dataStr == "" {
			dataStr = "0x"
		}
		auditLogger.LogRootWalletTransaction(userEmail, toAddr.Hex(), valueWei.String(), tx.Hash().Hex(), gasLimit, maxFeeStr, maxTipStr, dataStr)
	}

	// Return response
	response := SendTransactionResponse{
		TxHash: tx.Hash().Hex(),
		Nonce:  tx.Nonce(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
