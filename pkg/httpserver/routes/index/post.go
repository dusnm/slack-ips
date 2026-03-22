package index

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/dusnm/slack-ips/pkg/dto/commandresponse"
	"github.com/dusnm/slack-ips/pkg/dto/slack"
	"github.com/dusnm/slack-ips/pkg/httpserver"
	"github.com/dusnm/slack-ips/pkg/repositories"
	"github.com/rs/zerolog"
)

func POST(
	ctx context.Context,
	di *container.Container,
	logger zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if r.URL.Path != "/" {
		return httpserver.Err(http.StatusNotFound, w, httpserver.ErrNotFound)
	}

	if r.Body == nil {
		return httpserver.Err(http.StatusBadRequest, w, httpserver.ErrBadRequest)
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return httpserver.Err(http.StatusBadRequest, w, httpserver.ErrBadRequest)
	}

	timestamp, err := strconv.ParseInt(r.Header.Get("X-Slack-Request-Timestamp"), 10, 64)
	if err != nil {
		return httpserver.Err(http.StatusBadRequest, w, httpserver.ErrBadRequest)
	}

	auth := slack.AuthDetails{
		SigningSecret:    di.GetConfig().Slack.SigningSecret,
		Timestamp:        timestamp,
		RequestSignature: r.Header.Get("X-Slack-Signature"),
		RequestBody:      bodyBytes,
	}

	ok, err := di.GetRequestAuthService().Verify(auth)
	if err != nil {
		return httpserver.Err(http.StatusInternalServerError, w, httpserver.ErrInternalError)
	}

	if !ok {
		return httpserver.Err(http.StatusForbidden, w, httpserver.ErrForbidden)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err = r.ParseForm(); err != nil {
		return httpserver.Err(http.StatusBadRequest, w, httpserver.ErrBadRequest)
	}

	msg := slack.NewMessage(r.Form)
	response, err := di.GetMessageHandlerService().HandleMessage(ctx, msg)
	if err != nil {
		if !errors.Is(err, repositories.ErrNotFound) {
			logger.Error().Err(err).Send()
			return httpserver.JSON(http.StatusOK, w, commandresponse.Message{
				ResponseType: "ephemeral",
				Blocks: []any{
					commandresponse.Section{
						Type: "section",
						Text: commandresponse.Text{
							Type: "plain_text",
							Text: err.Error(),
						},
					},
				},
			})
		}

		return httpserver.JSON(http.StatusOK, w, commandresponse.Message{
			ResponseType: "ephemeral",
			Blocks: []any{
				commandresponse.Section{
					Type: "section",
					Text: commandresponse.Text{
						Type: "plain_text",
						Text: "Hmm, I can't seem to find your bank account details. Did you forget to save them with /ips init ?",
					},
				},
			},
		})
	}

	return httpserver.JSON(http.StatusOK, w, response)
}
