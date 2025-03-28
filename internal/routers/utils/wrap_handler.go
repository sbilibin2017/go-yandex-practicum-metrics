package utils

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func WrapHandler(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(w, r)
	}
}
