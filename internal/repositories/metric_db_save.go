package repositories

import (
	"context"
	"database/sql"
	"go-yandex-practicum-metrics/internal/types"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MetricDBSaveRepository struct {
	db *sql.DB
}

func NewMetricDBSaveRepository(db *sql.DB) *MetricDBSaveRepository {
	return &MetricDBSaveRepository{db: db}
}

// Save method to save a single metric (without transaction)
func (m *MetricDBSaveRepository) Save(ctx context.Context, metric *types.Metrics) bool {
	query := `INSERT INTO metrics (id, type, delta, value) 
			  VALUES ($1, $2, $3, $4) 
			  ON CONFLICT (id, type) 
			  DO UPDATE SET delta = EXCLUDED.delta, value = EXCLUDED.value`
	_, err := m.db.ExecContext(ctx, query, metric.ID, metric.Type, metric.Delta, metric.Value)
	return err == nil
}
