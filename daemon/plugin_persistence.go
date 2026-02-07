package daemon

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/daemon/db"
	"github.com/ethpandaops/spamoor/plugin"
)

// PluginStatus represents the combined status of a plugin from both DB and runtime.
type PluginStatus struct {
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
	CreatedAt          int64    `json:"created_at"`
	UpdatedAt          int64    `json:"updated_at"`
}

// PluginPersistence handles saving and restoring plugins from the database.
type PluginPersistence struct {
	logger       *logrus.Entry
	db           *db.Database
	pluginLoader *plugin.PluginLoader
}

// NewPluginPersistence creates a new PluginPersistence instance.
func NewPluginPersistence(
	logger logrus.FieldLogger,
	database *db.Database,
	pluginLoader *plugin.PluginLoader,
) *PluginPersistence {
	return &PluginPersistence{
		logger:       logger.WithField("component", "plugin-persistence"),
		db:           database,
		pluginLoader: pluginLoader,
	}
}

// SavePlugin persists a loaded plugin to the database.
// For non-local plugins, archiveData should be the original tar.gz bytes.
func (pp *PluginPersistence) SavePlugin(
	loaded *plugin.LoadedPlugin,
	archiveData []byte,
	sourcePath string,
) error {
	// Build scenario list
	allScenarios := loaded.Descriptor.GetAllScenarios()
	scenarioNames := make([]string, 0, len(allScenarios))
	for _, s := range allScenarios {
		scenarioNames = append(scenarioNames, s.Name)
	}

	// Encode archive data as base64 (skip for local plugins)
	var archiveBase64 string
	if loaded.SourceType != plugin.PluginSourceLocal && len(archiveData) > 0 {
		archiveBase64 = base64.StdEncoding.EncodeToString(archiveData)
	}

	dbPlugin := &db.Plugin{
		Name:               loaded.Descriptor.Name,
		SourceType:         loaded.SourceType.String(),
		SourcePath:         sourcePath,
		MetadataName:       loaded.Metadata.Name,
		MetadataBuildTime:  loaded.Metadata.BuildTime,
		MetadataGitVersion: loaded.Metadata.GitVersion,
		ArchiveData:        archiveBase64,
		Enabled:            true,
		LoadError:          "",
	}
	dbPlugin.SetScenariosFromSlice(scenarioNames)

	// Check if plugin already exists
	existing, err := pp.db.GetPlugin(loaded.Descriptor.Name)
	if err == nil && existing != nil {
		// Update existing plugin
		dbPlugin.ID = existing.ID
		dbPlugin.CreatedAt = existing.CreatedAt

		return pp.db.RunDBTransaction(func(tx *sqlx.Tx) error {
			return pp.db.UpdatePlugin(tx, dbPlugin)
		})
	}

	// Insert new plugin
	return pp.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return pp.db.InsertPlugin(tx, dbPlugin)
	})
}

// RestorePlugins restores all enabled plugins from the database on startup.
// Returns the number of successfully restored plugins.
func (pp *PluginPersistence) RestorePlugins() (int, error) {
	plugins, err := pp.db.GetEnabledPlugins()
	if err != nil {
		return 0, fmt.Errorf("failed to get enabled plugins: %w", err)
	}

	restored := 0

	for _, dbPlugin := range plugins {
		err := pp.restorePlugin(dbPlugin)
		if err != nil {
			pp.logger.Warnf("failed to restore plugin %s: %v", dbPlugin.Name, err)

			// Update load error in database
			updateErr := pp.db.RunDBTransaction(func(tx *sqlx.Tx) error {
				return pp.db.UpdatePluginLoadError(tx, dbPlugin.Name, err.Error())
			})
			if updateErr != nil {
				pp.logger.Warnf("failed to update load error for plugin %s: %v", dbPlugin.Name, updateErr)
			}

			continue
		}

		// Clear any previous load error
		clearErr := pp.db.RunDBTransaction(func(tx *sqlx.Tx) error {
			return pp.db.UpdatePluginLoadError(tx, dbPlugin.Name, "")
		})
		if clearErr != nil {
			pp.logger.Warnf("failed to clear load error for plugin %s: %v", dbPlugin.Name, clearErr)
		}

		restored++
	}

	return restored, nil
}

// restorePlugin restores a single plugin from the database.
// Always prefers stored archive for url/file/bytes types to ensure consistency.
func (pp *PluginPersistence) restorePlugin(dbPlugin *db.Plugin) error {
	var loaded *plugin.LoadedPlugin
	var err error

	switch dbPlugin.SourceType {
	case "url", "file", "bytes":
		// Always restore from stored archive for consistency
		if dbPlugin.ArchiveData == "" {
			return fmt.Errorf("no stored archive available for %s plugin", dbPlugin.SourceType)
		}

		loaded, err = pp.loadFromStoredArchive(dbPlugin)
		if err != nil {
			return fmt.Errorf("failed to load from stored archive: %w", err)
		}

	case "local":
		// Always re-read from local path (no archive stored)
		loaded, err = pp.pluginLoader.LoadFromLocalPath(dbPlugin.SourcePath)
		if err != nil {
			return fmt.Errorf("failed to load from local path: %w", err)
		}

	default:
		return fmt.Errorf("unknown source type: %s", dbPlugin.SourceType)
	}

	// Register plugin scenarios
	if err := pp.pluginLoader.RegisterPluginScenarios(loaded); err != nil {
		return fmt.Errorf("failed to register plugin scenarios: %w", err)
	}

	pp.logger.Infof("restored plugin %s with %d scenarios", loaded.Descriptor.Name, len(loaded.Descriptor.GetAllScenarios()))

	return nil
}

// loadFromStoredArchive loads a plugin from the base64-encoded archive data.
func (pp *PluginPersistence) loadFromStoredArchive(dbPlugin *db.Plugin) (*plugin.LoadedPlugin, error) {
	archiveData, err := base64.StdEncoding.DecodeString(dbPlugin.ArchiveData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode archive data: %w", err)
	}

	// The archive is gzip compressed
	return pp.pluginLoader.LoadFromBytes(archiveData, true)
}

// ReloadPluginFromURL re-downloads a URL plugin and updates the stored archive.
// This allows users to manually trigger an update from the original URL.
// Returns an error if the plugin has running spammers or is not a URL-based plugin.
func (pp *PluginPersistence) ReloadPluginFromURL(name string) (*plugin.LoadedPlugin, error) {
	// Get plugin from database
	dbPlugin, err := pp.db.GetPlugin(name)
	if err != nil {
		return nil, fmt.Errorf("plugin not found: %w", err)
	}

	// Check source type
	if dbPlugin.SourceType != "url" {
		return nil, fmt.Errorf("reload is only supported for URL plugins (plugin is %s type)", dbPlugin.SourceType)
	}

	// Check if plugin has running spammers
	loadedPlugin := pp.pluginLoader.GetPluginRegistry().Get(name)
	if loadedPlugin != nil && loadedPlugin.GetRunningCount() > 0 {
		return nil, fmt.Errorf("cannot reload plugin with %d running spammer(s)", loadedPlugin.GetRunningCount())
	}

	// Download fresh from URL
	resp, err := http.Get(dbPlugin.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download plugin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download plugin: HTTP %d", resp.StatusCode)
	}

	// Read all data for storage
	archiveData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin data: %w", err)
	}

	// Load the plugin
	loaded, err := pp.pluginLoader.LoadFromBytes(archiveData, true)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	loaded.SourceType = plugin.PluginSourceURL

	// Register scenarios (this will replace the old registration)
	if err := pp.pluginLoader.RegisterPluginScenarios(loaded); err != nil {
		return nil, fmt.Errorf("failed to register plugin scenarios: %w", err)
	}

	// Update the stored archive in database
	if err := pp.SavePlugin(loaded, archiveData, dbPlugin.SourcePath); err != nil {
		return nil, fmt.Errorf("failed to save plugin to database: %w", err)
	}

	pp.logger.Infof("reloaded plugin %s from URL with %d scenarios", loaded.Descriptor.Name, len(loaded.Descriptor.GetAllScenarios()))

	return loaded, nil
}

// DeletePlugin removes a plugin from the database.
// Returns an error if the plugin has running spammers.
func (pp *PluginPersistence) DeletePlugin(name string) error {
	// Check if plugin is loaded and has running spammers
	loadedPlugin := pp.pluginLoader.GetPluginRegistry().Get(name)
	if loadedPlugin != nil && loadedPlugin.GetRunningCount() > 0 {
		return fmt.Errorf("cannot delete plugin with %d running spammer(s)", loadedPlugin.GetRunningCount())
	}

	// Unregister plugin scenarios if loaded
	if loadedPlugin != nil {
		if err := pp.pluginLoader.UnregisterPluginScenarios(name); err != nil {
			pp.logger.Warnf("failed to unregister plugin scenarios: %v", err)
		}
	}

	// Delete from database
	return pp.db.RunDBTransaction(func(tx *sqlx.Tx) error {
		return pp.db.DeletePlugin(tx, name)
	})
}

// GetPluginStatuses returns the combined status of all plugins from DB and runtime.
func (pp *PluginPersistence) GetPluginStatuses() ([]*PluginStatus, error) {
	dbPlugins, err := pp.db.GetPlugins()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugins from database: %w", err)
	}

	statuses := make([]*PluginStatus, 0, len(dbPlugins))

	for _, dbPlugin := range dbPlugins {
		status := &PluginStatus{
			Name:               dbPlugin.Name,
			SourceType:         dbPlugin.SourceType,
			SourcePath:         dbPlugin.SourcePath,
			MetadataName:       dbPlugin.MetadataName,
			MetadataBuildTime:  dbPlugin.MetadataBuildTime,
			MetadataGitVersion: dbPlugin.MetadataGitVersion,
			Scenarios:          dbPlugin.GetScenariosAsSlice(),
			Enabled:            dbPlugin.Enabled,
			LoadError:          dbPlugin.LoadError,
			RunningCount:       0,
			IsLoaded:           false,
			CreatedAt:          dbPlugin.CreatedAt,
			UpdatedAt:          dbPlugin.UpdatedAt,
		}

		// Check runtime status
		loadedPlugin := pp.pluginLoader.GetPluginRegistry().Get(dbPlugin.Name)
		if loadedPlugin != nil {
			status.IsLoaded = !loadedPlugin.IsCleanedUp()
			status.RunningCount = loadedPlugin.GetRunningCount()

			// Update scenarios from runtime if loaded
			if status.IsLoaded {
				allScenarios := loadedPlugin.Descriptor.GetAllScenarios()
				status.Scenarios = make([]string, 0, len(allScenarios))
				for _, s := range allScenarios {
					status.Scenarios = append(status.Scenarios, s.Name)
				}

				// Update metadata from runtime
				if loadedPlugin.Metadata != nil {
					status.MetadataName = loadedPlugin.Metadata.Name
					status.MetadataBuildTime = loadedPlugin.Metadata.BuildTime
					status.MetadataGitVersion = loadedPlugin.Metadata.GitVersion
				}
			}
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// GetPluginStatus returns the status of a single plugin.
func (pp *PluginPersistence) GetPluginStatus(name string) (*PluginStatus, error) {
	dbPlugin, err := pp.db.GetPlugin(name)
	if err != nil {
		return nil, fmt.Errorf("plugin not found: %w", err)
	}

	status := &PluginStatus{
		Name:               dbPlugin.Name,
		SourceType:         dbPlugin.SourceType,
		SourcePath:         dbPlugin.SourcePath,
		MetadataName:       dbPlugin.MetadataName,
		MetadataBuildTime:  dbPlugin.MetadataBuildTime,
		MetadataGitVersion: dbPlugin.MetadataGitVersion,
		Scenarios:          dbPlugin.GetScenariosAsSlice(),
		Enabled:            dbPlugin.Enabled,
		LoadError:          dbPlugin.LoadError,
		RunningCount:       0,
		IsLoaded:           false,
		CreatedAt:          dbPlugin.CreatedAt,
		UpdatedAt:          dbPlugin.UpdatedAt,
	}

	// Check runtime status
	loadedPlugin := pp.pluginLoader.GetPluginRegistry().Get(name)
	if loadedPlugin != nil {
		status.IsLoaded = !loadedPlugin.IsCleanedUp()
		status.RunningCount = loadedPlugin.GetRunningCount()

		// Update scenarios from runtime if loaded
		if status.IsLoaded {
			allScenarios := loadedPlugin.Descriptor.GetAllScenarios()
			status.Scenarios = make([]string, 0, len(allScenarios))
			for _, s := range allScenarios {
				status.Scenarios = append(status.Scenarios, s.Name)
			}
		}
	}

	return status, nil
}

// RegisterPluginFromURL downloads and registers a plugin from a URL.
func (pp *PluginPersistence) RegisterPluginFromURL(url string) (*plugin.LoadedPlugin, error) {
	// Download the plugin
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download plugin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download plugin: HTTP %d", resp.StatusCode)
	}

	// Read all data for storage
	archiveData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin data: %w", err)
	}

	// Load the plugin
	loaded, err := pp.pluginLoader.LoadFromBytes(archiveData, true)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	loaded.SourceType = plugin.PluginSourceURL

	// Register scenarios
	if err := pp.pluginLoader.RegisterPluginScenarios(loaded); err != nil {
		return nil, fmt.Errorf("failed to register plugin scenarios: %w", err)
	}

	// Save to database
	if err := pp.SavePlugin(loaded, archiveData, url); err != nil {
		return nil, fmt.Errorf("failed to save plugin to database: %w", err)
	}

	return loaded, nil
}

// RegisterPluginFromFile loads and registers a plugin from a local file path.
func (pp *PluginPersistence) RegisterPluginFromFile(filePath string) (*plugin.LoadedPlugin, error) {
	// Read the file for storage
	archiveData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin file: %w", err)
	}

	// Load the plugin
	loaded, err := pp.pluginLoader.LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	// Register scenarios
	if err := pp.pluginLoader.RegisterPluginScenarios(loaded); err != nil {
		return nil, fmt.Errorf("failed to register plugin scenarios: %w", err)
	}

	// Save to database
	if err := pp.SavePlugin(loaded, archiveData, filePath); err != nil {
		return nil, fmt.Errorf("failed to save plugin to database: %w", err)
	}

	return loaded, nil
}

// RegisterPluginFromLocal loads and registers a plugin from a local directory path.
func (pp *PluginPersistence) RegisterPluginFromLocal(dirPath string) (*plugin.LoadedPlugin, error) {
	// Load the plugin
	loaded, err := pp.pluginLoader.LoadFromLocalPath(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	// Register scenarios
	if err := pp.pluginLoader.RegisterPluginScenarios(loaded); err != nil {
		return nil, fmt.Errorf("failed to register plugin scenarios: %w", err)
	}

	// Save to database (no archive data for local plugins)
	if err := pp.SavePlugin(loaded, nil, dirPath); err != nil {
		return nil, fmt.Errorf("failed to save plugin to database: %w", err)
	}

	return loaded, nil
}

// RegisterPluginFromUpload loads and registers a plugin from uploaded bytes.
func (pp *PluginPersistence) RegisterPluginFromUpload(data []byte, filename string) (*plugin.LoadedPlugin, error) {
	// Determine if data is compressed based on magic bytes or filename
	compressed := isGzipped(data)

	// Load the plugin
	loaded, err := pp.pluginLoader.LoadFromBytes(data, compressed)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	loaded.SourceType = plugin.PluginSourceBytes

	// Register scenarios
	if err := pp.pluginLoader.RegisterPluginScenarios(loaded); err != nil {
		return nil, fmt.Errorf("failed to register plugin scenarios: %w", err)
	}

	// Ensure data is gzipped for storage
	archiveData := data
	if !compressed {
		var buf bytes.Buffer
		gzWriter := gzip.NewWriter(&buf)
		if _, err := gzWriter.Write(data); err != nil {
			return nil, fmt.Errorf("failed to compress archive data: %w", err)
		}
		if err := gzWriter.Close(); err != nil {
			return nil, fmt.Errorf("failed to close gzip writer: %w", err)
		}
		archiveData = buf.Bytes()
	}

	// Save to database
	if err := pp.SavePlugin(loaded, archiveData, "upload:"+filename); err != nil {
		return nil, fmt.Errorf("failed to save plugin to database: %w", err)
	}

	return loaded, nil
}

// isGzipped checks if data starts with gzip magic bytes.
func isGzipped(data []byte) bool {
	return len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b
}
