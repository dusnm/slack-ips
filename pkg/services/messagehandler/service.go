package messagehandler

import (
	"context"
	"errors"
	"strings"

	"github.com/dusnm/slack-ips/pkg/config"
	"github.com/dusnm/slack-ips/pkg/dto/commandresponse"
	"github.com/dusnm/slack-ips/pkg/dto/slack"
	"github.com/dusnm/slack-ips/pkg/repositories/user"
	"github.com/rs/zerolog"
)

const (
	initCommand   = "init"
	sendCommand   = "send"
	deleteCommand = "delete"
	helpCommand   = "help"
)

var (
	ErrInvalidArguments = errors.New("invalid arguments for command")
	ErrUnknownCommand   = errors.New("unknown command")
)

type (
	Service struct {
		cfg      config.App
		userRepo *user.Repository
		logger   zerolog.Logger
	}
)

func New(
	cfg config.App,
	userRepo *user.Repository,
	logger zerolog.Logger,
) *Service {
	return &Service{
		cfg:      cfg,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *Service) HandleMessage(ctx context.Context, msg slack.Message) (commandresponse.Message, error) {
	text := strings.TrimSpace(msg.Text)

	switch {
	case strings.HasPrefix(text, initCommand):
		s.logger.
			Debug().
			Str("command", initCommand).
			Str("user", msg.UserName).
			Str("input", text).
			Send()
		return s.handleInitMessage(ctx, msg)
	case strings.HasPrefix(text, sendCommand):
		s.logger.
			Debug().
			Str("command", sendCommand).
			Str("user", msg.UserName).
			Str("input", text).
			Send()
		return s.handleSendMessage(ctx, msg)
	case strings.HasPrefix(text, deleteCommand):
		s.logger.
			Debug().
			Str("command", deleteCommand).
			Str("user", msg.UserName).
			Str("input", text).Send()
		return s.handleDeleteMessage(ctx, msg)
	case strings.HasPrefix(text, helpCommand):
		s.logger.
			Debug().
			Str("command", helpCommand).
			Str("user", msg.UserName).
			Str("input", text).
			Send()
		return s.handleHelpMessage(ctx, msg)
	default:
		return commandresponse.Message{}, ErrUnknownCommand
	}
}
