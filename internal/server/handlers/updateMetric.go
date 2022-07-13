package handlers

import (
	"errors"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UpdateMetricHandler struct {
	storage storage.MetricStorage
}

func NewUpdateMetricHandler(store storage.MetricStorage) *UpdateMetricHandler {
	return &UpdateMetricHandler{
		storage: store,
	}
}

func (handler *UpdateMetricHandler) Handle(c *gin.Context) {
	metricType := c.Param("type")
	name := c.Param("name")
	value := c.Param("value")

	if name == "" {
		c.AbortWithStatus(http.StatusNotFound)

		return
	}

	err := handler.storage.UpdateMetric(name, value, metricType)

	if err != nil && errors.Is(err, applicationerrors.ErrInvalidType) {
		c.AbortWithStatus(http.StatusNotImplemented)

		return
	} else if err != nil && errors.Is(err, applicationerrors.ErrInvalidValue) {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}
}
