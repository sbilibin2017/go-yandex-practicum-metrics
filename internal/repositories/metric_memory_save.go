package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/types"
	"sync"
)

type MetricMemorySaveRepository struct {
	data map[types.MetricID]*types.Metrics
	mu   sync.Mutex
}

func NewMetricMemorySaveRepository(data map[types.MetricID]*types.Metrics) *MetricMemorySaveRepository {
	return &MetricMemorySaveRepository{data: data}
}

// Save method for saving a single metric
func (m *MetricMemorySaveRepository) Save(ctx context.Context, metric *types.Metrics) error {
	if metric == nil {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[metric.MetricID] = metric
	return nil
}
