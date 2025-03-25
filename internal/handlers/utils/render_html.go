package utils

import "net/http"

func RenderHTML(w http.ResponseWriter, htmlContent string) error {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	_, err := w.Write([]byte(htmlContent))
	if err != nil {
		http.Error(w, "Failed to write HTML response", http.StatusInternalServerError)
		return err
	}
	return nil
}
