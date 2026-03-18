package httpserver

import "errors"

var (
	ErrBadRequest    = errors.New("bad request")
	Forbidden        = errors.New("forbidden")
	ErrNotFound      = errors.New("not found")
	ErrUnprocessable = errors.New("unprocessable")
	ErrInternalError = errors.New("internal server error")
)
