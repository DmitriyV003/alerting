package main

import (
	"net/http"

	_ "net/http/pprof"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type App struct {
	conf Config
	pool *pgxpool.Pool
}

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

// go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildTime=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git rev-parse HEAD)'" cmd/server/main.go

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

	if app.pool != nil {
		defer app.pool.Close()
	}

	log.Infof("server is starting at %s", app.conf.Address)
	srv := &http.Server{
		Addr:    app.conf.Address,
		Handler: app.routes(),
	}
	go http.ListenAndServe(":8082", nil)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
