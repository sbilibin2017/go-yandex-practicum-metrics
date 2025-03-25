package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-yandex-practicum-metrics/internal/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdatePathHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockMetricUpdatePathUsecase(ctrl)

	expectedResp := types.MetricUpdatePathResponse("Updated metric successfully")

	mockUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(&expectedResp, nil)

	req, err := http.NewRequest("GET", "/metric/some-name/some-type/some-value", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := MetricUpdatePathHandler(mockUsecase)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, string(expectedResp), rr.Body.String())
}

func TestMetricUpdatePathHandler_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockMetricUpdatePathUsecase(ctrl)
	mockUsecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))

	req, err := http.NewRequest("GET", "/metric/some-name/some-type/some-value", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	handler := MetricUpdatePathHandler(mockUsecase)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
