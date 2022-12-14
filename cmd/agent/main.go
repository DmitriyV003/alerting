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

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

// go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildTime=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git rev-parse HEAD)'" cmd/agent/main.go

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildTime == "" {
		buildTime = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	log.Printf("Build version=%s, Build date=%s\n, Build commit=%s\n", buildVersion, buildTime, buildCommit)
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
