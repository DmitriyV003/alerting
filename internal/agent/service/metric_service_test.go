package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/dmitriy/alerting/internal/agent/models"
	"github.com/dmitriy/alerting/internal/hasher"
	"github.com/stretchr/testify/assert"
)

func TestMetricService_gatherMetrics(t *testing.T) {
	type fields struct {
		Health *models.Health
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Gather metric test",
			fields: struct{ Health *models.Health }{Health: &models.Health{
				Metrics: &sync.Map{},
				Hasher:  hasher.New("fg"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricService := &MetricService{
				Health: *tt.fields.Health,
			}
			metricService.gatherMetrics()

			metric, isOk := metricService.Health.Metrics.Load("Alloc")
			metricCasted := metric.(models.Metric)

			assert.Equal(t, metricCasted.Name, "Alloc")
			assert.True(t, isOk)

			counterMetric, _ := metricService.Health.Metrics.Load("PollCount")
			counterMetricCasted := counterMetric.(models.Metric)

			assert.Equal(t, counterMetricCasted.Name, "PollCount")
		})
	}
}

func BenchmarkMetricService_GatherMetricsByInterval(b *testing.B) {
	h := models.Health{
		Metrics: &sync.Map{},
		Hasher:  hasher.New("fg"),
	}
	metricService := &MetricService{
		Health: h,
	}

	for i := 0; i < b.N; i++ {
		metricService.gatherMetrics()
	}
}

func ExampleMetricService_GatherMetricsByInterval() {
	metricService := NewMetricService("somesecretkey")
	d, err := time.ParseDuration("5s")
	if err != nil {
		return
	}
	go metricService.GatherMetricsByInterval(context.Background(), d)
}

func ExampleMetricService_GatherAdditionalMetricsByInterval() {
	metricService := NewMetricService("somesecretkey")
	d, err := time.ParseDuration("5s")
	if err != nil {
		return
	}
	go metricService.GatherAdditionalMetricsByInterval(context.Background(), d)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *MetricService
	}{
		{
			name: "Create Metric Service Test",
			want: &MetricService{Health: struct {
				Metrics *sync.Map
				Hasher  *hasher.Hasher
			}{
				Metrics: &sync.Map{},
				Hasher:  hasher.New("fg"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewMetricService("fg"), "NewMetricService()")
		})
	}
}
