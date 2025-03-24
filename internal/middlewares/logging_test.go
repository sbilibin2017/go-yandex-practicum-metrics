package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggingMiddleware(t *testing.T) {
	// Initialize a buffer to capture logs
	var buf bytes.Buffer
	writeSyncer := zapcore.AddSync(&buf)
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		writeSyncer,
		zap.NewAtomicLevelAt(zap.InfoLevel),
	))

	// Create a test HTTP handler to pass through the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Wrap the handler with the LoggingMiddleware using the custom logger
	loggingMiddleware := LoggingMiddleware(logger)(handler)

	// Create a test HTTP request
	req, err := http.NewRequest(http.MethodGet, "/test-uri", nil)
	require.NoError(t, err)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the middleware with the request and response recorder
	loggingMiddleware.ServeHTTP(rr, req)

	// Assert that the response status code is 200
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the response body is "Hello, World!"
	assert.Equal(t, "Hello, World!", rr.Body.String())

	// Assert that the log buffer contains the expected logs
	logs := buf.String()

	// Check that the logs contain the expected entries
	assert.Contains(t, logs, "Request received")
	assert.Contains(t, logs, "/test-uri") // Check if the correct URI is logged
	assert.Contains(t, logs, "GET")       // Check if the correct HTTP method is logged
	assert.Contains(t, logs, "duration")  // Ensure duration is logged

	assert.Contains(t, logs, "Response sent")
	assert.Contains(t, logs, "200")  // Check if the status code is logged
	assert.Contains(t, logs, "size") // Ensure the size of the response is logged
}
