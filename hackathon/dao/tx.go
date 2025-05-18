package dao

import "database/sql"

func DoInTx(db *sql.DB, f func(tx *sql.Tx) (interface{}, error)) (interface{}, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	v, err := f(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return v, nil
}
