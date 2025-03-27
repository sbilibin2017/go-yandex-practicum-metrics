package repositories

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/domain"
	"os"
	"sync"
)

type MetricFileSaveBatchRepository struct {
	file *os.File
	mu   sync.Mutex
}

func NewMetricFileSaveBatchRepository(file *os.File) *MetricFileSaveBatchRepository {
	return &MetricFileSaveBatchRepository{file: file}
}

func (repo *MetricFileSaveBatchRepository) SaveBatch(
	ctx context.Context, metrics []*domain.Metrics,
) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	encoder := json.NewEncoder(repo.file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metrics); err != nil {
		return err
	}
	return nil
}
