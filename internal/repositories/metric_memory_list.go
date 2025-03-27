package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"sync"
)

type MetricMemoryListAllRepository struct {
	data map[domain.MetricID]*domain.Metrics
	mu   sync.Mutex
}

func NewMetricMemoryListAllRepository(
	data map[domain.MetricID]*domain.Metrics,
) *MetricMemoryListAllRepository {
	return &MetricMemoryListAllRepository{data: data}
}

func (repo *MetricMemoryListAllRepository) List(
	ctx context.Context,
) (map[domain.MetricID]*domain.Metrics, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	result := make(map[domain.MetricID]*domain.Metrics)
	for metricID, metric := range repo.data {
		result[metricID] = metric
	}
	return result, nil
}
