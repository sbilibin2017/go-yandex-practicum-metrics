package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/types"
)

type MetricMemoryFilterRepository struct {
	data map[types.MetricID]*types.Metrics
}

func NewMetricMemoryFilterRepository(data map[types.MetricID]*types.Metrics) *MetricMemoryFilterRepository {
	return &MetricMemoryFilterRepository{data: data}
}

// Filter method to return a single metric based on the filter
func (m *MetricMemoryFilterRepository) Filter(ctx context.Context, filter types.MetricID) (*types.Metrics, error) {
	if metric, exists := m.data[filter]; exists {
		return metric, nil
	}
	return nil, nil
}
