package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/dmitriy/alerting/internal/server/transformers"
	log "github.com/sirupsen/logrus"
)

type GetMetricByTypeAndNameHandler struct {
	storage           storage.MetricStorage
	metricTransformer *transformers.MetricTransformer
}

type metricRequest struct {
	Name string           `json:"id"`
	Type model.MetricType `json:"type"`
}

func NewGetMetricByTypeAndNameHandler(store storage.MetricStorage, key string) *GetMetricByTypeAndNameHandler {
	return &GetMetricByTypeAndNameHandler{
		storage:           store,
		metricTransformer: transformers.NewTransformer(key),
	}
}

// Handle Get metric by type and name
func (h *GetMetricByTypeAndNameHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var name string
	var metricReq metricRequest

	if err := json.NewDecoder(r.Body).Decode(&metricReq); err != nil {
		applicationerrors.WriteHTTPError(&w, http.StatusBadRequest)

		return
	}

	name = metricReq.Name
	metric, err := h.storage.GetByNameAndType(context.Background(), name, string(metricReq.Type))

	if name == "" {
		applicationerrors.WriteHTTPError(&w, http.StatusNotFound)

		return
	}

	if err != nil {
		applicationerrors.SwitchError(err, &w)

		return
	}

	metric = h.metricTransformer.AddHash(metric)
	metricBytes, err := json.Marshal(metric)

	if err != nil {
		log.Error("Unknown error: ", err)
		applicationerrors.WriteHTTPError(&w, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(metricBytes)

	log.WithFields(log.Fields{
		"Metric": metricReq.Name,
		"Type":   metricReq.Type,
	}).Info("GET metric by Type and Name")

	if err != nil {
		applicationerrors.SwitchError(err, &w)

		return
	}
}
