package config

import "errors"

var (
	ErrPathEmpty = errors.New("path cannot be empty")
)

type (
	DB struct {
		Path string `toml:"path"`
	}
)

func (d DB) Validate() error {
	if len(d.Path) == 0 {
		return ErrPathEmpty
	}

	return nil
}
