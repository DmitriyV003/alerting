package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmitriy/alerting/internal/agent/models"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/go-chi/chi/v5"
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

func (handler *UpdateMetricHandler) Handle(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	metricType := chi.URLParam(r, "type")
	value := chi.URLParam(r, "value")
	var metricReq model.Metric

	if name == "" && metricType == "" && value == "" {
		if err := json.NewDecoder(r.Body).Decode(&metricReq); err != nil {
			applicationerrors.WriteHTTPError(&w, http.StatusBadRequest)

			return
		}

		name = metricReq.Name
		metricType = string(metricReq.Type)

		if metricType == models.GaugeType {
			value = fmt.Sprint(*metricReq.FloatValue)
		} else if metricType == models.CounterType {
			value = fmt.Sprint(*metricReq.IntValue)
		}
	}

	if name == "" {
		applicationerrors.WriteHTTPError(&w, http.StatusNotFound)

		return
	}

	err := handler.storage.UpdateMetric(name, value, metricType)

	if err != nil && errors.Is(err, applicationerrors.ErrInvalidType) {
		applicationerrors.WriteHTTPError(&w, http.StatusNotImplemented)

		return
	} else if err != nil && errors.Is(err, applicationerrors.ErrInvalidValue) {
		applicationerrors.WriteHTTPError(&w, http.StatusBadRequest)

		return
	}
}
