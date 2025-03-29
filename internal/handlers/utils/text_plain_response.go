package utils

import (
	"go-yandex-practicum-metrics/internal/logger"
	"net/http"
)

func MakeTextPlainResponse(w http.ResponseWriter, data []byte) {
	logger.Info("Sending plain text response", "status", http.StatusOK, "data", string(data))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
