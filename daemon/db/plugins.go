package db

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

// Plugin represents a database entity for storing plugin configuration and state.
// Maps to the "plugins" table with fields for source type, archive data,
// and scenario tracking.
type Plugin struct {
	ID                 int64  `db:"id"`                   // Unique identifier for the plugin
	Name               string `db:"name"`                 // Plugin name (unique)
	SourceType         string `db:"source_type"`          // Source type: "url", "file", "local", "bytes"
	SourcePath         string `db:"source_path"`          // URL, file path, or local directory
	MetadataName       string `db:"metadata_name"`        // Name from plugin.yaml
	MetadataBuildTime  string `db:"metadata_build_time"`  // Build time from plugin.yaml
	MetadataGitVersion string `db:"metadata_git_version"` // Git version from plugin.yaml
	ArchiveData        string `db:"archive_data"`         // Base64 encoded tar.gz (NULL for local)
	Scenarios          string `db:"scenarios"`            // JSON array of scenario names
	Enabled            bool   `db:"enabled"`              // Whether plugin is enabled
	LoadError          string `db:"load_error"`           // Error message if loading failed
	CreatedAt          int64  `db:"created_at"`           // Unix timestamp when created
	UpdatedAt          int64  `db:"updated_at"`           // Unix timestamp when last updated
}

// GetScenariosAsSlice parses the JSON scenarios field and returns a slice of strings.
func (p *Plugin) GetScenariosAsSlice() []string {
	if p.Scenarios == "" {
		return []string{}
	}

	var scenarios []string
	if err := json.Unmarshal([]byte(p.Scenarios), &scenarios); err != nil {
		return []string{}
	}

	return scenarios
}

// SetScenariosFromSlice converts a slice of scenario names to JSON for storage.
func (p *Plugin) SetScenariosFromSlice(scenarios []string) {
	if len(scenarios) == 0 {
		p.Scenarios = "[]"
		return
	}

	data, err := json.Marshal(scenarios)
	if err != nil {
		p.Scenarios = "[]"
		return
	}

	p.Scenarios = string(data)
}

// GetPlugin retrieves a single plugin by name from the database.
// Returns the plugin entity or an error if not found or database access fails.
func (d *Database) GetPlugin(name string) (*Plugin, error) {
	plugin := &Plugin{}
	err := d.ReaderDb.Get(plugin, `
		SELECT id, name, source_type, source_path, metadata_name, metadata_build_time,
		       metadata_git_version, archive_data, scenarios, enabled, load_error,
		       created_at, updated_at
		FROM plugins WHERE name = $1`, name)
	if err != nil {
		return nil, err
	}

	return plugin, nil
}

// GetPlugins retrieves all plugins from the database ordered by name.
// Returns a slice of plugin entities or an error if database access fails.
func (d *Database) GetPlugins() ([]*Plugin, error) {
	plugins := []*Plugin{}
	err := d.ReaderDb.Select(&plugins, `
		SELECT id, name, source_type, source_path, metadata_name, metadata_build_time,
		       metadata_git_version, archive_data, scenarios, enabled, load_error,
		       created_at, updated_at
		FROM plugins ORDER BY name ASC`)

	return plugins, err
}

// GetEnabledPlugins retrieves all enabled plugins from the database ordered by name.
// Returns a slice of plugin entities or an error if database access fails.
func (d *Database) GetEnabledPlugins() ([]*Plugin, error) {
	plugins := []*Plugin{}
	err := d.ReaderDb.Select(&plugins, `
		SELECT id, name, source_type, source_path, metadata_name, metadata_build_time,
		       metadata_git_version, archive_data, scenarios, enabled, load_error,
		       created_at, updated_at
		FROM plugins WHERE enabled = 1 ORDER BY name ASC`)

	return plugins, err
}

// InsertPlugin creates a new plugin record in the database within a transaction.
// Updates the plugin's ID field with the generated database ID after insertion.
// Returns an error if the insertion fails or transaction is invalid.
func (d *Database) InsertPlugin(tx *sqlx.Tx, plugin *Plugin) error {
	now := time.Now().Unix()
	plugin.CreatedAt = now
	plugin.UpdatedAt = now

	query := `
		INSERT INTO plugins (name, source_type, source_path, metadata_name, metadata_build_time,
		                     metadata_git_version, archive_data, scenarios, enabled, load_error,
		                     created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`

	return tx.QueryRow(query,
		plugin.Name,
		plugin.SourceType,
		plugin.SourcePath,
		plugin.MetadataName,
		plugin.MetadataBuildTime,
		plugin.MetadataGitVersion,
		plugin.ArchiveData,
		plugin.Scenarios,
		plugin.Enabled,
		plugin.LoadError,
		plugin.CreatedAt,
		plugin.UpdatedAt,
	).Scan(&plugin.ID)
}

// UpdatePlugin modifies an existing plugin record in the database within a transaction.
// Updates all mutable fields except ID and creation timestamp.
func (d *Database) UpdatePlugin(tx *sqlx.Tx, plugin *Plugin) error {
	plugin.UpdatedAt = time.Now().Unix()

	_, err := tx.Exec(`
		UPDATE plugins
		SET source_type = $1, source_path = $2, metadata_name = $3, metadata_build_time = $4,
		    metadata_git_version = $5, archive_data = $6, scenarios = $7, enabled = $8,
		    load_error = $9, updated_at = $10
		WHERE name = $11`,
		plugin.SourceType,
		plugin.SourcePath,
		plugin.MetadataName,
		plugin.MetadataBuildTime,
		plugin.MetadataGitVersion,
		plugin.ArchiveData,
		plugin.Scenarios,
		plugin.Enabled,
		plugin.LoadError,
		plugin.UpdatedAt,
		plugin.Name,
	)

	return err
}

// DeletePlugin removes a plugin record from the database within a transaction.
// Permanently deletes the plugin and all associated data.
// Returns an error if the deletion fails or the plugin doesn't exist.
func (d *Database) DeletePlugin(tx *sqlx.Tx, name string) error {
	_, err := tx.Exec(`DELETE FROM plugins WHERE name = $1`, name)
	return err
}

// UpdatePluginLoadError updates only the load_error field for a plugin.
// This is used to track loading failures without updating other fields.
func (d *Database) UpdatePluginLoadError(tx *sqlx.Tx, name string, loadError string) error {
	now := time.Now().Unix()

	_, err := tx.Exec(`
		UPDATE plugins
		SET load_error = $1, updated_at = $2
		WHERE name = $3`,
		loadError,
		now,
		name,
	)

	return err
}

// UpdatePluginEnabled updates only the enabled field for a plugin.
func (d *Database) UpdatePluginEnabled(tx *sqlx.Tx, name string, enabled bool) error {
	now := time.Now().Unix()

	_, err := tx.Exec(`
		UPDATE plugins
		SET enabled = $1, updated_at = $2
		WHERE name = $3`,
		enabled,
		now,
		name,
	)

	return err
}
