package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger имитирует интерфейс Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Infow(msg string, args ...any) {
	m.Called(msg, args)
}

func TestRegisterMetricUpdatePathRouter(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infow", "Request received", mock.Anything).Once()
	mockLogger.On("Infow", "Response sent", mock.Anything).Once()
	handler := httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Type: " + ps.ByName("type") + ", Name: " + ps.ByName("name") + ", Value: " + ps.ByName("value")))
	})
	router := httprouter.New()
	RegisterMetricUpdatePathRouter(router, handler, mockLogger)
	req, err := http.NewRequest("POST", "/update/some_type/some_name/some_value", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockLogger.AssertExpectations(t)
}
