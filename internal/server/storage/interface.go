package storage

type MetricData struct {
	Name  string
	Value interface{}
}

type MetricStorage interface {
	UpdateMetric(metric string, value string, metricType string) error
	GetAll() *[]MetricData
	GetByNameAndType(name string, metricType string) (interface{}, error)
}
