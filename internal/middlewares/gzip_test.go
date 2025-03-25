package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockGzipMIddlewareLogger struct {
	mock.Mock
}

func (m *MockGzipMIddlewareLogger) Infow(msg string, args ...any) {
	m.Called(msg, args)
}

func TestGzipMiddleware_ResponseCompression(t *testing.T) {
	logger := new(MockGzipMIddlewareLogger)
	logger.On("Infow", mock.Anything, mock.Anything).Return()

	data := "test response"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(data))
	})

	middleware := GzipMiddleware(logger)(handler)

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

	logger.AssertCalled(t, "Infow", "Response compression enabled", mock.Anything)
	logger.AssertCalled(t, "Infow", "Response sent with compression", mock.Anything)
}

func TestGzipMiddleware_WriteHeader(t *testing.T) {
	logger := new(MockGzipMIddlewareLogger)
	logger.On("Infow", mock.Anything, mock.Anything).Return()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	middleware := GzipMiddleware(logger)(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	logger.AssertCalled(t, "Infow", "Setting response status", mock.Anything)
}

func TestGzipMiddleware_InvalidGzipRequest(t *testing.T) {
	logger := new(MockGzipMIddlewareLogger)
	logger.On("Infow", mock.Anything, mock.Anything).Return()

	invalidData := []byte("invalid gzip data")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := GzipMiddleware(logger)(handler)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(invalidData))
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, ErrFailedToReadGzipRequest.Error()+"\n", w.Body.String())

	logger.AssertCalled(t, "Infow", "Failed to read gzipped request body", mock.Anything)
}

func TestGzipMiddleware_RequestDecompression_Logging(t *testing.T) {
	data := []byte("test payload")
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)
	_, err := gzipWriter.Write(data)
	require.NoError(t, err)
	require.NoError(t, gzipWriter.Close())

	mockLogger := new(MockGzipMIddlewareLogger)

	mockLogger.On("Infow", "Request received", mock.Anything).Once()
	mockLogger.On("Infow", "Request body decompressed", mock.Anything).Once()
	mockLogger.On("Infow", "Response sent without compression", mock.Anything).Once()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, data, body)
		w.WriteHeader(http.StatusOK)
	})

	middleware := GzipMiddleware(mockLogger)(handler)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(compressedData.Bytes()))
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockLogger.AssertExpectations(t)
}
