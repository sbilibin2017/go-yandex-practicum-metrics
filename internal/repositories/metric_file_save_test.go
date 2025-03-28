package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveMetricsWithNilWriter(t *testing.T) {
	repo := NewMetricFileSaveBatchRepository(nil)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Value: ptrFloat64(100),
		},
	}
	result := repo.Save(context.Background(), metrics)
	assert.False(t, result)
}

func TestSaveFileMetricsSuccessfully(t *testing.T) {
	buf := new(bytes.Buffer)
	repo := NewMetricFileSaveBatchRepository(buf)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Value: ptrFloat64(100),
		},
		{
			ID:    "metric2",
			Type:  "gauge",
			Value: ptrFloat64(200),
		},
	}
	result := repo.Save(context.Background(), metrics)
	assert.True(t, result)

	// Verify the content in the buffer
	var savedMetrics []*domain.Metrics
	decoder := json.NewDecoder(buf)
	for {
		var metric domain.Metrics
		if err := decoder.Decode(&metric); err != nil {
			break
		}
		savedMetrics = append(savedMetrics, &metric)
	}
	assert.Equal(t, 2, len(savedMetrics))
	assert.Equal(t, "metric1", savedMetrics[0].ID)
	assert.Equal(t, "metric2", savedMetrics[1].ID)
}

func TestSaveMetricsWithWriteError(t *testing.T) {
	errorWriter := &errorWriter{}
	repo := NewMetricFileSaveBatchRepository(errorWriter)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  "counter",
			Value: ptrFloat64(100),
		},
	}
	result := repo.Save(context.Background(), metrics)
	assert.False(t, result)
}

type errorWriter struct{}

func (ew *errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}
