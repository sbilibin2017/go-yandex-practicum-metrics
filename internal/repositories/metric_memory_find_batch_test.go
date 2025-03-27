package repositories_test

import (
	"context"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/repositories"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindBatch(t *testing.T) {
	// Setup
	data := map[domain.MetricID]*domain.Metrics{
		{ID: "metric1", Type: "counter"}: {
			ID:    "metric1",
			Type:  "counter",
			Delta: ptrInt64(10),
		},
		{ID: "metric2", Type: "gauge"}: {
			ID:    "metric2",
			Type:  "gauge",
			Value: ptrFloat64(20.5),
		},
	}
	repo := repositories.NewMetricMemoryFindBatchRepository(data)

	// Test with existing metric IDs
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
		{ID: "metric2", Type: "gauge"},
	}
	result, err := repo.FindBatch(context.Background(), filters)

	// Assert no error
	require.NoError(t, err)

	// Assert both metrics are found
	assert.Len(t, result, 2)
	assert.Equal(t, data[domain.MetricID{ID: "metric1", Type: "counter"}], result[domain.MetricID{ID: "metric1", Type: "counter"}])
	assert.Equal(t, data[domain.MetricID{ID: "metric2", Type: "gauge"}], result[domain.MetricID{ID: "metric2", Type: "gauge"}])
}

func TestFindBatch_MissingMetrics(t *testing.T) {
	// Setup
	data := map[domain.MetricID]*domain.Metrics{
		{ID: "metric1", Type: "counter"}: {
			ID:    "metric1",
			Type:  "counter",
			Delta: ptrInt64(10),
		},
	}
	repo := repositories.NewMetricMemoryFindBatchRepository(data)

	// Test with a non-existing metric ID
	filters := []domain.MetricID{
		{ID: "metric2", Type: "gauge"},
	}
	result, err := repo.FindBatch(context.Background(), filters)

	// Assert no error
	require.NoError(t, err)

	// Assert that no metrics were found
	assert.Len(t, result, 0)
}

func TestFindBatch_EmptyFilters(t *testing.T) {
	// Setup
	data := map[domain.MetricID]*domain.Metrics{
		{ID: "metric1", Type: "counter"}: {
			ID:    "metric1",
			Type:  "counter",
			Delta: ptrInt64(10),
		},
		{ID: "metric2", Type: "gauge"}: {
			ID:    "metric2",
			Type:  "gauge",
			Value: ptrFloat64(20.5),
		},
	}
	repo := repositories.NewMetricMemoryFindBatchRepository(data)

	// Test with an empty filters slice
	filters := []domain.MetricID{}
	result, err := repo.FindBatch(context.Background(), filters)

	// Assert no error
	require.NoError(t, err)

	// Assert that no metrics were returned
	assert.Len(t, result, 0)
}

func TestFindBatch_Concurrency(t *testing.T) {
	// Setup
	data := map[domain.MetricID]*domain.Metrics{
		{ID: "metric1", Type: "counter"}: {
			ID:    "metric1",
			Type:  "counter",
			Delta: ptrInt64(10),
		},
		{ID: "metric2", Type: "gauge"}: {
			ID:    "metric2",
			Type:  "gauge",
			Value: ptrFloat64(20.5),
		},
	}
	repo := repositories.NewMetricMemoryFindBatchRepository(data)

	// Test with existing metric IDs
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
		{ID: "metric2", Type: "gauge"},
	}

	// Run FindBatch concurrently in different goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := repo.FindBatch(context.Background(), filters)
			require.NoError(t, err)

			// Assert both metrics are found
			assert.Len(t, result, 2)
			assert.Equal(t, data[domain.MetricID{ID: "metric1", Type: "counter"}], result[domain.MetricID{ID: "metric1", Type: "counter"}])
			assert.Equal(t, data[domain.MetricID{ID: "metric2", Type: "gauge"}], result[domain.MetricID{ID: "metric2", Type: "gauge"}])
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
}

// Helper functions for pointers to values
func ptrInt64(v int64) *int64 {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}
