package domain

type Metrics struct {
	MetricID
	MetricValue
}

func NewMetrics(mtype string, name string, value string) (*Metrics, error) {
	metricID, err := NewMetricID(name, mtype)
	if err != nil {
		return nil, err
	}
	metricValue, err := NewMetricValue(mtype, value)
	if err != nil {
		return nil, err
	}
	return &Metrics{
		MetricID:    *metricID,
		MetricValue: *metricValue,
	}, nil
}
