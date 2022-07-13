package handlers

import (
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetMetricByTypeAndNameHandler struct {
	storage storage.MetricStorage
}

func NewGetMetricByTypeAndNameHandler(store storage.MetricStorage) *GetMetricByTypeAndNameHandler {
	return &GetMetricByTypeAndNameHandler{
		storage: store,
	}
}

func (h *GetMetricByTypeAndNameHandler) Handle(c *gin.Context) {
	name := c.Param("name")
	metricType := c.Param("type")
	metric, err := h.storage.GetByNameAndType(name, metricType)

	if err != nil && err.Error() == "not found" {
		c.AbortWithStatus(http.StatusNotFound)

		return
	} else if err != nil && err.Error() == "unknown type" {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusOK, metric)
}
