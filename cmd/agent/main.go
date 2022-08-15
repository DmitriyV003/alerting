package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/dmitriy/alerting/internal/agent/client"
	"github.com/dmitriy/alerting/internal/agent/service"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	var conf config
	err := env.Parse(&conf)
	if err != nil {
		log.Error("Unable to parse ENV: ", err)
	}
	log.Infof("Agent starting. Poll interval: %s; Report interval: %s", fmt.Sprint(conf.PollInterval), fmt.Sprint(conf.ReportInterval))

	metricService := service.New()
	go metricService.GatherMetricsByInterval(conf.PollInterval)

	sender := client.New()
	go sender.SendWithInterval(fmt.Sprintf("http://%s/update", conf.Address), &metricService.Health, conf.ReportInterval)

	select {}
}
