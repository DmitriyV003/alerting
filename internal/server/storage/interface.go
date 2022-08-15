package storage

import "github.com/dmitriy/alerting/internal/server/model"

type MetricStorage interface {
	UpdateMetric(metric string, value string, metricType string) error
	GetAll() *[]model.Metric
	GetByNameAndType(name string, metricType string) (*model.Metric, error)
	SaveAllMetricsData(metrics *[]model.Metric)
	AddOnUpdateListener(fn func())
}
