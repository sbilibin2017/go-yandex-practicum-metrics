package transaction

import (
	"context"
	"database/sql"
)

type contextKey string

var dbTransactionKey = contextKey("db_transaction_key")

type DB interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) // Меняем *sql.Tx на Tx
}

type Tx interface {
	Commit() error
	Rollback() error
}

type DBTransactionEngine struct {
	db DB
}

func (e *DBTransactionEngine) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, dbTransactionKey, tx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = fn(ctx)
	return err
}
