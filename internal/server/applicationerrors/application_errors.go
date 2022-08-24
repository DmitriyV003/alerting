package applicationerrors

import (
	"errors"
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
		log.Error("Not Found: ", err)
		WriteHTTPError(w, http.StatusNotFound)
	case errors.Is(err, ErrUnknownType):
		log.Error("Unknown metric type: ", err)
		WriteHTTPError(w, http.StatusBadRequest)
	case errors.Is(err, ErrInvalidType):
		log.Error("Type does not supported: ", err)
		WriteHTTPError(w, http.StatusNotImplemented)
	case errors.Is(err, ErrInvalidValue):
		log.Error("Invalid metric value: ", err)
		WriteHTTPError(w, http.StatusBadRequest)
	default:
		log.Error("Unknown error: ", err)
		WriteHTTPError(w, http.StatusInternalServerError)
	}
}
