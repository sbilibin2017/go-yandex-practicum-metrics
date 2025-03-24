package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTimeoutMiddlewareSuccess tests the case where the request completes within the timeout
func TestTimeoutMiddlewareSuccess(t *testing.T) {
	// Handler that sleeps for 1 second to simulate processing
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.Write([]byte("Success"))
	})

	// Wrap the handler with the TimeoutMiddleware
	timeoutMiddleware := TimeoutMiddleware(handler)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// Record the response
	rr := httptest.NewRecorder()
	timeoutMiddleware.ServeHTTP(rr, req)

	// Check that the status code is 200 OK because the request finished within the timeout
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	assert.Equal(t, "Success", rr.Body.String())
}

// TestTimeoutMiddlewareTimeout tests the case where the request takes too long and times out
func TestTimeoutMiddlewareTimeout(t *testing.T) {
	// Handler that sleeps for 10 seconds to simulate long processing
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second) // Simulate a long-running request
		w.Write([]byte("Success"))
	})

	// Wrap the handler with the TimeoutMiddleware
	timeoutMiddleware := TimeoutMiddleware(handler)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// Record the response
	rr := httptest.NewRecorder()
	timeoutMiddleware.ServeHTTP(rr, req)

	// Check that the status code is 408 Request Timeout because the request took too long
	assert.Equal(t, http.StatusRequestTimeout, rr.Code)

	// Check the response body
	assert.Equal(t, "Request Timeout\n", rr.Body.String())
}
