package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestRegisterMetricUpdatePathRouter(t *testing.T) {
	handler := httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Type: " + ps.ByName("type") + ", Name: " + ps.ByName("name") + ", Value: " + ps.ByName("value")))
	})
	router := httprouter.New()
	RegisterMetricUpdatePathRouter(router, handler)
	req, err := http.NewRequest("POST", "/update/some_type/some_name/some_value", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	expectedBody := "Type: some_type, Name: some_name, Value: some_value"
	assert.Equal(t, expectedBody, rr.Body.String())
}
