package repositories

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/types"
	"sync"
)

type MetricMemoryListRepository struct {
	data map[types.MetricID]*types.Metrics
	mu   sync.Mutex
}

func NewMetricMemoryListRepository(data map[types.MetricID]*types.Metrics) *MetricMemoryListRepository {
	return &MetricMemoryListRepository{data: data}
}

// ListAll method to return all metrics from memory
func (m *MetricMemoryListRepository) List(ctx context.Context) ([]*types.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data == nil {
		return nil, errors.New("data is not initialized") // Return error if data is nil
	}

	var result []*types.Metrics
	for _, metric := range m.data {
		result = append(result, metric)
	}

	return result, nil
}
