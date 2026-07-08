package mysql

import (
	"context"
	"database/sql"
	"hackathon/repository"
)

type TxExecutor struct{}

var _ repository.TransactionExecutor = (*TxExecutor)(nil)

func NewTxExecutor() *TxExecutor {
	return &TxExecutor{}
}

func (te *TxExecutor) DoInTx(ctx context.Context, db *sql.DB, f func(ctx context.Context, tx *sql.Tx) (interface{}, error)) (interface{}, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	v, err := f(ctx, tx)
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
