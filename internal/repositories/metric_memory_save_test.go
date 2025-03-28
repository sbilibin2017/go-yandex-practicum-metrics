package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveMetricsSuccessfully(t *testing.T) {
	data := make(map[domain.MetricID]*domain.Metrics)
	repo := NewMetricMemorySaveBatchRepository(data)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Value: ptrFloat64(100),
		},
		{
			ID:    "metric2",
			Type:  "gauge",
			Value: ptrFloat64(200),
		},
	}
	result := repo.Save(context.Background(), metrics)
	assert.True(t, result)
	assert.Equal(t, 2, len(data))
	assert.Contains(t, data, domain.MetricID{ID: "metric1", Type: "counter"})
	assert.Contains(t, data, domain.MetricID{ID: "metric2", Type: "gauge"})
}

func ptrFloat64(value float64) *float64 {
	return &value
}
