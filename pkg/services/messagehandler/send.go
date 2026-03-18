package messagehandler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dusnm/slack-ips/pkg/config"
	"github.com/dusnm/slack-ips/pkg/dto/commandresponse"
	"github.com/dusnm/slack-ips/pkg/dto/slack"
	"github.com/dusnm/slack-ips/pkg/models"
)

var (
	ErrInvalidAmount = errors.New("invalid amount, please enter a positive decimal value, example: 99.99")
)

func (s *Service) handleSendMessage(
	ctx context.Context,
	msg slack.Message,
) (commandresponse.Message, error) {
	user, err := s.userRepo.FindByID(ctx, msg.UserID)
	if err != nil {
		return commandresponse.Message{}, err
	}

	text, _ := strings.CutPrefix(msg.Text, "send")
	text = strings.TrimSpace(text)

	// Since amount is the only argument
	// this is fine, but it'll need to be
	// changed in the future if I want to support
	// multiple arguments.
	amount := 0.0
	if len(text) > 0 {
		amount, err = strconv.ParseFloat(text, 64)
		if err != nil {
			return commandresponse.Message{}, ErrInvalidAmount
		}

		if amount <= 0 {
			return commandresponse.Message{}, ErrInvalidAmount
		}
	}

	return constructSuccessfulSendResponse(s.cfg, user, amount), nil
}

func constructSuccessfulSendResponse(cfg config.App, user models.User, amount float64) commandresponse.Message {
	return commandresponse.Message{
		ResponseType: "in_channel",
		Blocks: []any{
			commandresponse.Section{
				Type: "section",
				Text: commandresponse.Text{
					Type: "mrkdwn",
					Text: fmt.Sprintf(
						"Name: *%s*\nIBAN: *%s*\nPlace: *%s*\nAmount: *%.2fRSD*",
						user.Name,
						user.BankAccountNumber,
						user.City,
						amount,
					),
				},
			},
			commandresponse.Image{
				Type:     "image",
				ImageURL: user.QRCodeURL(cfg.Domain, cfg.Secure, amount),
				AltText:  "IPS QR Code",
			},
		},
	}
}
