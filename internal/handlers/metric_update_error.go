package handlers

import (
	"errors"
	e "go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/handlers/utils"
	"net/http"
)

func MetricErrorResponse(w http.ResponseWriter, err error) {
	var statusCode int
	var errMsg string

	switch {
	case errors.Is(err, e.ErrMetricTypeRequired):
		statusCode = http.StatusBadRequest
		errMsg = err.Error()
	case errors.Is(err, e.ErrMetricInvalidType):
		statusCode = http.StatusBadRequest
		errMsg = err.Error()
	case errors.Is(err, e.ErrMetricIDRequired):
		statusCode = http.StatusNotFound
		errMsg = err.Error()
	case errors.Is(err, e.ErrMetricValueRequired):
		statusCode = http.StatusBadRequest
		errMsg = err.Error()
	case errors.Is(err, e.ErrMetricInvalidDelta):
		statusCode = http.StatusBadRequest
		errMsg = err.Error()
	case errors.Is(err, e.ErrMetricInvalidValue):
		statusCode = http.StatusBadRequest
		errMsg = err.Error()
	case errors.Is(err, e.ErrMetricNotFound):
		statusCode = http.StatusNotFound
		errMsg = err.Error()
	default:
		statusCode = http.StatusInternalServerError
		errMsg = "internal server error"
	}
	utils.ErrorResponse(w, errMsg, statusCode)
}
