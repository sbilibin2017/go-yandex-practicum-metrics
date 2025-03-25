package errors

import "errors"

var (
	ErrMetricBodyNotProvider = errors.New("metric body not provided")
	ErrMetricInternal        = errors.New("internal error")
	ErrMetricNotFound        = errors.New("metric not found")
	ErrMetricIDRequired      = errors.New("metric ID required")
	ErrMetricTypeRequired    = errors.New("metric type required")
	ErrMetricValueRequired   = errors.New("metric value required")
	ErrMetricInvalidType     = errors.New("invalid metric type")
	ErrMetricInvalidDelta    = errors.New("invalid metric delta")
	ErrMetricInvalidValue    = errors.New("invalid metric value")
)
