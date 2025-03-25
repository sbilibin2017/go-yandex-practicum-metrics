package middlewares

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"strings"
)

var (
	ErrFailedToReadGzipRequest = errors.New("failed to read gzipped request body")
)

type GzipMiddlewareLogger interface {
	Infow(msg string, args ...any)
}

func GzipMiddleware(logger GzipMiddlewareLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infow("Request received", "method", r.Method, "uri", r.URL.Path)
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				gzipReader, err := gzip.NewReader(r.Body)
				if err != nil {
					logger.Infow("Failed to read gzipped request body", "error", err)
					http.Error(w, ErrFailedToReadGzipRequest.Error(), http.StatusBadRequest)
					return
				}
				defer gzipReader.Close()
				r.Body = io.NopCloser(gzipReader)
				logger.Infow("Request body decompressed")
			}
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				w.Header().Set("Content-Encoding", "gzip")
				gzipWriter := gzip.NewWriter(w)
				defer gzipWriter.Close()

				gzipResponseWriter := &gzipResponseWriter{
					ResponseWriter: w,
					Writer:         gzipWriter,
					logger:         logger,
				}
				logger.Infow("Response compression enabled")
				next.ServeHTTP(gzipResponseWriter, r)
				logger.Infow("Response sent with compression")
				return
			}
			logger.Infow("Response sent without compression")
			next.ServeHTTP(w, r)
		})
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
	logger GzipMiddlewareLogger
}

func (rw *gzipResponseWriter) Write(p []byte) (int, error) {
	rw.logger.Infow("Writing response body", "size", len(p))
	return rw.Writer.Write(p)
}

func (rw *gzipResponseWriter) WriteHeader(statusCode int) {
	rw.logger.Infow("Setting response status", "status", statusCode)
	rw.ResponseWriter.WriteHeader(statusCode)
}
