package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmitriy/alerting/internal/server/applicationerrors"
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/go-chi/chi/v5"
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
	name := chi.URLParam(r, "name")
	metricType := chi.URLParam(r, "type")
	metric, err := h.storage.GetByNameAndType(name, metricType)

	if name == "" {
		log.Info(fmt.Printf("Metric Not Found "))
		applicationerrors.WriteHTTPError(&w, http.StatusNotFound)

		return
	}

	if err != nil {
		switch {
		case errors.Is(err, applicationerrors.ErrNotFound):
			log.Info(fmt.Printf("Not Found"))
			applicationerrors.WriteHTTPError(&w, http.StatusNotFound)
		case errors.Is(err, applicationerrors.ErrUnknownType):
			log.Info("Unknown metric type")
			applicationerrors.WriteHTTPError(&w, http.StatusBadRequest)
		default:
			log.Info("Unknown error")
			applicationerrors.WriteHTTPError(&w, http.StatusInternalServerError)
		}

		return
	}

	metricBytes, err := json.Marshal(metric)

	if err != nil {
		log.Info("Unknown error")
		applicationerrors.WriteHTTPError(&w, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(metricBytes)
	if err != nil {
		log.Info("Unknown error")

		return
	}
}
