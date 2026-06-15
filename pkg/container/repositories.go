package container

import (
	"github.com/dusnm/slack-ips/pkg/repositories/settings"
	"github.com/dusnm/slack-ips/pkg/repositories/user"
)

func (c *Container) GetUserRepository() *user.Repository {
	if c.userRepo == nil {
		c.userRepo = user.New(
			c.GetDB(),
			c.logger.
				With().
				Str("component", "repository:user").
				Logger(),
		)
	}

	return c.userRepo
}

func (c *Container) GetSettingsRepository() *settings.Repository {
	if c.settingsRepo == nil {
		c.settingsRepo = settings.New(
			c.GetDB(),
			c.logger.
				With().
				Str("component", "repository:settings").
				Logger(),
		)
	}

	return c.settingsRepo
}
