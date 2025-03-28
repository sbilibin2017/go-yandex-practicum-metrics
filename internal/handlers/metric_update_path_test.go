package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	e "go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/responses"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricUpdatePathHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUseCase := NewMockMetricUpdatePathUsecase(ctrl)
	tests := []struct {
		name           string
		mockResponse   *responses.MetricUpdatePathResponse
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful case",
			mockResponse:   &responses.MetricUpdatePathResponse{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Metric updated successfully",
		},
		{
			name:           "Missing Metric ID",
			mockResponse:   nil,
			mockError:      e.ErrMissingMetricID,
			expectedStatus: http.StatusNotFound,
			expectedBody:   e.ErrMissingMetricID.Error(),
		},
		{
			name:           "Invalid Metric Type",
			mockResponse:   nil,
			mockError:      e.ErrInvalidMetricType,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   e.ErrInvalidMetricType.Error(),
		},
		{
			name:           "Invalid Metric Value",
			mockResponse:   nil,
			mockError:      e.ErrInvalidMetricValue,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   e.ErrInvalidMetricValue.Error(),
		},
		{
			name:           "General error",
			mockResponse:   nil,
			mockError:      errors.New("internal error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase.EXPECT().
				Execute(gomock.Any(), gomock.Any()).
				Return(tt.mockResponse, tt.mockError).
				Times(1)
			router := httprouter.New()
			handler := MetricUpdatePathHandler(mockUseCase)
			router.POST("/update/:type/:name/:value", handler)
			req, err := http.NewRequest("POST", "/update/some_type/some_name/some_value", nil)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, strings.TrimSpace(tt.expectedBody), strings.TrimSpace(rr.Body.String()))
		})
	}
}
