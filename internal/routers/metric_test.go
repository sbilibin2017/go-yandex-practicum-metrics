package routers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Infow(msg string, args ...any) {
	m.Called(msg, args)
}

func TestRegisterUpdateMetricPathRoute(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infow", mock.Anything, mock.Anything).Return()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Update Metric Path Route"))
	}
	r := chi.NewRouter()
	RegisterUpdateMetricPathRoute(r, handler, mockLogger)
	req, err := http.NewRequest("POST", "/update/someType/123/456", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Update Metric Path Route", rr.Body.String())
	mockLogger.AssertExpectations(t)
}

func TestRegisterUpdatesMetricBodyRoute(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infow", mock.Anything, mock.Anything).Return()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Metrics updates received"))
	})
	r := chi.NewRouter()
	RegisterUpdatesMetricBodyRoute(r, handler, mockLogger)
	reqBody := []byte(`{"type":"someType","id":"123","value":"456"}`)
	req, err := http.NewRequest("POST", "/updates/", bytes.NewReader(reqBody))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Metrics updates received", rr.Body.String())
	mockLogger.AssertExpectations(t)
}

func TestRegisterGetMetricByTypeAndIDBodyRoute(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infow", mock.Anything, mock.Anything).Return()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Metric received"))
	})
	r := chi.NewRouter()
	RegisterGetMetricByByTypeAndIDBodyRoute(r, handler, mockLogger)
	reqBody := []byte(`{"type":"someType","id":"123","value":"456"}`)
	req, err := http.NewRequest("POST", "/value/", bytes.NewReader(reqBody))
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Metric received", rr.Body.String())
	mockLogger.AssertExpectations(t)
}

func TestRegisterUpdateMetricBodyRoute(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infow", mock.Anything, mock.Anything).Return()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Update Metric Body Route"))
	}
	r := chi.NewRouter()
	RegisterUpdateMetricBodyRoute(r, handler, mockLogger)
	req, err := http.NewRequest("POST", "/update/", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Update Metric Body Route", rr.Body.String())
	mockLogger.AssertExpectations(t)
}

func TestRegisterGetMetricByTypeAndIDPathRoute(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infow", mock.Anything, mock.Anything).Return()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get Metric By Type and ID Path Route"))
	}
	r := chi.NewRouter()
	RegisterGetMetricByTypeAndIDPathRoute(r, handler, mockLogger)
	req, err := http.NewRequest("GET", "/value/someType/123", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Get Metric By Type and ID Path Route", rr.Body.String())
	mockLogger.AssertExpectations(t)
}

func TestRegisterListMetricsHTMLRoute(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infow", mock.Anything, mock.Anything).Return()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("List Metrics HTML Route"))
	}
	r := chi.NewRouter()
	RegisterListMetricsHTMLRoute(r, handler, mockLogger)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "List Metrics HTML Route", rr.Body.String())
	mockLogger.AssertExpectations(t)
}

func TestMetricMiddlewares(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infow", mock.Anything, mock.Anything).Return()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Middleware Test"))
	})
	r := chi.NewRouter()
	useMetricMiddlewares(r, mockLogger)
	r.Get("/", handler)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Middleware Test", rr.Body.String())
	mockLogger.AssertExpectations(t)
}

func TestLogger(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Infow", mock.Anything, mock.Anything).Return()
	mockLogger.Infow("Some message", "key", "value")
	mockLogger.AssertExpectations(t)
}
