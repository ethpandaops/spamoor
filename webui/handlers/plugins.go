package handlers

import (
	"context"
	"net/http"

	"github.com/ethpandaops/spamoor/webui/server"
)

// PluginsPage represents the data passed to the plugins template.
type PluginsPage struct {
	Plugins     []*PluginsPageEntry `json:"plugins"`
	PluginCount uint64              `json:"plugin_count"`
}

// PluginsPageEntry represents a single plugin entry for the frontend.
type PluginsPageEntry struct {
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

// Plugins renders the plugins management page.
func (fh *FrontendHandler) Plugins(w http.ResponseWriter, r *http.Request) {
	templateFiles := append(server.LayoutTemplateFiles,
		"plugins/plugins.html",
	)

	pageTemplate := server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "plugins", "/plugins", "Plugins", templateFiles)

	var pageError error
	data.Data, pageError = fh.getPluginsPageData(r.Context())
	if pageError != nil {
		server.HandlePageError(w, r, pageError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "plugins.go", "Plugins", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return // an error has occurred and was processed
	}
}

func (fh *FrontendHandler) getPluginsPageData(_ context.Context) (*PluginsPage, error) {
	pageData := &PluginsPage{
		Plugins: []*PluginsPageEntry{},
	}

	persistence := fh.daemon.GetPluginPersistence()
	if persistence == nil {
		// Plugin persistence not available, return empty page
		return pageData, nil
	}

	statuses, err := persistence.GetPluginStatuses()
	if err != nil {
		return nil, err
	}

	for _, status := range statuses {
		entry := &PluginsPageEntry{
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
		pageData.Plugins = append(pageData.Plugins, entry)
	}

	pageData.PluginCount = uint64(len(pageData.Plugins))

	return pageData, nil
}
