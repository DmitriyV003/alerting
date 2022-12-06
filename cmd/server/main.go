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
