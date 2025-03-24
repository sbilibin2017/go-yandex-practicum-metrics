package validators

import (
	"go-yandex-practicum-metrics/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetricName(t *testing.T) {
	// Test cases for ValidateName function
	tests := []struct {
		name     string
		expected error
	}{
		{"validName", nil},                 // Valid name
		{"", errors.ErrMetricNameRequired}, // Empty name, should return the required error
	}

	for _, tt := range tests {
		t.Run("ValidateName", func(t *testing.T) {
			result := ValidateMetricName(tt.name)
			assert.Equal(t, tt.expected, result)
		})
	}
}
