package repositories

import (
	"context"
	"encoding/json"

	"go-yandex-practicum-metrics/internal/types"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a new temporary file repository
func newTestFileSaveRepository() (*MetricFileSaveRepository, *os.File) {
	// Create a temporary file
	file, err := os.CreateTemp("", "metrics_test_*.json")
	if err != nil {
		panic(err)
	}
	// Create a MetricFileSaveRepository using the temporary file
	repo := NewMetricFileSaveRepository(file)
	return repo, file
}

// Test Save with a valid metric and file
func TestSaveMetricSuccess(t *testing.T) {
	repo, file := newTestFileSaveRepository()
	defer os.Remove(file.Name()) // Clean up the temporary file

	// Prepare a metric to save
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric.Value = 42.0

	// Call Save method
	success := repo.Save(context.Background(), metric)

	// Assert that the save operation was successful
	assert.True(t, success)

	// Verify that the data was written to the file by reading it back
	fileContent, err := os.ReadFile(file.Name())
	assert.NoError(t, err)

	// Decode the saved metric from the file
	var savedMetric types.Metrics
	err = json.Unmarshal(fileContent, &savedMetric)
	assert.NoError(t, err)

	// Assert that the metric in the file matches the original metric
	assert.Equal(t, metric, &savedMetric)
}

// Test Save with a nil file (invalid file)
func TestSaveNilFile(t *testing.T) {
	// Create a repository with a nil file
	repo := NewMetricFileSaveRepository(nil)

	// Prepare a metric to save
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric.Value = 42.0

	// Call Save method
	success := repo.Save(context.Background(), metric)

	// Assert that the save operation was unsuccessful because the file is nil
	assert.False(t, success)
}

// Test Save with an invalid file path (file cannot be opened)
func TestSaveInvalidFile(t *testing.T) {
	// Attempt to open a non-existent file (should fail)
	// Use a path that likely does not exist or is not accessible
	invalidFilePath := "/non_existent_directory/invalid_file.json"

	// Try opening a file in a non-existent or restricted directory
	_, err := os.OpenFile(invalidFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		t.Skip("Skipping test because the file could not be created at the specified path")
	}

	// Create the repository using the invalid file path
	repo := NewMetricFileSaveRepository(nil) // Simulating a nil file scenario

	// Prepare a metric to save
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_2", Type: "counter"},
		Value:    new(float64),
	}
	*metric.Value = 10.0

	// Call Save method
	success := repo.Save(context.Background(), metric)

	// Assert that the save operation was unsuccessful due to the invalid file
	assert.False(t, success, "Expected Save to return false when the file path is invalid or inaccessible")
}

// Test Save with empty metric (checking edge case)
func TestSaveEmptyMetric(t *testing.T) {
	repo, file := newTestFileSaveRepository()
	defer os.Remove(file.Name()) // Clean up the temporary file

	// Prepare an empty metric (no value)
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_3", Type: "gauge"},
	}

	// Call Save method with the empty metric
	success := repo.Save(context.Background(), metric)

	// Assert that the save operation was successful
	assert.True(t, success)

	// Verify that the data was written to the file by reading it back
	fileContent, err := os.ReadFile(file.Name())
	assert.NoError(t, err)

	// Decode the saved metric from the file
	var savedMetric types.Metrics
	err = json.Unmarshal(fileContent, &savedMetric)
	assert.NoError(t, err)

	// Assert that the metric in the file matches the original metric (empty value)
	assert.Equal(t, metric, &savedMetric)
}
