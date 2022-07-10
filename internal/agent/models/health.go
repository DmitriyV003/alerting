package models

import "sync"

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

func NewGauge(name string, value float64) Gauge {
	return Gauge{
		Metric: Metric{Name: name},
		Value:  value,
	}
}

func NewCounter(name string, value int64) Counter {
	return Counter{
		Metric: Metric{Name: name},
		Value:  value,
	}
}

func NewHealth() *Health {
	return &Health{
		Gauges:   sync.Map{},
		Counters: sync.Map{},
	}
}

type Health struct {
	Gauges   sync.Map
	Counters sync.Map
}
