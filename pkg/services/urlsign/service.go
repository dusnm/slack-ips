package urlsign

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

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
// signing based on data of the http.Request and a provided signing key.
func (s *Service) Sign(request *http.Request) ([]byte, error) {
	buff := bytes.Buffer{}
	buff.WriteString(request.Method)
	buff.WriteString(request.URL.Path)
	buff.WriteString(request.URL.Query().Encode())

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

// Verify
//
// Verifies the signature in constant time.
func (s *Service) Verify(request *http.Request, providedSignature []byte) error {
	computedSignature, err := s.Sign(request)
	if err != nil {
		return err
	}

	if !hmac.Equal(computedSignature, providedSignature) {
		return ErrSignatureMismatch
	}

	return nil
}
