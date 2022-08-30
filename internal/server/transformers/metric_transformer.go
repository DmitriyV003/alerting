package transformers

import (
	"fmt"
	"github.com/dmitriy/alerting/internal/hasher"
	"github.com/dmitriy/alerting/internal/server/model"
)

type MetricTransformer struct {
	hasher *hasher.Hasher
}

func NewTransformer(key string) *MetricTransformer {
	return &MetricTransformer{
		hasher: hasher.New(key),
	}
}

func (tr *MetricTransformer) AddHash(metric *model.Metric) *model.Metric {
	str := ""
	if metric.Type == "counter" {
		str = fmt.Sprintf("%s:%s:%d", metric.Name, metric.Type, *metric.IntValue)
	} else if metric.Type == "gauge" {
		str = fmt.Sprintf("%s:%s:%f", metric.Name, metric.Type, *metric.FloatValue)
	}
	hash := tr.hasher.Hash(str)
	metric.Hash = hash

	return metric
}
