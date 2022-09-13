package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

type App struct {
	conf Config
}

func main() {
	app := App{
		conf: Config{},
	}
	app.config()
	log.Info().Msgf(
		"Agent starting. Poll interval: %s; Report interval: %s",
		fmt.Sprint(app.conf.PollInterval),
		fmt.Sprint(app.conf.ReportInterval),
	)

	app.run()
}
