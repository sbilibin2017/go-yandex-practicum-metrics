package repositories

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/domain"
	"io"
	"sync"
)

type MetricFileSaveBatchRepository struct {
	writer io.Writer
	mu     sync.Mutex
}

func NewMetricFileSaveBatchRepository(
	writer io.Writer,
) *MetricFileSaveBatchRepository {
	return &MetricFileSaveBatchRepository{writer: writer}
}

func (repo *MetricFileSaveBatchRepository) Save(
	ctx context.Context, metrics []*domain.Metrics,
) bool {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.writer == nil {
		return false
	}
	encoder := json.NewEncoder(repo.writer)
	for _, metric := range metrics {
		if err := encoder.Encode(metric); err != nil {
			return false
		}
	}
	return true
}
