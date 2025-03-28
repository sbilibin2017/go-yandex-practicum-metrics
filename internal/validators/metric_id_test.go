package validators

import (
	"go-yandex-practicum-metrics/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMetricID(t *testing.T) {
	tests := []struct {
		name      string
		metricID  string
		expectErr error
	}{
		{
			name:      "empty metricID",
			metricID:  "",
			expectErr: errors.ErrMissingMetricID,
		},
		{
			name:      "non-empty metricID",
			metricID:  "valid-metric-id",
			expectErr: nil,
		},
		{
			name:      "another non-empty metricID",
			metricID:  "another-valid-id",
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetricID(tt.metricID)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}
