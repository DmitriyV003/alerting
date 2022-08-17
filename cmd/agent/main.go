package main

import (
	"fmt"
	"github.com/dmitriy/alerting/cmd/agent/config"
	"github.com/dmitriy/alerting/internal/agent/client"
	"github.com/dmitriy/alerting/internal/agent/service"
	log "github.com/sirupsen/logrus"
	"os"
)

var conf config.Config

func init() {
	config.ParseEnv(&conf)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Infof("Agent starting. Poll interval: %s; Report interval: %s", fmt.Sprint(conf.PollInterval), fmt.Sprint(conf.ReportInterval))

	metricService := service.New()
	go metricService.GatherMetricsByInterval(conf.PollInterval)

	sender := client.New()
	go sender.SendWithInterval(fmt.Sprintf("http://%s/update", conf.Address), &metricService.Health, conf.ReportInterval)

	select {}
}
