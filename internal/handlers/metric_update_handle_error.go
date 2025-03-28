package handlers

import (
	"errors"
	"net/http"

	e "go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/handlers/utils"
)

func metricUpdateHandleError(w http.ResponseWriter, err error) {
	if errors.Is(err, e.ErrMissingMetricID) {
		utils.MakeNotFoundResponse(w, err)
		return
	} else if errors.Is(err, e.ErrInvalidMetricType) {
		utils.MakeBadRequestResponse(w, err)
		return
	} else if errors.Is(err, e.ErrInvalidMetricValue) {
		utils.MakeBadRequestResponse(w, err)
		return
	}
	utils.MakeInternalServerErrorResponse(w, err)
}
