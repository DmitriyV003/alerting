package model

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
