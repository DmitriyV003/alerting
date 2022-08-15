package models

import "sync"

const CounterType = "counter"
const GaugeType = "gauge"

type MetricType string

func NewGauge(name string, value float64) Metric {
	return Metric{
		Name:       name,
		Type:       GaugeType,
		FloatValue: &value,
	}
}

func NewCounter(name string, value int64) Metric {
	return Metric{
		Name:     name,
		Type:     CounterType,
		IntValue: &value,
	}
}

func NewHealth() *Health {
	return &Health{
		Metrics: &sync.Map{},
	}
}

type Health struct {
	Metrics *sync.Map
}

type Metric struct {
	Name       string     `json:"id"`
	Type       MetricType `json:"type"`
	IntValue   *int64     `json:"delta,omitempty"`
	FloatValue *float64   `json:"value,omitempty"`
}
