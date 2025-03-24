package middlewares

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipDecompressionOnRequest(t *testing.T) {
	// Create a gzipped request body
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write([]byte("Hello, World!"))
	assert.NoError(t, err)
	err = gzipWriter.Close()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/test", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	// Create a basic handler to test the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Wrap the handler with the GzipMiddleware
	gzipMiddleware := GzipMiddleware(handler)

	// Record the response
	rr := httptest.NewRecorder()
	gzipMiddleware.ServeHTTP(rr, req)

	// Check that the response is correct
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Hello, World!", rr.Body.String())
}

func TestGzipCompressionOnResponse(t *testing.T) {
	// Create a request that accepts gzip encoding
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	// Create a basic handler to test the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Wrap the handler with the GzipMiddleware
	gzipMiddleware := GzipMiddleware(handler)

	// Record the response
	rr := httptest.NewRecorder()
	gzipMiddleware.ServeHTTP(rr, req)

	// Check that the response status is OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that the response body is gzipped
	contentEncoding := rr.Header().Get("Content-Encoding")
	assert.Equal(t, "gzip", contentEncoding)

	// Try to decompress the response body to ensure it's valid gzip
	gzipReader, err := gzip.NewReader(rr.Body)
	assert.NoError(t, err)
	defer gzipReader.Close()

	// Read the decompressed body
	var decompressedBuf bytes.Buffer
	_, err = decompressedBuf.ReadFrom(gzipReader)
	assert.NoError(t, err)

	// Check the decompressed content
	assert.Equal(t, "Hello, World!", decompressedBuf.String())
}

func TestNoGzipCompressionIfNotAccepted(t *testing.T) {
	// Create a request that doesn't accept gzip encoding
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "identity") // no gzip

	// Create a basic handler to test the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Wrap the handler with the GzipMiddleware
	gzipMiddleware := GzipMiddleware(handler)

	// Record the response
	rr := httptest.NewRecorder()
	gzipMiddleware.ServeHTTP(rr, req)

	// Check that the response is not gzipped
	contentEncoding := rr.Header().Get("Content-Encoding")
	assert.Empty(t, contentEncoding) // Should not have the Content-Encoding header
	assert.Equal(t, "Hello, World!", rr.Body.String())
}

// TestInvalidGzippedRequestBody tests the case where gzip.NewReader fails due to invalid gzipped data
func TestInvalidGzippedRequestBody(t *testing.T) {
	// Create a handler that we will pass through the GzipMiddleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler won't be called if there is an error in decompression
		w.Write([]byte("Success"))
	})

	// Wrap the handler with the GzipMiddleware
	gzipMiddleware := GzipMiddleware(handler)

	// Create a request with invalid gzip data in the body (e.g., random bytes that are not valid gzip)
	invalidGzipData := []byte("This is not a valid gzip data")
	var buf bytes.Buffer
	// Simulating invalid gzip encoding by writing random data
	_, err := buf.Write(invalidGzipData)
	assert.NoError(t, err)

	// Create a request with the invalid gzipped body
	req := httptest.NewRequest(http.MethodPost, "/test", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	// Record the response
	rr := httptest.NewRecorder()
	gzipMiddleware.ServeHTTP(rr, req)

	// Verify that the status code is 400 Bad Request (since the gzip decompression failed)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Verify the error message in the response body
	assert.Equal(t, "failed to read gzipped request body\n", rr.Body.String())
}

func TestWriteWithoutGzipCompression(t *testing.T) {
	// Create a request that does not accept gzip encoding
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "identity") // No gzip

	// Create a basic handler to test the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The middleware should pass through the response without compression
		w.Write([]byte("Hello, World!"))
	})

	// Wrap the handler with the GzipMiddleware
	gzipMiddleware := GzipMiddleware(handler)

	// Record the response
	rr := httptest.NewRecorder()
	gzipMiddleware.ServeHTTP(rr, req)

	// Check that the response status is OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that the response body is not gzipped
	contentEncoding := rr.Header().Get("Content-Encoding")
	assert.Empty(t, contentEncoding) // Should not have the Content-Encoding header

	// Check the response body is as expected
	assert.Equal(t, "Hello, World!", rr.Body.String())
}
