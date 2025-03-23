package repositories

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/types"
)

// Extend ReadFile interface to include Seek
type ReadFilterFile interface {
	Read(p []byte) (n int, err error)
	Seek(offset int64, whence int) (int64, error)
}

type MetricFileFilterRepository struct {
	file ReadFilterFile
}

func NewMetricFileFilterRepository(file ReadFilterFile) *MetricFileFilterRepository {
	return &MetricFileFilterRepository{file: file}
}

// Filter method to read the file and return a metric based on the single filter
func (m *MetricFileFilterRepository) Filter(ctx context.Context, filter types.MetricID) (*types.Metrics, bool) {
	if m.file == nil {
		return nil, false
	}
	_, err := m.file.Seek(0, 0)
	if err != nil {
		return nil, false
	}
	decoder := json.NewDecoder(m.file)
	for {
		var metric types.Metrics
		if err := decoder.Decode(&metric); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, false
		}
		if metric.MetricID == filter {
			return &metric, true
		}
	}
	return nil, false
}
