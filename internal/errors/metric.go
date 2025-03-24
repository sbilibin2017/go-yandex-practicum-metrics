package errors

import "errors"

var (
	ErrMetricInternal = errors.New("internal error")
	ErrMetricNotFound = errors.New("metric not found")
)
