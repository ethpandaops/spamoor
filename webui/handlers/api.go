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

	// Convert to a simpler format for JSON
	type LogEntry struct {
		Time    time.Time `json:"time"`
		Level   string    `json:"level"`
		Message string    `json:"message"`
	}

	logs := make([]LogEntry, len(entries))
	for i, entry := range entries {
		logs[i] = LogEntry{
			Time:    entry.Time,
			Level:   entry.Level.String(),
			Message: entry.Message,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
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
