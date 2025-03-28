package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"sync"
)

type MetricMemorySaveBatchRepository struct {
	data map[domain.MetricID]*domain.Metrics
	mu   sync.Mutex
}

func NewMetricMemorySaveBatchRepository(
	data map[domain.MetricID]*domain.Metrics,
) *MetricMemorySaveBatchRepository {
	return &MetricMemorySaveBatchRepository{
		data: data,
	}
}

func (repo *MetricMemorySaveBatchRepository) Save(
	ctx context.Context, metrics []*domain.Metrics,
) bool {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, metric := range metrics {
		repo.data[domain.MetricID{ID: metric.ID, Type: metric.Type}] = metric
	}
	return true
}
