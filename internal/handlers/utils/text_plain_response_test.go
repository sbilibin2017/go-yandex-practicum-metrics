package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeTextPlainResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	data := []byte("Hello, world!")
	MakeTextPlainResponse(rr, data)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "text/plain", rr.Header().Get("Content-Type"))
	assert.Equal(t, string(data), rr.Body.String())
}
