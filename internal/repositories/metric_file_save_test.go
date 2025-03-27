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

func TestFileSaveBatch(t *testing.T) {
	// Setup: Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Create the repository with the temp file
	repo := repositories.NewMetricFileSaveBatchRepository(tmpFile)

	// Create test data
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

	// Call SaveBatch
	err = repo.SaveBatch(context.Background(), metrics)
	require.NoError(t, err)

	// Open the file to verify the written data
	tmpFile.Seek(0, 0) // Rewind the file pointer to the beginning
	var writtenMetrics []*domain.Metrics
	decoder := json.NewDecoder(tmpFile)
	err = decoder.Decode(&writtenMetrics)
	require.NoError(t, err)

	// Assert that the written metrics are the same as the original metrics
	assert.Len(t, writtenMetrics, 2)
	assert.Equal(t, metrics[0], writtenMetrics[0])
	assert.Equal(t, metrics[1], writtenMetrics[1])
}

func TestFileSaveBatch_ErrorOnFileWrite(t *testing.T) {
	// Setup: Create a read-only temporary file for testing
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	tmpFile.Close()                      // Close the file immediately
	err = os.Chmod(tmpFile.Name(), 0444) // Change permissions to read-only
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Create the repository with the read-only file
	repo := repositories.NewMetricFileSaveBatchRepository(tmpFile)

	// Create test data
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

	// Call SaveBatch, expecting an error because the file is read-only
	err = repo.SaveBatch(context.Background(), metrics)
	require.Error(t, err)
}

func TestFileSaveBatch_EmptyMetrics(t *testing.T) {
	// Setup: Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "metrics_test_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name()) // Clean up after test

	// Create the repository with the temp file
	repo := repositories.NewMetricFileSaveBatchRepository(tmpFile)

	// Call SaveBatch with an empty slice
	err = repo.SaveBatch(context.Background(), []*domain.Metrics{})
	require.NoError(t, err)

	// Open the file to verify the written data
	tmpFile.Seek(0, 0) // Rewind the file pointer to the beginning
	var writtenMetrics []*domain.Metrics
	decoder := json.NewDecoder(tmpFile)
	err = decoder.Decode(&writtenMetrics)
	require.NoError(t, err)

	// Assert that the file is empty (no metrics written)
	assert.Len(t, writtenMetrics, 0)
}
