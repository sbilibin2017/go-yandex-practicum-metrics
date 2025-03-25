package repositories

import (
	"context"
	"database/sql"
	"go-yandex-practicum-metrics/internal/types"
)

// DBExecutor - интерфейс для работы с базой данных
type DBListQuerier interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type MetricDBListRepository struct {
	db DBListQuerier
}

func NewMetricDBListRepository(db DBListQuerier) *MetricDBListRepository {
	return &MetricDBListRepository{db: db}
}

// ListAll method to return all metrics from the database as a slice
func (m *MetricDBListRepository) List(ctx context.Context) ([]*types.Metrics, error) {
	result := make(map[types.MetricID]*types.Metrics)
	query := `SELECT id, type, delta, value FROM metrics`
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err // Return the error if there is an issue executing the query.
	}
	defer rows.Close()

	for rows.Next() {
		var metric types.Metrics
		if err := rows.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value); err != nil {
			return nil, err // Return the error if there is an issue scanning a row.
		}
		metricID := types.MetricID{ID: metric.ID, Type: metric.Type}
		result[metricID] = &metric
	}

	if err := rows.Err(); err != nil {
		return nil, err // Return the error if there was an issue with the rows.
	}

	var resultSlice []*types.Metrics
	for _, metric := range result {
		resultSlice = append(resultSlice, metric)
	}

	return resultSlice, nil // Return the slice of metrics and nil for the error when everything is successful.
}
