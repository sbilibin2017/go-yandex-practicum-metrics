package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzipMiddleware_ResponseCompression(t *testing.T) {
	data := "test response"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(data))
	})
	middleware := GzipMiddleware()(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	gzipReader, err := gzip.NewReader(w.Body)
	require.NoError(t, err)
	decompressedData, err := io.ReadAll(gzipReader)
	require.NoError(t, err)
	assert.Equal(t, data, string(decompressedData))
}

func TestGzipMiddleware_WriteHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	middleware := GzipMiddleware()(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGzipMiddleware_InvalidGzipRequest(t *testing.T) {
	invalidData := []byte("invalid gzip data")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := GzipMiddleware()(handler)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(invalidData))
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, ErrFailedToReadGzipRequest.Error()+"\n", w.Body.String())
}

func TestGzipMiddleware_RequestDecompression(t *testing.T) {
	data := []byte("test payload")
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	_, err := gzipWriter.Write(data)
	require.NoError(t, err)
	require.NoError(t, gzipWriter.Close())
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, data, body)
		w.WriteHeader(http.StatusOK)
	})
	middleware := GzipMiddleware()(handler)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(compressedData.Bytes()))
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
