package unitofwork

import (
	"context"
	"go-yandex-practicum-metrics/internal/logger"
)

type DB interface {
	BeginTx(ctx context.Context) error
	Commit() error
	Rollback() error
}

type UnitOfWork struct {
	db DB
}

func NewUnitOfWork(db DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (w *UnitOfWork) Do(ctx context.Context, fn func() error) error {
	logger.Info("Executing unit of work")
	var err error
	if w.db != nil {
		logger.Info("Beginning transaction")
		err = w.db.BeginTx(ctx)
		if err != nil {
			logger.Error("Error beginning transaction", "error", err)
			return err
		}
	}
	defer func() {
		if err != nil {
			logger.Info("Rolling back transaction due to error")
			_ = w.db.Rollback()
		}
	}()
	err = fn()
	if err != nil {
		logger.Error("Error during unit of work execution", "error", err)
		if err := w.db.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := w.db.Commit(); err != nil {
		logger.Error("Error committing transaction", "error", err)
		return err
	}
	logger.Info("Unit of work executed successfully")
	return nil
}
