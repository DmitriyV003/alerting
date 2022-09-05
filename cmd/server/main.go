package main

import (
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type App struct {
	conf Config
	pool *pgxpool.Pool
}

func main() {
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
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
