package utils

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestGetPathParam(t *testing.T) {
	params := httprouter.Params{
		{Key: "type", Value: "some_type"},
		{Key: "name", Value: "some_name"},
		{Key: "value", Value: "some_value"},
	}

	tests := []struct {
		paramName     string
		expectedValue string
	}{
		{"type", "some_type"},
		{"name", "some_name"},
		{"value", "some_value"},
		{"nonexistent", ""},
	}
	for _, tt := range tests {
		t.Run(tt.paramName, func(t *testing.T) {
			result := GetPathParam(params, tt.paramName)
			assert.Equal(t, tt.expectedValue, result)
		})
	}
}
