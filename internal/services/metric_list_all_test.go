package services

import (
	"context"
	"testing"

	"go-yandex-practicum-metrics/internal/errors"

	"go-yandex-practicum-metrics/internal/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricListAllService_ListAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock repository
	mockListRepo := NewMockMetricListAllRepository(ctrl)

	// Create the service instance
	service := MetricListAllService{
		list: mockListRepo,
	}

	// Define the expected metrics
	expectedMetrics := []*types.Metrics{
		{
			MetricID: types.MetricID{ID: "1", Type: "gauge"},
			Value:    listAllFloat64Ptr(10.0),
		},
		{
			MetricID: types.MetricID{ID: "2", Type: "counter"},
			Delta:    listAllInt64Ptr(5),
		},
	}

	// Mock the ListAll method to return the expected metrics
	mockListRepo.EXPECT().ListAll(gomock.Any()).Return(expectedMetrics, true)

	// Call ListAll
	metrics, err := service.ListAll(context.Background())

	// Validate that no error occurred
	require.NoError(t, err)

	// Assert that the returned metrics are the expected ones
	assert.Equal(t, expectedMetrics, metrics)
}

func TestMetricListAllService_ListAll_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create the mock repository
	mockListRepo := NewMockMetricListAllRepository(ctrl)

	// Create the service instance
	service := MetricListAllService{
		list: mockListRepo,
	}

	// Mock the ListAll method to return false (failure)
	mockListRepo.EXPECT().ListAll(gomock.Any()).Return(nil, false)

	// Call ListAll
	metrics, err := service.ListAll(context.Background())

	// Validate that the error is ErrMetricInternal
	require.Error(t, err)
	assert.Equal(t, errors.ErrMetricInternal, err)

	// Assert that the returned metrics are nil
	assert.Nil(t, metrics)
}

// Helper function to return a pointer to a float64 value
func listAllFloat64Ptr(f float64) *float64 {
	return &f
}

// Helper function to return a pointer to an int64 value
func listAllInt64Ptr(i int64) *int64 {
	return &i
}
