package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeRequestBody[T any](w http.ResponseWriter, r *http.Request, req T) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return err
	}
	return nil
}

func EncodeResponseBody[T any](w http.ResponseWriter, resp T) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}
