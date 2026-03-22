package urlsign

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"

	"github.com/dusnm/slack-ips/pkg/config"
	"github.com/rs/zerolog"
)

var (
	ErrSignatureMismatch = errors.New("signature mismatch")
)

type (
	Service struct {
		cfg    config.App
		logger zerolog.Logger
	}
)

func New(
	cfg config.App,
	logger zerolog.Logger,
) *Service {
	return &Service{
		cfg:    cfg,
		logger: logger,
	}
}

// Sign
//
// Incredibly rudimentary function for HMAC-SHA256 URL
// signing based on url.Values and a provided signing key.
func (s *Service) Sign(
	values url.Values,
) ([]byte, error) {
	if len(values) == 0 {
		s.logger.Fatal().Msg("no values provided")
	}

	buff := bytes.Buffer{}
	buff.WriteString(values.Encode())

	key, err := hex.DecodeString(s.cfg.SigningSecret)
	if err != nil {
		return nil, err
	}

	hasher := hmac.New(sha256.New, key)
	_, err = hasher.Write(buff.Bytes())
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

func (s *Service) Verify(values url.Values, providedSignature []byte) error {
	computedSignature, err := s.Sign(values)
	if err != nil {
		return err
	}

	if !hmac.Equal(computedSignature, providedSignature) {
		return ErrSignatureMismatch
	}

	return nil
}
