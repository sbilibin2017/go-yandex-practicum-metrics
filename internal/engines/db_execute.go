package engines

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBExecuteEngine struct {
	db *sql.DB
}

func NewDBExecuteEngine(db *sql.DB) (*DBExecuteEngine, error) {
	return &DBExecuteEngine{
		db: db,
	}, nil
}

func (e *DBExecuteEngine) Execute(ctx context.Context, query string, args ...any) bool {
	_, err := e.db.ExecContext(ctx, query, args...)
	return err == nil
}
