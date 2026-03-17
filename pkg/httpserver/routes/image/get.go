package image

import (
	"context"
	"errors"
	"net/http"

	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/dusnm/slack-ips/pkg/httpserver"
	"github.com/dusnm/slack-ips/pkg/repositories"
	"github.com/rs/zerolog"
	"github.com/skip2/go-qrcode"
)

func GET(
	ctx context.Context,
	di *container.Container,
	_ zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if r.URL.Path != "/image" {
		return httpserver.Err(http.StatusNotFound, w, httpserver.ErrNotFound)
	}

	userID := r.URL.Query().Get("userId")
	user, err := di.GetUserRepository().FindByID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repositories.ErrNotFound) {
			return httpserver.Err(http.StatusInternalServerError, w, err)
		}

		return httpserver.Err(http.StatusNotFound, w, httpserver.ErrNotFound)
	}

	image, err := qrcode.Encode(user.IPSString, qrcode.Medium, 350)
	if err != nil {
		return httpserver.Err(http.StatusInternalServerError, w, err)
	}

	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(image)
	return err
}
