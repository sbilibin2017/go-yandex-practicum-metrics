package engines

import (
	"context"
	"database/sql"
)

type DBTransactionEngine struct {
	tx *sql.Tx
	db *sql.DB
}

func (e *DBTransactionEngine) Begin(ctx context.Context) bool {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return false
	}
	e.tx = tx
	return true
}

func (e *DBTransactionEngine) Commit() bool {
	if e.tx == nil {
		return false
	}
	err := e.tx.Commit()
	return err == nil
}

func (e *DBTransactionEngine) Rollback() bool {
	if e.tx == nil {
		return false
	}
	err := e.tx.Rollback()
	return err == nil
}
