package utils

import "net/http"

func MakeTextPlainResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
