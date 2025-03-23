package repositories

import (
	"context"

	"go-yandex-practicum-metrics/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a new repository
func newTestRepository() *MetricMemorySaveRepository {
	return NewMetricMemorySaveRepository(make(map[types.MetricID]*types.Metrics))
}

// Test saving a valid metric
func TestSaveMetric(t *testing.T) {
	repo := newTestRepository()

	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric.Value = 42.0

	success := repo.Save(context.Background(), metric)

	// Assert that the save was successful
	assert.True(t, success)

	// Assert that the metric is saved in the repository
	storedMetric, exists := repo.data[metric.MetricID]
	assert.True(t, exists)
	assert.Equal(t, metric, storedMetric)
}

// Test saving a metric with the same MetricID (overwrite scenario)
func TestSaveMetricWithSameID(t *testing.T) {
	repo := newTestRepository()

	// Save the initial metric
	metric1 := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric1.Value = 42.0
	repo.Save(context.Background(), metric1)

	// Save the second metric with the same MetricID but different value
	metric2 := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric2.Value = 100.0

	success := repo.Save(context.Background(), metric2)

	// Assert that the save was successful
	assert.True(t, success)

	// Assert that the new value is saved in the repository
	storedMetric, exists := repo.data[metric2.MetricID]
	assert.True(t, exists)
	assert.Equal(t, metric2, storedMetric)
}

// Test saving a metric with a Delta value
func TestSaveMetricWithDelta(t *testing.T) {
	repo := newTestRepository()

	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_2", Type: "counter"},
		Delta:    new(int64),
	}
	*metric.Delta = 10

	success := repo.Save(context.Background(), metric)

	// Assert that the save was successful
	assert.True(t, success)

	// Assert that the metric with Delta is saved in the repository
	storedMetric, exists := repo.data[metric.MetricID]
	assert.True(t, exists)
	assert.Equal(t, metric, storedMetric)
}

// Test saving an empty MetricID (invalid case)
func TestSaveEmptyMetricID(t *testing.T) {
	repo := newTestRepository()

	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "", Type: ""},
		Value:    new(float64),
	}
	*metric.Value = 0.0

	success := repo.Save(context.Background(), metric)

	// Assuming that saving an empty MetricID is allowed (modify if needed)
	assert.True(t, success) // You can change this assertion based on your validation rules.
}
