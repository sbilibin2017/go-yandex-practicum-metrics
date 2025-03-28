package validators

import (
	"go-yandex-practicum-metrics/internal/errors"
	"strconv"
)

func ValidateMetricGaugeStringValue(value string) error {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return errors.ErrInvalidMetricValue
	}
	return nil

}

func ValidateMetricCounterStringValue(value string) error {
	_, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return errors.ErrInvalidMetricValue
	}
	return nil

}

func ValidateMetricGaugeValue(value *float64) error {
	if value == nil {
		return errors.ErrInvalidMetricValue
	}
	return nil
}

func ValidateMetricCounterValue(value *int64) error {
	if value == nil {
		return errors.ErrInvalidMetricValue
	}
	return nil

}
