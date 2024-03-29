package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmitriy/alerting/internal/agent/models"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type GetMetricValueByTypeAndNameHandler struct {
	storage storage.MetricStorage
}

func NewGetMetricValueByTypeAndNameHandler(store storage.MetricStorage) *GetMetricValueByTypeAndNameHandler {
	return &GetMetricValueByTypeAndNameHandler{
		storage: store,
	}
}

// Handle Get metric value by type and name
func (h *GetMetricValueByTypeAndNameHandler) Handle(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	metricType := chi.URLParam(r, "type")
	metric, err := h.storage.GetByNameAndType(context.Background(), name, metricType)

	if name == "" {
		log.Info(fmt.Printf("Metric Not Found "))
		applicationerrors.WriteHTTPError(&w, http.StatusNotFound)

		return
	}

	if err != nil {
		applicationerrors.SwitchError(err, &w)

		return
	}

	var metricBytes []byte

	if metric.Type == models.GaugeType {
		metricBytes, err = json.Marshal(metric.FloatValue)
	} else if metric.Type == models.CounterType {
		metricBytes, err = json.Marshal(metric.IntValue)
	}

	if err != nil {
		applicationerrors.SwitchError(err, &w)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(metricBytes)
	if err != nil {
		applicationerrors.SwitchError(err, &w)

		return
	}
}
