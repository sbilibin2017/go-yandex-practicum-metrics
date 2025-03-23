package repositories

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/types"
	"os"
)

type MetricFileSaveRepository struct {
	file *os.File
}

func NewMetricFileSaveRepository(file *os.File) *MetricFileSaveRepository {
	return &MetricFileSaveRepository{file: file}
}

// Save method to save a single metric to the file
func (m *MetricFileSaveRepository) Save(ctx context.Context, metric *types.Metrics) bool {
	if m.file == nil {
		return false
	}
	encoder := json.NewEncoder(m.file)
	err := encoder.Encode(metric)
	return err == nil
}
