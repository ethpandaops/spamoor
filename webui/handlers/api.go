package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/ethpandaops/spamoor/scenarios"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func (fh *FrontendHandler) StartSpammer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	spammer := fh.daemon.GetSpammer(id)
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

func (fh *FrontendHandler) PauseSpammer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	spammer := fh.daemon.GetSpammer(id)
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

func (fh *FrontendHandler) DeleteSpammer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	err = fh.daemon.DeleteSpammer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type ScenarioEntries struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (fh *FrontendHandler) GetScenarios(w http.ResponseWriter, r *http.Request) {
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

func (fh *FrontendHandler) GetScenarioConfig(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/yaml")
	fmt.Fprintf(w, "# wallet settings\n")
	fmt.Fprintf(w, "seed: %s-%v # seed for the wallet\n", scenarioName, rand.Intn(1000000))
	fmt.Fprintf(w, "refill_amount: 5000000000000000000 # refill 5 ETH when\n")
	fmt.Fprintf(w, "refill_balance: 1000000000000000000 # balance drops below 1 ETH\n")
	fmt.Fprintf(w, "refill_interval: 600 # check every 10 minutes\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "# scenario: %s\n", scenarioName)
	fmt.Fprintf(w, "%s\n", string(configYaml))
}

func (fh *FrontendHandler) CreateSpammer(w http.ResponseWriter, r *http.Request) {
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

	spammer, err := fh.daemon.NewSpammer(req.Scenario, req.Config, req.Name, req.Description, req.StartImmediately)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spammer.GetID())
}

func (fh *FrontendHandler) GetSpammerDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	spammer := fh.daemon.GetSpammer(id)
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

func (fh *FrontendHandler) UpdateSpammer(w http.ResponseWriter, r *http.Request) {
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

	spammer := fh.daemon.GetSpammer(id)
	if spammer == nil {
		http.Error(w, "Spammer not found", http.StatusNotFound)
		return
	}

	err = fh.daemon.UpdateSpammer(id, req.Name, req.Description, req.Config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type LogEntry struct {
	Time    string            `json:"time"`
	Level   string            `json:"level"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

func (fh *FrontendHandler) GetSpammerLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	spammer := fh.daemon.GetSpammer(id)
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

func (fh *FrontendHandler) StreamSpammerLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid spammer ID", http.StatusBadRequest)
		return
	}

	spammer := fh.daemon.GetSpammer(id)
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
