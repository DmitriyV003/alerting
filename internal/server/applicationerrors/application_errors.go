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
var ErrInternalServer = errors.New("internal server error")

func WriteHTTPError(w *http.ResponseWriter, status int) {
	http.Error(*w, http.StatusText(status), status)
}

func SwitchError(err error, w *http.ResponseWriter) {
	switch {
	case errors.Is(err, ErrNotFound):
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Not Found")
		WriteHTTPError(w, http.StatusNotFound)
	case errors.Is(err, ErrUnknownType):
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Unknown metric type")
		WriteHTTPError(w, http.StatusBadRequest)
	case errors.Is(err, ErrInvalidType):
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Type does not supported")
		WriteHTTPError(w, http.StatusNotImplemented)
	case errors.Is(err, ErrInvalidValue):
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Invalid metric value")
		WriteHTTPError(w, http.StatusBadRequest)
	default:
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Unknown error")
		WriteHTTPError(w, http.StatusInternalServerError)
	}
}
