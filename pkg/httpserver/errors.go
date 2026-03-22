package httpserver

import "errors"

var (
	ErrBadRequest    = errors.New("bad request")
	ErrForbidden     = errors.New("forbidden")
	ErrNotFound      = errors.New("not found")
	ErrUnprocessable = errors.New("unprocessable")
	ErrInternalError = errors.New("internal server error")
)
