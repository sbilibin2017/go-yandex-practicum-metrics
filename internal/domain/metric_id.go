package domain

import (
	"go-yandex-practicum-metrics/internal/errors"
	"go-yandex-practicum-metrics/internal/logger"
)

type MetricID struct {
	ID   string     `json:"id"`
	Type MetricType `json:"type"`
}

func NewMetricID(id string, mtype string) (*MetricID, error) {
	logger.Info("Creating new MetricID", "id", id, "type", mtype)
	if id == "" {
		logger.Error("Invalid MetricID: ID is empty", "id", id)
		return nil, errors.ErrInvalidMetricID
	}
	if mtype != string(Counter) && mtype != string(Gauge) {
		logger.Error("Invalid MetricID: Invalid type", "type", mtype)
		return nil, errors.ErrInvalidMetricType
	}
	metricID := &MetricID{
		ID:   id,
		Type: MetricType(mtype),
	}
	logger.Info("Successfully created MetricID", "id", id, "type", mtype)
	return metricID, nil
}
