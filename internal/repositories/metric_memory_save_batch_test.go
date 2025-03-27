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

func TestSaveBatch(t *testing.T) {
	// Setup
	data := make(map[domain.MetricID]*domain.Metrics)
	repo := repositories.NewMetricMemorySaveBatchRepository(data)

	// Create some test metrics
	delta1 := int64(10)
	value2 := float64(20.5)

	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: &delta1,
		},
		{
			ID:    "metric2",
			Type:  "gauge",
			Value: &value2,
		},
	}

	// Call SaveBatch method
	err := repo.SaveBatch(context.Background(), metrics)

	// Assert no error
	require.NoError(t, err)

	// Assert metrics were saved correctly
	assert.Len(t, data, 2)
	assert.Equal(t, metrics[0], data[domain.MetricID{ID: "metric1", Type: "counter"}])
	assert.Equal(t, metrics[1], data[domain.MetricID{ID: "metric2", Type: "gauge"}])
}

func TestSaveBatch_EmptyMetrics(t *testing.T) {
	// Setup
	data := make(map[domain.MetricID]*domain.Metrics)
	repo := repositories.NewMetricMemorySaveBatchRepository(data)

	// Call SaveBatch with an empty slice
	err := repo.SaveBatch(context.Background(), []*domain.Metrics{})

	// Assert no error
	require.NoError(t, err)

	// Assert that data is still empty
	assert.Len(t, data, 0)
}

func TestSaveBatch_Concurrency(t *testing.T) {
	// Setup
	data := make(map[domain.MetricID]*domain.Metrics)
	repo := repositories.NewMetricMemorySaveBatchRepository(data)

	// Create some test metrics
	delta1 := int64(10)
	value2 := float64(20.5)

	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: &delta1,
		},
		{
			ID:    "metric2",
			Type:  "gauge",
			Value: &value2,
		},
	}

	// Run SaveBatch concurrently in different goroutines
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.SaveBatch(context.Background(), metrics)
			require.NoError(t, err)
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Assert metrics are still saved correctly
	assert.Len(t, data, 2)
	assert.Equal(t, metrics[0], data[domain.MetricID{ID: "metric1", Type: "counter"}])
	assert.Equal(t, metrics[1], data[domain.MetricID{ID: "metric2", Type: "gauge"}])
}
