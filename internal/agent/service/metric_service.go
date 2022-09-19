package service

import (
	"fmt"
	"github.com/dmitriy/alerting/internal/agent/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

type MetricService struct {
	Health models.Health
}

func NewMetricService(key string) *MetricService {
	return &MetricService{Health: *models.NewHealth(key)}
}

func (metricService *MetricService) gatherMetrics() {
	var runtimeStats runtime.MemStats
	rand.Seed(time.Now().UnixNano())

	runtime.ReadMemStats(&runtimeStats)

	currentCounter, ok := metricService.Health.Metrics.Load("PollCount")
	if !ok {
		metricService.Health.Store("PollCount", models.CounterType, fmt.Sprint(0))
	} else {
		metric := currentCounter.(models.Metric)
		metricService.Health.Store("PollCount", models.CounterType, strconv.FormatInt(*metric.IntValue+1, 10))
	}

	metricService.Health.Store("Alloc", models.GaugeType, strconv.FormatUint(runtimeStats.Alloc, 10))
	metricService.Health.Store("Frees", models.GaugeType, strconv.FormatUint(runtimeStats.Frees, 10))
	metricService.Health.Store("BuckHashSys", models.GaugeType, strconv.FormatUint(runtimeStats.BuckHashSys, 10))
	metricService.Health.Store("GCCPUFraction", models.GaugeType, strconv.FormatUint(uint64(runtimeStats.GCCPUFraction), 10))
	metricService.Health.Store("GCSys", models.GaugeType, strconv.FormatUint(runtimeStats.GCSys, 10))
	metricService.Health.Store("HeapAlloc", models.GaugeType, strconv.FormatUint(runtimeStats.HeapAlloc, 10))
	metricService.Health.Store("HeapObjects", models.GaugeType, strconv.FormatUint(runtimeStats.HeapObjects, 10))
	metricService.Health.Store("HeapReleased", models.GaugeType, strconv.FormatUint(runtimeStats.HeapReleased, 10))
	metricService.Health.Store("HeapIdle", models.GaugeType, strconv.FormatUint(runtimeStats.HeapIdle, 10))
	metricService.Health.Store("HeapInuse", models.GaugeType, strconv.FormatUint(runtimeStats.HeapInuse, 10))
	metricService.Health.Store("HeapSys", models.GaugeType, strconv.FormatUint(runtimeStats.HeapSys, 10))
	metricService.Health.Store("LastGC", models.GaugeType, strconv.FormatUint(runtimeStats.LastGC, 10))
	metricService.Health.Store("Lookups", models.GaugeType, strconv.FormatUint(runtimeStats.Lookups, 10))
	metricService.Health.Store("MCacheInuse", models.GaugeType, strconv.FormatUint(runtimeStats.MCacheInuse, 10))
	metricService.Health.Store("MCacheSys", models.GaugeType, strconv.FormatUint(runtimeStats.MCacheSys, 10))
	metricService.Health.Store("MSpanInuse", models.GaugeType, strconv.FormatUint(runtimeStats.MSpanInuse, 10))
	metricService.Health.Store("MSpanSys", models.GaugeType, strconv.FormatUint(runtimeStats.MSpanSys, 10))
	metricService.Health.Store("Mallocs", models.GaugeType, strconv.FormatUint(runtimeStats.Mallocs, 10))
	metricService.Health.Store("NextGC", models.GaugeType, strconv.FormatUint(runtimeStats.NextGC, 10))
	metricService.Health.Store("NumForcedGC", models.GaugeType, strconv.FormatUint(uint64(runtimeStats.NumForcedGC), 10))
	metricService.Health.Store("NumGC", models.GaugeType, strconv.FormatUint(uint64(runtimeStats.NumGC), 10))
	metricService.Health.Store("OtherSys", models.GaugeType, strconv.FormatUint(runtimeStats.OtherSys, 10))
	metricService.Health.Store("PauseTotalNs", models.GaugeType, strconv.FormatUint(runtimeStats.PauseTotalNs, 10))
	metricService.Health.Store("StackInuse", models.GaugeType, strconv.FormatUint(runtimeStats.StackInuse, 10))
	metricService.Health.Store("StackSys", models.GaugeType, strconv.FormatUint(runtimeStats.StackSys, 10))
	metricService.Health.Store("Sys", models.GaugeType, strconv.FormatUint(runtimeStats.Sys, 10))
	metricService.Health.Store("TotalAlloc", models.GaugeType, strconv.FormatUint(runtimeStats.TotalAlloc, 10))
	metricService.Health.Store("RandomValue", models.GaugeType, fmt.Sprint(rand.Float64()))
}

func (metricService *MetricService) gatherAdditionalMetrics() {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		log.Error("Unable to gather virtual memory metrics: ", err)
		return
	}

	cpuUsage, err := cpu.Percent(3*time.Second, false)
	if err != nil {
		log.Error("Unable to gather CPU usage: ", err)
		return
	}

	metricService.Health.Store("TotalMemory", models.GaugeType, strconv.FormatUint(virtualMemory.Total, 10))
	metricService.Health.Store("FreeMemory", models.GaugeType, strconv.FormatUint(virtualMemory.Free, 10))
	metricService.Health.Store("CPUutilization1", models.GaugeType, strconv.FormatFloat(cpuUsage[0], 'f', 6, 64))
}

func (metricService *MetricService) GatherMetricsByInterval(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {
		metricService.gatherMetrics()
	}
}

func (metricService *MetricService) GatherAdditionalMetricsByInterval(duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {
		metricService.gatherAdditionalMetrics()
	}
}
