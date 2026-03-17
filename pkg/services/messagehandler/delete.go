package messagehandler

import (
	"context"

	"github.com/dusnm/slack-ips/pkg/dto/commandresponse"
	"github.com/dusnm/slack-ips/pkg/dto/slack"
)

func (s *Service) handleDeleteMessage(
	ctx context.Context,
	msg slack.Message,
) (commandresponse.Message, error) {
	if err := s.userRepo.DeleteByID(ctx, msg.UserID); err != nil {
		return commandresponse.Message{}, err
	}

	return constructSuccessfulDeleteResponse(), nil
}

func constructSuccessfulDeleteResponse() commandresponse.Message {
	return commandresponse.Message{
		ResponseType: "ephemeral",
		Blocks: []any{
			commandresponse.Text{
				Type: "markdown",
				Text: "Success ✅\nThe details of your bank account have been deleted.",
			},
		},
	}
}
