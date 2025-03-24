package usecases

import (
	"context"
	"testing"

	e "go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdatePathUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the MetricUpdatePathService
	mockService := NewMockMetricUpdatePathService(ctrl)

	// Create the usecase instance with the mocked service
	usecase := &MetricUpdatePathUsecase{
		svc: mockService,
	}

	// Test cases for successful execution
	tests := []struct {
		name           string
		req            *types.MetricUpdatePathRequest
		mockUpdateResp *types.Metrics
		expectedResp   *types.MetricUpdatePathResponse
	}{
		{
			name: "valid counter update",
			req: &types.MetricUpdatePathRequest{
				Type:  "counter",
				Name:  "metric1",
				Value: "10",
			},
			mockUpdateResp: &types.Metrics{
				MetricID: types.MetricID{
					ID:   "metric1",
					Type: types.Counter,
				},
				Delta: metricUpdatePtrInt64(10),
			},
			expectedResp: &types.MetricUpdatePathResponse{
				MetricID: types.MetricID{
					Type: types.Counter,
					ID:   "metric1",
				},
				Delta: metricUpdatePtrInt64(10),
				Value: nil,
			},
		},
		{
			name: "valid gauge update",
			req: &types.MetricUpdatePathRequest{
				Type:  "gauge",
				Name:  "metric2",
				Value: "3.14",
			},
			mockUpdateResp: &types.Metrics{
				MetricID: types.MetricID{
					ID:   "metric2",
					Type: types.Gauge,
				},
				Value: metricUpdatePtrFloat64(3.14),
			},
			expectedResp: &types.MetricUpdatePathResponse{
				MetricID: types.MetricID{
					Type: types.Gauge,
					ID:   "metric2",
				},
				Delta: nil,
				Value: metricUpdatePtrFloat64(3.14),
			},
		},
	}

	// Run successful test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock behavior
			mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tt.mockUpdateResp, nil)

			// Call Execute
			resp, err := usecase.Execute(context.Background(), tt.req)

			// Assert the response and error
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}

// Helper functions for setting up mock services and testing logic

func TestMetricUpdatePathUsecase_Execute_ValidCounterUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the MetricUpdatePathService
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := &MetricUpdatePathUsecase{
		svc: mockService,
	}

	// Test case where the counter update is valid
	req := &types.MetricUpdatePathRequest{
		Type:  "counter",
		Name:  "metric1",
		Value: "10",
	}
	mockUpdateResp := &types.Metrics{
		MetricID: types.MetricID{
			ID:   "metric1",
			Type: types.Counter,
		},
		Delta: metricUpdatePtrInt64(10),
	}
	expectedResp := &types.MetricUpdatePathResponse{
		MetricID: types.MetricID{
			Type: types.Counter,
			ID:   "metric1",
		},
		Delta: metricUpdatePtrInt64(10),
		Value: nil,
	}

	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(mockUpdateResp, nil)

	// Execute the use case
	resp, err := usecase.Execute(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}

func TestMetricUpdatePathUsecase_Execute_ValidGaugeUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the MetricUpdatePathService
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := &MetricUpdatePathUsecase{
		svc: mockService,
	}

	// Test case where the gauge update is valid
	req := &types.MetricUpdatePathRequest{
		Type:  "gauge",
		Name:  "metric2",
		Value: "3.14",
	}
	mockUpdateResp := &types.Metrics{
		MetricID: types.MetricID{
			ID:   "metric2",
			Type: types.Gauge,
		},
		Value: metricUpdatePtrFloat64(3.14),
	}
	expectedResp := &types.MetricUpdatePathResponse{
		MetricID: types.MetricID{
			Type: types.Gauge,
			ID:   "metric2",
		},
		Delta: nil,
		Value: metricUpdatePtrFloat64(3.14),
	}

	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(mockUpdateResp, nil)

	// Execute the use case
	resp, err := usecase.Execute(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
}

func TestMetricUpdatePathUsecase_Execute_InvalidCounterValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the MetricUpdatePathService
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := &MetricUpdatePathUsecase{
		svc: mockService,
	}

	// Test case where the counter value is invalid
	req := &types.MetricUpdatePathRequest{
		Type:  "counter",
		Name:  "metric1",
		Value: "invalid", // Invalid value
	}

	// We don't expect the Update function to be called
	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

	// Execute the use case
	resp, err := usecase.Execute(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, e.ErrInvalidCounterValue, err)
}

func TestMetricUpdatePathUsecase_Execute_InvalidGaugeValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the MetricUpdatePathService
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := &MetricUpdatePathUsecase{
		svc: mockService,
	}

	// Test case where the gauge value is invalid
	req := &types.MetricUpdatePathRequest{
		Type:  "gauge",
		Name:  "metric2",
		Value: "invalid", // Invalid value for gauge
	}

	// We don't expect the Update function to be called
	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

	// Execute the use case
	resp, err := usecase.Execute(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, e.ErrInvalidGaugeValue, err)
}

func TestMetricUpdatePathUsecase_Execute_UpdateServiceFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the MetricUpdatePathService
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := &MetricUpdatePathUsecase{
		svc: mockService,
	}

	// Test case where the update service fails
	req := &types.MetricUpdatePathRequest{
		Type:  "gauge",
		Name:  "metric3",
		Value: "5.67",
	}

	mockUpdateErr := e.ErrMetricInternal
	// We expect the Update method to be called
	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, mockUpdateErr)

	// Execute the use case
	resp, err := usecase.Execute(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, mockUpdateErr, err)
}

// Helper functions to create pointers to basic types
func metricUpdatePtrInt64(i int64) *int64 {
	return &i
}

func metricUpdatePtrFloat64(f float64) *float64 {
	return &f
}

func TestMetricUpdatePathUsecase_Execute_ValidateTypeFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the MetricUpdatePathService
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := &MetricUpdatePathUsecase{
		svc: mockService,
	}

	// Test case where the ValidateType fails
	// We assume "invalid_type" is an invalid type
	req := &types.MetricUpdatePathRequest{
		Type:  "invalid_type", // Invalid type
		Name:  "metric1",
		Value: "10",
	}

	// No expectation for Update call, because validation fails
	// Execute the use case
	resp, err := usecase.Execute(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, e.ErrInvalidMetricType, err) // We expect the validation error to be returned
}

func TestMetricUpdatePathUsecase_Execute_ValidateNameFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the MetricUpdatePathService
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := &MetricUpdatePathUsecase{
		svc: mockService,
	}

	// Test case where the ValidateName fails
	// We assume an empty name will fail the validation
	req := &types.MetricUpdatePathRequest{
		Type:  "counter",
		Name:  "", // Invalid name (empty string)
		Value: "10",
	}

	// No expectation for Update call, because validation fails
	// Execute the use case
	resp, err := usecase.Execute(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, e.ErrMetricNameRequired, err) // We expect the validation error to be returned
}

func TestMetricUpdatePathUsecase_Execute_ValidateValueFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the MetricUpdatePathService
	mockService := NewMockMetricUpdatePathService(ctrl)
	usecase := &MetricUpdatePathUsecase{
		svc: mockService,
	}

	// Test case where the ValidateValue fails
	// We assume an empty value will fail the validation
	req := &types.MetricUpdatePathRequest{
		Type:  "counter",
		Name:  "metric1",
		Value: "", // Invalid value (empty string)
	}

	// No expectation for Update call, because validation fails
	// Execute the use case
	resp, err := usecase.Execute(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, e.ErrMetricValueRequired, err) // We expect the validation error to be returned
}
