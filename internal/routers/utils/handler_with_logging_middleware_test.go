package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the logger
type MockLoggingMiddlewareLogger struct {
	mock.Mock
}

func (m *MockLoggingMiddlewareLogger) Infow(msg string, args ...any) {
	m.Called(msg, args)
}

func TestHandlerWithLoggingMiddleware(t *testing.T) {
	logger := new(MockLoggingMiddlewareLogger)
	logger.On("Infow", mock.Anything, mock.Anything).Return()
	handler := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Test Handler"))
	}
	handlerWithLogging := HandlerWithLoggingMiddleware(handler, logger)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	handlerWithLogging.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Test Handler", rr.Body.String())
	logger.AssertCalled(t, "Infow", "Request received", mock.Anything)
	logger.AssertCalled(t, "Infow", "Response sent", mock.Anything)
}
