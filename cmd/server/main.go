package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/dmitriy/alerting/internal/server/handlers"
	"github.com/dmitriy/alerting/internal/server/service"
	"github.com/dmitriy/alerting/internal/server/storage/memory"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type config struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func main() {
	var conf config
	err := env.Parse(&conf)
	if err != nil {
		log.Error("Unable to parse ENV: ", err)
	}
	if conf.StoreFile == "" {
		conf.StoreFile = "/tmp/devops-metrics-db.json"
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)

	store := memory.New()
	updateMetricHandler := handlers.NewUpdateMetricHandler(store)
	getAllMetricsHandler := handlers.NewGetAllMetricHandler(store)
	getMetricValueByTypeAndNameHandler := handlers.NewGetMetricValueByTypeAndNameHandler(store)
	getMetricByTypeAndNameHandler := handlers.NewGetMetricByTypeAndNameHandler(store)

	fileSaver := service.NewFileSaver(conf.StoreFile, conf.StoreInterval, conf.Restore, store)
	log.Info("ENV VARS: ", conf.StoreFile, conf.StoreInterval, conf.Restore)
	fileSaver.Restore()
	if conf.StoreInterval == 0 {
		store.AddOnUpdateListener(fileSaver.StoreAllData)
	} else {
		go fileSaver.StoreAllDataWithInterval()
	}

	router.Get("/", getAllMetricsHandler.Handle)
	router.Get("/value/{type}/{name}", getMetricValueByTypeAndNameHandler.Handle)
	router.Post("/update/{type}/{name}/{value}", updateMetricHandler.Handle)
	router.Post("/value", getMetricByTypeAndNameHandler.Handle)
	router.Post("/update", updateMetricHandler.Handle)

	log.Infof("server is starting at %s", conf.Address)
	log.Fatal(http.ListenAndServe(conf.Address, router))
}
