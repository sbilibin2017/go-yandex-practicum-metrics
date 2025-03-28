package requests

import (
	"go-yandex-practicum-metrics/internal/converters"
	"go-yandex-practicum-metrics/internal/domain"
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/validators"
)

type MetricUpdatePathRequest struct {
	Type  string
	Name  string
	Value string
}

func (r *MetricUpdatePathRequest) Validate() error {
	if err := validators.ValidateMetricID(r.Name); err != nil {
		return errors.ErrMissingMetricID
	}
	if err := validators.ValidateMetricType(r.Type); err != nil {
		return errors.ErrInvalidMetricType
	}
	switch r.Type {
	case string(domain.Gauge):
		if err := validators.ValidateMetricGaugeStringValue(r.Value); err != nil {
			return errors.ErrInvalidMetricValue
		}
	case string(domain.Counter):
		if err := validators.ValidateMetricCounterStringValue(r.Value); err != nil {
			return errors.ErrInvalidMetricValue
		}
	}
	return nil
}

func (r *MetricUpdatePathRequest) ToDomain() *domain.Metrics {
	var delta *int64
	var value *float64
	switch r.Type {
	case string(domain.Gauge):
		value, _ = converters.ConvertToFloat64(r.Value)
	case string(domain.Counter):
		delta, _ = converters.ConvertToInt64(r.Value)
	}
	return &domain.Metrics{
		ID:    r.Name,
		Type:  domain.MetricType(r.Type),
		Delta: delta,
		Value: value,
	}
}
