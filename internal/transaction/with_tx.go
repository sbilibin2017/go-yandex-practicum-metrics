package transaction

import (
	"context"
)

type Tx interface {
	Commit() error
	Rollback() error
}

type DB interface {
	Begin(ctx context.Context) (Tx, error)
}

type WithTx struct {
	db DB
}

func NewWithTx(db DB) *WithTx {
	return &WithTx{
		db: db,
	}
}

func (w *WithTx) Begin(ctx context.Context) (Tx, error) {
	if w.db != nil {
		return w.db.Begin(ctx)
	}
	return nil, nil
}

func (w *WithTx) Do(ctx context.Context, fn func(tx Tx) error) error {
	tx, err := w.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil && tx != nil {
			_ = tx.Rollback()
		}
	}()
	err = fn(tx)
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		return err
	}
	if tx != nil {
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
