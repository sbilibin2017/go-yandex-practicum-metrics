package repositories

import (
	"context"

	"go-yandex-practicum-metrics/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a new repository
func newTestListAllRepository() *MetricMemoryListAllRepository {
	return NewMetricMemoryListAllRepository(make(map[types.MetricID]*types.Metrics))
}

// Test ListAll with metrics in the repository
func TestListAllWithMetrics(t *testing.T) {
	repo := newTestListAllRepository()

	// Add some metrics to the repository
	metric1 := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric1.Value = 42.0
	repo.data[metric1.MetricID] = metric1

	metric2 := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_2", Type: "counter"},
		Value:    new(float64),
	}
	*metric2.Value = 10.0
	repo.data[metric2.MetricID] = metric2

	// List all metrics
	metrics, success := repo.ListAll(context.Background())

	// Assert that the ListAll was successful and the result is correct
	assert.True(t, success)
	assert.Len(t, metrics, 2)
	assert.Contains(t, metrics, metric1)
	assert.Contains(t, metrics, metric2)
}

// Test ListAll with an empty repository
func TestListAllEmptyRepository(t *testing.T) {
	repo := newTestListAllRepository()

	// List all metrics from an empty repository
	metrics, success := repo.ListAll(context.Background())

	// Assert that the ListAll was successful and the result is an empty slice
	assert.True(t, success)
	assert.Empty(t, metrics)
}

// Test ListAll with only one metric in the repository
func TestListAllWithOneMetric(t *testing.T) {
	repo := newTestListAllRepository()

	// Add a single metric to the repository
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric.Value = 42.0
	repo.data[metric.MetricID] = metric

	// List all metrics
	metrics, success := repo.ListAll(context.Background())

	// Assert that the ListAll was successful and the result contains one metric
	assert.True(t, success)
	assert.Len(t, metrics, 1)
	assert.Contains(t, metrics, metric)
}
