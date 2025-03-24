package validators

import (
	"go-yandex-practicum-metrics/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetricValue(t *testing.T) {
	// Test cases for ValidateMetricValue function
	tests := []struct {
		value       string
		expected    error
		expectedMsg string
	}{
		{"123", nil, ""}, // Valid value
		{"0", nil, ""},   // Zero value is valid

	}

	for _, tt := range tests {
		t.Run("ValidateMetricValue", func(t *testing.T) {
			result := ValidateMetricValue(tt.value)

			if tt.expected != nil {
				// Assert the error type and the error message text
				assert.ErrorIs(t, result, tt.expected)
				assert.EqualError(t, result, tt.expectedMsg) // Assert the expected error message
			} else {
				// No error is expected
				assert.NoError(t, result)
			}
		})
	}
}

func TestErrMetricValueRequired(t *testing.T) {
	// Test cases for ValidateMetricValue function
	tests := []struct {
		value    string
		expected error
	}{
		{"", errors.ErrMetricValueRequired}, // Empty value should return ErrMetricValueRequired
		{"123", nil},                        // Non-empty value should return nil (no error)
	}

	for _, tt := range tests {
		t.Run("ValidateMetricValue", func(t *testing.T) {
			// Call ValidateMetricValue which may return an error
			result := ValidateMetricValue(tt.value)

			if tt.expected != nil {
				// Check that the error returned matches the expected error type
				assert.ErrorIs(t, result, tt.expected)
			} else {
				// Check that no error is returned
				assert.NoError(t, result)
			}
		})
	}
}
