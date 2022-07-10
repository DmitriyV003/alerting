package service

import (
	"github.com/dmitriy/alerting/internal/agent/models"
	"math/rand"
	"runtime"
	"time"
)

type MetricService struct {
	Health models.Health
}

func New() *MetricService {
	return &MetricService{Health: *models.NewHealth()}
}

func (metricService *MetricService) gatherMetrics() {
	var runtimeStats runtime.MemStats
	rand.Seed(time.Now().UnixNano())

	runtime.ReadMemStats(&runtimeStats)

	currentCounter, ok := metricService.Health.Counters.Load("PollCount")
	if !ok {
		metricService.Health.Counters.Store("PollCount", models.NewCounter("PollCount", 0))
	} else {
		metric := currentCounter.(models.Counter)
		metricService.Health.Counters.Store("PollCount", models.NewCounter("PollCount", metric.Value+1))
	}

	metricService.Health.Gauges.Store("Alloc", models.NewGauge("Alloc", float64(runtimeStats.Alloc)))
	metricService.Health.Gauges.Store("BuckHashSys", models.NewGauge("BuckHashSys", float64(runtimeStats.BuckHashSys)))
	metricService.Health.Gauges.Store("GCCPUFraction", models.NewGauge("GCCPUFraction", runtimeStats.GCCPUFraction))
	metricService.Health.Gauges.Store("GCSys", models.NewGauge("GCSys", float64(runtimeStats.GCSys)))
	metricService.Health.Gauges.Store("HeapAlloc", models.NewGauge("HeapAlloc", float64(runtimeStats.HeapAlloc)))
	metricService.Health.Gauges.Store("HeapIdle", models.NewGauge("HeapIdle", float64(runtimeStats.HeapIdle)))
	metricService.Health.Gauges.Store("HeapInuse", models.NewGauge("HeapInuse", float64(runtimeStats.HeapInuse)))
	metricService.Health.Gauges.Store("HeapSys", models.NewGauge("HeapSys", float64(runtimeStats.HeapSys)))
	metricService.Health.Gauges.Store("LastGC", models.NewGauge("LastGC", float64(runtimeStats.LastGC)))
	metricService.Health.Gauges.Store("Lookups", models.NewGauge("Lookups", float64(runtimeStats.Lookups)))
	metricService.Health.Gauges.Store("MCacheInuse", models.NewGauge("MCacheInuse", float64(runtimeStats.MCacheInuse)))
	metricService.Health.Gauges.Store("MCacheSys", models.NewGauge("MCacheSys", float64(runtimeStats.MCacheSys)))
	metricService.Health.Gauges.Store("MSpanInuse", models.NewGauge("MSpanInuse", float64(runtimeStats.MSpanInuse)))
	metricService.Health.Gauges.Store("MSpanSys", models.NewGauge("MSpanSys", float64(runtimeStats.MSpanSys)))
	metricService.Health.Gauges.Store("Mallocs", models.NewGauge("Mallocs", float64(runtimeStats.Mallocs)))
	metricService.Health.Gauges.Store("NextGC", models.NewGauge("NextGC", float64(runtimeStats.NextGC)))
	metricService.Health.Gauges.Store("NumForcedGC", models.NewGauge("NumForcedGC", float64(runtimeStats.NumForcedGC)))
	metricService.Health.Gauges.Store("NumGC", models.NewGauge("NumGC", float64(runtimeStats.NumGC)))
	metricService.Health.Gauges.Store("OtherSys", models.NewGauge("OtherSys", float64(runtimeStats.OtherSys)))
	metricService.Health.Gauges.Store("PauseTotalNs", models.NewGauge("PauseTotalNs", float64(runtimeStats.PauseTotalNs)))
	metricService.Health.Gauges.Store("StackInuse", models.NewGauge("StackInuse", float64(runtimeStats.StackInuse)))
	metricService.Health.Gauges.Store("StackSys", models.NewGauge("StackSys", float64(runtimeStats.StackSys)))
	metricService.Health.Gauges.Store("Sys", models.NewGauge("Sys", float64(runtimeStats.Sys)))
	metricService.Health.Gauges.Store("TotalAlloc", models.NewGauge("TotalAlloc", float64(runtimeStats.TotalAlloc)))
	metricService.Health.Gauges.Store("RandomValue", models.NewGauge("RandomValue", rand.Float64()))
}

func (metricService *MetricService) GatherMetricsByInterval(seconds int) {
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)

	for range ticker.C {
		metricService.gatherMetrics()
	}
}
