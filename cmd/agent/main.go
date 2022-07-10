package main

import (
	"github.com/dmitriy/alerting/internal/agent/client"
	"github.com/dmitriy/alerting/internal/agent/service"
)

func main() {

	metricService := service.New()
	go metricService.GatherMetricsByInterval(2)

	sender := client.New()
	go sender.SendWithInterval("http://localhost:8080", &metricService.Health, 10)

	select {}
}
