package services

import (
	"context"
	"testing"

	"go-yandex-practicum-metrics/internal/errors"

	"go-yandex-practicum-metrics/internal/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdateService_Update(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories
	mockSaveRepo := NewMockMetricUpdateSaveRepository(ctrl)
	mockFilterRepo := NewMockMetricUpdateFilterRepository(ctrl)

	// Create a MetricUpdateService instance with the mocked repositories
	service := MetricUpdateService{
		save:   mockSaveRepo,
		filter: mockFilterRepo,
	}

	// Define test cases
	tests := []struct {
		name           string
		metric         *types.Metrics
		filterReturn   *types.Metrics
		filterFound    bool
		saveReturn     bool
		expectedResult *types.Metrics
		expectedError  error
	}{
		{
			name: "Metric does not exist, save new metric",
			metric: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "gauge"},
				Value:    updateFloat64Ptr(123.45),
			},
			filterReturn:   nil,
			filterFound:    false,
			saveReturn:     true,
			expectedResult: &types.Metrics{MetricID: types.MetricID{ID: "1", Type: "gauge"}, Value: updateFloat64Ptr(123.45)},
			expectedError:  nil,
		},
		{
			name: "Metric exists, update value",
			metric: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "gauge"},
				Value:    updateFloat64Ptr(456.78),
			},
			filterReturn: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "gauge"},
				Value:    updateFloat64Ptr(123.45),
			},
			filterFound: true,
			saveReturn:  true,
			expectedResult: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "gauge"},
				Value:    updateFloat64Ptr(456.78),
			},
			expectedError: nil,
		},
		{
			name: "Metric exists, but save fails",
			metric: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "gauge"},
				Value:    updateFloat64Ptr(456.78),
			},
			filterReturn: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "gauge"},
				Value:    updateFloat64Ptr(123.45),
			},
			filterFound:    true,
			saveReturn:     false,
			expectedResult: nil,
			expectedError:  errors.ErrMetricInternal,
		},
		{
			name: "Counter type metric, update delta",
			metric: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "counter"},
				Delta:    updateInt64Ptr(10),
			},
			filterReturn: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "counter"},
				Delta:    updateInt64Ptr(5),
			},
			filterFound: true,
			saveReturn:  true,
			expectedResult: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "counter"},
				Delta:    updateInt64Ptr(15),
			},
			expectedError: nil,
		},
		{
			name: "Counter type metric, save fails",
			metric: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "counter"},
				Delta:    updateInt64Ptr(10),
			},
			filterReturn: &types.Metrics{
				MetricID: types.MetricID{ID: "1", Type: "counter"},
				Delta:    updateInt64Ptr(5),
			},
			filterFound:    true,
			saveReturn:     false,
			expectedResult: nil,
			expectedError:  errors.ErrMetricInternal,
		},
	}

	// Loop through each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up expectations for the filter repository
			mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(tt.filterReturn, tt.filterFound)

			// Set up expectations for the save repository
			mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(tt.saveReturn)

			// Call the Update method
			result, err := service.Update(context.Background(), tt.metric)

			// Assert the results
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestMetricUpdateService_Update_SaveFails(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories
	mockSaveRepo := NewMockMetricUpdateSaveRepository(ctrl)
	mockFilterRepo := NewMockMetricUpdateFilterRepository(ctrl)

	// Create a MetricUpdateService instance with the mocked repositories
	service := MetricUpdateService{
		save:   mockSaveRepo,
		filter: mockFilterRepo,
	}

	// Define test case where Save fails
	metric := &types.Metrics{
		MetricID: types.MetricID{ID: "1", Type: "gauge"},
		Value:    updateFloat64Ptr(100.0),
	}

	// Set up mock expectations
	// Mock Filter method to return false (meaning metric not found)
	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(nil, false)

	// Mock Save method to return false (simulating a save failure)
	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(false)

	// Call the Update method
	result, err := service.Update(context.Background(), metric)

	// Assert that the result is nil and the expected error is returned
	assert.Nil(t, result)
	assert.Equal(t, errors.ErrMetricInternal, err)
}

// Helper function to create *float64 for test cases
func updateFloat64Ptr(value float64) *float64 {
	return &value
}

// Helper function to create *int64 for test cases
func updateInt64Ptr(value int64) *int64 {
	return &value
}
