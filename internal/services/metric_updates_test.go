package services

import (
	"context"
	"go-yandex-practicum-metrics/internal/errors"

	"go-yandex-practicum-metrics/internal/types"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create *float64 for test cases
func updatesFloat64Ptr(value float64) *float64 {
	return &value
}

// Helper function to create *int64 for test cases
func updatesInt64Ptr(value int64) *int64 {
	return &value
}

func TestMetricUpdatesService_Updates_SuccessfulUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories and transaction
	mockSaveRepo := NewMockMetricUpdatesSaveRepository(ctrl)
	mockFilterRepo := NewMockMetricUpdatesFilterRepository(ctrl)
	mockTx := NewMockTransaction(ctrl)

	// Create service instance
	service := MetricUpdatesService{
		save:        mockSaveRepo,
		filter:      mockFilterRepo,
		transaction: mockTx,
	}

	// Test metrics
	metrics := []*types.Metrics{
		{
			MetricID: types.MetricID{ID: "1", Type: "gauge"},
			Value:    updatesFloat64Ptr(100.0),
		},
		{
			MetricID: types.MetricID{ID: "2", Type: "counter"},
			Delta:    updatesInt64Ptr(10),
		},
	}

	// Set up mocks
	mockTx.EXPECT().Rollback().AnyTimes()         // Since Rollback is deferred
	mockTx.EXPECT().Commit().Return(nil).Times(1) // Commit should succeed

	// Mock Filter and Save methods
	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, false).Times(2) // Metric not found
	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(true).Times(2)           // Save succeeds

	// Call Update
	result, err := service.Updates(context.Background(), metrics)

	// Validate results
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestMetricUpdatesService_Updates_SaveFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories and transaction
	mockSaveRepo := NewMockMetricUpdatesSaveRepository(ctrl)
	mockFilterRepo := NewMockMetricUpdatesFilterRepository(ctrl)
	mockTx := NewMockTransaction(ctrl)

	// Create service instance
	service := MetricUpdatesService{
		save:        mockSaveRepo,
		filter:      mockFilterRepo,
		transaction: mockTx,
	}

	// Test metrics
	metrics := []*types.Metrics{
		{
			MetricID: types.MetricID{ID: "1", Type: "gauge"},
			Value:    updatesFloat64Ptr(100.0),
		},
	}

	// Set up mocks
	mockTx.EXPECT().Rollback().Times(1) // Rollback should always be called
	// Commit should not be called since we have a failure during save
	mockTx.EXPECT().Commit().Times(0)

	// Mock Filter method
	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, false) // Metric not found

	// Mock Save method to fail
	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(false) // Save fails

	// Call Update and assert error
	result, err := service.Updates(context.Background(), metrics)

	// Validate results
	require.Error(t, err)
	assert.Equal(t, errors.ErrMetricInternal, err)
	assert.Nil(t, result)
}

func TestMetricUpdatesService_Updates_TransactionCommitFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories and transaction
	mockSaveRepo := NewMockMetricUpdatesSaveRepository(ctrl)
	mockFilterRepo := NewMockMetricUpdatesFilterRepository(ctrl)
	mockTx := NewMockTransaction(ctrl)

	// Create service instance
	service := MetricUpdatesService{
		save:        mockSaveRepo,
		filter:      mockFilterRepo,
		transaction: mockTx,
	}

	// Test metrics
	metrics := []*types.Metrics{
		{
			MetricID: types.MetricID{ID: "1", Type: "gauge"},
			Value:    updatesFloat64Ptr(100.0),
		},
	}

	// Set up mocks
	mockTx.EXPECT().Rollback().AnyTimes()                              // Rollback should always be called
	mockTx.EXPECT().Commit().Return(errors.ErrMetricInternal).Times(1) // Commit fails

	// Mock Filter method
	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, false) // Metric not found

	// Mock Save method
	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(true) // Save succeeds

	// Call Update and assert error
	result, err := service.Updates(context.Background(), metrics)

	require.Error(t, err)
	assert.Equal(t, errors.ErrMetricInternal, err)
	assert.Nil(t, result)
}

func TestMetricUpdatesService_Updates_MetricNotFoundAndSaved(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories and transaction
	mockSaveRepo := NewMockMetricUpdatesSaveRepository(ctrl)
	mockFilterRepo := NewMockMetricUpdatesFilterRepository(ctrl)
	mockTx := NewMockTransaction(ctrl)

	// Create service instance
	service := MetricUpdatesService{
		save:        mockSaveRepo,
		filter:      mockFilterRepo,
		transaction: mockTx,
	}

	// Test metrics
	metrics := []*types.Metrics{
		{
			MetricID: types.MetricID{ID: "1", Type: "gauge"},
			Value:    updatesFloat64Ptr(100.0),
		},
	}

	// Set up mocks
	mockTx.EXPECT().Rollback().AnyTimes()         // Rollback should always be called
	mockTx.EXPECT().Commit().Return(nil).Times(1) // Commit should succeed

	// Mock Filter method to return nil (metric not found)
	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, false)

	// Mock Save method to return true (save succeeds)
	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)

	// Call Update
	result, err := service.Updates(context.Background(), metrics)

	// Validate results
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "1", result[0].MetricID.ID)
}

func TestMetricUpdatesService_Updates_DeferRollback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories and transaction
	mockSaveRepo := NewMockMetricUpdatesSaveRepository(ctrl)
	mockFilterRepo := NewMockMetricUpdatesFilterRepository(ctrl)
	mockTx := NewMockTransaction(ctrl)

	// Create service instance
	service := MetricUpdatesService{
		save:        mockSaveRepo,
		filter:      mockFilterRepo,
		transaction: mockTx,
	}

	// Test metrics
	metrics := []*types.Metrics{
		{
			MetricID: types.MetricID{ID: "1", Type: "gauge"},
			Value:    updatesFloat64Ptr(100.0),
		},
	}

	// Expect Rollback to be called once even if no error occurs
	mockTx.EXPECT().Rollback().Times(1)
	mockTx.EXPECT().Commit().Times(1) // Commit will be called if everything is fine

	// Mock Filter method to not find the metric (simulating a "not found" scenario)
	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, false)

	// Mock Save method to succeed
	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)

	// Call Update and assert no error
	result, err := service.Updates(context.Background(), metrics)

	// Validate results
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, *metrics[0], *result[0])
}

func TestMetricUpdatesService_Updates_GaugeAndCounter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories and transaction
	mockSaveRepo := NewMockMetricUpdatesSaveRepository(ctrl)
	mockFilterRepo := NewMockMetricUpdatesFilterRepository(ctrl)
	mockTx := NewMockTransaction(ctrl)

	// Create service instance
	service := MetricUpdatesService{
		save:        mockSaveRepo,
		filter:      mockFilterRepo,
		transaction: mockTx,
	}

	// Test metrics
	gaugeMetric := &types.Metrics{
		MetricID: types.MetricID{ID: "1", Type: "gauge"},
		Value:    updatesFloat64Ptr(100.0),
	}

	counterMetric := &types.Metrics{
		MetricID: types.MetricID{ID: "2", Type: "counter"},
		Delta:    updateInt64Ptr(10),
	}

	// Mock Filter to return an existing metric for both metrics
	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any()).Return(gaugeMetric, true).Times(1)
	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any()).Return(counterMetric, true).Times(1)

	// Mock Save method to succeed for both metrics
	mockSaveRepo.EXPECT().Save(gomock.Any(), gaugeMetric, gomock.Any()).Return(true).Times(1)
	mockSaveRepo.EXPECT().Save(gomock.Any(), counterMetric, gomock.Any()).Return(true).Times(1)

	// Expect Commit to be called after all operations
	mockTx.EXPECT().Commit().Times(1)

	// Expect Rollback to be called in case of error in the middle (Rollback happens at the end)
	mockTx.EXPECT().Rollback().Times(1)

	// Call Update for both metrics
	metrics := []*types.Metrics{gaugeMetric, counterMetric}
	result, err := service.Updates(context.Background(), metrics)

	// Validate results
	require.NoError(t, err)
	assert.Len(t, result, 2)

	// Check the values of the updated metrics
	assert.Equal(t, *result[0].Value, 100.0)     // Gauge metric value should be set
	assert.Equal(t, *result[1].Delta, int64(20)) // Counter metric delta should be updated (10 + 10)
}

func TestMetricUpdatesService_Updates_RollbackFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repositories and transaction
	mockSaveRepo := NewMockMetricUpdatesSaveRepository(ctrl)
	mockFilterRepo := NewMockMetricUpdatesFilterRepository(ctrl)
	mockTx := NewMockTransaction(ctrl)

	// Create service instance
	service := MetricUpdatesService{
		save:        mockSaveRepo,
		filter:      mockFilterRepo,
		transaction: mockTx,
	}

	// Test metrics
	metrics := []*types.Metrics{
		{
			MetricID: types.MetricID{ID: "1", Type: "gauge"},
			Value:    updatesFloat64Ptr(100.0),
		},
	}

	// Mock the Rollback to fail
	mockTx.EXPECT().Rollback().Return(errors.ErrMetricInternal).Times(1) // Simulate rollback failure

	// Mock Commit to succeed
	mockTx.EXPECT().Commit().Return(nil).Times(1)

	// Mock Filter method to return nil (metric not found)
	mockFilterRepo.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, false)

	// Mock Save method to succeed
	mockSaveRepo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)

	// Call Update and assert no error
	result, err := service.Updates(context.Background(), metrics)

	// Validate results
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, *metrics[0], *result[0])
}
