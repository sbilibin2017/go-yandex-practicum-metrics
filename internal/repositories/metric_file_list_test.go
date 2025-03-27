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

func TestMetricFileListRepository_EmptyFile(t *testing.T) {
	// Setup: Create an empty temporary file
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Create the repository with the empty file
	repo := repositories.NewMetricFileListRepository(tmpFile)

	// Call List on the empty file
	result, err := repo.List(context.Background())
	require.NoError(t, err)

	// Assert that the result is empty since the file has no data
	assert.Len(t, result, 0)
}

func TestMetricFileListRepository_MultipleMetrics(t *testing.T) {
	// Setup: Create a temporary file with some metrics data
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
	repo := repositories.NewMetricFileListRepository(tmpFile)

	// Call List to retrieve all metrics
	result, err := repo.List(context.Background())
	require.NoError(t, err)

	// Assert that the result contains the correct metrics
	assert.Len(t, result, 2)
	assert.Equal(t, metrics[0], result[domain.MetricID{ID: "metric1", Type: "counter"}])
	assert.Equal(t, metrics[1], result[domain.MetricID{ID: "metric2", Type: "gauge"}])
}

func TestMetricFileListRepository_MalformedJSON(t *testing.T) {
	// Setup: Create a temporary file with malformed JSON
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Write malformed JSON to the file
	_, err = tmpFile.Write([]byte(`{"ID":"metric1","Type":"counter","Delta":10,`)) // Invalid JSON (missing closing brace)
	require.NoError(t, err)

	// Close the file before reading
	err = tmpFile.Close()
	require.NoError(t, err)

	// Create the repository with the malformed file
	repo := repositories.NewMetricFileListRepository(tmpFile)

	// Call List and expect an error due to malformed JSON
	result, err := repo.List(context.Background())
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestMetricFileListRepository_OneMetric(t *testing.T) {
	// Setup: Create a temporary file with a single metric
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Write a single metric to the file
	metric := &domain.Metrics{
		ID:    "metric1",
		Type:  "counter",
		Delta: ptrInt64(10),
	}
	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(metric)
	require.NoError(t, err)

	// Rewind the file pointer to the beginning before reading
	tmpFile.Seek(0, 0)

	// Create repository with the file
	repo := repositories.NewMetricFileListRepository(tmpFile)

	// Call List to retrieve the metrics
	result, err := repo.List(context.Background())
	require.NoError(t, err)

	// Assert that the result contains the correct metric
	assert.Len(t, result, 1)
	assert.Equal(t, metric, result[domain.MetricID{ID: "metric1", Type: "counter"}])
}
