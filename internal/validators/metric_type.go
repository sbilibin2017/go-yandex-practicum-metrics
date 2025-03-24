package validators

import (
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/types"
)

// ValidateType checks if the metric type is valid.
func ValidateMetricType(metricType string) error {
	if metricType == "" {
		return errors.ErrMetricTypeRequired
	}
	if metricType != string(types.Counter) && metricType != string(types.Gauge) {
		return errors.ErrInvalidMetricType
	}
	return nil
}
