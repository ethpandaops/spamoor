package db

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

// SpamoorState represents a key-value pair for storing application-wide configuration state.
// Maps to the "spamoor_state" table for persisting settings like scenario counters,
// first launch flags, and other global application state that survives restarts.
type SpamoorState struct {
	Key   string `db:"key"`   // Unique identifier for the state entry
	Value string `db:"value"` // JSON-serialized value for the state data
}

// GetSpamoorState retrieves and deserializes a state value by key.
// The returnValue parameter should be a pointer to the target type for unmarshaling.
// Returns the deserialized value or an error if the key doesn't exist or JSON parsing fails.
func (d *Database) GetSpamoorState(key string, returnValue interface{}) (interface{}, error) {
	entry := SpamoorState{}
	err := d.ReaderDb.Get(&entry, `SELECT key, value FROM spamoor_state WHERE key = $1`, key)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(entry.Value), returnValue)
	if err != nil {
		return nil, err
	}
	return returnValue, nil
}

// SetSpamoorState stores a state value by key with JSON serialization.
// If tx is nil, creates and manages its own transaction automatically.
// Uses INSERT OR REPLACE to handle both new entries and updates atomically.
// Returns an error if JSON marshaling or database operation fails.
func (d *Database) SetSpamoorState(tx *sqlx.Tx, key string, value interface{}) error {
	if tx == nil {
		return d.RunDBTransaction(func(tx *sqlx.Tx) error {
			return d.SetSpamoorState(tx, key, value)
		})
	}

	valueMarshal, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		`INSERT OR REPLACE INTO spamoor_state (key, value) VALUES ($1, $2)`,
		key, valueMarshal,
	)
	if err != nil {
		return err
	}
	return nil
}
