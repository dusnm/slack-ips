package container

import "github.com/dusnm/slack-ips/pkg/repositories/user"

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
