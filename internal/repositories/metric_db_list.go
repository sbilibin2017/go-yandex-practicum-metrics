package repositories

import (
	"context"
	"database/sql"
	"go-yandex-practicum-metrics/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MetricDBListRepository struct {
	db *sql.DB
}

func NewMetricDBListRepository(db *sql.DB) *MetricDBListRepository {
	return &MetricDBListRepository{db: db}
}

var listAllMetricsQuery = "SELECT id, type, delta, value FROM metrics"

func (repo *MetricDBListRepository) List(
	ctx context.Context,
) (map[domain.MetricID]*domain.Metrics, error) {
	rows, err := repo.db.QueryContext(ctx, listAllMetricsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make(map[domain.MetricID]*domain.Metrics)
	for rows.Next() {
		var metric domain.Metrics
		if err := rows.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value); err != nil {
			return nil, err
		}
		results[domain.MetricID{ID: metric.ID, Type: metric.Type}] = &metric
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
