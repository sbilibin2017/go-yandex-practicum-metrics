package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeNotFoundResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	err := fmt.Errorf("not found")
	MakeNotFoundResponse(rr, err)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "not found\n", rr.Body.String())
}

func TestMakeBadRequestResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	err := fmt.Errorf("bad request")
	MakeBadRequestResponse(rr, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "bad request\n", rr.Body.String())
}

func TestMakeInternalServerErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	err := fmt.Errorf("internal error")
	MakeInternalServerErrorResponse(rr, err)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "internal error\n", rr.Body.String())
}
