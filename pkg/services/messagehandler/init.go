package messagehandler

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"github.com/dusnm/slack-ips/pkg/dto/command"
	"github.com/dusnm/slack-ips/pkg/dto/commandresponse"
	"github.com/dusnm/slack-ips/pkg/dto/slack"
	"github.com/dusnm/slack-ips/pkg/repositories"
)

func (s *Service) handleInitMessage(
	ctx context.Context,
	msg slack.Message,
) (commandresponse.Message, error) {
	payload, err := constructInitPayload(msg)
	if err != nil {
		return commandresponse.Message{}, err
	}

	// Check if the user is already there.
	// For some reason I find this more elegant than allowing
	// the unique constraint to be triggered on the database level.
	_, err = s.userRepo.FindByID(ctx, payload.UserID)
	if err == nil {
		// Found user, do nothing
		return constructSuccessfulInitResponse(), nil
	}

	// Some other error happened, aborting
	if !errors.Is(err, repositories.ErrNotFound) {
		return commandresponse.Message{}, err
	}

	if err = s.userRepo.Insert(ctx, payload); err != nil {
		return commandresponse.Message{}, err
	}

	return constructSuccessfulInitResponse(), nil
}

func constructInitPayload(msg slack.Message) (command.Init, error) {
	if len(msg.Text) == 0 {
		return command.Init{}, ErrInvalidArguments
	}

	text, _ := strings.CutPrefix(msg.Text, "init")
	text = strings.TrimSpace(text)

	// Is using a csv reader here overkill? Possibly.
	records := make([][]string, 0, 1)
	reader := csv.NewReader(strings.NewReader(text))

	for {
		record, err := reader.Read()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return command.Init{}, err
			}

			break
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		return command.Init{}, ErrInvalidArguments
	}

	record := records[0]
	if len(record) != 3 {
		return command.Init{}, ErrInvalidArguments
	}

	for i, r := range record {
		record[i] = strings.TrimSpace(r)
	}

	init := command.Init{
		Name:              record[0],
		BankAccountNumber: record[1],
		City:              record[2],
		UserID:            msg.UserID,
		UserName:          msg.UserName,
	}

	if err := init.Validate(); err != nil {
		return command.Init{}, err
	}

	return init.Format(), nil
}

func constructSuccessfulInitResponse() commandresponse.Message {
	return commandresponse.Message{
		ResponseType: "ephemeral",
		Blocks: []any{
			commandresponse.Text{
				Type: "markdown",
				Text: "# Success ✅\nThe details of your bank account have been saved.",
			},
		},
	}
}
