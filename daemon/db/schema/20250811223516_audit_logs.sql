-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Create audit_logs table for tracking all actions
CREATE TABLE IF NOT EXISTS audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_email TEXT NOT NULL, -- User who performed the action (from header)
    action_type TEXT NOT NULL, -- Type of action (e.g., 'spammer_create', 'client_update')
    entity_type TEXT NOT NULL, -- Type of entity ('spammer', 'client')
    entity_id TEXT NOT NULL, -- ID or identifier of the entity
    entity_name TEXT, -- Human-readable name of the entity
    diff TEXT, -- Unified diff of configuration changes
    metadata TEXT, -- Additional metadata (JSON)
    timestamp INTEGER NOT NULL -- Unix timestamp
);

-- Create index for faster lookups by timestamp
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp DESC);

-- Create index for filtering by user
CREATE INDEX idx_audit_logs_user ON audit_logs(user_email);

-- Create index for filtering by entity
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);

-- Create index for filtering by action type
CREATE INDEX idx_audit_logs_action ON audit_logs(action_type);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_entity;
DROP INDEX IF EXISTS idx_audit_logs_user;
DROP INDEX IF EXISTS idx_audit_logs_timestamp;
DROP TABLE IF EXISTS audit_logs;