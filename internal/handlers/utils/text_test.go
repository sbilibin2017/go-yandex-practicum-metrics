package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSendTextResponse tests the SendTextResponse function with valid data.
func TestSendTextResponse(t *testing.T) {
	// Setup the recorder and request
	w := httptest.NewRecorder()
	data := "Hello, World!"

	// Call the function
	SendTextResponse(w, data)

	// Assert the status code is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the Content-Type header is set correctly
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))

	// Assert the response body is the expected text
	assert.Equal(t, data, w.Body.String())
}

// TestSendTextResponseEmpty tests the SendTextResponse function with an empty string as the data.
func TestSendTextResponseEmpty(t *testing.T) {
	// Setup the recorder and empty data
	w := httptest.NewRecorder()
	data := ""

	// Call the function
	SendTextResponse(w, data)

	// Assert the status code is OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the Content-Type header is set correctly
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))

	// Assert the response body is an empty string
	assert.Equal(t, data, w.Body.String())
}
