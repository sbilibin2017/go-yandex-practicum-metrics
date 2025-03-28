package validators

import (
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetricType(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		expectErr  error
	}{
		{
			name:       "invalid metricType",
			metricType: "invalid_type",
			expectErr:  errors.ErrInvalidMetricType,
		},
		{
			name:       "valid metricType - Gauge",
			metricType: string(domain.Gauge),
			expectErr:  nil,
		},
		{
			name:       "valid metricType - Counter",
			metricType: string(domain.Counter),
			expectErr:  nil,
		},
		{
			name:       "empty metricType",
			metricType: "",
			expectErr:  errors.ErrInvalidMetricType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricType(tt.metricType)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}
