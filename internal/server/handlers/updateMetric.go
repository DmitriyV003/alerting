package handlers

import (
	"fmt"
	"github.com/dmitriy/alerting/internal/server/storage"
	"github.com/dmitriy/alerting/internal/server/storage/memory"
	"net/http"
	"regexp"
)

type UpdateMetricHandler struct {
	storage storage.MetricStorage
}

func NewUpdateMetricHandler() *UpdateMetricHandler {
	return &UpdateMetricHandler{
		storage: memory.New(),
	}
}

func (handler *UpdateMetricHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	regex, regErr := regexp.Compile("^/update/(\\w+)/(\\w+|[\\-\\d.,]+)/(\\w+|[\\-\\d.,]+)$")
	if regErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte("Bad regex"))
		return
	}
	matches := regex.FindSubmatch([]byte(r.URL.Path))

	if matches == nil {
		http.Error(rw, "Error", http.StatusNotFound)
		return
	}

	fmt.Println(fmt.Sprintf("Got message with: %s, %s, %s", string(matches[2]), string(matches[3]), string(matches[1])))
	err := handler.storage.UpdateMetric(string(matches[2]), string(matches[3]), string(matches[1]))

	if err != nil && err.Error() == "invalid type" {
		http.Error(rw, "Error", http.StatusNotImplemented)
		return
	} else if err != nil && err.Error() == "invalid value" {
		http.Error(rw, "Error", http.StatusBadRequest)
		return
	}
}
