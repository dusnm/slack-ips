package messagehandler

import (
	"context"
	"fmt"

	"github.com/dusnm/slack-ips/pkg/config"
	"github.com/dusnm/slack-ips/pkg/dto/commandresponse"
	"github.com/dusnm/slack-ips/pkg/dto/slack"
	"github.com/dusnm/slack-ips/pkg/models"
)

func (s *Service) handleSendMessage(
	ctx context.Context,
	msg slack.Message,
) (commandresponse.Message, error) {
	user, err := s.userRepo.FindByID(ctx, msg.UserID)
	if err != nil {
		return commandresponse.Message{}, err
	}

	return constructSuccessfulSendResponse(s.cfg, user), nil
}

func constructSuccessfulSendResponse(cfg config.App, user models.User) commandresponse.Message {
	return commandresponse.Message{
		ResponseType: "in_channel",
		Blocks: []any{
			commandresponse.Section{
				Type: "section",
				Text: commandresponse.Text{
					Type: "mrkdwn",
					Text: fmt.Sprintf(
						"*%s*\n* IBAN: %s\n*Place: %s",
						user.Name,
						user.BankAccountNumber,
						user.City,
					),
				},
			},
			commandresponse.Image{
				Type:     "image",
				ImageURL: user.QRCodeURL(cfg.Domain, cfg.Secure),
				AltText:  "IPS QR Code",
			},
		},
	}
}
