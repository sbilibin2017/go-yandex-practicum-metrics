package validators

import (
	"go-yandex-practicum-metrics/internal/errors"
)

// ValidateName checks if the metric name is provided.
func ValidateMetricName(name string) error {
	if name == "" {
		return errors.ErrMetricNameRequired
	}
	return nil
}
