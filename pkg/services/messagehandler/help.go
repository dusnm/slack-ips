package messagehandler

import (
	"context"

	"github.com/dusnm/slack-ips/pkg/dto/commandresponse"
	"github.com/dusnm/slack-ips/pkg/dto/slack"
)

func (s *Service) handleHelpMessage(ctx context.Context, msg slack.Message) (commandresponse.Message, error) {
	return commandresponse.Message{
		ResponseType: "ephemeral",
		Blocks: []any{
			commandresponse.Section{
				Type: "header",
				Text: commandresponse.Text{
					Type: "plain_text",
					Text: "IPS QR - Help",
				},
			},
			commandresponse.Section{
				Type: "section",
				Text: commandresponse.Text{
					Type: "plain_text",
					Text: "Copyright © 2026 Dušan Mitrović <dusan@dusanmitrovic.rs>\nEasily share your bank account details with others.",
				},
			},
			commandresponse.Section{
				Type: "section",
				Text: commandresponse.Text{
					Type: "mrkdwn",
					Text: "Use `/ips <command>` to interact with the IPS QR app.\n\n*Available commands:*",
				},
			},
			commandresponse.Section{
				Type: "section",
				Fields: []any{
					commandresponse.Text{
						Type: "mrkdwn",
						Text: "*init*\nInitialize your profile and bank details.",
					},
					commandresponse.Text{
						Type: "mrkdwn",
						Text: "*send*\nGenerate a QR code for receiving a payment with optional amount.",
					},
					commandresponse.Text{
						Type: "mrkdwn",
						Text: "*delete*\nRemove your saved profile and data.",
					},
					commandresponse.Text{
						Type: "mrkdwn",
						Text: "*help*\nShow this help message.",
					},
				},
			},
			commandresponse.Section{
				Type: "section",
				Text: commandresponse.Text{
					Type: "mrkdwn",
					Text: "*Examples:*\n• `/ips init Malina Vojvodić,260-0056010016113-79,Beograd`\n• `/ips send 1500`\n• `/ips delete`\n• `/ips help`",
				},
			},
		},
	}, nil
}
