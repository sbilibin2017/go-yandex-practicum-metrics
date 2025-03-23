package repositories

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/types"
	"io"
)

type FileListAllReader interface {
	io.Reader
	io.Seeker
}

type MetricFileListAllRepository struct {
	file FileListAllReader
}

func NewMetricFileListAllRepository(file FileListAllReader) *MetricFileListAllRepository {
	return &MetricFileListAllRepository{file: file}
}

// ListAll method to return all metrics from the file, now using a map instead of a slice
func (m *MetricFileListAllRepository) ListAll(ctx context.Context) ([]*types.Metrics, bool) {
	if m.file == nil {
		return nil, false
	}
	_, err := m.file.Seek(0, 0)
	if err != nil {
		return nil, false
	}
	result := make(map[types.MetricID]*types.Metrics)
	decoder := json.NewDecoder(m.file)
	for {
		var metric types.Metrics
		if err := decoder.Decode(&metric); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, false
		}
		result[metric.MetricID] = &metric
	}
	var metricsSlice []*types.Metrics
	for _, metric := range result {
		metricsSlice = append(metricsSlice, metric)
	}
	return metricsSlice, true
}
