package memory

import (
	"fmt"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type metricStorage struct {
	counters sync.Map
	gauges   sync.Map
}

func New() *metricStorage {
	return &metricStorage{
		counters: sync.Map{},
		gauges:   sync.Map{},
	}
}

func (s *metricStorage) GetByNameAndType(name string, metricType string) (interface{}, error) {
	if metricType == model.GaugeType {
		metric, ok := s.gauges.Load(name)

		if !ok {
			return nil, applicationerrors.ErrNotFound
		}

		gauge := metric.(model.Gauge)

		return gauge.Value, nil
	} else if metricType == model.CounterType {
		metric, ok := s.counters.Load(name)

		if !ok {
			return nil, applicationerrors.ErrNotFound
		}

		counter := metric.(model.Counter)

		return counter.Value, nil
	}

	return nil, applicationerrors.ErrUnknownType
}

func (s *metricStorage) GetAll() *[]storage.MetricData {
	var metrics []storage.MetricData

	s.gauges.Range(func(key, value interface{}) bool {
		metric := value.(model.Gauge)
		metricData := storage.MetricData{
			Name:  key.(string),
			Value: fmt.Sprint(metric.Value),
		}
		metrics = append(metrics, metricData)

		return true
	})

	s.counters.Range(func(key, value interface{}) bool {
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

func (s *metricStorage) UpdateMetric(metric string, value string, metricType string) error {
	if metricType == model.GaugeType {
		foundedMetric := model.NewGauge(metric)
		val, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return applicationerrors.ErrInvalidValue
		}

		foundedMetric.Value = val
		foundedMetric.Name = metric
		s.gauges.Store(metric, foundedMetric)
	} else if metricType == model.CounterType {
		foundedMetric, ok := s.counters.Load(metric)

		val, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			log.Info("Invalid metric Value")

			return applicationerrors.ErrInvalidValue
		}

		if !ok {
			newCounter := model.NewCounter(metric)
			newCounter.Value = val
			newCounter.Name = metric
			s.counters.Store(metric, newCounter)
		} else {
			foundedMetric := foundedMetric.(model.Counter)
			foundedMetric.Value += val
			s.counters.Store(metric, foundedMetric)
		}
	} else {
		log.Info("Invalid metric Type")

		return applicationerrors.ErrInvalidType
	}

	return nil
}
