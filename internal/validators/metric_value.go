package validators

import (
	"go-yandex-practicum-metrics/internal/errors"
)

// ValidateValue checks if the metric value is provided.
func ValidateMetricValue(value string) error {
	if value == "" {
		return errors.ErrMetricValueRequired
	}
	return nil
}
