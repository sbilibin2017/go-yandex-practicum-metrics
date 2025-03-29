package middlewares

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"go-yandex-practicum-metrics/internal/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateUUID()
		start := time.Now()
		logger.Info("Request received",
			"request_id", requestID,
			"method", r.Method,
			"uri", r.URL.Path,
		)
		responseRecorder := &responseLoggingWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		w = responseRecorder
		next.ServeHTTP(w, r)
		logger.Info("Response sent",
			"request_id", requestID,
			"status", responseRecorder.StatusCode(),
			"size", responseRecorder.Size(),
			"duration", time.Since(start),
		)
	})
}

func generateUUID() string {
	var uuid [16]byte
	rand.Read(uuid[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

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
