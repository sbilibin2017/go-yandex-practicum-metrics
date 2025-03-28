package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestWrapHandler(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Handler called"))
	})
	wrappedHandler := WrapHandler(handler)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/dummy-path", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := httprouter.Params{}
	wrappedHandler(rr, req, params)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Handler called", rr.Body.String())
}
