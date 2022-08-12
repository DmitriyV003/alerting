package main

import (
	"github.com/dmitriy/alerting/internal/server/handlers"
	"github.com/dmitriy/alerting/internal/server/storage/memory"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	store := memory.New()
	updateMetricHandler := handlers.NewUpdateMetricHandler(store)
	getAllMetricsHandler := handlers.NewGetAllMetricHandler(store)
	getMetricByTypeAndNameHandler := handlers.NewGetMetricByTypeAndNameHandler(store)

	router.Get("/", getAllMetricsHandler.Handle)
	router.Get("/value/{type}/{name}", getMetricByTypeAndNameHandler.Handle)
	router.Post("/update/{type}/{name}/{value}", updateMetricHandler.Handle)

	log.Info("server is starting at http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8090", router))
}
