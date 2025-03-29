package utils

import (
	"go-yandex-practicum-metrics/internal/logger"
	"net/http"
)

func MakeNotFoundResponse(w http.ResponseWriter, err error) {
	logger.Error("Not found error", "error", err)
	http.Error(w, err.Error(), http.StatusNotFound)
}

func MakeBadRequestResponse(w http.ResponseWriter, err error) {
	logger.Error("Bad request error", "error", err)
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func MakeInternalServerErrorResponse(w http.ResponseWriter, err error) {
	logger.Error("Internal server error", "error", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
