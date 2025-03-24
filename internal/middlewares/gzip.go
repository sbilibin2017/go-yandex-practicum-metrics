package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "failed to read gzipped request body", http.StatusBadRequest)
				return
			}
			defer gzipReader.Close()

			r.Body = io.NopCloser(gzipReader)
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")

			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			w = &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gzipWriter,
			}
		}

		next.ServeHTTP(w, r)
	})
}

// gzipResponseWriter wraps http.ResponseWriter to write the gzipped response
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write compresses the response body and writes it to the gzip.Writer
func (rw *gzipResponseWriter) Write(p []byte) (int, error) {
	return rw.Writer.Write(p)
}
