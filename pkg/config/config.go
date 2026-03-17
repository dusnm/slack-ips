package config

import (
	"io"

	"github.com/pelletier/go-toml/v2"
)

type (
	Config struct {
		App   App   `toml:"app"`
		Slack Slack `toml:"slack"`
		DB    DB    `toml:"db"`
	}
)

func New(r io.Reader) (*Config, error) {
	cfg := &Config{}

	if err := toml.NewDecoder(r).Decode(cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	validators := []func() error{
		c.App.Validate,
		c.Slack.Validate,
		c.DB.Validate,
	}

	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}

	return nil
}
