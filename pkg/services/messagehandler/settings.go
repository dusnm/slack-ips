package messagehandler

import (
	"context"
	"encoding/hex"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dusnm/slack-ips/pkg/config"
	"github.com/dusnm/slack-ips/pkg/dto/commandresponse"
	"github.com/dusnm/slack-ips/pkg/dto/slack"
)

func (s *Service) handleSettingsMessage(
	ctx context.Context,
	msg slack.Message,
) (commandresponse.Message, error) {
	expirationDatetime := time.
		Now().
		UTC().
		Add(10 * time.Minute).
		Unix()

	query := url.Values{}
	query.Add("userId", msg.UserID)
	query.Add("expiresAt", strconv.FormatInt(expirationDatetime, 10))

	requestForSigning, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/settings",
		nil,
	)

	if err != nil {
		return commandresponse.Message{}, err
	}

	requestForSigning.URL.RawQuery = query.Encode()

	signature, err := s.urlSignService.Sign(requestForSigning)
	if err != nil {
		return commandresponse.Message{}, err
	}

	query.Add("sig", hex.EncodeToString(signature))
	uri := settingsURI(s.cfg, query)

	return constructSuccessfulSettingsResponse(uri), nil
}

func settingsURI(cfg config.App, query url.Values) *url.URL {
	uri := new(url.URL)
	if cfg.Secure {
		uri.Scheme = "https"
	} else {
		uri.Scheme = "http"
	}

	if cfg.Port > 0 && !cfg.BehindProxy {
		uri.Host = net.JoinHostPort(cfg.Domain, strconv.FormatUint(uint64(cfg.Port), 10))
	} else {
		uri.Host = cfg.Domain
	}

	uri.Path = "/settings"
	uri.RawQuery = query.Encode()

	return uri
}

func constructSuccessfulSettingsResponse(
	settingsURI *url.URL,
) commandresponse.Message {
	return commandresponse.Message{
		ResponseType: "ephemeral",
		Blocks: []any{
			commandresponse.Section{
				Type: "actions",
				Elements: []any{
					commandresponse.Button{
						Type: "button",
						Text: commandresponse.Text{
							Type: "plain_text",
							Text: "Open settings",
						},
						URL: settingsURI.String(),
					},
				},
			},
		},
	}
}
