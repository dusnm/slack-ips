package config

import "errors"

var (
	ErrAppIDEmpty         = errors.New("app id cannot be empty")
	ErrClientIDEmpty      = errors.New("client id cannot be empty")
	ErrClientSecretEmpty  = errors.New("client secret cannot be empty")
	ErrSigningSecretEmpty = errors.New("signing secret cannot be empty")
)

type (
	Slack struct {
		AppID         string `toml:"app_id"`
		ClientID      string `toml:"client_id"`
		ClientSecret  string `toml:"client_secret"`
		SigningSecret string `toml:"signing_secret"`
	}
)

func (s Slack) Validate() error {
	if len(s.AppID) == 0 {
		return ErrAppIDEmpty
	}

	if len(s.ClientID) == 0 {
		return ErrClientIDEmpty
	}

	if len(s.ClientSecret) == 0 {
		return ErrClientSecretEmpty
	}

	if len(s.SigningSecret) == 0 {
		return ErrSigningSecretEmpty
	}

	return nil
}
