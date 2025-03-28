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

type Spammer struct {
	ID          int64  `db:"id"`
	Scenario    string `db:"scenario"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Config      string `db:"config"`
	Status      int    `db:"status"`
	CreatedAt   int64  `db:"created_at"`
	State       string `db:"state"`
}

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

func (d *Database) GetSpammers() ([]*Spammer, error) {
	spammers := []*Spammer{}
	err := d.ReaderDb.Select(&spammers, `SELECT id, scenario, name, description, config, status, created_at, state FROM spammers ORDER BY created_at DESC`)
	return spammers, err
}

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

func (d *Database) DeleteSpammer(tx *sqlx.Tx, id int64) error {
	_, err := tx.Exec(`DELETE FROM spammers WHERE id = $1`, id)
	return err
}
