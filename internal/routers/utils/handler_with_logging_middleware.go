package utils

import (
	"go-yandex-practicum-metrics/internal/middlewares"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandlerWithLoggingMiddleware(h httprouter.Handle, logger middlewares.LoggingMiddlewareLogger) http.Handler {
	return middlewares.LoggingMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, httprouter.ParamsFromContext(r.Context()))
	}))
}
