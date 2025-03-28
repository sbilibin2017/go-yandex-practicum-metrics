package validators

import (
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
)

func ValidateMetricType(metricType string) error {
	if !(metricType == string(domain.Gauge) || metricType == string(domain.Counter)) {
		return errors.ErrInvalidMetricType
	}
	return nil
}
