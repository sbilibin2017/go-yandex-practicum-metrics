package repositories

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/logger"
	"io"
	"sync"
)

type MetricFileSaveRepository struct {
	writer io.Writer
	mu     sync.Mutex
}

func NewMetricFileSaveRepository(
	writer io.Writer,
) *MetricFileSaveRepository {
	return &MetricFileSaveRepository{writer: writer}
}

func (repo *MetricFileSaveRepository) Save(
	ctx context.Context, metrics []*domain.Metrics,
) error {
	logger.Info("Saving metrics to file", "metrics_count", len(metrics))
	repo.mu.Lock()
	defer repo.mu.Unlock()
	encoder := json.NewEncoder(repo.writer)
	for _, metric := range metrics {
		if err := encoder.Encode(metric); err != nil {
			logger.Error("Error encoding metric", "error", err)
			return err
		}
	}
	logger.Info("Metrics saved successfully", "metrics_count", len(metrics))
	return nil
}
