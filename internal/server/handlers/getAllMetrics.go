package handlers

import (
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetAllMetricHandler struct {
	storage storage.MetricStorage
}

func NewGetAllMetricHandler(store storage.MetricStorage) *GetAllMetricHandler {
	return &GetAllMetricHandler{
		storage: store,
	}
}

func (h *GetAllMetricHandler) Handle(c *gin.Context) {
	metrics := h.storage.GetAll()
	c.JSON(http.StatusOK, *metrics)
}
