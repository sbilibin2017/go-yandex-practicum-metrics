package routers

import (
	"go-yandex-practicum-metrics/internal/logger"
	"go-yandex-practicum-metrics/internal/router"

	"github.com/julienschmidt/httprouter"
)

func RegisterMetricUpdatePathRouter(
	r *router.Router,
	h httprouter.Handle,
) {
	logger.Info("Registering MetricUpdatePath router", "method", "POST", "path", "/update/:type/:name/:value")
	r.AddHandler(
		"POST",
		"/update/:type/:name/:value",
		h,
	)
}
