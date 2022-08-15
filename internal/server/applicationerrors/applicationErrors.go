package applicationerrors

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var ErrNotFound = errors.New("not found")
var ErrUnknownType = errors.New("unknown type")
var ErrInvalidValue = errors.New("invalid value")
var ErrInvalidType = errors.New("invalid type")

func WriteHTTPError(w *http.ResponseWriter, status int) {
	http.Error(*w, http.StatusText(status), status)
}

func SwitchError(err error, w *http.ResponseWriter) {
	switch {
	case errors.Is(err, ErrNotFound):
		log.Info(fmt.Printf("Not Found"))
		WriteHTTPError(w, http.StatusNotFound)
	case errors.Is(err, ErrUnknownType):
		log.Info("Unknown metric type")
		WriteHTTPError(w, http.StatusBadRequest)
	default:
		log.Info("Unknown error")
		WriteHTTPError(w, http.StatusInternalServerError)
	}
}
