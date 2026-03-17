package container

import (
	"github.com/dusnm/slack-ips/pkg/services/messagehandler"
	"github.com/dusnm/slack-ips/pkg/services/requestauth"
)

func (c *Container) GetMessageHandlerService() *messagehandler.Service {
	if c.messageHandlerService == nil {
		c.messageHandlerService = messagehandler.New(
			c.GetConfig().App,
			c.GetUserRepository(),
			c.logger.
				With().
				Str("component", "service:message_handler").
				Logger(),
		)
	}

	return c.messageHandlerService
}

func (c *Container) GetRequestAuthService() *requestauth.Service {
	if c.requestAuthService == nil {
		c.requestAuthService = requestauth.New(
			c.logger.
				With().
				Str("component", "service:request_auth").
				Logger(),
		)
	}

	return c.requestAuthService
}
