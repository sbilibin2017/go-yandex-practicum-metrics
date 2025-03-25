package utils

import "net/http"

func SendTextResponse(w http.ResponseWriter, data string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
