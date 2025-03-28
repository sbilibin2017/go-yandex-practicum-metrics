package utils

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandlerToRouterHandle(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}
