package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	msg := "Bad Request"
	code := http.StatusBadRequest

	ErrorResponse(w, msg, code)

	assert.Equal(t, code, w.Code)
	assert.Equal(t, msg+"\n", w.Body.String())
}

func TestErrorResponseNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	msg := "Not Found"
	code := http.StatusNotFound

	ErrorResponse(w, msg, code)

	assert.Equal(t, code, w.Code)
	assert.Equal(t, msg+"\n", w.Body.String())
}

func TestErrorResponseInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	msg := "Internal Server Error"
	code := http.StatusInternalServerError

	ErrorResponse(w, msg, code)

	assert.Equal(t, code, w.Code)
	assert.Equal(t, msg+"\n", w.Body.String())
}
