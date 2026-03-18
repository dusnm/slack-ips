package requestauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/dusnm/slack-ips/pkg/dto/slack"
	"github.com/rs/zerolog"
)

type (
	Service struct {
		logger zerolog.Logger
	}
)

func New(
	logger zerolog.Logger,
) *Service {
	return &Service{
		logger: logger,
	}
}

func (s *Service) Verify(details slack.AuthDetails) (bool, error) {
	// Copy the body bytes buffer, just in case
	bodyBytes := make([]byte, len(details.RequestBody))
	copy(bodyBytes, details.RequestBody)

	if time.Since(time.Unix(details.Timestamp, 0)) >= 5*time.Minute {
		s.logger.
			Warn().
			Bytes("payload", bodyBytes).
			Msg("possible replay attack")

		return false, nil
	}

	timestampStr := strconv.FormatInt(details.Timestamp, 10)
	buff := bytes.Buffer{}
	buff.Grow(len(timestampStr) + len(bodyBytes) + 4) // 4 comes from the length of v0 and 2 ":" characters
	buff.WriteString("v0:")
	buff.WriteString(timestampStr)
	buff.WriteString(":")
	buff.Write(bodyBytes)

	hasher := hmac.New(sha256.New, []byte(details.SigningSecret))
	_, err := hasher.Write(buff.Bytes())
	if err != nil {
		return false, err
	}

	providedSignatureStr, _ := strings.CutPrefix(details.RequestSignature, "v0=")
	providedSignature, err := hex.DecodeString(providedSignatureStr)
	if err != nil {
		return false, err
	}

	computedSignature := hasher.Sum(nil)
	verified := hmac.Equal(providedSignature, computedSignature)
	if !verified {
		s.logger.
			Warn().
			Str("provided_signature", providedSignatureStr).
			Str("computed_signature", hex.EncodeToString(computedSignature)).
			Bytes("payload", bodyBytes).
			Msg("signatures don't match")
	}

	return verified, nil
}
