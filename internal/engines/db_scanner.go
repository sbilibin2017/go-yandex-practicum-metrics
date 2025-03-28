package engines

import (
	"context"
	"database/sql"
)

type DBScannerEngine struct {
	db *sql.DB
}

func NewDBScannerEngine(db *sql.DB) (*DBScannerEngine, error) {
	return &DBScannerEngine{db: db}, nil
}

func (e *DBScannerEngine) Scan(ctx context.Context, query string, args ...any) (map[any]any, bool) {
	rows, err := e.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, false
	}
	defer rows.Close()
	results := make(map[any]any)
	columns, err := rows.Columns()
	if err != nil {
		return nil, false
	}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, false
		}
		rowData := make(map[string]any)
		for i, colName := range columns {
			rowData[colName] = values[i]
		}
		results[len(results)] = rowData
	}
	if err := rows.Err(); err != nil {
		return nil, false
	}
	return results, true
}
