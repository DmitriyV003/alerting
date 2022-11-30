package models

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/dmitriy/alerting/internal/hasher"
	log "github.com/sirupsen/logrus"
)

const CounterType = "counter"
const GaugeType = "gauge"

type MetricType string

func NewGauge(name string, value float64, hash string) Metric {
	return Metric{
		Name:       name,
		Type:       GaugeType,
		FloatValue: &value,
		Hash:       hash,
	}
}

func NewCounter(name string, value int64, hash string) Metric {
	return Metric{
		Name:     name,
		Type:     CounterType,
		IntValue: &value,
		Hash:     hash,
	}
}

func NewHealth(hashKey string) *Health {
	return &Health{
		Metrics: &sync.Map{},
		Hasher:  hasher.New(hashKey),
	}
}

type Health struct {
	Metrics *sync.Map
	Hasher  *hasher.Hasher
}

func (h *Health) Store(id string, metricType MetricType, value string) {
	var metric interface{}
	switch metricType {
	case GaugeType:
		{
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Error("Unable to parse value: ", err)
				return
			}
			hash := h.Hasher.Hash(fmt.Sprintf("%s:%s:%f", id, metricType, val))

			metric = NewGauge(id, val, hash)
		}
	case CounterType:
		{
			val, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Error("Unable to parse value: ", err)
				return
			}
			hash := h.Hasher.Hash(fmt.Sprintf("%s:%s:%d", id, metricType, val))

			metric = NewCounter(id, val, hash)
		}
	default:
		{
			log.Error("Unknown Metric Type")
			return
		}
	}

	h.Metrics.Store(id, metric)
}

type Metric struct {
	Name       string     `json:"id"`
	Type       MetricType `json:"type"`
	IntValue   *int64     `json:"delta,omitempty"`
	FloatValue *float64   `json:"value,omitempty"`
	Hash       string     `json:"hash,omitempty"`
}
