package memory

import (
	"errors"
	"github.com/dmitriy/alerting/internal/server/model"
	"strconv"
	"sync"
)

type Storage struct {
	Counters sync.Map
	Gauges   sync.Map
}

func New() *Storage {
	return &Storage{
		Counters: sync.Map{},
		Gauges:   sync.Map{},
	}
}

func NewGauge(name string) model.Gauge {
	return model.Gauge{
		Metric: model.Metric{
			Name: name,
		},
		Value: 0,
	}
}
func NewCounter(name string) model.Counter {
	return model.Counter{
		Metric: model.Metric{
			Name: name,
		},
		Value: 0,
	}
}

func (s *Storage) UpdateMetric(metric string, value string, metricType string) error {
	if metricType == "gauge" {
		foundedMetric := NewGauge(metric)
		val, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return errors.New("invalid value")
		}

		foundedMetric.Value = val
		foundedMetric.Name = metric
		s.Gauges.Store(metric, foundedMetric)
	} else if metricType == "counter" {
		foundedMetric, ok := s.Counters.Load(metric)

		val, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			return errors.New("invalid value")
		}

		if !ok {
			newCounter := NewCounter(metric)
			newCounter.Value = val
			newCounter.Name = metric
			s.Counters.Store(metric, newCounter)
		} else {
			foundedMetric := foundedMetric.(model.Counter)
			foundedMetric.Value += val
			s.Counters.Store(metric, foundedMetric)
		}
	} else {
		return errors.New("invalid type")
	}

	return nil
}
