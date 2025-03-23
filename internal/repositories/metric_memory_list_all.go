package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/types"
	"sync"
)

type MetricMemoryListAllRepository struct {
	data map[types.MetricID]*types.Metrics
	mu   sync.Mutex
}

func NewMetricMemoryListAllRepository(data map[types.MetricID]*types.Metrics) *MetricMemoryListAllRepository {
	return &MetricMemoryListAllRepository{data: data}
}

// ListAll method to return all metrics from memory
func (m *MetricMemoryListAllRepository) ListAll(ctx context.Context) ([]*types.Metrics, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []*types.Metrics
	for _, metric := range m.data {
		result = append(result, metric)
	}
	return result, true
}
