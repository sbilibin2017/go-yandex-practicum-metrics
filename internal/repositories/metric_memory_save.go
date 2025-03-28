package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"sync"
)

type MemoryExecutorEngine[K comparable, V any] interface {
	Execute(ctx context.Context, query string, args ...any) bool
}

type MetricMemorySaveBatchRepository struct {
	engine MemoryExecutorEngine[domain.MetricID, *domain.Metrics]
	mu     sync.Mutex
}

func NewMetricMemorySaveBatchRepository(
	engine MemoryExecutorEngine[domain.MetricID, *domain.Metrics],
) *MetricMemorySaveBatchRepository {
	return &MetricMemorySaveBatchRepository{engine: engine}
}

func (repo *MetricMemorySaveBatchRepository) Save(
	ctx context.Context, metrics []*domain.Metrics,
) bool {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, metric := range metrics {
		ok := repo.engine.Execute(ctx, "", domain.MetricID{ID: metric.ID, Type: metric.Type}, metric)
		if !ok {
			return false
		}
	}
	return true
}
