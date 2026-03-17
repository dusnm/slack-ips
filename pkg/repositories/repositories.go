package repositories

import "errors"

var (
	ErrNotFound = errors.New("the requested resource was not found")
)
