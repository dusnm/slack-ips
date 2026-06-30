package container

import (
	"github.com/dusnm/slack-ips/pkg/services/messagehandler"
	"github.com/dusnm/slack-ips/pkg/services/qr"
	"github.com/dusnm/slack-ips/pkg/services/qrcaption"
	"github.com/dusnm/slack-ips/pkg/services/requestauth"
	"github.com/dusnm/slack-ips/pkg/services/templating"
	"github.com/dusnm/slack-ips/pkg/services/urlsign"
)

func (c *Container) GetMessageHandlerService() *messagehandler.Service {
	if c.messageHandlerService == nil {
		c.messageHandlerService = messagehandler.New(
			c.GetConfig().App,
			c.GetUserRepository(),
			c.GetURLSignService(),
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

func (c *Container) GetURLSignService() *urlsign.Service {
	if c.urlSignService == nil {
		c.urlSignService = urlsign.New(
			c.GetConfig().App,
			c.logger.
				With().
				Str("component", "service:urlsign").
				Logger(),
		)
	}

	return c.urlSignService
}

func (c *Container) GetTemplateService() *templating.Service {
	if c.templateService == nil {
		c.templateService = templating.New(c.templatesFS)
	}

	return c.templateService
}

func (c *Container) GetQRService() *qr.Service {
	if c.qrService == nil {
		c.qrService = qr.New(c.GetQRCaptionService())
	}

	return c.qrService
}

func (c *Container) GetQRCaptionService() *qrcaption.Service {
	if c.qrCaptionService == nil {
		fontBytes, err := c.AssetsFS.ReadFile("assets/fonts/LibreBaskerville-Regular.ttf")
		if err != nil {
			c.logger.
				Fatal().
				Err(err).
				Str("component", "service:qrcaption").
				Msg("could not load font")
		}

		service, err := qrcaption.New(fontBytes)
		if err != nil {
			c.logger.
				Fatal().
				Err(err).
				Str("component", "service:qrcaption").
				Msg("could not instantiate")
		}

		c.qrCaptionService = service
	}

	return c.qrCaptionService
}
