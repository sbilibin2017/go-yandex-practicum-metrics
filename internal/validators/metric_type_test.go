package validators

import (
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetricType(t *testing.T) {
	// Test cases for ValidateMetricType function
	tests := []struct {
		metricType string
		expected   error
	}{
		{string(types.Counter), nil},                 // Valid metric type "Counter"
		{string(types.Gauge), nil},                   // Valid metric type "Gauge"
		{"", errors.ErrMetricTypeRequired},           // Empty metric type, should return the required error
		{"invalidType", errors.ErrInvalidMetricType}, // Invalid metric type, should return the invalid type error
		{"123", errors.ErrInvalidMetricType},         // Another invalid type, should return the invalid type error
		// Valid metric type "Gauge" (as a string)
	}

	for _, tt := range tests {
		t.Run("ValidateMetricType", func(t *testing.T) {
			result := ValidateMetricType(tt.metricType)

			// Compare error types using assert.ErrorIs
			if tt.expected != nil {
				assert.ErrorIs(t, result, tt.expected) // Check if the error is of the expected type
			} else {
				assert.NoError(t, result) // No error is expected
			}
		})
	}
}
