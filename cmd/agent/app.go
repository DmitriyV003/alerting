package main

import (
	"fmt"
	"github.com/dmitriy/alerting/internal/agent/client"
	"github.com/dmitriy/alerting/internal/agent/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
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
	log.Info("Application config: ", app.conf)
	metricService := service.New(app.conf.Key)

	ticker := time.NewTicker(time.Second)
	clientPing := http.Client{}
	for range ticker.C {
		err := app.ping(&clientPing)
		if err == nil {
			ticker.Stop()
			break
		}
	}

	go metricService.GatherMetricsByInterval(app.conf.PollInterval)

	sender := client.New()
	go sender.SendWithInterval(fmt.Sprintf("http://%s/update", app.conf.Address), &metricService.Health, app.conf.ReportInterval)

	select {}
}

func (app *App) ping(client *http.Client) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/heartbeat", app.conf.Address), nil)
	if err != nil {
		log.Error("Request failed: ", err)
		return err
	}
	res, err := client.Do(request)
	log.Info("Ping Server: ", err)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
