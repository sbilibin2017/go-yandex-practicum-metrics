package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToFloat64(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
		expected  *float64
	}{
		{
			name:      "valid float64 value",
			value:     "123.45",
			expectErr: false,
			expected:  floatPointer(123.45),
		},
		{
			name:      "valid integer as string",
			value:     "100",
			expectErr: false,
			expected:  floatPointer(100),
		},
		{
			name:      "invalid value (non-numeric string)",
			value:     "invalid",
			expectErr: true,
			expected:  nil,
		},
		{
			name:      "empty string",
			value:     "",
			expectErr: true,
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertToFloat64(tt.value)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func floatPointer(v float64) *float64 {
	return &v
}
