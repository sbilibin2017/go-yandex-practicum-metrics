package engines

import (
	"context"
	"database/sql"
)

type DBQueryEngine struct {
	db *sql.DB
}

func NewDBQueryEngine(db *sql.DB) (*DBQueryEngine, error) {
	return &DBQueryEngine{
		db: db,
	}, nil
}

func (e *DBQueryEngine) Query(ctx context.Context, query string, args ...any) ([]any, bool) {
	var results []any
	rows, err := e.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, false
	}
	defer rows.Close()
	for rows.Next() {
		var result any
		if err := rows.Scan(&result); err != nil {
			return nil, false
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, false
	}
	return results, true
}
