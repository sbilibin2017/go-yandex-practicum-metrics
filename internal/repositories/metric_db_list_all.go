package repositories

import (
	"context"
	"database/sql"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricDBListAllRepository struct {
	db *sql.DB
}

func NewMetricDBListAllRepository(db *sql.DB) *MetricDBListAllRepository {
	return &MetricDBListAllRepository{db: db}
}

// ListAll method to return all metrics from the database as a slice
func (m *MetricDBListAllRepository) ListAll(ctx context.Context) ([]*types.Metrics, bool) {
	result := make(map[types.MetricID]*types.Metrics)
	query := `SELECT id, type, delta, value FROM metrics`
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, false
	}
	defer rows.Close()
	for rows.Next() {
		var metric types.Metrics
		if err := rows.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value); err != nil {
			return nil, false
		}
		metricID := types.MetricID{ID: metric.ID, Type: metric.Type}
		result[metricID] = &metric
	}
	if err := rows.Err(); err != nil {
		return nil, false
	}
	var resultSlice []*types.Metrics
	for _, metric := range result {
		resultSlice = append(resultSlice, metric)
	}
	return resultSlice, true
}
