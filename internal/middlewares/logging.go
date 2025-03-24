package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggingMiddleware logs incoming HTTP requests and their corresponding responses
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseRecorder := &responseLoggingWriter{ResponseWriter: w}
			start := time.Now()
			logger.Info("Request received",
				zap.String("method", r.Method),
				zap.String("uri", r.URL.Path),
				zap.Duration("duration", time.Since(start)),
			)

			next.ServeHTTP(responseRecorder, r)

			logger.Info("Response sent",
				zap.Int("status", responseRecorder.StatusCode()),
				zap.Int("size", responseRecorder.Size()),
			)
		})
	}
}

// Custom response writer to capture status code and response size
type responseLoggingWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseLoggingWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseLoggingWriter) Write(p []byte) (n int, err error) {
	n, err = rw.ResponseWriter.Write(p)
	rw.size += n
	return n, err
}

func (rw *responseLoggingWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *responseLoggingWriter) Size() int {
	return rw.size
}
