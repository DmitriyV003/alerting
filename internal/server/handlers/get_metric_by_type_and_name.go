package handlers

import (
	"encoding/json"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
	log "github.com/sirupsen/logrus"
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

func (h *GetMetricByTypeAndNameHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var name string
	var metricReq model.Metric

	if err := json.NewDecoder(r.Body).Decode(&metricReq); err != nil {
		applicationerrors.WriteHTTPError(&w, http.StatusBadRequest)

		return
	}

	name = metricReq.Name
	metric, err := h.storage.GetByNameAndType(name, string(metricReq.Type))

	if name == "" {
		log.Infof("Metric Not Found: %s", name)
		applicationerrors.WriteHTTPError(&w, http.StatusNotFound)

		return
	}

	if err != nil {
		applicationerrors.SwitchError(err, &w)

		return
	}

	metricBytes, err := json.Marshal(metric)

	if err != nil {
		log.Info("Unknown error: ", err)
		applicationerrors.WriteHTTPError(&w, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(metricBytes)

	log.WithFields(log.Fields{
		"Metric": metricReq.Name,
		"Type":   metricReq.Type,
	}).Info("Got metric by Type and Name")

	if err != nil {
		applicationerrors.SwitchError(err, &w)

		return
	}
}
