package repositories_test

import (
	"context"
	"encoding/json"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/repositories"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileFindBatch(t *testing.T) {
	// Setup: Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Write some test data to the file
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: ptrInt64(10),
		},
		{
			ID:    "metric2",
			Type:  "gauge",
			Value: ptrFloat64(20.5),
		},
	}
	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "  ")
	for _, metric := range metrics {
		err = encoder.Encode(metric)
		require.NoError(t, err)
	}

	// Rewind the file pointer to the beginning before reading
	tmpFile.Seek(0, 0)

	// Create repository with the file
	repo := repositories.NewMetricFileFindBatchRepository(tmpFile)

	// Define filters (we will filter for "metric1")
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}

	// Call FindBatch
	result, err := repo.FindBatch(context.Background(), filters)
	require.NoError(t, err)

	// Assert that the result contains only the filtered metric
	assert.Len(t, result, 1)
	assert.Equal(t, metrics[0], result[domain.MetricID{ID: "metric1", Type: "counter"}])
}

func TestFileFindBatch_EmptyFile(t *testing.T) {
	// Setup: Create an empty temporary file
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Create the repository with the empty file
	repo := repositories.NewMetricFileFindBatchRepository(tmpFile)

	// Define some filters
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}

	// Call FindBatch on the empty file
	result, err := repo.FindBatch(context.Background(), filters)
	require.NoError(t, err)

	// Assert that the result is empty since the file has no data
	assert.Len(t, result, 0)
}

func TestFileFindBatch_NoMatchingFilters(t *testing.T) {
	// Setup: Create a temporary file with some data
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Write some test data to the file
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Delta: ptrInt64(10),
		},
		{
			ID:    "metric2",
			Type:  "gauge",
			Value: ptrFloat64(20.5),
		},
	}
	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "  ")
	for _, metric := range metrics {
		err = encoder.Encode(metric)
		require.NoError(t, err)
	}

	// Rewind the file pointer to the beginning before reading
	tmpFile.Seek(0, 0)

	// Create repository with the file
	repo := repositories.NewMetricFileFindBatchRepository(tmpFile)

	// Define filters that do not match any metrics
	filters := []domain.MetricID{
		{ID: "nonexistent", Type: "counter"},
	}

	// Call FindBatch with non-matching filters
	result, err := repo.FindBatch(context.Background(), filters)
	require.NoError(t, err)

	// Assert that no metrics were found
	assert.Len(t, result, 0)
}

func TestFileFindBatch_ErrorOnFileRead(t *testing.T) {
	// Setup: Create a temporary file and write malformed JSON to it
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Write malformed JSON to the file
	_, err = tmpFile.Write([]byte(`{"ID":"metric1","Type":"counter","Delta":10,`)) // Invalid JSON (missing closing brace)
	require.NoError(t, err)

	// Close the file before reading
	err = tmpFile.Close()
	require.NoError(t, err)

	// Create the repository with the file
	repo := repositories.NewMetricFileFindBatchRepository(tmpFile)

	// Define filters
	filters := []domain.MetricID{
		{ID: "metric1", Type: "counter"},
	}

	// Call FindBatch and expect an error due to malformed JSON
	result, err := repo.FindBatch(context.Background(), filters)
	require.Error(t, err)
	assert.Nil(t, result)
}
