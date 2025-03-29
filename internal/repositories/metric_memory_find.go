package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/logger"
	"sync"
)

type MetricMemoryFindRepository struct {
	data map[domain.MetricID]*domain.Metrics
	mu   sync.Mutex
}

func NewMetricMemoryFindRepository(
	data map[domain.MetricID]*domain.Metrics,
) *MetricMemoryFindRepository {
	return &MetricMemoryFindRepository{data: data}
}

func (repo *MetricMemoryFindRepository) Find(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, error) {
	logger.Info("Finding metrics from memory", "filters_count", len(filters))
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
			logger.Info("Metric found", "id", metricID.ID, "type", metricID.Type)
		}
	}
	logger.Info("Metrics found", "result_count", len(result))
	return result, nil
}
