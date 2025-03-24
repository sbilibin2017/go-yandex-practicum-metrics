package errors

import "errors"

var (
	ErrMetricInternal      = errors.New("internal error")
	ErrMetricNotFound      = errors.New("metric not found")
	ErrMetricTypeRequired  = errors.New("metric type required")
	ErrInvalidMetricType   = errors.New("metric counter or gauge required")
	ErrMetricNameRequired  = errors.New("metric name required")
	ErrMetricValueRequired = errors.New("metric value required")
	ErrInvalidCounterValue = errors.New("metric invalid counter value")
	ErrInvalidGaugeValue   = errors.New("metric invalid gauge value")
)
