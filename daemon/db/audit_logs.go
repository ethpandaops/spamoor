package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// AuditLog represents an audit log entry in the database
type AuditLog struct {
	ID         int64  `db:"id"`
	UserEmail  string `db:"user_email"`
	ActionType string `db:"action_type"`
	EntityType string `db:"entity_type"`
	EntityID   string `db:"entity_id"`
	EntityName string `db:"entity_name"`
	Diff       string `db:"diff"`
	Metadata   string `db:"metadata"`
	Timestamp  int64  `db:"timestamp"`
}

// AuditActionType defines the types of actions that can be logged
type AuditActionType string

const (
	// Spammer actions
	AuditActionSpammerCreate  AuditActionType = "spammer_create"
	AuditActionSpammerUpdate  AuditActionType = "spammer_update"
	AuditActionSpammerDelete  AuditActionType = "spammer_delete"
	AuditActionSpammerStart   AuditActionType = "spammer_start"
	AuditActionSpammerPause   AuditActionType = "spammer_pause"
	AuditActionSpammerReclaim AuditActionType = "spammer_reclaim"

	// Client actions
	AuditActionClientUpdate      AuditActionType = "client_update"
	AuditActionClientGroupUpdate AuditActionType = "client_group_update"
	AuditActionClientNameUpdate  AuditActionType = "client_name_update"
	AuditActionClientTypeUpdate  AuditActionType = "client_type_update"
	AuditActionClientToggle      AuditActionType = "client_toggle"

	// Import/Export actions
	AuditActionSpammersImport AuditActionType = "spammers_import"
	AuditActionSpammersExport AuditActionType = "spammers_export"

	// Root wallet actions
	AuditActionRootWalletSend AuditActionType = "root_wallet_send"

	// Plugin actions
	AuditActionPluginRegister AuditActionType = "plugin_register"
	AuditActionPluginDelete   AuditActionType = "plugin_delete"
	AuditActionPluginReload   AuditActionType = "plugin_reload"
)

// AuditEntityType defines the types of entities that can be logged
type AuditEntityType string

const (
	AuditEntitySpammer    AuditEntityType = "spammer"
	AuditEntityClient     AuditEntityType = "client"
	AuditEntitySystem     AuditEntityType = "system"
	AuditEntityRootWallet AuditEntityType = "root_wallet"
	AuditEntityPlugin     AuditEntityType = "plugin"
)

// AuditMetadata represents additional metadata for an audit log entry
type AuditMetadata map[string]interface{}

// InsertAuditLog creates a new audit log entry in the database within a transaction
func (d *Database) InsertAuditLog(tx *sqlx.Tx, log *AuditLog) error {
	if log.Timestamp == 0 {
		log.Timestamp = time.Now().Unix()
	}

	query := `
		INSERT INTO audit_logs (user_email, action_type, entity_type, entity_id, entity_name, diff, metadata, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	return tx.QueryRow(query,
		log.UserEmail,
		log.ActionType,
		log.EntityType,
		log.EntityID,
		log.EntityName,
		log.Diff,
		log.Metadata,
		log.Timestamp,
	).Scan(&log.ID)
}

// GetAuditLogs retrieves audit logs with optional filtering
func (d *Database) GetAuditLogs(filters AuditLogFilters) ([]*AuditLog, error) {
	query := `SELECT id, user_email, action_type, entity_type, entity_id, entity_name, diff, metadata, timestamp FROM audit_logs WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	// Add filters
	if filters.UserEmail != "" {
		argCount++
		query += fmt.Sprintf(" AND user_email = $%d", argCount)
		args = append(args, filters.UserEmail)
	}

	if filters.ActionType != "" {
		argCount++
		query += fmt.Sprintf(" AND action_type = $%d", argCount)
		args = append(args, filters.ActionType)
	}

	if filters.EntityType != "" {
		argCount++
		query += fmt.Sprintf(" AND entity_type = $%d", argCount)
		args = append(args, filters.EntityType)
	}

	if filters.EntityID != "" {
		argCount++
		query += fmt.Sprintf(" AND entity_id = $%d", argCount)
		args = append(args, filters.EntityID)
	}

	if filters.StartTime > 0 {
		argCount++
		query += fmt.Sprintf(" AND timestamp >= $%d", argCount)
		args = append(args, filters.StartTime)
	}

	if filters.EndTime > 0 {
		argCount++
		query += fmt.Sprintf(" AND timestamp <= $%d", argCount)
		args = append(args, filters.EndTime)
	}

	// Add ordering and limit
	query += " ORDER BY timestamp DESC"

	if filters.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.Limit)
	}

	if filters.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filters.Offset)
	}

	var logs []*AuditLog
	err := d.ReaderDb.Select(&logs, query, args...)
	return logs, err
}

// GetAuditLogByID retrieves a single audit log entry by ID
func (d *Database) GetAuditLogByID(id int64) (*AuditLog, error) {
	var log AuditLog
	err := d.ReaderDb.Get(&log, `
		SELECT id, user_email, action_type, entity_type, entity_id, entity_name, diff, metadata, timestamp
		FROM audit_logs
		WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// CountAuditLogs returns the total count of audit logs matching the filters
func (d *Database) CountAuditLogs(filters AuditLogFilters) (int64, error) {
	query := `SELECT COUNT(*) FROM audit_logs WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	// Add filters (same as GetAuditLogs but without limit/offset)
	if filters.UserEmail != "" {
		argCount++
		query += fmt.Sprintf(" AND user_email = $%d", argCount)
		args = append(args, filters.UserEmail)
	}

	if filters.ActionType != "" {
		argCount++
		query += fmt.Sprintf(" AND action_type = $%d", argCount)
		args = append(args, filters.ActionType)
	}

	if filters.EntityType != "" {
		argCount++
		query += fmt.Sprintf(" AND entity_type = $%d", argCount)
		args = append(args, filters.EntityType)
	}

	if filters.EntityID != "" {
		argCount++
		query += fmt.Sprintf(" AND entity_id = $%d", argCount)
		args = append(args, filters.EntityID)
	}

	if filters.StartTime > 0 {
		argCount++
		query += fmt.Sprintf(" AND timestamp >= $%d", argCount)
		args = append(args, filters.StartTime)
	}

	if filters.EndTime > 0 {
		argCount++
		query += fmt.Sprintf(" AND timestamp <= $%d", argCount)
		args = append(args, filters.EndTime)
	}

	var count int64
	err := d.ReaderDb.Get(&count, query, args...)
	return count, err
}

// AuditLogFilters represents filters for querying audit logs
type AuditLogFilters struct {
	UserEmail  string
	ActionType string
	EntityType string
	EntityID   string
	StartTime  int64
	EndTime    int64
	Limit      int
	Offset     int
}

// Helper function to create metadata JSON
func MarshalAuditMetadata(metadata AuditMetadata) string {
	if metadata == nil {
		return "{}"
	}
	data, err := json.Marshal(metadata)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// Helper function to parse metadata JSON
func UnmarshalAuditMetadata(metadataStr string) AuditMetadata {
	var metadata AuditMetadata
	if metadataStr == "" {
		return metadata
	}
	json.Unmarshal([]byte(metadataStr), &metadata)
	return metadata
}
