package settings

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/dusnm/slack-ips/pkg/httpserver"
	"github.com/dusnm/slack-ips/pkg/models"
	"github.com/dusnm/slack-ips/pkg/repositories"
	"github.com/dusnm/slack-ips/pkg/types"
	"github.com/rs/zerolog"
)

type (
	Data struct {
		User      models.User
		QR        models.QR
		Signature string
		ExpiresAt string
	}
)

func GET(
	_ context.Context,
	di *container.Container,
	_ zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if r.URL.Path != "/settings" {
		return httpserver.Err(http.StatusNotFound, w, httpserver.ErrNotFound)
	}

	providedSignatureStr := r.URL.Query().Get("sig")

	if len(providedSignatureStr) == 0 {
		return httpserver.Err(http.StatusUnprocessableEntity, w, httpserver.ErrUnprocessable)
	}

	providedSignature, err := hex.DecodeString(providedSignatureStr)
	if err != nil {
		return httpserver.Err(http.StatusUnprocessableEntity, w, fmt.Errorf("%w: malformed signature", httpserver.ErrUnprocessable))
	}

	values := r.URL.Query()
	values.Del("sig")
	r.URL.RawQuery = values.Encode()

	urlSignService := di.GetURLSignService()
	if err = urlSignService.Verify(r, providedSignature); err != nil {
		return httpserver.Err(http.StatusForbidden, w, fmt.Errorf("%w: %w", httpserver.ErrForbidden, err))
	}

	userId := r.URL.Query().Get("userId")
	expiresAtStr := r.URL.Query().Get("expiresAt")
	expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
	if err != nil {
		return httpserver.Err(http.StatusUnprocessableEntity, w, httpserver.ErrUnprocessable)
	}

	if time.Now().UTC().After(time.Unix(expiresAt, 0).UTC()) {
		return httpserver.Err(http.StatusUnprocessableEntity, w, fmt.Errorf("%w: link expired", httpserver.ErrUnprocessable))
	}

	user, err := di.GetUserRepository().FindByID(r.Context(), userId)
	if err != nil {
		if !errors.Is(err, repositories.ErrNotFound) {
			return httpserver.Err(http.StatusInternalServerError, w, err)
		}

		return httpserver.Err(http.StatusNotFound, w, httpserver.ErrNotFound)
	}

	ipsString := strings.Replace(user.IPSString, "RSD0,00", "RSD1000,00", -1)
	qr, err := di.GetQRService().Generate(user, ipsString)
	if err != nil {
		return httpserver.Err(http.StatusInternalServerError, w, err)
	}

	requestForSigning, _ := http.NewRequest(http.MethodPost, "/settings", nil)
	newExpiration := strconv.FormatInt(time.Now().UTC().Add(10*time.Minute).Unix(), 10)
	values = url.Values{}
	values.Add("userId", userId)
	values.Add("expiresAt", newExpiration)
	requestForSigning.URL.RawQuery = values.Encode()
	newSignature, err := urlSignService.Sign(requestForSigning)
	if err != nil {
		return httpserver.Err(http.StatusInternalServerError, w, err)
	}

	return httpserver.HTML(
		http.StatusOK,
		w,
		di.GetTemplateService(),
		types.PageSettings,
		Data{
			User:      user,
			QR:        qr,
			Signature: hex.EncodeToString(newSignature),
			ExpiresAt: newExpiration,
		},
	)
}
