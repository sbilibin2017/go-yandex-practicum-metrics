package repositories

import (
	"context"
	"fmt"
	"go-yandex-practicum-metrics/internal/domain"
	"strings"
)

type DBScannerEngine interface {
	Scan(ctx context.Context, query string, args ...any) (map[any]any, bool)
}

type DBQuerierEngine interface {
	Query(ctx context.Context, query string, args ...any) (DBScannerEngine, bool)
}

type MetricDBFindBatchRepository struct {
	q DBQuerierEngine
	s DBScannerEngine
}

func NewMetricDBFindBatchRepository(q DBQuerierEngine, s DBScannerEngine) *MetricDBFindBatchRepository {
	return &MetricDBFindBatchRepository{q: q, s: s}
}

var findBatchQueryTemplate = "SELECT id, type, delta, value FROM metrics WHERE %s"

func buildFindQuery(filters []domain.MetricID) (string, []any) {
	var args []any
	conditions := []string{}
	if len(filters) > 0 {
		for i, filter := range filters {
			conditions = append(conditions, fmt.Sprintf("(id = $%d AND type = $%d)", i*2+1, i*2+2))
			args = append(args, filter.ID, filter.Type)
		}
	} else {
		args = []any{}
	}
	query := fmt.Sprintf(findBatchQueryTemplate, strings.Join(conditions, " OR "))
	return query, args
}

func (repo *MetricDBFindBatchRepository) FindBatch(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, bool) {
	query, args := buildFindQuery(filters)
	scannerEngine, ok := repo.q.Query(ctx, query, args...)
	if !ok {
		return nil, false
	}
	results := make(map[domain.MetricID]*domain.Metrics)
	scannedResults, ok := scannerEngine.Scan(ctx, query, args...)
	if !ok {
		return nil, false
	}
	for _, row := range scannedResults {
		row, _ := row.(map[string]any)
		var m domain.Metrics
		m.ID = row["id"].(string)
		m.Type = domain.MetricType(row["type"].(string))
		delta := row["delta"].(int64)
		m.Delta = &delta
		value := row["value"].(float64)
		m.Value = &value
		results[domain.MetricID{ID: m.ID, Type: m.Type}] = &m
	}
	return results, true
}
