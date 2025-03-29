package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/logger"
	"strings"
)

type DBSave interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type MetricDBSaveRepository struct {
	db DBSave
}

var saveBatchQueryTemplate = `
INSERT INTO metrics (id, type, delta, value) 
VALUES %s
ON CONFLICT (id, type) 
DO UPDATE SET 
	delta = EXCLUDED.delta, 
	value = EXCLUDED.value
`

func buildSaveBatchQuery(metrics []*domain.Metrics) (string, []any) {
	var args []any
	placeholders := []string{}
	for i, metric := range metrics {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		args = append(args, metric.ID, metric.Type, metric.Delta, metric.Value)
	}
	query := fmt.Sprintf(saveBatchQueryTemplate, strings.Join(placeholders, ", "))
	return query, args
}

func NewMetricDBSaveRepository(db DBSave) *MetricDBSaveRepository {
	return &MetricDBSaveRepository{db: db}
}

func (repo *MetricDBSaveRepository) Save(
	ctx context.Context, metrics []*domain.Metrics,
) error {
	if len(metrics) == 0 {
		logger.Info("No metrics to save")
		return nil
	}
	query, args := buildSaveBatchQuery(metrics)
	logger.Info("Executing batch save query", "query", query, "args_count", len(args))
	_, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error("Error executing batch save query", "error", err)
		return err
	}
	logger.Info("Batch save completed successfully", "metrics_count", len(metrics))
	return nil
}
