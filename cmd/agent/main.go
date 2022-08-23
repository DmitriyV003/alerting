package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type App struct {
	conf Config
}

func main() {
	app := App{
		conf: Config{},
	}
	app.config()
	log.Infof("Agent starting. Poll interval: %s; Report interval: %s", fmt.Sprint(app.conf.PollInterval), fmt.Sprint(app.conf.ReportInterval))

	app.run()
}
