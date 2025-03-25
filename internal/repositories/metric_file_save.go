package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/types"
)

type FileEncoder interface {
	Encode(v interface{}) error
}

type MetricFileSaveRepository struct {
	encoder FileEncoder
}

func NewMetricFileSaveRepository(encoder FileEncoder) *MetricFileSaveRepository {
	return &MetricFileSaveRepository{encoder: encoder}
}

// Save method to save a single metric to the file
func (m *MetricFileSaveRepository) Save(ctx context.Context, metric *types.Metrics) error {
	err := m.encoder.Encode(metric)
	if err != nil {
		return err
	}
	return nil
}
