package main

import (
	"github.com/dmitriy/alerting/internal/hasher"
	"github.com/dmitriy/alerting/internal/server/handlers"
	"github.com/dmitriy/alerting/internal/server/service"
	"github.com/dmitriy/alerting/internal/server/storage/memory"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
)

func (app *App) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Heartbeat("/ping"))

	store := memory.New()
	mHasher := hasher.New(app.conf.Key)
	updateMetricHandler := handlers.NewUpdateMetricHandler(store, mHasher)
	getAllMetricsHandler := handlers.NewGetAllMetricHandler(store)
	getMetricValueByTypeAndNameHandler := handlers.NewGetMetricValueByTypeAndNameHandler(store)
	getMetricByTypeAndNameHandler := handlers.NewGetMetricByTypeAndNameHandler(store, app.conf.Key)

	restore, _ := strconv.ParseBool(app.conf.Restore)
	fileSaver := service.NewFileSaver(app.conf.StoreFile, app.conf.StoreInterval, restore, store)
	fileSaver.Restore()

	if app.conf.StoreInterval.String() == "0s" {
		store.AddOnUpdateListener(fileSaver.StoreAllData)
	} else {
		go fileSaver.StoreAllDataWithInterval()
	}

	router.Get("/", getAllMetricsHandler.Handle)
	router.Get("/value/{type}/{name}", getMetricValueByTypeAndNameHandler.Handle)
	router.Post("/update/{type}/{name}/{value}", updateMetricHandler.Handle)
	router.Post("/value", getMetricByTypeAndNameHandler.Handle)
	router.Post("/update", updateMetricHandler.Handle)

	return router
}

func (app *App) config() {
	app.conf.parseEnv()
}
