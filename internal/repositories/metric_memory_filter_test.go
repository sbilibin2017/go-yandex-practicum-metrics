package repositories

import (
	"context"

	"go-yandex-practicum-metrics/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a new repository
func newTestFilterRepository() *MetricMemoryFilterRepository {
	return NewMetricMemoryFilterRepository(make(map[types.MetricID]*types.Metrics))
}

// Test filtering a metric that exists in the repository
func TestFilterMetricExists(t *testing.T) {
	repo := newTestFilterRepository()

	// Create and add a metric to the repository
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric.Value = 42.0
	repo.data[metric.MetricID] = metric

	// Filter by the metric ID
	filter := types.MetricID{ID: "metric_1", Type: "gauge"}
	storedMetric, exists := repo.Filter(context.Background(), filter)

	// Assert that the metric exists and is correctly retrieved
	assert.True(t, exists)
	assert.Equal(t, metric, storedMetric)
}

// Test filtering a metric that does not exist in the repository
func TestFilterMetricNotExists(t *testing.T) {
	repo := newTestFilterRepository()

	// Try to filter by a non-existing metric ID
	filter := types.MetricID{ID: "metric_nonexistent", Type: "gauge"}
	storedMetric, exists := repo.Filter(context.Background(), filter)

	// Assert that the metric does not exist
	assert.False(t, exists)
	assert.Nil(t, storedMetric)
}

// Test filtering an empty MetricID (invalid case)
func TestFilterEmptyMetricID(t *testing.T) {
	repo := newTestFilterRepository()

	// Add a metric with a non-empty MetricID
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric.Value = 42.0
	repo.data[metric.MetricID] = metric

	// Try filtering with an empty MetricID (assuming empty IDs are invalid in this case)
	filter := types.MetricID{ID: "", Type: ""}
	storedMetric, exists := repo.Filter(context.Background(), filter)

	// Assert that no metric is found with an empty MetricID
	assert.False(t, exists)
	assert.Nil(t, storedMetric)
}

// Test filtering a metric with a valid but different MetricType
func TestFilterMetricWithDifferentType(t *testing.T) {
	repo := newTestFilterRepository()

	// Add a metric with a different type
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric.Value = 42.0
	repo.data[metric.MetricID] = metric

	// Try filtering with the same MetricID but different Type
	filter := types.MetricID{ID: "metric_1", Type: "counter"}
	storedMetric, exists := repo.Filter(context.Background(), filter)

	// Assert that the metric is not found due to type mismatch
	assert.False(t, exists)
	assert.Nil(t, storedMetric)
}
