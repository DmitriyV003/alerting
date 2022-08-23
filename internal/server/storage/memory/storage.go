package memory

import (
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type metricStorage struct {
	metrics  *sync.Map
	events   map[string][]chan func()
	onUpdate []func()
}

func New() *metricStorage {
	metricStore := metricStorage{
		metrics: &sync.Map{},
		events: map[string][]chan func(){
			"OnUpdate": {
				make(chan func()),
			},
		},
		onUpdate: []func(){},
	}
	go func() {
		chs := metricStore.events["OnUpdate"]
		for _, ch := range chs {
			for {
				handler := <-ch
				handler()
				log.Info("Event: OnUpdate")
			}
		}
	}()

	return &metricStore
}

func (s *metricStorage) AddOnUpdateListener(fn func()) {
	s.onUpdate = append(s.onUpdate, fn)
}

func (s *metricStorage) GetByNameAndType(name string, metricType string) (*model.Metric, error) {
	metric, ok := s.metrics.Load(name)

	if !ok {
		return nil, applicationerrors.ErrNotFound
	}
	castedMetric := metric.(model.Metric)

	if metricType == model.GaugeType {
		return &castedMetric, nil
	} else if metricType == model.CounterType {
		return &castedMetric, nil
	}

	return nil, applicationerrors.ErrUnknownType
}

func (s *metricStorage) GetAll() *[]model.Metric {
	var metrics []model.Metric

	s.metrics.Range(func(key, value interface{}) bool {
		metric := value.(model.Metric)

		metrics = append(metrics, metric)

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
	s.emit("OnUpdate")

	return nil
}

func (s *metricStorage) SaveAllMetricsData(metrics *[]model.Metric) {
	for _, metric := range *metrics {
		s.metrics.Store(metric.Name, metric)
	}
}

func (s *metricStorage) emit(event string) {
	if _, ok := s.events[event]; ok {
		for _, handler := range s.events[event] {
			go func(handler chan func()) {
				for _, h := range s.onUpdate {
					handler <- h
				}
			}(handler)
		}
	}
}
