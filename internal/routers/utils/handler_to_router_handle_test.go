package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestHandlerToRouterHandle(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test Handler"))
	})
	routerHandle := HandlerToRouterHandle(handler)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	routerHandle(rr, req, httprouter.Params{})
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Test Handler", rr.Body.String())
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "/test", req.URL.Path)
}
