package repositories

import (
	"context"
	"database/sql"
	"go-yandex-practicum-metrics/internal/types"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DBExecutor — интерфейс для выполнения запросов с контекстом
type DBSaveExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type MetricDBSaveRepository struct {
	db DBSaveExecutor
}

func NewMetricDBSaveRepository(db DBSaveExecutor) *MetricDBSaveRepository {
	return &MetricDBSaveRepository{db: db}
}

// Save method to save a single metric (without transaction)
func (m *MetricDBSaveRepository) Save(ctx context.Context, metric *types.Metrics) error {
	query := `INSERT INTO metrics (id, type, delta, value) 
			  VALUES ($1, $2, $3, $4) 
			  ON CONFLICT (id, type) 
			  DO UPDATE SET delta = EXCLUDED.delta, value = EXCLUDED.value`
	_, err := m.db.ExecContext(ctx, query, metric.ID, metric.Type, metric.Delta, metric.Value)
	if err != nil {
		return err // Return the error if there was an issue executing the query.
	}
	return nil // Return nil when the operation is successful.
}
