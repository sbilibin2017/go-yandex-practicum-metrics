package repositories

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/domain"
	"os"
	"sync"
)

type MetricFileFindBatchRepository struct {
	file *os.File
	mu   sync.Mutex
}

func NewMetricFileFindBatchRepository(file *os.File) *MetricFileFindBatchRepository {
	return &MetricFileFindBatchRepository{file: file}
}

func (repo *MetricFileFindBatchRepository) FindBatch(
	ctx context.Context, filters []domain.MetricID,
) (map[domain.MetricID]*domain.Metrics, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	filterMap := make(map[domain.MetricID]struct{})
	for _, filter := range filters {
		filterMap[filter] = struct{}{}
	}
	var result = make(map[domain.MetricID]*domain.Metrics)
	decoder := json.NewDecoder(repo.file)
	for {
		var metric domain.Metrics
		if err := decoder.Decode(&metric); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if _, found := filterMap[domain.MetricID{ID: metric.ID, Type: metric.Type}]; found {
			result[domain.MetricID{ID: metric.ID, Type: metric.Type}] = &metric
		}
	}
	return result, nil
}
