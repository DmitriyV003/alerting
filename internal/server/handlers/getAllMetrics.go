package handlers

import (
	"encoding/json"
	"github.com/dmitriy/alerting/internal/server/storage"
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

func (h *GetAllMetricHandler) Handle(w http.ResponseWriter, r *http.Request) {
	metrics := h.storage.GetAll()
	metricsBytes, _ := json.Marshal(metrics)

	w.WriteHeader(http.StatusOK)
	_, err := w.Write(metricsBytes)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}
