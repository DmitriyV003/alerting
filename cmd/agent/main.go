package main

import (
	"fmt"

	"net/http"
	_ "net/http/pprof"

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
	go http.ListenAndServe(":8083", nil)
	app.run()
}
