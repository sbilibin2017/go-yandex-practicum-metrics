package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"go-yandex-practicum-metrics/internal/types"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a new temporary file repository
func newTestFileFilterRepository() (*MetricFileFilterRepository, *os.File) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "metrics_test_*.json")
	if err != nil {
		panic(err)
	}
	// Create the MetricFileFilterRepository using the temporary file
	repo := NewMetricFileFilterRepository(file)
	return repo, file
}

// Test Filter with a valid metric in the file
func TestFilterWithValidMetric(t *testing.T) {
	repo, file := newTestFileFilterRepository()
	defer os.Remove(file.Name()) // Clean up the temporary file after the test

	// Prepare a metric to save
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric.Value = 42.0

	// Write the metric to the file
	encoder := json.NewEncoder(file)
	err := encoder.Encode(metric)
	assert.NoError(t, err)

	// Reset the file pointer to the beginning before calling Filter
	_, err = file.Seek(0, 0)
	assert.NoError(t, err)

	// Call Filter method with the metric filter
	filter := types.MetricID{ID: "metric_1", Type: "gauge"}
	result, success := repo.Filter(context.Background(), filter)

	// Assert that the metric was found in the file
	assert.True(t, success)
	assert.Equal(t, metric, result)
}

// Test Filter with a non-existing metric in the file
func TestFilterWithNonExistingMetric(t *testing.T) {
	repo, file := newTestFileFilterRepository()
	defer os.Remove(file.Name()) // Clean up the temporary file after the test

	// Prepare a metric to save
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "metric_1", Type: "gauge"},
		Value:    new(float64),
	}
	*metric.Value = 42.0

	// Write the metric to the file
	encoder := json.NewEncoder(file)
	err := encoder.Encode(metric)
	assert.NoError(t, err)

	// Reset the file pointer to the beginning before calling Filter
	_, err = file.Seek(0, 0)
	assert.NoError(t, err)

	// Call Filter method with a non-existing filter
	filter := types.MetricID{ID: "metric_2", Type: "counter"}
	result, success := repo.Filter(context.Background(), filter)

	// Assert that the metric was not found in the file
	assert.False(t, success)
	assert.Nil(t, result)
}

// Test Filter with an empty file
func TestFilterWithEmptyFile(t *testing.T) {
	repo, file := newTestFileFilterRepository()
	defer os.Remove(file.Name()) // Clean up the temporary file after the test

	// No metrics are written to the file, so it's empty

	// Reset the file pointer to the beginning before calling Filter
	_, err := file.Seek(0, 0)
	assert.NoError(t, err)

	// Call Filter method with any filter
	filter := types.MetricID{ID: "metric_1", Type: "gauge"}
	result, success := repo.Filter(context.Background(), filter)

	// Assert that the result is nil and success is false
	assert.False(t, success)
	assert.Nil(t, result)
}

// Test Filter with nil file (invalid repository)
func TestFilterWithNilFile(t *testing.T) {
	// Create the repository with a nil file
	repo := NewMetricFileFilterRepository(nil)

	// Prepare a filter
	filter := types.MetricID{ID: "metric_1", Type: "gauge"}

	// Call Filter method
	result, success := repo.Filter(context.Background(), filter)

	// Assert that the result is nil and success is false
	assert.False(t, success)
	assert.Nil(t, result)
}

// Mocking a reader that always returns an error
type errorReader struct{}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF // Simulate a read error
}

func (er *errorReader) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("seek error") // Simulate seek error
}

// Test Filter with file read error
func TestFilterWithFileReadError(t *testing.T) {
	// Use the errorReader to simulate a read failure
	errorReader := &errorReader{}

	// Create the repository with the errorReader
	repo := NewMetricFileFilterRepository(errorReader)

	// Call Filter with any metric ID; the read will fail
	filter := types.MetricID{ID: "metric_1", Type: "gauge"}
	result, success := repo.Filter(context.Background(), filter)

	// Assert that the result is nil and success is false due to the read error
	assert.False(t, success, "Expected Filter to return false due to read error")
	assert.Nil(t, result, "Expected Filter to return nil due to read error")
}

// Mock для симуляции ошибки при JSON-декодировании
type errorDecodeReader struct{}

func (e *errorDecodeReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("decode error") // Всегда возвращает ошибку
}

func (e *errorDecodeReader) Seek(offset int64, whence int) (int64, error) {
	return 0, nil // Seek работает нормально
}

// Тест обработки ошибки в decoder.Decode()
func TestFilter_DecodeError(t *testing.T) {
	// Используем кастомный reader, который всегда возвращает ошибку при Read()
	errorReader := &errorDecodeReader{}

	// Создаем репозиторий
	repo := NewMetricFileFilterRepository(errorReader)

	// Вызываем Filter с любым фильтром
	filter := types.MetricID{ID: "metric_1", Type: "gauge"}
	result, success := repo.Filter(context.Background(), filter)

	// Проверяем, что вернулась ошибка
	assert.False(t, success, "Expected Filter to return false due to decode error")
	assert.Nil(t, result, "Expected Filter to return nil due to decode error")
}
