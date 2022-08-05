package handlers

import (
	"errors"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
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
	metricType := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

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
