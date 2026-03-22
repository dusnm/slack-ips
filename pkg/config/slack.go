package config

import "errors"

var (
	ErrSlackAppIDEmpty         = errors.New("app id cannot be empty")
	ErrSlackClientIDEmpty      = errors.New("client id cannot be empty")
	ErrSlackClientSecretEmpty  = errors.New("client secret cannot be empty")
	ErrSlackSigningSecretEmpty = errors.New("signing secret cannot be empty")
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
		return ErrSlackAppIDEmpty
	}

	if len(s.ClientID) == 0 {
		return ErrSlackClientIDEmpty
	}

	if len(s.ClientSecret) == 0 {
		return ErrSlackClientSecretEmpty
	}

	if len(s.SigningSecret) == 0 {
		return ErrSlackSigningSecretEmpty
	}

	return nil
}
