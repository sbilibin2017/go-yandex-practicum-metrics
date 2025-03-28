package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"sync"
)

type MetricMemoryFindBatchRepository struct {
	data map[domain.MetricID]*domain.Metrics
	mu   sync.Mutex
}

func NewMetricMemoryFindBatchRepository(
	data map[domain.MetricID]*domain.Metrics,
) *MetricMemoryFindBatchRepository {
	return &MetricMemoryFindBatchRepository{data: data}
}

func (repo *MetricMemoryFindBatchRepository) Find(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	filterMap := make(map[domain.MetricID]struct{})
	for _, filter := range filters {
		filterMap[filter] = struct{}{}
	}
	result := make(map[domain.MetricID]*domain.Metrics)
	for metricID, metric := range repo.data {
		if _, exists := filterMap[metricID]; exists {
			result[metricID] = metric
		}
	}
	return result, true
}
