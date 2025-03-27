package db

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type SpamoorState struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

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
