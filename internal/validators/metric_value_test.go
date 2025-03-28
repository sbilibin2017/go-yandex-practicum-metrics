package validators

import (
	"go-yandex-practicum-metrics/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetricGaugeStringValue(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr error
	}{
		{
			name:      "valid gauge value",
			value:     "123.45",
			expectErr: nil,
		},
		{
			name:      "invalid gauge value (non-numeric)",
			value:     "invalid",
			expectErr: errors.ErrInvalidMetricValue,
		},
		{
			name:      "empty string",
			value:     "",
			expectErr: errors.ErrInvalidMetricValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricGaugeStringValue(tt.value)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestValidateMetricCounterStringValue(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr error
	}{
		{
			name:      "valid counter value",
			value:     "123",
			expectErr: nil,
		},
		{
			name:      "invalid counter value (non-numeric)",
			value:     "invalid",
			expectErr: errors.ErrInvalidMetricValue,
		},
		{
			name:      "empty string",
			value:     "",
			expectErr: errors.ErrInvalidMetricValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricCounterStringValue(tt.value)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestValidateMetricGaugeValue(t *testing.T) {
	tests := []struct {
		name      string
		value     *float64
		expectErr error
	}{
		{
			name:      "valid gauge value",
			value:     new(float64),
			expectErr: nil,
		},
		{
			name:      "nil gauge value",
			value:     nil,
			expectErr: errors.ErrInvalidMetricValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricGaugeValue(tt.value)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}

func TestValidateMetricCounterValue(t *testing.T) {
	tests := []struct {
		name      string
		value     *int64
		expectErr error
	}{
		{
			name:      "valid counter value",
			value:     new(int64),
			expectErr: nil,
		},
		{
			name:      "nil counter value",
			value:     nil,
			expectErr: errors.ErrInvalidMetricValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricCounterValue(tt.value)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}
