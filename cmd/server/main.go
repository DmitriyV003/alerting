package main

import (
	"github.com/dmitriy/alerting/internal/server/handlers"
	"github.com/dmitriy/alerting/internal/server/storage/memory"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.Default()
	//router.LoadHTMLGlob("./internal/server/templates/*")
	store := memory.New()
	updateMetricHandler := handlers.NewUpdateMetricHandler(store)
	getAllMetricsHandler := handlers.NewGetAllMetricHandler(store)
	getMetricByTypeAndNameHandler := handlers.NewGetMetricByTypeAndNameHandler(store)

	router.GET("/", getAllMetricsHandler.Handle)
	router.GET("/value/:type/:name", getMetricByTypeAndNameHandler.Handle)
	router.POST("/update/:type/:name/:value", updateMetricHandler.Handle)

	log.Fatal(router.Run(":8080"))
}
