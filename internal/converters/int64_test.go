package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToInt64(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
		expected  *int64
	}{
		{
			name:      "valid int64 value",
			value:     "12345",
			expectErr: false,
			expected:  int64Pointer(12345),
		},
		{
			name:      "valid negative int64 value",
			value:     "-98765",
			expectErr: false,
			expected:  int64Pointer(-98765),
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
			result, err := ConvertToInt64(tt.value)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func int64Pointer(v int64) *int64 {
	return &v
}
