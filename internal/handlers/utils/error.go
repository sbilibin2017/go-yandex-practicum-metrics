package utils

import "net/http"

func ErrorResponse(w http.ResponseWriter, msg string, code int) {
	http.Error(w, msg, code)
}
