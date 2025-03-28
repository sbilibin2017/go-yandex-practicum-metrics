package routers

import (
	"go-yandex-practicum-metrics/internal/routers/utils"

	"github.com/julienschmidt/httprouter"
)

type Logger interface {
	Infow(msg string, args ...any)
}

func RegisterMetricUpdatePathRouter(
	r *httprouter.Router,
	h httprouter.Handle,
	logger Logger,
) {
	r.POST(
		"/update/:type/:name/:value",
		utils.HandlerToRouterHandle(utils.HandlerWithLoggingMiddleware(h, logger)),
	)
}
