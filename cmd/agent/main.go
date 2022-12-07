package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	_ "net/http/pprof"
)

type App struct {
	conf Config
}

var (
	BuildVersion string
	BuildTime    string
	BuildCommit  string
)

func main() {
	if BuildVersion == "" {
		BuildVersion = "N/A"
	}
	if BuildTime == "" {
		BuildTime = "N/A"
	}
	if BuildCommit == "" {
		BuildCommit = "N/A"
	}
	log.Printf("Build version=%s, Build date=%s\n, Build commit=%s\n", BuildVersion, BuildTime, BuildCommit)
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
