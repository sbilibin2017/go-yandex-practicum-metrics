package services

import (
	"context"
	"testing"

	"go-yandex-practicum-metrics/internal/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-yandex-practicum-metrics/internal/types"
)

func TestMetricGetByTypeAndIDService_GetByTypeAndID_Found(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock repository
	mockFilterRepo := NewMockMetricGetByTypeAndIDFilterRepository(ctrl)

	// Create service instance
	service := MetricGetByTypeAndIDService{
		filter: mockFilterRepo,
	}

	// Define the test MetricID
	metricID := types.MetricID{ID: "1", Type: "gauge"}

	// Mock the Filter method to return a metric
	expectedMetric := &types.Metrics{
		MetricID: metricID,
		Value:    getFloat64Ptr(100.0),
	}
	mockFilterRepo.EXPECT().Filter(gomock.Any(), metricID).Return(expectedMetric, true)

	// Call GetByTypeAndID
	metric, err := service.GetByTypeAndID(context.Background(), metricID)

	// Validate that no error occurred
	require.NoError(t, err)

	// Assert that the returned metric is the expected one
	assert.Equal(t, expectedMetric, metric)
}

func TestMetricGetByTypeAndIDService_GetByTypeAndID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock repository
	mockFilterRepo := NewMockMetricGetByTypeAndIDFilterRepository(ctrl)

	// Create service instance
	service := MetricGetByTypeAndIDService{
		filter: mockFilterRepo,
	}

	// Define the test MetricID
	metricID := types.MetricID{ID: "1", Type: "gauge"}

	// Mock the Filter method to return nothing (not found)
	mockFilterRepo.EXPECT().Filter(gomock.Any(), metricID).Return(nil, false)

	// Call GetByTypeAndID
	metric, err := service.GetByTypeAndID(context.Background(), metricID)

	// Validate that the error is ErrMetricNotFound
	require.Error(t, err)
	assert.Equal(t, errors.ErrMetricNotFound, err)

	// Assert that the returned metric is nil
	assert.Nil(t, metric)
}

// Helper function to return a pointer to a float64 value
func getFloat64Ptr(f float64) *float64 {
	return &f
}
