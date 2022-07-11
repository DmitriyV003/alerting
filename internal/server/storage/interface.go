package storage

type MetricStorage interface {
	UpdateMetric(metric string, value string, metricType string) error
}
