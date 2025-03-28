package usecases

import (
	"context"
	"errors"
	"go-yandex-practicum-metrics/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdatePathUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := NewMockMetricUpdatePathService(ctrl)
	mockRequester := NewMockMetricUpdateRequester(ctrl)
	mockResponser := NewMockMetricUpdateResponser(ctrl)
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			Type:  domain.Gauge,
			Value: new(float64),
		},
	}
	mockRequester.EXPECT().Validate().Return(nil)
	mockRequester.EXPECT().ToDomain().Return(metrics)
	mockSvc.EXPECT().Update(gomock.Any(), metrics).Return(metrics, nil)
	expectedResponse := []byte("Metric updated successfully")
	mockResponser.EXPECT().ToResponse().Return(&expectedResponse)
	uc := NewMetricUpdatePathUsecase(mockSvc, mockRequester, mockResponser)
	resp, err := uc.Execute(context.Background())
	require.NoError(t, err)
	assert.Equal(t, &expectedResponse, resp)
}

func TestMetricUpdatePathUsecase_Execute_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := NewMockMetricUpdatePathService(ctrl)
	mockRequester := NewMockMetricUpdateRequester(ctrl)
	mockResponser := NewMockMetricUpdateResponser(ctrl)
	mockRequester.EXPECT().Validate().Return(errors.New("validation error"))
	uc := NewMetricUpdatePathUsecase(mockSvc, mockRequester, mockResponser)
	resp, err := uc.Execute(context.Background())
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestMetricUpdatePathUsecase_Execute_UpdateServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := NewMockMetricUpdatePathService(ctrl)
	mockRequester := NewMockMetricUpdateRequester(ctrl)
	mockResponser := NewMockMetricUpdateResponser(ctrl)
	mockRequester.EXPECT().Validate().Return(nil)
	mockRequester.EXPECT().ToDomain().Return([]*domain.Metrics{{}})
	mockSvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, errors.New("service update error"))
	uc := NewMetricUpdatePathUsecase(mockSvc, mockRequester, mockResponser)
	resp, err := uc.Execute(context.Background())
	require.Error(t, err)
	assert.Nil(t, resp)
}
