package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmitriy/alerting/internal/agent/models"
	"github.com/dmitriy/alerting/internal/hasher"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/service"
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type UpdateMetricHandler struct {
	storage   storage.MetricStorage
	fileSaver *service.FileSaver
	hasher    *hasher.Hasher
}

func NewUpdateMetricHandler(store storage.MetricStorage, hasher *hasher.Hasher) *UpdateMetricHandler {
	return &UpdateMetricHandler{
		storage: store,
		hasher:  hasher,
	}
}

// Handle Update metric with given params
func (handler *UpdateMetricHandler) Handle(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	metricType := chi.URLParam(r, "type")
	value := chi.URLParam(r, "value")
	var metricReq model.Metric
	var calculatedHash string
	var stringToHash string

	if name == "" && metricType == "" && value == "" {
		if err := json.NewDecoder(r.Body).Decode(&metricReq); err != nil {
			applicationerrors.WriteHTTPError(&w, http.StatusBadRequest)

			return
		}

		name = metricReq.Name
		metricType = string(metricReq.Type)

		if metricType == models.GaugeType {
			value = fmt.Sprint(*metricReq.FloatValue)
			stringToHash = fmt.Sprintf("%s:%s:%f", metricReq.Name, metricType, *metricReq.FloatValue)
		} else if metricType == models.CounterType {
			value = fmt.Sprint(*metricReq.IntValue)
			stringToHash = fmt.Sprintf("%s:%s:%d", metricReq.Name, metricType, *metricReq.IntValue)
		}
	}

	if metricReq.Hash != "" {
		calculatedHash = handler.hasher.Hash(stringToHash)

		if !handler.hasher.IsEqual(metricReq.Hash, calculatedHash) {
			applicationerrors.WriteHTTPError(&w, http.StatusBadRequest)

			return
		}
	}

	if name == "" {
		applicationerrors.WriteHTTPError(&w, http.StatusNotFound)

		return
	}

	err := handler.storage.UpdateOrCreate(context.Background(), name, value, metricType)

	if err != nil {
		applicationerrors.SwitchError(err, &w)

		return
	}

	log.WithFields(log.Fields{
		"Name":  name,
		"Type":  metricType,
		"Value": value,
		"Hash":  calculatedHash,
	}).Info("Updated metric!")

	res, _ := json.Marshal(updateMetricResponse{Hash: calculatedHash})
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}

type updateMetricResponse struct {
	Hash string `json:"hash,omitempty"`
}
