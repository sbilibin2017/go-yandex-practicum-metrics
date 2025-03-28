package repositories

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/domain"
	"io"
	"sync"
)

type MetricFileFindBatchRepository struct {
	reader io.Reader
	seeker io.Seeker
	mu     sync.Mutex
}

func NewMetricFileFindBatchRepository(
	reader io.Reader,
	seeker io.Seeker,
) *MetricFileFindBatchRepository {
	return &MetricFileFindBatchRepository{
		reader: reader,
		seeker: seeker,
	}
}

func (repo *MetricFileFindBatchRepository) Find(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	filterMap := make(map[domain.MetricID]struct{})
	for _, filter := range filters {
		filterMap[filter] = struct{}{}
	}
	if repo.reader == nil {
		return nil, false
	}
	_, err := repo.seeker.Seek(0, io.SeekStart)
	if err != nil {
		return nil, false
	}
	var allMetrics []*domain.Metrics
	decoder := json.NewDecoder(repo.reader)
	for {
		var metric domain.Metrics
		if err := decoder.Decode(&metric); err != nil {
			if err.Error() == "EOF" {
				break
			}
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
	return result, true
}
