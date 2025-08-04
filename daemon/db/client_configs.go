package db

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

/*
CREATE TABLE IF NOT EXISTS "client_configs"
(
    "rpc_url" TEXT NOT NULL,
    "name" TEXT NOT NULL DEFAULT '',
    "tags" TEXT NOT NULL DEFAULT '',
    "enabled" INTEGER NOT NULL DEFAULT 1,
    "created_at" INTEGER NOT NULL,
    "updated_at" INTEGER NOT NULL,
    CONSTRAINT "client_configs_pkey" PRIMARY KEY("rpc_url")
);
*/

// ClientConfig represents a database entity for storing client configuration.
// Maps to the "client_configs" table with RPC URL as primary key and persistent
// settings for name, tags, and enabled state.
type ClientConfig struct {
	RpcUrl     string `db:"rpc_url"`     // RPC URL used as primary key identifier
	Name       string `db:"name"`        // Human-readable name for the client
	Tags       string `db:"tags"`        // Comma-separated tags for client categorization
	ClientType string `db:"client_type"` // Type of client: 'client' or 'builder'
	Enabled    bool   `db:"enabled"`     // Whether the client is enabled for use
	CreatedAt  int64  `db:"created_at"`  // Unix timestamp when the config was created
	UpdatedAt  int64  `db:"updated_at"`  // Unix timestamp when the config was last modified
}

// GetClientConfig retrieves a single client config by RPC URL from the database.
// Returns the client config entity or an error if not found or database access fails.
func (d *Database) GetClientConfig(rpcUrl string) (*ClientConfig, error) {
	config := &ClientConfig{}
	err := d.ReaderDb.Get(config, `
		SELECT rpc_url, name, tags, client_type, enabled, created_at, updated_at
		FROM client_configs WHERE rpc_url = $1`, rpcUrl)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// GetClientConfigs retrieves all client configs from the database ordered by creation time.
// Returns a slice of client config entities or an error if database access fails.
func (d *Database) GetClientConfigs() ([]*ClientConfig, error) {
	configs := []*ClientConfig{}
	err := d.ReaderDb.Select(&configs, `
		SELECT rpc_url, name, tags, client_type, enabled, created_at, updated_at 
		FROM client_configs ORDER BY created_at ASC`)
	return configs, err
}

// InsertClientConfig creates a new client config record in the database within a transaction.
// Sets creation and update timestamps automatically.
// Returns an error if the insertion fails or transaction is invalid.
func (d *Database) InsertClientConfig(tx *sqlx.Tx, config *ClientConfig) error {
	now := time.Now().Unix()
	config.CreatedAt = now
	config.UpdatedAt = now

	_, err := tx.Exec(`
		INSERT INTO client_configs (rpc_url, name, tags, client_type, enabled, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		config.RpcUrl,
		config.Name,
		config.Tags,
		config.ClientType,
		config.Enabled,
		config.CreatedAt,
		config.UpdatedAt,
	)
	return err
}

// UpdateClientConfig modifies an existing client config record in the database within a transaction.
// Updates name, tags, enabled state, and update timestamp. RPC URL and creation time remain unchanged.
// Returns an error if the update fails or the config doesn't exist.
func (d *Database) UpdateClientConfig(tx *sqlx.Tx, config *ClientConfig) error {
	config.UpdatedAt = time.Now().Unix()

	_, err := tx.Exec(`
		UPDATE client_configs 
		SET name = $1, tags = $2, client_type = $3, enabled = $4, updated_at = $5
		WHERE rpc_url = $6`,
		config.Name,
		config.Tags,
		config.ClientType,
		config.Enabled,
		config.UpdatedAt,
		config.RpcUrl,
	)
	return err
}

// UpsertClientConfig creates or updates a client config record in the database within a transaction.
// If the RPC URL exists, updates the record; otherwise creates a new one.
// Returns an error if the operation fails.
func (d *Database) UpsertClientConfig(tx *sqlx.Tx, config *ClientConfig) error {
	now := time.Now().Unix()
	config.UpdatedAt = now

	// Try to update first
	result, err := tx.Exec(`
		UPDATE client_configs 
		SET name = $1, tags = $2, client_type = $3, enabled = $4, updated_at = $5
		WHERE rpc_url = $6`,
		config.Name,
		config.Tags,
		config.ClientType,
		config.Enabled,
		config.UpdatedAt,
		config.RpcUrl,
	)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If no rows were updated, insert a new record
	if rowsAffected == 0 {
		config.CreatedAt = now
		_, err = tx.Exec(`
			INSERT INTO client_configs (rpc_url, name, tags, client_type, enabled, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			config.RpcUrl,
			config.Name,
			config.Tags,
			config.ClientType,
			config.Enabled,
			config.CreatedAt,
			config.UpdatedAt,
		)
	}

	return err
}

// DeleteClientConfig removes a client config record from the database within a transaction.
// Permanently deletes the client config for the specified RPC URL.
// Returns an error if the deletion fails or the config doesn't exist.
func (d *Database) DeleteClientConfig(tx *sqlx.Tx, rpcUrl string) error {
	_, err := tx.Exec(`DELETE FROM client_configs WHERE rpc_url = $1`, rpcUrl)
	return err
}

// GetTagsAsSlice parses the comma-separated tags string into a slice of strings.
// Returns an empty slice if tags is empty or contains only whitespace.
func (c *ClientConfig) GetTagsAsSlice() []string {
	if c.Tags == "" {
		return []string{}
	}

	tags := strings.Split(c.Tags, ",")
	result := make([]string, 0, len(tags))

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}

	return result
}

// SetTagsFromSlice converts a slice of tags into a comma-separated string.
// Trims whitespace from each tag and filters out empty strings.
func (c *ClientConfig) SetTagsFromSlice(tags []string) {
	filtered := make([]string, 0, len(tags))

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			filtered = append(filtered, tag)
		}
	}

	c.Tags = strings.Join(filtered, ",")
}
