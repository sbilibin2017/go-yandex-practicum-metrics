package usecases

import (
	"context"
	"errors"
	"testing"

	"go-yandex-practicum-metrics/internal/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdatesBodyUsecase_Execute_EmptyRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMetricUpdatesBodyService(ctrl)
	usecase := NewMetricUpdatesBodyUsecase(mockService)

	// Test case where request is empty
	req := &MetricUpdatesBodyRequest{}
	_, err := usecase.Execute(context.Background(), req)

	// Assert error is expected
	assert.EqualError(t, err, ErrMetricUpdatesBodyNotProvider.Error())
}

func TestMetricUpdatesBodyUsecase_Execute_MissingMetricID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMetricUpdatesBodyService(ctrl)
	usecase := NewMetricUpdatesBodyUsecase(mockService)

	// Test case where metric ID is missing
	req := &MetricUpdatesBodyRequest{
		&MetricUpdateBodyRequest{ID: "", Type: "gauge", Value: nil},
	}
	_, err := usecase.Execute(context.Background(), req)

	// Assert error is expected
	assert.EqualError(t, err, ErrMetricUpdatesBodyIDRequired.Error())
}

func TestMetricUpdatesBodyUsecase_Execute_MissingMetricType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMetricUpdatesBodyService(ctrl)
	usecase := NewMetricUpdatesBodyUsecase(mockService)

	// Test case where metric type is missing
	req := &MetricUpdatesBodyRequest{
		&MetricUpdateBodyRequest{ID: "test_id", Type: "", Value: nil},
	}
	_, err := usecase.Execute(context.Background(), req)

	// Assert error is expected
	assert.EqualError(t, err, ErrMetricUpdatesBodyTypeRequired.Error())
}

func TestMetricUpdatesBodyUsecase_Execute_InvalidMetricType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMetricUpdatesBodyService(ctrl)
	usecase := NewMetricUpdatesBodyUsecase(mockService)

	// Test case where metric type is invalid
	req := &MetricUpdatesBodyRequest{
		&MetricUpdateBodyRequest{ID: "test_id", Type: "invalid", Value: nil},
	}
	_, err := usecase.Execute(context.Background(), req)

	// Assert error is expected
	assert.EqualError(t, err, ErrMetricUpdatesBodyInvalidType.Error())
}

func TestMetricUpdatesBodyUsecase_Execute_InvalidCounterValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMetricUpdatesBodyService(ctrl)
	usecase := NewMetricUpdatesBodyUsecase(mockService)

	// Test case where counter type is missing Delta value
	req := &MetricUpdatesBodyRequest{
		&MetricUpdateBodyRequest{ID: "test_id", Type: "counter", Delta: nil},
	}
	_, err := usecase.Execute(context.Background(), req)

	// Assert error is expected
	assert.EqualError(t, err, ErrMetricUpdatesBodyInvalidValue.Error())
}

func TestMetricUpdatesBodyUsecase_Execute_InvalidGaugeValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMetricUpdatesBodyService(ctrl)
	usecase := NewMetricUpdatesBodyUsecase(mockService)

	// Test case where gauge type is missing Value
	req := &MetricUpdatesBodyRequest{
		&MetricUpdateBodyRequest{ID: "test_id", Type: "gauge", Value: nil},
	}
	_, err := usecase.Execute(context.Background(), req)

	// Assert error is expected
	assert.EqualError(t, err, ErrMetricUpdatesBodyInvalidValue.Error())
}

func TestMetricUpdatesBodyUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMetricUpdatesBodyService(ctrl)
	usecase := NewMetricUpdatesBodyUsecase(mockService)

	// Test case where everything is correct
	req := &MetricUpdatesBodyRequest{
		&MetricUpdateBodyRequest{ID: "test_id", Type: "gauge", Value: float64Ptr(10.5)},
	}

	// Mock the service response
	metrics := []*domain.Metrics{
		{
			MetricID: domain.MetricID{
				ID:   "test_id",
				Type: domain.Gauge,
			},
			Value: float64Ptr(10.5),
		},
	}

	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(metrics, nil)

	response, err := usecase.Execute(context.Background(), req)

	// Assert no error and correct response
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Len(t, *response, 1)
	assert.Equal(t, "test_id", (*response)[0].ID)
	assert.Equal(t, "gauge", (*response)[0].Type)
	assert.Equal(t, 10.5, *(*response)[0].Value)
}

func TestMetricUpdatesBodyUsecase_Execute_AssignDeltaToMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMetricUpdatesBodyService(ctrl)
	usecase := NewMetricUpdatesBodyUsecase(mockService)

	// Test case where Delta is assigned to the metric
	delta := int64(42)
	req := &MetricUpdatesBodyRequest{
		&MetricUpdateBodyRequest{
			ID:    "test_id",
			Type:  "counter",
			Delta: &delta,
		},
	}

	// Create the expected metric to check Delta assignment
	expectedMetric := &domain.Metrics{
		MetricID: domain.MetricID{
			ID:   "test_id",
			Type: domain.Counter,
		},

		Delta: &delta, // This is the value we're checking for assignment
	}

	// Mock the service response
	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return([]*domain.Metrics{expectedMetric}, nil)

	// Call the method
	response, err := usecase.Execute(context.Background(), req)

	// Assert no error and correct response
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Len(t, *response, 1)

	// Check that Delta has been assigned correctly
	assert.Equal(t, delta, *(*response)[0].Delta)
}

func TestMetricUpdatesBodyUsecase_Execute_UpdateServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := NewMockMetricUpdatesBodyService(ctrl)
	usecase := NewMetricUpdatesBodyUsecase(mockService)

	// Test case where service returns an error
	req := &MetricUpdatesBodyRequest{
		&MetricUpdateBodyRequest{
			ID:    "test_id",
			Type:  "gauge",
			Value: float64Ptr(10),
		},
	}

	// Mock the service to return an error
	mockService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))

	// Call the method
	_, err := usecase.Execute(context.Background(), req)

	// Assert the error is the internal error we expect
	assert.EqualError(t, err, ErrMetricUpdatesBodyInternal.Error())
}

func float64Ptr(val float64) *float64 {
	return &val
}
