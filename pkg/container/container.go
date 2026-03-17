package container

import (
	"database/sql"
	"errors"

	"github.com/dusnm/slack-ips/pkg/config"
	"github.com/dusnm/slack-ips/pkg/repositories/user"
	"github.com/dusnm/slack-ips/pkg/services/messagehandler"
	"github.com/dusnm/slack-ips/pkg/services/requestauth"
	"github.com/rs/zerolog"
)

type (
	Container struct {
		cfg                   *config.Config
		db                    *sql.DB
		userRepo              *user.Repository
		messageHandlerService *messagehandler.Service
		requestAuthService    *requestauth.Service
		logger                zerolog.Logger
	}
)

func New(
	logger zerolog.Logger,
) *Container {
	return &Container{
		logger: logger,
	}
}

func (c *Container) Close() error {
	c.logger.Info().Msg("closing")

	var err error
	if c.userRepo != nil {
		err = errors.Join(err, c.userRepo.Close())
	}

	if c.db != nil {
		err = errors.Join(err, c.db.Close())
	}

	return err
}
