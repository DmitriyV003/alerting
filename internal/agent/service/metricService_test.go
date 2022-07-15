package service

import (
	"github.com/dmitriy/alerting/internal/agent/models"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMetricService_gatherMetrics(t *testing.T) {
	type fields struct {
		Health models.Health
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Gather metric test",
			fields: struct{ Health models.Health }{Health: models.Health{
				Gauges:   sync.Map{},
				Counters: sync.Map{},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricService := &MetricService{
				Health: tt.fields.Health,
			}
			metricService.gatherMetrics()

			metric, _ := metricService.Health.Gauges.Load("Alloc")
			metricCasted := metric.(models.Gauge)
			assert.Equal(t, metricCasted.Name, "Alloc")

			counterMetric, _ := metricService.Health.Counters.Load("PollCount")
			counterMetricCasted := counterMetric.(models.Counter)
			assert.Equal(t, counterMetricCasted.Name, "PollCount")
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *MetricService
	}{
		{
			name: "Create Metric Service Test",
			want: &MetricService{Health: struct {
				Gauges   sync.Map
				Counters sync.Map
			}{Gauges: sync.Map{}, Counters: sync.Map{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, New(), "New()")
		})
	}
}
