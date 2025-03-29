package middlewares

import (
	"compress/gzip"
	"errors"
	"go-yandex-practicum-metrics/internal/logger"
	"io"
	"net/http"
	"strings"
)

var (
	ErrFailedToReadGzipRequest = errors.New("failed to read gzipped request body")
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Error("Failed to read gzipped request body", "error", err)
				http.Error(w, ErrFailedToReadGzipRequest.Error(), http.StatusBadRequest)
				return
			}
			defer gzipReader.Close()
			r.Body = io.NopCloser(gzipReader)
			logger.Info("Gzip request body decompressed")
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()
			gzipResponseWriter := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gzipWriter,
			}
			logger.Info("Sending gzipped response")
			next.ServeHTTP(gzipResponseWriter, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (rw *gzipResponseWriter) Write(p []byte) (int, error) {
	return rw.Writer.Write(p)
}

func (rw *gzipResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
}
