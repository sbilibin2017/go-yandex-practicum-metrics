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

func TestList(t *testing.T) {
	// Setup: Initialize the repository with some test data
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
	repo := repositories.NewMetricMemoryListAllRepository(data)

	// Call List method
	result, err := repo.List(context.Background())

	// Assert no error
	require.NoError(t, err)

	// Assert that the result contains all the metrics
	assert.Len(t, result, 2)
	assert.Equal(t, data[domain.MetricID{ID: "metric1", Type: "counter"}], result[domain.MetricID{ID: "metric1", Type: "counter"}])
	assert.Equal(t, data[domain.MetricID{ID: "metric2", Type: "gauge"}], result[domain.MetricID{ID: "metric2", Type: "gauge"}])
}

func TestList_EmptyRepository(t *testing.T) {
	// Setup: Initialize the repository with no data
	data := make(map[domain.MetricID]*domain.Metrics)
	repo := repositories.NewMetricMemoryListAllRepository(data)

	// Call List method on an empty repository
	result, err := repo.List(context.Background())

	// Assert no error
	require.NoError(t, err)

	// Assert that the result is an empty map
	assert.Len(t, result, 0)
}

func TestList_Concurrency(t *testing.T) {
	// Setup: Initialize the repository with some test data
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
	repo := repositories.NewMetricMemoryListAllRepository(data)

	// Run List method concurrently in different goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := repo.List(context.Background())
			require.NoError(t, err)

			// Assert that the result contains all the metrics
			assert.Len(t, result, 2)
			assert.Equal(t, data[domain.MetricID{ID: "metric1", Type: "counter"}], result[domain.MetricID{ID: "metric1", Type: "counter"}])
			assert.Equal(t, data[domain.MetricID{ID: "metric2", Type: "gauge"}], result[domain.MetricID{ID: "metric2", Type: "gauge"}])
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
}
