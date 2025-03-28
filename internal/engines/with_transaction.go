package engines

import (
	"context"
)

type Transaction interface {
	Begin(ctx context.Context) error
	Commit() error
	Rollback() error
}

type WithTransaction struct {
	transaction Transaction
}

func NewWithTransaction(transaction Transaction) *WithTransaction {
	return &WithTransaction{
		transaction: transaction,
	}
}

func (wt *WithTransaction) Do(ctx context.Context, f func() error) bool {
	if err := wt.transaction.Begin(ctx); err != nil {
		return false
	}
	defer func() {
		if p := recover(); p != nil {
			wt.transaction.Rollback()
			panic(p)
		}
	}()
	if err := f(); err != nil {
		wt.transaction.Rollback()
		return false
	}
	if err := wt.transaction.Commit(); err != nil {
		return false
	}
	return true
}
