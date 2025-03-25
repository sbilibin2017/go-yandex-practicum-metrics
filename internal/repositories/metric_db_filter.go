package repositories

import (
	"context"
	"database/sql"
	"go-yandex-practicum-metrics/internal/types"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DBQuerier - интерфейс для выполнения запроса, который возвращает одну строку
type DBFilterQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type MetricDBFilterRepository struct {
	db DBFilterQuerier
}

func NewMetricDBFilterRepository(db DBFilterQuerier) *MetricDBFilterRepository {
	return &MetricDBFilterRepository{db: db}
}

func (m *MetricDBFilterRepository) Filter(
	ctx context.Context, filter types.MetricID,
) (*types.Metrics, error) {
	query := `SELECT id, type, delta, value FROM metrics WHERE id = $1 AND type = $2`
	var metric types.Metrics
	err := m.db.QueryRowContext(
		ctx, query, filter.ID, filter.Type,
	).Scan(
		&metric.ID, &metric.Type, &metric.Delta, &metric.Value,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No result found, returning nil for metrics and nil for error.
		}
		return nil, err // Return the error if it's not sql.ErrNoRows.
	}
	return &metric, nil // Return the metric and nil for error when successful.
}
