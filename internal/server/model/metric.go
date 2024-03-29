package model

const CounterType = "counter"
const GaugeType = "gauge"

type MetricType string

func newMetric(name string) *Metric {
	return &Metric{
		Name: name,
	}
}

func NewGauge(name string, val float64) *Metric {
	metric := newMetric(name)
	metric.Type = GaugeType
	metric.FloatValue = &val

	return metric
}

func NewCounter(name string, val int64) *Metric {
	metric := newMetric(name)
	metric.Type = CounterType
	metric.IntValue = &val

	return metric
}

type Metric struct {
	ID         *int64     `json:"-"`
	Name       string     `json:"id"`
	Type       MetricType `json:"type"`
	IntValue   *int64     `json:"delta,omitempty"`
	FloatValue *float64   `json:"value,omitempty"`
	Hash       string     `json:"hash,omitempty"`
}
