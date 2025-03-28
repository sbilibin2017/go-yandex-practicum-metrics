package repositories

import (
	"context"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
	"strings"
)

type DBExecutorEngine interface {
	Execute(ctx context.Context, query string, args ...any) bool
}

type MetricDBSaveBatchRepository struct {
	engine DBExecutorEngine
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

func NewMetricDBSaveBatchRepository(engine DBExecutorEngine) *MetricDBSaveBatchRepository {
	return &MetricDBSaveBatchRepository{engine: engine}
}

func (repo *MetricDBSaveBatchRepository) SaveBatch(
	ctx context.Context, metrics []*domain.Metrics,
) bool {
	query, args := buildSaveBatchQuery(metrics)
	return repo.engine.Execute(ctx, query, args...)
}
