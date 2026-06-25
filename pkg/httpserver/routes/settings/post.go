package settings

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dusnm/slack-ips/pkg/constants"
	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/dusnm/slack-ips/pkg/dto/command"
	"github.com/dusnm/slack-ips/pkg/httpserver"
	"github.com/dusnm/slack-ips/pkg/repositories"
	"github.com/dusnm/slack-ips/pkg/types"
	"github.com/rs/zerolog"
)

func POST(
	_ context.Context,
	di *container.Container,
	logger zerolog.Logger,
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

	if err = r.ParseMultipartForm(10240); err != nil {
		logger.Debug().Err(err).Send()
		return httpserver.Err(http.StatusBadRequest, w, httpserver.ErrBadRequest)
	}

	logoFile, logoFileHeader, err := r.FormFile("logo")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			logger.Debug().Err(err).Send()
			return httpserver.Err(http.StatusInternalServerError, w, httpserver.ErrInternalError)
		}
	}

	var logoBytes []byte
	if logoFile != nil {
		sizeLimit := math.Round(di.GetConfig().App.UploadedFileSizeLimit * constants.MiB)
		if logoFileHeader.Size > int64(sizeLimit) {
			return httpserver.Err(http.StatusUnprocessableEntity, w, fmt.Errorf("%w: logo file too large", httpserver.ErrUnprocessable))
		}

		logoBytes, err = io.ReadAll(logoFile)
		if err != nil {
			logger.Debug().Err(err).Send()
			return httpserver.Err(http.StatusInternalServerError, w, httpserver.ErrInternalError)
		}
	}

	payload := command.Settings{
		Init: command.Init{
			Name:              r.Form.Get("name"),
			BankAccountNumber: r.Form.Get("bank_account_number"),
			City:              r.Form.Get("city"),
		},
		QRFGColor:  r.Form.Get("fg_color"),
		QRBGColor:  r.Form.Get("bg_color"),
		QRShape:    r.Form.Get("shape"),
		QRLogo:     logoBytes,
		QRShowLogo: r.Form.Get("show_logo") == "on",
	}

	if err = payload.Validate(); err != nil {
		return httpserver.Err(http.StatusUnprocessableEntity, w, err)
	}

	payload = payload.Format()

	if err = di.GetUserRepository().Update(r.Context(), user.ID, payload.Init); err != nil {
		logger.Debug().Err(err).Send()
		return httpserver.Err(http.StatusInternalServerError, w, httpserver.ErrInternalError)
	}

	user.Name = payload.Name
	user.BankAccountNumber = payload.BankAccountNumber
	user.City = payload.City

	settings, err := di.GetSettingsRepository().UpsertByUserID(r.Context(), user.ID, payload)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return httpserver.Err(http.StatusNotFound, w, httpserver.ErrNotFound)
		}

		logger.Debug().Err(err).Send()
		return httpserver.Err(http.StatusInternalServerError, w, httpserver.ErrInternalError)
	}

	user.Settings = settings
	ipsString := strings.Replace(user.IPSString, "RSD0,00", "RSD1000,00", -1)
	qrModel, err := di.GetQRService().Generate(user, ipsString)
	if err != nil {
		logger.Debug().Err(err).Send()
		return httpserver.Err(http.StatusInternalServerError, w, httpserver.ErrInternalError)
	}

	requestForSigning, _ := http.NewRequest(http.MethodPost, "/settings", nil)
	newExpiration := strconv.FormatInt(time.Now().UTC().Add(10*time.Minute).Unix(), 10)
	values = url.Values{}
	values.Add("userId", userId)
	values.Add("expiresAt", newExpiration)
	requestForSigning.URL.RawQuery = values.Encode()
	newSignature, err := urlSignService.Sign(requestForSigning)
	if err != nil {
		logger.Debug().Err(err).Send()
		return httpserver.Err(http.StatusInternalServerError, w, err)
	}

	return httpserver.HTML(
		http.StatusOK,
		w,
		di.GetTemplateService(),
		types.PageSettings,
		Data{
			User:      user,
			QR:        qrModel,
			Signature: hex.EncodeToString(newSignature),
			ExpiresAt: newExpiration,
		},
	)
}
