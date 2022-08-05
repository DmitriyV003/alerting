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
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	}

	err := handler.storage.UpdateMetric(name, value, metricType)

	if err != nil && errors.Is(err, applicationerrors.ErrInvalidType) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)

		return
	} else if err != nil && errors.Is(err, applicationerrors.ErrInvalidValue) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}
}
