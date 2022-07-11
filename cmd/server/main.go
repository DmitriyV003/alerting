package main

import (
	"github.com/dmitriy/alerting/internal/server/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	updateMetricHandler := handlers.NewUpdateMetricHandler()
	mux.HandleFunc("/update/", updateMetricHandler.Handle)

	log.Println("Server started on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
