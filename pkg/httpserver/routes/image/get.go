package image

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	amount := r.URL.Query().Get("amount")

	if len(userID) == 0 {
		return httpserver.Err(http.StatusUnprocessableEntity, w, httpserver.ErrUnprocessable)
	}

	user, err := di.GetUserRepository().FindByID(ctx, userID)
	if err != nil {
		if !errors.Is(err, repositories.ErrNotFound) {
			return httpserver.Err(http.StatusInternalServerError, w, err)
		}

		return httpserver.Err(http.StatusNotFound, w, httpserver.ErrNotFound)
	}

	ipsString := user.IPSString
	if len(amount) > 0 {
		amountFloat, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return httpserver.Err(http.StatusUnprocessableEntity, w, httpserver.ErrUnprocessable)
		}

		if amountFloat <= 0 {
			return httpserver.Err(http.StatusUnprocessableEntity, w, httpserver.ErrUnprocessable)
		}

		amount = strings.Replace(fmt.Sprintf("RSD%.2f", amountFloat), ".", ",", -1)
		ipsString = strings.Replace(ipsString, "RSD0,00", amount, -1)
	}

	image, err := qrcode.Encode(ipsString, qrcode.Medium, 350)
	if err != nil {
		return httpserver.Err(http.StatusInternalServerError, w, err)
	}

	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(image)
	return err
}
