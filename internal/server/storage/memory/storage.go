package memory

import (
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type metricStorage struct {
	metrics *sync.Map
}

func New() *metricStorage {
	return &metricStorage{
		metrics: &sync.Map{},
	}
}

func (s *metricStorage) GetByNameAndType(name string, metricType string) (interface{}, error) {
	metric, ok := s.metrics.Load(name)

	if !ok {
		return nil, applicationerrors.ErrNotFound
	}
	castedMetric := metric.(model.Metric)

	if metricType == model.GaugeType {
		return castedMetric.FloatValue, nil
	} else if metricType == model.CounterType {
		return castedMetric.IntValue, nil
	}

	return nil, applicationerrors.ErrUnknownType
}

func (s *metricStorage) GetAll() *[]storage.MetricData {
	var metrics []storage.MetricData

	s.metrics.Range(func(key, value interface{}) bool {
		metric := value.(model.Metric)
		var val interface{}

		if metric.Type == model.GaugeType {
			val = *metric.FloatValue
		} else if metric.Type == model.CounterType {
			val = *metric.IntValue
		}

		metricData := storage.MetricData{
			Name:  key.(string),
			Value: val,
		}
		metrics = append(metrics, metricData)

		return true
	})

	return &metrics
}

func (s *metricStorage) UpdateMetric(name string, value string, metricType string) error {
	var metric interface{}

	if metricType == model.GaugeType {
		val, err := strconv.ParseFloat(value, 64)

		if err != nil {
			return applicationerrors.ErrInvalidValue
		}

		metric = model.NewGauge(name, val)
	} else if metricType == model.CounterType {
		var ok bool
		metric, ok = s.metrics.Load(name)
		val, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			log.Info("Invalid metric Value")

			return applicationerrors.ErrInvalidValue
		}

		if !ok {
			metric = model.NewCounter(name, val)
		} else {
			*metric.(model.Metric).IntValue += val
		}
	} else {
		log.Info("Invalid metric Type")

		return applicationerrors.ErrInvalidType
	}

	s.metrics.Store(name, metric)

	return nil
}
