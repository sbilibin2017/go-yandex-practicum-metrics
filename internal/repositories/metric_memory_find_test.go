package repositories

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindMetricsSuccessfully(t *testing.T) {
	data := make(map[domain.MetricID]*domain.Metrics)
	repo := NewMetricMemoryFindBatchRepository(data)
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
	data[domain.MetricID{ID: "metric1", Type: "counter"}] = metrics[0]
	data[domain.MetricID{ID: "metric2", Type: "gauge"}] = metrics[1]
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}
	result, success := repo.Find(context.Background(), filters)
	assert.True(t, success)
	assert.Equal(t, 1, len(result))
	assert.Contains(t, result, domain.MetricID{ID: "metric1", Type: "counter"})
}

func TestFindNoMatchingMetrics(t *testing.T) {
	data := make(map[domain.MetricID]*domain.Metrics)
	repo := NewMetricMemoryFindBatchRepository(data)
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
	data[domain.MetricID{ID: "metric1", Type: "counter"}] = metrics[0]
	data[domain.MetricID{ID: "metric2", Type: "gauge"}] = metrics[1]
	filters := []domain.MetricID{
		{ID: "metric3", Type: "counter"},
	}
	result, success := repo.Find(context.Background(), filters)
	assert.True(t, success)
	assert.Equal(t, 0, len(result))
}

func TestFindMultipleMetrics(t *testing.T) {
	data := make(map[domain.MetricID]*domain.Metrics)
	repo := NewMetricMemoryFindBatchRepository(data)
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
		{
			ID:    "metric3",
			Type:  "counter",
			Value: ptrFloat64(300),
		},
	}
	data[domain.MetricID{ID: "metric1", Type: "counter"}] = metrics[0]
	data[domain.MetricID{ID: "metric2", Type: "gauge"}] = metrics[1]
	data[domain.MetricID{ID: "metric3", Type: "counter"}] = metrics[2]
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
		{ID: "metric2", Type: "gauge"},
	}
	result, success := repo.Find(context.Background(), filters)
	assert.True(t, success)
	assert.Equal(t, 2, len(result))
	assert.Contains(t, result, domain.MetricID{ID: "metric1", Type: "counter"})
	assert.Contains(t, result, domain.MetricID{ID: "metric2", Type: "gauge"})
}
