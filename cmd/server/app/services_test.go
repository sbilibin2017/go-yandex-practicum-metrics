package app_test

import (
	"go-yandex-practicum-metrics/cmd/server/app"
	"go-yandex-practicum-metrics/internal/services"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock RepoProvider
	mockRepoProvider := app.NewMockRepoProvider(ctrl)
	// Mock Transaction
	mockTx := app.NewMockTransaction(ctrl)

	// Mock the methods of the RepoProvider interface
	mockSaveRepo := app.NewMockServiceSaveRepository(ctrl)
	mockFilterRepo := app.NewMockServiceFilterRepository(ctrl)
	mockListRepo := app.NewMockServiceListRepository(ctrl)

	// Set expectations for the mocks
	mockRepoProvider.EXPECT().GetSaveRepository().Return(mockSaveRepo).Times(2)     // GetSaveRepository() called twice
	mockRepoProvider.EXPECT().GetFilterRepository().Return(mockFilterRepo).Times(3) // GetFilterRepository() called three times
	mockRepoProvider.EXPECT().GetListRepository().Return(mockListRepo).Times(1)     // GetListRepository() called once

	// Call the service creation function
	metricUpdateService, metricUpdatesService, metricGetService, metricListService, err := app.NewServices(mockRepoProvider, mockTx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, metricUpdateService)
	assert.IsType(t, &services.MetricUpdateService{}, metricUpdateService)
	assert.NotNil(t, metricUpdatesService)
	assert.IsType(t, &services.MetricUpdatesService{}, metricUpdatesService)
	assert.NotNil(t, metricGetService)
	assert.IsType(t, &services.MetricGetService{}, metricGetService)
	assert.NotNil(t, metricListService)
	assert.IsType(t, &services.MetricListService{}, metricListService)
}
