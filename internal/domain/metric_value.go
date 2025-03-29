package domain

import (
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/logger"
	"strconv"
)

type MetricValue struct {
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func NewMetricValue(mtype string, value string) (*MetricValue, error) {
	logger.Info("Creating new MetricValue", "type", mtype, "value", value)
	var d *int64
	var v *float64
	switch MetricType(mtype) {
	case Gauge:
		_v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logger.Error("Invalid MetricValue for Gauge type", "value", value, "error", err)
			return nil, errors.ErrInvalidMetricValue
		}
		v = &_v
	case Counter:
		_v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logger.Error("Invalid MetricValue for Counter type", "value", value, "error", err)
			return nil, errors.ErrInvalidMetricValue
		}
		d = &_v
	default:
		logger.Error("Invalid MetricType", "type", mtype)
		return nil, errors.ErrInvalidMetricType
	}
	metricValue := &MetricValue{
		Delta: d,
		Value: v,
	}
	logger.Info("Successfully created MetricValue", "type", mtype)
	return metricValue, nil
}
