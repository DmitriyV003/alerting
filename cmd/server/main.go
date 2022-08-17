package main

import (
	"github.com/dmitriy/alerting/cmd/server/config"
	"github.com/dmitriy/alerting/internal/server/handlers"
	"github.com/dmitriy/alerting/internal/server/service"
	"github.com/dmitriy/alerting/internal/server/storage/memory"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var conf config.Config

func init() {
	config.ParseEnv(&conf)
}

func main() {
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

	restore, _ := strconv.ParseBool(conf.Restore)
	fileSaver := service.NewFileSaver(conf.StoreFile, conf.StoreInterval, restore, store)
	fileSaver.Restore()

	if conf.StoreInterval.String() == "0s" {
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
