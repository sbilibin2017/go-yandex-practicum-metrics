package routers

import (
	"github.com/julienschmidt/httprouter"
)

func RegisterMetricUpdatePathRouter(r *httprouter.Router, h httprouter.Handle) {
	r.POST("/update/:type/:name/:value", h)
}
