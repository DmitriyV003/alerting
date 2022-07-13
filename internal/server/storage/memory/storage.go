package memory

import (
	"errors"
	"fmt"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
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

func (s *Storage) GetByNameAndType(name string, metricType string) (interface{}, error) {
	if metricType == "gauge" {
		metric, ok := s.Gauges.Load(name)

		if !ok {
			return nil, errors.New("not found")
		}

		gauge := metric.(model.Gauge)

		return gauge.Value, nil
	} else if metricType == "counter" {
		metric, ok := s.Counters.Load(name)

		if !ok {
			return nil, errors.New("not found")
		}

		counter := metric.(model.Counter)

		return counter.Value, nil
	}

	return nil, errors.New("unknown type")
}

func (s *Storage) GetAll() *[]storage.MetricData {
	var metrics []storage.MetricData

	s.Gauges.Range(func(key, value interface{}) bool {
		metric := value.(model.Gauge)
		metricData := storage.MetricData{
			Name:  key.(string),
			Value: fmt.Sprint(metric.Value),
		}
		metrics = append(metrics, metricData)

		return true
	})

	s.Counters.Range(func(key, value interface{}) bool {
		metric := value.(model.Counter)
		metricData := storage.MetricData{
			Name:  key.(string),
			Value: fmt.Sprint(metric.Value),
		}
		metrics = append(metrics, metricData)

		return true
	})

	return &metrics

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
