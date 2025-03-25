package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDecodeRequestBodyValid tests the DecodeRequestBody function with a valid request body.
func TestDecodeRequestBodyValid(t *testing.T) {
	// Mock the request and response
	reqBody := `{"name":"John", "age":30}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	var requestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	err := DecodeRequestBody(w, req, &requestData)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the requestData has been correctly populated
	assert.Equal(t, "John", requestData.Name)
	assert.Equal(t, 30, requestData.Age)
}

// TestDecodeRequestBodyInvalid tests the DecodeRequestBody function with an invalid request body.
func TestDecodeRequestBodyInvalid(t *testing.T) {
	// Mock the request and response
	reqBody := `{"name":"John", "age":"invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	var requestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	err := DecodeRequestBody(w, req, &requestData)

	// Assert that an error occurred
	assert.Error(t, err)

	// Assert the correct status code is returned for invalid body
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestEncodeResponseBodyValid tests the EncodeResponseBody function with a valid response body.
func TestEncodeResponseBodyValid(t *testing.T) {
	// Mock the response and request
	respData := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}

	w := httptest.NewRecorder()

	err := EncodeResponseBody(w, respData)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the response body is correctly encoded as JSON
	expectedResponse := `{"name":"John","age":30}`
	assert.JSONEq(t, expectedResponse, w.Body.String())

	// Assert that the correct content type header is set
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

// TestEncodeResponseBodyInvalid tests the EncodeResponseBody function with an invalid response body.
func TestEncodeResponseBodyInvalid(t *testing.T) {
	// Mock the response and request
	w := httptest.NewRecorder()

	// Create a channel that will cause an error when trying to encode
	invalidRespData := make(chan int)

	// Try to encode a channel, which will result in an error
	err := EncodeResponseBody(w, invalidRespData)

	// Assert that an error occurred
	assert.Error(t, err)

	// Assert the correct status code is returned for internal server error
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
