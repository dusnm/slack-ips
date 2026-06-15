package container

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/dusnm/slack-ips/pkg/config"
	"github.com/dusnm/slack-ips/pkg/repositories/settings"
	"github.com/dusnm/slack-ips/pkg/repositories/user"
	"github.com/dusnm/slack-ips/pkg/services/messagehandler"
	"github.com/dusnm/slack-ips/pkg/services/qr"
	"github.com/dusnm/slack-ips/pkg/services/requestauth"
	"github.com/dusnm/slack-ips/pkg/services/templating"
	"github.com/dusnm/slack-ips/pkg/services/urlsign"
	"github.com/rs/zerolog"
)

type (
	Container struct {
		cfg                   *config.Config
		db                    *sql.DB
		userRepo              *user.Repository
		settingsRepo          *settings.Repository
		messageHandlerService *messagehandler.Service
		requestAuthService    *requestauth.Service
		urlSignService        *urlsign.Service
		templateService       *templating.Service
		qrService             *qr.Service
		AssetsFS              embed.FS
		templatesFS           embed.FS
		logger                zerolog.Logger
	}
)

func New(
	assetsFS embed.FS,
	templatesFS embed.FS,
	logger zerolog.Logger,
) *Container {
	return &Container{
		AssetsFS:    assetsFS,
		templatesFS: templatesFS,
		logger:      logger,
	}
}

func (c *Container) Close() error {
	c.logger.Info().Msg("closing")

	var err error
	if c.userRepo != nil {
		err = errors.Join(err, c.userRepo.Close())
	}

	if c.settingsRepo != nil {
		err = errors.Join(err, c.settingsRepo.Close())
	}

	if c.db != nil {
		err = errors.Join(err, c.db.Close())
	}

	return err
}
