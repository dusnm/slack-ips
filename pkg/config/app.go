package config

import (
	"errors"
	"net"
	"strconv"
)

var (
	ErrBindEmpty   = errors.New("bind cannot be empty")
	ErrPortEmpty   = errors.New("port must be greater than zero")
	ErrDomainEmpty = errors.New("domain cannot be empty")
)

type (
	App struct {
		Bind   string `toml:"bind"`
		Port   uint16 `toml:"port"`
		Domain string `toml:"domain"`
		Secure bool   `toml:"secure"`
	}
)

func (a App) Validate() error {
	if len(a.Bind) == 0 {
		return ErrBindEmpty
	}

	if a.Port == 0 {
		return ErrPortEmpty
	}

	if len(a.Domain) == 0 {
		return ErrDomainEmpty
	}

	return nil
}

func (a App) Socket() string {
	return net.JoinHostPort(a.Bind, strconv.FormatUint(uint64(a.Port), 10))
}
