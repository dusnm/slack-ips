package container

import "github.com/rs/zerolog"

func (c *Container) GetLogger() zerolog.Logger {
	return c.logger
}
