package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/ethpandaops/spamoor/daemon"
	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/plugin"
	"github.com/gorilla/mux"
)

// PluginEntry represents a plugin in the API response.
type PluginEntry struct {
	Name               string   `json:"name"`
	SourceType         string   `json:"source_type"`
	SourcePath         string   `json:"source_path"`
	MetadataName       string   `json:"metadata_name,omitempty"`
	MetadataBuildTime  string   `json:"metadata_build_time,omitempty"`
	MetadataGitVersion string   `json:"metadata_git_version,omitempty"`
	Scenarios          []string `json:"scenarios"`
	Enabled            bool     `json:"enabled"`
	LoadError          string   `json:"load_error,omitempty"`
	RunningCount       int32    `json:"running_count"`
	IsLoaded           bool     `json:"is_loaded"`
	Deprecated         bool     `json:"deprecated"`
	CreatedAt          int64    `json:"created_at"`
	UpdatedAt          int64    `json:"updated_at"`
}

// RegisterPluginRequest represents the request body for registering a plugin.
type RegisterPluginRequest struct {
	Type string `json:"type"` // "url", "local", or "upload" (upload uses multipart)
	Path string `json:"path"` // URL or local path
}

// RegisterPluginResponse represents the response after registering a plugin.
type RegisterPluginResponse struct {
	Name      string   `json:"name"`
	Scenarios []string `json:"scenarios"`
}

// GetPlugins godoc
// @Id getPlugins
// @Summary Get all plugins
// @Tags Plugin
// @Description Returns a list of all registered plugins with their status
// @Produce json
// @Success 200 {array} PluginEntry "Success"
// @Failure 500 {string} string "Server Error"
// @Router /api/plugins [get]
func (ah *APIHandler) GetPlugins(w http.ResponseWriter, r *http.Request) {
	persistence := ah.daemon.GetPluginPersistence()
	if persistence == nil {
		http.Error(w, "Plugin persistence not available", http.StatusServiceUnavailable)
		return
	}

	statuses, err := persistence.GetPluginStatuses()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	entries := make([]*PluginEntry, len(statuses))
	for i, status := range statuses {
		entries[i] = convertStatusToEntry(status)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// RegisterPlugin godoc
// @Id registerPlugin
// @Summary Register a new plugin
// @Tags Plugin
// @Description Registers a plugin from URL, local path, or file upload
// @Accept multipart/form-data
// @Produce json
// @Param type formData string true "Source type: url, local, or upload"
// @Param path formData string false "URL or local path (required for url/local types)"
// @Param file formData file false "Plugin archive file (required for upload type)"
// @Success 200 {object} RegisterPluginResponse "Success"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Server Error"
// @Router /api/plugins [post]
func (ah *APIHandler) RegisterPlugin(w http.ResponseWriter, r *http.Request) {
	persistence := ah.daemon.GetPluginPersistence()
	if persistence == nil {
		http.Error(w, "Plugin persistence not available", http.StatusServiceUnavailable)
		return
	}

	userEmail := ah.getPluginUserEmail(r)

	// Check content type for multipart
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle multipart form upload
		err := r.ParseMultipartForm(100 << 20) // 100 MB max
		if err != nil {
			http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
			return
		}

		sourceType := r.FormValue("type")
		sourcePath := r.FormValue("path")

		switch sourceType {
		case "url":
			if sourcePath == "" {
				http.Error(w, "URL is required for type=url", http.StatusBadRequest)
				return
			}

			loaded, err := persistence.RegisterPluginFromURL(sourcePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ah.auditPluginAction(userEmail, db.AuditActionPluginRegister, loaded.Descriptor.Name, buildPluginMetadata(loaded, "url", sourcePath))
			sendPluginResponse(w, loaded)

		case "local":
			if sourcePath == "" {
				http.Error(w, "Path is required for type=local", http.StatusBadRequest)
				return
			}

			loaded, err := persistence.RegisterPluginFromLocal(sourcePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ah.auditPluginAction(userEmail, db.AuditActionPluginRegister, loaded.Descriptor.Name, buildPluginMetadata(loaded, "local", sourcePath))
			sendPluginResponse(w, loaded)

		case "upload":
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "File upload is required for type=upload: "+err.Error(), http.StatusBadRequest)
				return
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Failed to read uploaded file: "+err.Error(), http.StatusInternalServerError)
				return
			}

			loaded, err := persistence.RegisterPluginFromUpload(data, header.Filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ah.auditPluginAction(userEmail, db.AuditActionPluginRegister, loaded.Descriptor.Name, buildPluginMetadata(loaded, "upload", header.Filename))
			sendPluginResponse(w, loaded)

		default:
			http.Error(w, "Invalid type. Must be 'url', 'local', or 'upload'", http.StatusBadRequest)
		}
	} else {
		// Handle JSON request
		var req RegisterPluginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch req.Type {
		case "url":
			if req.Path == "" {
				http.Error(w, "URL is required for type=url", http.StatusBadRequest)
				return
			}

			loaded, err := persistence.RegisterPluginFromURL(req.Path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ah.auditPluginAction(userEmail, db.AuditActionPluginRegister, loaded.Descriptor.Name, buildPluginMetadata(loaded, "url", req.Path))
			sendPluginResponse(w, loaded)

		case "local":
			if req.Path == "" {
				http.Error(w, "Path is required for type=local", http.StatusBadRequest)
				return
			}

			loaded, err := persistence.RegisterPluginFromLocal(req.Path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ah.auditPluginAction(userEmail, db.AuditActionPluginRegister, loaded.Descriptor.Name, buildPluginMetadata(loaded, "local", req.Path))
			sendPluginResponse(w, loaded)

		default:
			http.Error(w, "Invalid type for JSON request. Must be 'url' or 'local'. Use multipart/form-data for file uploads.", http.StatusBadRequest)
		}
	}
}

// DeletePlugin godoc
// @Id deletePlugin
// @Summary Delete a plugin
// @Tags Plugin
// @Description Deletes a plugin by name. Fails if the plugin has running spammers.
// @Param name path string true "Plugin name"
// @Success 200 "Success"
// @Failure 400 {string} string "Plugin has running spammers"
// @Failure 404 {string} string "Plugin not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/plugins/{name} [delete]
func (ah *APIHandler) DeletePlugin(w http.ResponseWriter, r *http.Request) {
	persistence := ah.daemon.GetPluginPersistence()
	if persistence == nil {
		http.Error(w, "Plugin persistence not available", http.StatusServiceUnavailable)
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		http.Error(w, "Plugin name is required", http.StatusBadRequest)
		return
	}

	userEmail := ah.getPluginUserEmail(r)

	err := persistence.DeletePlugin(name)
	if err != nil {
		if strings.Contains(err.Error(), "running spammer") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	ah.auditPluginAction(userEmail, db.AuditActionPluginDelete, name, nil)

	w.WriteHeader(http.StatusOK)
}

// ReloadPlugin godoc
// @Id reloadPlugin
// @Summary Reload a plugin from URL
// @Tags Plugin
// @Description Re-downloads a URL plugin and updates the stored archive
// @Param name path string true "Plugin name"
// @Success 200 {object} RegisterPluginResponse "Success"
// @Failure 400 {string} string "Plugin has running spammers or is not a URL plugin"
// @Failure 404 {string} string "Plugin not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/plugins/{name}/reload [post]
func (ah *APIHandler) ReloadPlugin(w http.ResponseWriter, r *http.Request) {
	persistence := ah.daemon.GetPluginPersistence()
	if persistence == nil {
		http.Error(w, "Plugin persistence not available", http.StatusServiceUnavailable)
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		http.Error(w, "Plugin name is required", http.StatusBadRequest)
		return
	}

	userEmail := ah.getPluginUserEmail(r)

	loaded, err := persistence.ReloadPluginFromURL(name)
	if err != nil {
		if strings.Contains(err.Error(), "running spammer") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if strings.Contains(err.Error(), "only supported for URL") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	ah.auditPluginAction(userEmail, db.AuditActionPluginReload, loaded.Descriptor.Name, buildPluginMetadata(loaded, "url", ""))
	sendPluginResponse(w, loaded)
}

// GetPlugin godoc
// @Id getPlugin
// @Summary Get a specific plugin
// @Tags Plugin
// @Description Returns details of a specific plugin by name
// @Produce json
// @Param name path string true "Plugin name"
// @Success 200 {object} PluginEntry "Success"
// @Failure 404 {string} string "Plugin not found"
// @Failure 500 {string} string "Server Error"
// @Router /api/plugins/{name} [get]
func (ah *APIHandler) GetPlugin(w http.ResponseWriter, r *http.Request) {
	persistence := ah.daemon.GetPluginPersistence()
	if persistence == nil {
		http.Error(w, "Plugin persistence not available", http.StatusServiceUnavailable)
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		http.Error(w, "Plugin name is required", http.StatusBadRequest)
		return
	}

	status, err := persistence.GetPluginStatus(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	entry := convertStatusToEntry(status)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

// getPluginUserEmail extracts the user email from the request for audit logging.
func (ah *APIHandler) getPluginUserEmail(r *http.Request) string {
	userEmail := "user"
	if auditLogger := ah.daemon.GetAuditLogger(); auditLogger != nil {
		userEmail = auditLogger.GetUserFromRequest(r.Header)
	}

	return userEmail
}

// auditPluginAction logs a plugin action to the audit log.
func (ah *APIHandler) auditPluginAction(
	userEmail string,
	action db.AuditActionType,
	pluginName string,
	metadata db.AuditMetadata,
) {
	auditLogger := ah.daemon.GetAuditLogger()
	if auditLogger == nil || userEmail == "" {
		return
	}

	//nolint:errcheck // best-effort audit logging
	auditLogger.LogPluginAction(userEmail, action, pluginName, metadata)
}

// buildPluginMetadata builds audit metadata from a loaded plugin.
func buildPluginMetadata(
	loaded *plugin.LoadedPlugin,
	sourceType string,
	sourcePath string,
) db.AuditMetadata {
	allScenarios := loaded.Descriptor.GetAllScenarios()
	scenarioNames := make([]string, 0, len(allScenarios))
	for _, s := range allScenarios {
		scenarioNames = append(scenarioNames, s.Name)
	}

	return db.AuditMetadata{
		"source_type":    sourceType,
		"source_path":    sourcePath,
		"scenario_count": len(scenarioNames),
		"scenarios":      scenarioNames,
	}
}

// sendPluginResponse sends a successful plugin registration response.
func sendPluginResponse(w http.ResponseWriter, loaded *plugin.LoadedPlugin) {
	allScenarios := loaded.Descriptor.GetAllScenarios()
	scenarios := make([]string, 0, len(allScenarios))
	for _, s := range allScenarios {
		scenarios = append(scenarios, s.Name)
	}

	response := RegisterPluginResponse{
		Name:      loaded.Descriptor.Name,
		Scenarios: scenarios,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// convertStatusToEntry converts a PluginStatus to a PluginEntry.
func convertStatusToEntry(status *daemon.PluginStatus) *PluginEntry {
	return &PluginEntry{
		Name:               status.Name,
		SourceType:         status.SourceType,
		SourcePath:         status.SourcePath,
		MetadataName:       status.MetadataName,
		MetadataBuildTime:  status.MetadataBuildTime,
		MetadataGitVersion: status.MetadataGitVersion,
		Scenarios:          status.Scenarios,
		Enabled:            status.Enabled,
		LoadError:          status.LoadError,
		RunningCount:       status.RunningCount,
		IsLoaded:           status.IsLoaded,
		Deprecated:         status.Deprecated,
		CreatedAt:          status.CreatedAt,
		UpdatedAt:          status.UpdatedAt,
	}
}
