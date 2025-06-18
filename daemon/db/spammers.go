package db

import (
	"github.com/jmoiron/sqlx"
)

/*
CREATE TABLE IF NOT EXISTS "spammers"
(
    "id" INTEGER NOT NULL,
    "scenario" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "config" TEXT NOT NULL,
    "status" INTEGER NOT NULL DEFAULT 0,
    "created_at" INTEGER NOT NULL,
    "state" TEXT NOT NULL,
    CONSTRAINT "spammers_pkey" PRIMARY KEY("id")
);
*/

// Spammer represents a database entity for storing spammer configuration and state.
// Maps to the "spammers" table with fields for scenario definition, execution status,
// and persistent state management. The status field tracks execution state using SpammerStatus constants.
type Spammer struct {
	ID          int64  `db:"id"`          // Unique identifier for the spammer instance
	Scenario    string `db:"scenario"`    // Name of the scenario to execute
	Name        string `db:"name"`        // Human-readable name for the spammer
	Description string `db:"description"` // Detailed description of what this spammer does
	Config      string `db:"config"`      // YAML configuration for scenario and wallet settings
	Status      int    `db:"status"`      // Current execution status (see SpammerStatus constants)
	CreatedAt   int64  `db:"created_at"`  // Unix timestamp when the spammer was created
	State       string `db:"state"`       // Persistent state data for scenario execution
}

// GetSpammer retrieves a single spammer by ID from the database.
// Returns the spammer entity or an error if not found or database access fails.
func (d *Database) GetSpammer(id int64) (*Spammer, error) {
	spammer := &Spammer{}
	err := d.ReaderDb.Get(spammer, `
		SELECT id, scenario, name, description, config, status, created_at, state 
		FROM spammers WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return spammer, nil
}

// GetSpammers retrieves all spammers from the database ordered by creation time (newest first).
// Returns a slice of spammer entities or an error if database access fails.
func (d *Database) GetSpammers() ([]*Spammer, error) {
	spammers := []*Spammer{}
	err := d.ReaderDb.Select(&spammers, `SELECT id, scenario, name, description, config, status, created_at, state FROM spammers ORDER BY created_at DESC`)
	return spammers, err
}

// InsertSpammer creates a new spammer record in the database within a transaction.
// Updates the spammer's ID field with the generated database ID after insertion.
// Returns an error if the insertion fails or transaction is invalid.
func (d *Database) InsertSpammer(tx *sqlx.Tx, spammer *Spammer) error {
	query := `
		INSERT INTO spammers (id, scenario, name, description, config, status, created_at, state)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	return tx.QueryRow(query,
		spammer.ID,
		spammer.Scenario,
		spammer.Name,
		spammer.Description,
		spammer.Config,
		spammer.Status,
		spammer.CreatedAt,
		spammer.State,
	).Scan(&spammer.ID)
}

// UpdateSpammer modifies an existing spammer record in the database within a transaction.
// Updates all mutable fields: name, description, config, status, and state.
// The ID and creation timestamp remain unchanged.
func (d *Database) UpdateSpammer(tx *sqlx.Tx, spammer *Spammer) error {
	_, err := tx.Exec(`
		UPDATE spammers 
		SET name = $1, description = $2, config = $3, status = $4, state = $5
		WHERE id = $6`,
		spammer.Name,
		spammer.Description,
		spammer.Config,
		spammer.Status,
		spammer.State,
		spammer.ID,
	)
	return err
}

// DeleteSpammer removes a spammer record from the database within a transaction.
// Permanently deletes the spammer and all associated data.
// Returns an error if the deletion fails or the spammer doesn't exist.
func (d *Database) DeleteSpammer(tx *sqlx.Tx, id int64) error {
	_, err := tx.Exec(`DELETE FROM spammers WHERE id = $1`, id)
	return err
}
