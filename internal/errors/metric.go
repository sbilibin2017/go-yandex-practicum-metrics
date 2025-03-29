package errors

import "errors"

var (
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricID    = errors.New("invalid metric id")
	ErrInvalidMetricValue = errors.New("invalid metric value")
	ErrMetricNotFound     = errors.New("metric not found")
	ErrInternal           = errors.New("internal error")
)
