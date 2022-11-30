package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
	log "github.com/sirupsen/logrus"
)

type UpdateMetricsCollectionHandler struct {
	storage storage.MetricStorage
}

func NewUpdateMetricsCollectionHandler(store storage.MetricStorage) *UpdateMetricsCollectionHandler {
	return &UpdateMetricsCollectionHandler{
		storage: store,
	}
}

func (h *UpdateMetricsCollectionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var metricRequest []model.Metric
	if err := json.NewDecoder(r.Body).Decode(&metricRequest); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("error to decode request")
		applicationerrors.WriteHTTPError(&w, http.StatusBadRequest)

		return
	}

	err := h.storage.SaveCollection(context.Background(), &metricRequest)
	if err != nil {
		applicationerrors.WriteHTTPError(&w, http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
}
