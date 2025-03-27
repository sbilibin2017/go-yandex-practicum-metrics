package repositories

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/domain"
	"os"
	"sync"
)

type MetricFileListRepository struct {
	file *os.File
	mu   sync.Mutex
}

func NewMetricFileListRepository(file *os.File) *MetricFileListRepository {
	return &MetricFileListRepository{file: file}
}

func (repo *MetricFileListRepository) List(
	ctx context.Context,
) (map[domain.MetricID]*domain.Metrics, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	result := make(map[domain.MetricID]*domain.Metrics)
	decoder := json.NewDecoder(repo.file)
	for {
		var metric domain.Metrics
		if err := decoder.Decode(&metric); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		result[domain.MetricID{ID: metric.ID, Type: metric.Type}] = &metric
	}
	return result, nil
}
