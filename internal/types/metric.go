package types

type Metrics struct {
	MetricID
	Delta *int64
	Value *float64
}
