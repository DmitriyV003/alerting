package applicationerrors

import (
	"errors"
	"net/http"
)

var ErrNotFound = errors.New("not found")
var ErrUnknownType = errors.New("unknown type")
var ErrInvalidValue = errors.New("invalid value")
var ErrInvalidType = errors.New("invalid type")

func WriteHTTPError(w *http.ResponseWriter, status int) {
	http.Error(*w, http.StatusText(status), status)
}
