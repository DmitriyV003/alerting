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

	currentCounter, ok := metricService.Health.Metrics.Load("PollCount")
	if !ok {
		metricService.Health.Metrics.Store("PollCount", models.NewCounter("PollCount", 0))
	} else {
		metric := currentCounter.(models.Metric)
		metricService.Health.Metrics.Store("PollCount", models.NewCounter("PollCount", *metric.IntValue+1))
	}

	metricService.Health.Metrics.Store("Alloc", models.NewGauge("Alloc", float64(runtimeStats.Alloc)))
	metricService.Health.Metrics.Store("Frees", models.NewGauge("Frees", float64(runtimeStats.Frees)))
	metricService.Health.Metrics.Store("BuckHashSys", models.NewGauge("BuckHashSys", float64(runtimeStats.BuckHashSys)))
	metricService.Health.Metrics.Store("GCCPUFraction", models.NewGauge("GCCPUFraction", runtimeStats.GCCPUFraction))
	metricService.Health.Metrics.Store("GCSys", models.NewGauge("GCSys", float64(runtimeStats.GCSys)))
	metricService.Health.Metrics.Store("HeapAlloc", models.NewGauge("HeapAlloc", float64(runtimeStats.HeapAlloc)))
	metricService.Health.Metrics.Store("HeapObjects", models.NewGauge("HeapObjects", float64(runtimeStats.HeapObjects)))
	metricService.Health.Metrics.Store("HeapReleased", models.NewGauge("HeapReleased", float64(runtimeStats.HeapReleased)))
	metricService.Health.Metrics.Store("HeapIdle", models.NewGauge("HeapIdle", float64(runtimeStats.HeapIdle)))
	metricService.Health.Metrics.Store("HeapInuse", models.NewGauge("HeapInuse", float64(runtimeStats.HeapInuse)))
	metricService.Health.Metrics.Store("HeapSys", models.NewGauge("HeapSys", float64(runtimeStats.HeapSys)))
	metricService.Health.Metrics.Store("LastGC", models.NewGauge("LastGC", float64(runtimeStats.LastGC)))
	metricService.Health.Metrics.Store("Lookups", models.NewGauge("Lookups", float64(runtimeStats.Lookups)))
	metricService.Health.Metrics.Store("MCacheInuse", models.NewGauge("MCacheInuse", float64(runtimeStats.MCacheInuse)))
	metricService.Health.Metrics.Store("MCacheSys", models.NewGauge("MCacheSys", float64(runtimeStats.MCacheSys)))
	metricService.Health.Metrics.Store("MSpanInuse", models.NewGauge("MSpanInuse", float64(runtimeStats.MSpanInuse)))
	metricService.Health.Metrics.Store("MSpanSys", models.NewGauge("MSpanSys", float64(runtimeStats.MSpanSys)))
	metricService.Health.Metrics.Store("Mallocs", models.NewGauge("Mallocs", float64(runtimeStats.Mallocs)))
	metricService.Health.Metrics.Store("NextGC", models.NewGauge("NextGC", float64(runtimeStats.NextGC)))
	metricService.Health.Metrics.Store("NumForcedGC", models.NewGauge("NumForcedGC", float64(runtimeStats.NumForcedGC)))
	metricService.Health.Metrics.Store("NumGC", models.NewGauge("NumGC", float64(runtimeStats.NumGC)))
	metricService.Health.Metrics.Store("OtherSys", models.NewGauge("OtherSys", float64(runtimeStats.OtherSys)))
	metricService.Health.Metrics.Store("PauseTotalNs", models.NewGauge("PauseTotalNs", float64(runtimeStats.PauseTotalNs)))
	metricService.Health.Metrics.Store("StackInuse", models.NewGauge("StackInuse", float64(runtimeStats.StackInuse)))
	metricService.Health.Metrics.Store("StackSys", models.NewGauge("StackSys", float64(runtimeStats.StackSys)))
	metricService.Health.Metrics.Store("Sys", models.NewGauge("Sys", float64(runtimeStats.Sys)))
	metricService.Health.Metrics.Store("TotalAlloc", models.NewGauge("TotalAlloc", float64(runtimeStats.TotalAlloc)))
	metricService.Health.Metrics.Store("RandomValue", models.NewGauge("RandomValue", rand.Float64()))
}

func (metricService *MetricService) GatherMetricsByInterval(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {
		metricService.gatherMetrics()
	}
}
