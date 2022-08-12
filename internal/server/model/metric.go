package model

const CounterType = "counter"
const GaugeType = "gauge"

type MetricType string

func newMetric(name string) *Metric {
	return &Metric{
		Name: name,
	}
}

func NewGauge(name string, val float64) Metric {
	metric := newMetric(name)
	metric.Type = GaugeType
	metric.FloatValue = &val

	return *metric
}

func NewCounter(name string, val int64) Metric {
	metric := newMetric(name)
	metric.Type = CounterType
	metric.IntValue = &val

	return *metric
}

type Metric struct {
	Name       string
	Type       MetricType
	IntValue   *int64
	FloatValue *float64
}
