package main

import (
	"fmt"
	"github.com/dmitriy/alerting/internal/agent/client"
	"github.com/dmitriy/alerting/internal/agent/service"
	log "github.com/sirupsen/logrus"
	"os"
)

func (app *App) config() {
	app.conf.parseEnv()
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func (app *App) run() {
	metricService := service.New()
	go metricService.GatherMetricsByInterval(app.conf.PollInterval)

	sender := client.New()
	go sender.SendWithInterval(fmt.Sprintf("http://%s/update", app.conf.Address), &metricService.Health, app.conf.ReportInterval)

	select {}
}
