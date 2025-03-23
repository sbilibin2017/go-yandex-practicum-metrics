package repositories

import (
	"context"
	"database/sql"
	"go-yandex-practicum-metrics/internal/types"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

type MetricDBFilterRepository struct {
	db *sql.DB
}

func NewMetricDBFilterRepository(db *sql.DB) *MetricDBFilterRepository {
	// The connection is assumed to be already open
	return &MetricDBFilterRepository{db: db}
}

// Filter method to retrieve a metric based on the given MetricID filter
func (m *MetricDBFilterRepository) Filter(ctx context.Context, filter types.MetricID) (*types.Metrics, bool) {
	query := `SELECT id, type, delta, value FROM metrics WHERE id = $1 AND type = $2`
	var metric types.Metrics
	err := m.db.QueryRowContext(
		ctx, query, filter.ID, filter.Type,
	).Scan(
		&metric.ID, &metric.Type, &metric.Delta, &metric.Value,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		return nil, false
	}
	return &metric, true
}
