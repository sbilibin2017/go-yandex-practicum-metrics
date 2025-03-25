package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderHTML(t *testing.T) {
	w := httptest.NewRecorder()
	htmlContent := "<html><body><h1>Hello, World!</h1></body></html>"

	err := RenderHTML(w, htmlContent)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=UTF-8", w.Header().Get("Content-Type"))
	assert.Equal(t, htmlContent, w.Body.String())
}

func TestRenderHTMLEmptyContent(t *testing.T) {
	w := httptest.NewRecorder()
	htmlContent := ""

	err := RenderHTML(w, htmlContent)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=UTF-8", w.Header().Get("Content-Type"))
	assert.Equal(t, htmlContent, w.Body.String())
}

// TestRenderHTMLWriteError tests the scenario where an error occurs while writing the HTML response.
func TestRenderHTMLWriteError(t *testing.T) {
	w := &mockWriter{}
	htmlContent := "<html><body><h1>Error Test</h1></body></html>"
	err := RenderHTML(w, htmlContent)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, w.statusCode)
	assert.Equal(t, "Failed to write HTML response", w.body)
}

type mockWriter struct {
	statusCode int
	body       string
}

func (m *mockWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	m.body = "Failed to write HTML response"
	m.statusCode = http.StatusInternalServerError
	return 0, assert.AnError
}

func (m *mockWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}
