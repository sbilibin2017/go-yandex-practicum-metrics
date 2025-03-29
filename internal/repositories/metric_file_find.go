package repositories

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/logger"
	"io"
	"sync"
)

type ReaderSeaker interface {
	io.Reader
	io.Seeker
}

type MetricFileFindRepository struct {
	file ReaderSeaker
	mu   sync.Mutex
}

func NewMetricFileFindRepository(
	file ReaderSeaker,
) *MetricFileFindRepository {
	return &MetricFileFindRepository{
		file: file,
	}
}

func (repo *MetricFileFindRepository) Find(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, error) {
	logger.Info("Finding metrics from file", "filters_count", len(filters))
	repo.mu.Lock()
	defer repo.mu.Unlock()
	filterMap := make(map[domain.MetricID]struct{})
	for _, filter := range filters {
		filterMap[filter] = struct{}{}
	}
	_, err := repo.file.Seek(0, io.SeekStart)
	if err != nil {
		logger.Error("Error seeking file", "error", err)
		return nil, err
	}
	var allMetrics []*domain.Metrics
	decoder := json.NewDecoder(repo.file)
	for {
		var metric domain.Metrics
		if err := decoder.Decode(&metric); err != nil {
			if err.Error() == "EOF" {
				break
			}
			logger.Error("Error decoding metric", "error", err)
			return nil, err
		}
		allMetrics = append(allMetrics, &metric)
	}
	result := make(map[domain.MetricID]*domain.Metrics)
	for _, metric := range allMetrics {
		metricID := domain.MetricID{ID: metric.ID, Type: metric.Type}
		if _, exists := filterMap[metricID]; exists {
			result[metricID] = metric
		}
	}
	logger.Info("Metrics found", "result_count", len(result))
	return result, nil
}
