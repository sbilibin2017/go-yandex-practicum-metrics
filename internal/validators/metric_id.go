package validators

import "go-yandex-practicum-metrics/internal/errors"

func ValidateMetricID(metricID string) error {
	if metricID == "" {
		return errors.ErrMissingMetricID
	}
	return nil
}
