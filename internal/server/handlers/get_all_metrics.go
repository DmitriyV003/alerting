package handlers

import (
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/model"
	"github.com/dmitriy/alerting/internal/server/storage"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
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

func (h *GetAllMetricHandler) Handle(w http.ResponseWriter, r *http.Request) {
	metrics := h.storage.GetAll()
	viewMetrics := viewData{Metrics: *metrics}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ts := template.Must(template.ParseFiles("./internal/server/templates/index.gohtml"))
	err := ts.Execute(w, viewMetrics)

	if err != nil {
		log.Error("Exec html: ", err)
		applicationerrors.WriteHTTPError(&w, http.StatusInternalServerError)

		return
	}
}