package storage

import (
	"context"
	"github.com/dmitriy/alerting/internal/server/model"
)

type MetricStorage interface {
	UpdateOrCreate(ctx context.Context, metric string, value string, metricType string) error
	GetAll(ctx context.Context) *[]model.Metric
	GetByNameAndType(ctx context.Context, name string, metricType string) (*model.Metric, error)
	SaveAllMetricsData(ctx context.Context, metrics *[]model.Metric)
	AddOnUpdateListener(fn func())
	Emit(event string)
}
