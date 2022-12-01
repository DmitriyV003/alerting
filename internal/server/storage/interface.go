package storage

import (
	"context"

	"github.com/dmitriy/alerting/internal/server/model"
)

type MetricStorage interface {
	// UpdateOrCreate Update or create metric
	UpdateOrCreate(ctx context.Context, metric string, value string, metricType string) error

	// GetAll Get all metrics from store
	GetAll(ctx context.Context) *[]model.Metric

	// GetByNameAndType Get metric by name and type
	GetByNameAndType(ctx context.Context, name string, metricType string) (*model.Metric, error)

	// SaveCollection Save metric collection
	SaveCollection(ctx context.Context, metrics *[]model.Metric) error

	// RestoreCollection Resave metric collection
	RestoreCollection(ctx context.Context, metrics *[]model.Metric)
	AddOnUpdateListener(fn func())
	Emit(event string)
}
