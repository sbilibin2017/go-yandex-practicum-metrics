package engines

import (
	"context"
	"database/sql"
	"fmt"
	"go-yandex-practicum-metrics/internal/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DatabaseDSNGetter interface {
	GetDatabaseDSN() string
}

type DBEngine struct {
	*sql.DB
	tx *sql.Tx
}

func NewDBEngine(g DatabaseDSNGetter) (*DBEngine, error) {
	con, err := sql.Open("pgx", g.GetDatabaseDSN())
	if err != nil {
		logger.Error("Error opening database connection", "error", err)
		return nil, err
	}
	if err := con.Ping(); err != nil {
		logger.Error("Error pinging database", "error", err)
		return nil, err
	}
	return &DBEngine{DB: con}, nil
}

func (d *DBEngine) BeginTx(ctx context.Context) error {
	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("Error starting transaction", "error", err)
		return err
	}
	d.tx = tx
	logger.Info("Transaction started")
	return nil
}

func (d *DBEngine) Commit() error {
	if d.tx == nil {
		logger.Error("No transaction to commit")
		return fmt.Errorf("no transaction to commit")
	}
	err := d.tx.Commit()
	if err != nil {
		logger.Error("Error committing transaction", "error", err)
		return err
	}
	logger.Info("Transaction committed successfully")
	return nil
}

func (d *DBEngine) Rollback() error {
	if d.tx == nil {
		logger.Error("No transaction to rollback")
		return fmt.Errorf("no transaction to rollback")
	}
	err := d.tx.Rollback()
	if err != nil {
		logger.Error("Error rolling back transaction", "error", err)
		return err
	}
	logger.Info("Transaction rolled back successfully")
	return nil
}
