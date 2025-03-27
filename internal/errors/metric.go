package errors

import "errors"

var (
	InvalidMetricTypeError  = errors.New("invalid metric type")
	MissingMetricNameError  = errors.New("missing metric name")
	InvalidMetricValueError = errors.New("invalid metric value")
)
