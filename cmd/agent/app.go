package main

import (
	"fmt"
	"os"

	"github.com/dmitriy/alerting/internal/agent/client"
	"github.com/dmitriy/alerting/internal/agent/service"
	"github.com/dmitriy/alerting/internal/helpers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (app *App) config() {
	app.conf.parseEnv()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func (app *App) run() {
	publicKey, err := helpers.ImportPublicKeyFromFile(app.conf.PublicKey)
	if err != nil {
		log.Warn().Err(err).Msg("error to get public key from file")
	}

	metricService := service.NewMetricService(app.conf.Key)
	pingService := service.NewPingService()

	pingService.Ping(fmt.Sprintf("http://%s/heartbeat", app.conf.Address))

	go metricService.GatherMetricsByInterval(app.conf.PollInterval)
	go metricService.GatherAdditionalMetricsByInterval(app.conf.PollInterval)

	sender := client.New(publicKey)
	go sender.SendWithInterval(fmt.Sprintf("http://%s/update", app.conf.Address), &metricService.Health, app.conf.ReportInterval)

	select {}
}
