package middlewares

import (
	"io"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLoggingMiddlewareLogger struct {
	mock.Mock
}

func (m *MockLoggingMiddlewareLogger) Infow(msg string, args ...any) {
	m.Called(msg, args)
}

func TestLoggingMiddleware(t *testing.T) {
	mockLogger := new(MockLoggingMiddlewareLogger)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		_, _ = io.WriteString(w, "Hello, world!")
	})
	middleware := LoggingMiddleware(mockLogger)
	wrappedHandler := middleware(handler)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	mockLogger.On("Infow", "Request received", mock.Anything).Once()
	mockLogger.On("Infow", "Response sent", mock.Anything).Once()
	wrappedHandler.ServeHTTP(w, req)
	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
	assert.Equal(t, "Hello, world!", string(body))
	mockLogger.AssertExpectations(t)
}
