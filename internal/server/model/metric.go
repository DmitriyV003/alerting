package model

const CounterType = "counter"
const GaugeType = "gauge"

func NewGauge(name string) Gauge {
	return Gauge{
		Metric: Metric{
			Name: name,
		},
		Value: 0,
	}
}

func NewCounter(name string) Counter {
	return Counter{
		Metric: Metric{
			Name: name,
		},
		Value: 0,
	}
}

type Metric struct {
	Name string
}

type Gauge struct {
	Metric
	Value float64
}

type Counter struct {
	Metric
	Value int64
}
