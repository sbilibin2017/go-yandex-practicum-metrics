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

func (repo *MetricMemoryFindBatchRepository) FindBatch(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	result := make(map[domain.MetricID]*domain.Metrics)
	for _, metricID := range filters {
		if metric, exists := repo.data[metricID]; exists {
			result[metricID] = metric
		}
	}
	return result, nil
}
