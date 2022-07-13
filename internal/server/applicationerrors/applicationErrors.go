package applicationerrors

import "errors"

var ErrNotFound = errors.New("not found")
var ErrUnknownType = errors.New("unknown type")
var ErrInvalidValue = errors.New("invalid value")
var ErrInvalidType = errors.New("invalid type")
