package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (app *App) config() {
	app.conf.parseEnv()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
