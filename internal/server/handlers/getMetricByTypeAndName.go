package handlers

import (
	"encoding/json"
	"errors"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/go-chi/chi/v5"
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
	name := chi.URLParam(r, "name")
	metricType := chi.URLParam(r, "type")
	metric, err := h.storage.GetByNameAndType(name, metricType)

	metricBytes, _ := json.Marshal(metric)

	if err != nil && errors.Is(err, applicationerrors.ErrNotFound) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

		return
	} else if err != nil && errors.Is(err, applicationerrors.ErrUnknownType) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(metricBytes)
	if err != nil {
		return
	}
}
