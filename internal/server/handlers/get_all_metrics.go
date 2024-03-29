package handlers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
)

type GetAllMetricHandler struct {
	storage storage.MetricStorage
}

type viewData struct {
	Metrics []model.Metric
}

func NewGetAllMetricHandler(store storage.MetricStorage) *GetAllMetricHandler {
	return &GetAllMetricHandler{
		storage: store,
	}
}

// Handle Get all Metrics
func (h *GetAllMetricHandler) Handle(w http.ResponseWriter, r *http.Request) {
	metrics := h.storage.GetAll(context.Background())
	viewMetrics := viewData{Metrics: *metrics}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ts := template.Must(template.ParseFiles("./internal/server/templates/index.gohtml"))
	err := ts.Execute(w, viewMetrics)

	if err != nil {
		applicationerrors.SwitchError(err, &w)

		return
	}
}
