package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/logger"
	"sync"
)

type MetricMemorySaveRepository struct {
	data map[domain.MetricID]*domain.Metrics
	mu   sync.Mutex
}

func NewMetricMemorySaveRepository(
	data map[domain.MetricID]*domain.Metrics,
) *MetricMemorySaveRepository {
	return &MetricMemorySaveRepository{
		data: data,
	}
}

func (repo *MetricMemorySaveRepository) Save(
	ctx context.Context, metrics []*domain.Metrics,
) error {
	logger.Info("Saving metrics to memory")
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, metric := range metrics {
		repo.data[domain.MetricID{ID: metric.ID, Type: metric.Type}] = metric
		logger.Info("Metric saved", "id", metric.ID, "type", metric.Type)
	}
	logger.Info("Metrics saved successfully", "metrics_count", len(metrics))
	return nil
}
