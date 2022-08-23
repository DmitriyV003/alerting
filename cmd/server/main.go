package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type App struct {
	conf Config
}

func main() {
	app := App{
		conf: Config{},
	}
	app.config()

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
