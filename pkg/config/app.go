package config

import (
	"errors"
	"net"
	"strconv"
)

var (
	ErrAppBindEmpty          = errors.New("bind cannot be empty")
	ErrAppPortEmpty          = errors.New("port must be greater than zero")
	ErrAppDomainEmpty        = errors.New("domain cannot be empty")
	ErrAppSigningSecretEmpty = errors.New("signing secret cannot be empty")
)

type (
	App struct {
		Bind          string `toml:"bind"`
		Port          uint16 `toml:"port"`
		Domain        string `toml:"domain"`
		Secure        bool   `toml:"secure"`
		BehindProxy   bool   `toml:"behind_proxy"`
		SigningSecret string `toml:"signing_secret"`
	}
)

func (a App) Validate() error {
	if len(a.Bind) == 0 {
		return ErrAppBindEmpty
	}

	if a.Port == 0 {
		return ErrAppPortEmpty
	}

	if len(a.Domain) == 0 {
		return ErrAppDomainEmpty
	}

	if len(a.SigningSecret) == 0 {
		return ErrAppSigningSecretEmpty
	}

	return nil
}

func (a App) Socket() string {
	return net.JoinHostPort(a.Bind, strconv.FormatUint(uint64(a.Port), 10))
}
