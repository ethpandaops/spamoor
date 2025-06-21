package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/ethpandaops/spamoor/scenarios"
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

type ScenarioEntries struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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
// @Description Returns a list of all scenarios
// @Produce json
// @Success 200 {array} ScenarioEntries "Success"
// @Failure 400 {string} string "Failure"
// @Failure 500 {string} string "Server Error"
// @Router /api/scenarios [get]
func (ah *APIHandler) GetScenarios(w http.ResponseWriter, r *http.Request) {
	scenarioNames := scenarios.GetScenarioNames()
	entries := make([]*ScenarioEntries, len(scenarioNames))
	for i, scenarioName := range scenarioNames {
		entries[i] = &ScenarioEntries{
			Name:        scenarioName,
			Description: scenarios.GetScenario(scenarioName).Description,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
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
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer [post]
func (ah *APIHandler) CreateSpammer(w http.ResponseWriter, r *http.Request) {
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

	spammer, err := ah.daemon.NewSpammer(req.Scenario, req.Config, req.Name, req.Description, req.StartImmediately)
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
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id}/start [post]
func (ah *APIHandler) StartSpammer(w http.ResponseWriter, r *http.Request) {
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

	err = spammer.Start()
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
// @Failure 400 {object} Response "Invalid request"
// @Failure 404 {object} Response "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id} [put]
func (ah *APIHandler) UpdateSpammer(w http.ResponseWriter, r *http.Request) {
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

	err = ah.daemon.UpdateSpammer(id, req.Name, req.Description, req.Config)
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
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id}/pause [post]
func (ah *APIHandler) PauseSpammer(w http.ResponseWriter, r *http.Request) {
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

	err = spammer.Pause()
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
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id} [delete]
func (ah *APIHandler) DeleteSpammer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	err = ah.daemon.DeleteSpammer(id)
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
// @Failure 404 {string} string "Spammer not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/spammer/{id}/reclaim [post]
func (ah *APIHandler) ReclaimFunds(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	err = ah.daemon.ReclaimSpammer(id)
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
// @Failure 400 {object} Response "Invalid spammer ID"
// @Failure 404 {object} Response "Spammer not found"
// @Router /api/spammer/{id} [get]
func (ah *APIHandler) GetSpammerDetails(w http.ResponseWriter, r *http.Request) {
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
// @Failure 400 {object} Response "Invalid spammer ID"
// @Failure 404 {object} Response "Spammer not found"
// @Router /api/spammer/{id}/logs [get]
func (ah *APIHandler) GetSpammerLogs(w http.ResponseWriter, r *http.Request) {
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
// @Failure 400 {object} Response "Invalid spammer ID"
// @Failure 404 {object} Response "Spammer not found"
// @Failure 500 {string} string "Streaming unsupported"
// @Router /api/spammer/{id}/logs/stream [get]
func (ah *APIHandler) StreamSpammerLogs(w http.ResponseWriter, r *http.Request) {
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

// ClientEntry represents a client in the API response
type ClientEntry struct {
	Index        int      `json:"index"`
	Name         string   `json:"name"`
	Group        string   `json:"group"`  // First group for backward compatibility
	Groups       []string `json:"groups"` // All groups
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

// GetClients godoc
// @Id getClients
// @Summary Get all clients
// @Tags Client
// @Description Returns a list of all clients with their details
// @Produce json
// @Success 200 {array} ClientEntry "Success"
// @Router /api/clients [get]
func (ah *APIHandler) GetClients(w http.ResponseWriter, r *http.Request) {
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

		response[i] = ClientEntry{
			Index:        i,
			Name:         client.GetName(),
			Group:        client.GetClientGroup(),
			Groups:       client.GetClientGroups(),
			Version:      version,
			BlockHeight:  blockHeight,
			IsReady:      slices.Contains(goodClients, client),
			RpcHost:      client.GetRPCHost(),
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
// @Failure 404 {string} string "Client not found"
// @Router /api/client/{index}/group [put]
func (ah *APIHandler) UpdateClientGroup(w http.ResponseWriter, r *http.Request) {
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

	// Handle both single group (backward compatibility) and multiple groups
	if len(req.Groups) > 0 {
		client.SetClientGroups(req.Groups)
	} else if req.Group != "" {
		client.SetClientGroups([]string{req.Group})
	} else {
		http.Error(w, "Either 'group' or 'groups' must be provided", http.StatusBadRequest)
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
// @Failure 404 {string} string "Client not found"
// @Router /api/client/{index}/enabled [put]
func (ah *APIHandler) UpdateClientEnabled(w http.ResponseWriter, r *http.Request) {
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
	client.SetEnabled(req.Enabled)

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
// @Failure 404 {string} string "Client not found"
// @Router /api/client/{index}/name [put]
func (ah *APIHandler) UpdateClientName(w http.ResponseWriter, r *http.Request) {
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
	client.SetNameOverride(req.NameOverride)

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
// @Failure 500 {string} string "Server Error"
// @Router /api/spammers/export [post]
func (ah *APIHandler) ExportSpammers(w http.ResponseWriter, r *http.Request) {
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
// @Failure 500 {string} string "Server Error"
// @Router /api/spammers/import [post]
func (ah *APIHandler) ImportSpammers(w http.ResponseWriter, r *http.Request) {
	var req ImportSpammersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Input == "" {
		http.Error(w, "input is required", http.StatusBadRequest)
		return
	}

	result, err := ah.daemon.ImportSpammers(req.Input)
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
