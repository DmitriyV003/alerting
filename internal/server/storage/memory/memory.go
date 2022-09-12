package memory

import (
	"context"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage/commonstorage"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type metricStorage struct {
	metrics *sync.Map
	*commonstorage.CommonStorage
}

func New() *metricStorage {
	metricStore := metricStorage{
		metrics:       &sync.Map{},
		CommonStorage: commonstorage.New(),
	}
	metricStore.ListenEvents()

	return &metricStore
}

func (s *metricStorage) GetByNameAndType(ctx context.Context, name string, metricType string) (*model.Metric, error) {
	metric, ok := s.metrics.Load(name)

	if !ok {
		return nil, applicationerrors.ErrNotFound
	}
	castedMetric := metric.(*model.Metric)

	if metricType == model.GaugeType || metricType == model.CounterType {
		return castedMetric, nil
	}

	return nil, applicationerrors.ErrUnknownType
}

func (s *metricStorage) GetAll(ctx context.Context) *[]model.Metric {
	var metrics []model.Metric

	s.metrics.Range(func(key, value interface{}) bool {
		metric := value.(model.Metric)

		metrics = append(metrics, metric)

		return true
	})

	return &metrics
}

func (s *metricStorage) UpdateOrCreate(ctx context.Context, name string, value string, metricType string) error {
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
			*metric.(*model.Metric).IntValue += val
		}
	} else {
		log.Info("Invalid metric Type")

		return applicationerrors.ErrInvalidType
	}

	s.metrics.Store(name, metric)
	s.Emit("OnUpdate")

	return nil
}

func (s *metricStorage) SaveCollection(ctx context.Context, metrics *[]model.Metric) error {
	for _, metric := range *metrics {
		s.metrics.Store(metric.Name, metric)
	}

	return nil
}

func (s *metricStorage) RestoreCollection(ctx context.Context, metrics *[]model.Metric) {
	for _, metric := range *metrics {
		s.metrics.Store(metric.Name, metric)
	}
}
