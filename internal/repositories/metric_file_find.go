package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"sync"
)

type FileQuerierEngine interface {
	Query(ctx context.Context, query string, args ...any) ([]any, bool)
}

type MetricFileFindBatchRepository struct {
	engine FileQuerierEngine
	mu     sync.Mutex
}

func NewMetricFileFindBatchRepository(
	engine FileQuerierEngine,
) *MetricFileFindBatchRepository {
	return &MetricFileFindBatchRepository{engine: engine}
}

func (repo *MetricFileFindBatchRepository) Find(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	filterMap := make(map[domain.MetricID]struct{})
	for _, filter := range filters {
		filterMap[filter] = struct{}{}
	}
	allMetrics, ok := repo.engine.Query(ctx, "")
	if !ok {
		return nil, false
	}
	result := make(map[domain.MetricID]*domain.Metrics)
	for _, metric := range allMetrics {
		if m, ok := metric.(*domain.Metrics); ok {
			metricID := domain.MetricID{ID: m.ID, Type: m.Type}
			if _, exists := filterMap[metricID]; exists {
				result[metricID] = m
			}
		}
	}
	return result, true
}
