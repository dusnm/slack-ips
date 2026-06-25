package command

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/png"
	"slices"
	"strings"

	"github.com/dusnm/slack-ips/pkg/utils"
)

var (
	ErrInvalidHexColor     = errors.New("invalid hex color")
	ErrInvalidValues       = errors.New("invalid values")
	ErrDecodingLogo        = errors.New("decoding logo failed")
	ErrInvalidLogoEncoding = errors.New("invalid logo encoding")
)

type (
	Settings struct {
		Init
		QRFGColor  string
		QRBGColor  string
		QRShape    string
		QRLogo     []byte
		QRShowLogo bool
	}
)

func (s Settings) Validate() error {
	if err := s.ValidateInit(); err != nil {
		return err
	}

	if s.QRFGColor != "" {
		color, _ := strings.CutPrefix(s.QRFGColor, "#")
		_, err := hex.DecodeString(color)
		if err != nil {
			return fmt.Errorf("%w for qr foreground", ErrInvalidHexColor)
		}
	}

	if s.QRBGColor != "" {
		color, _ := strings.CutPrefix(s.QRBGColor, "#")
		_, err := hex.DecodeString(color)
		if err != nil {
			return fmt.Errorf("%w for qr background", ErrInvalidHexColor)
		}
	}

	if s.QRShape != "" {
		allowedShapes := []string{
			"square",
			"circle",
			"liquid",
		}

		shape := strings.ToLower(strings.TrimSpace(s.QRShape))
		if !slices.Contains(allowedShapes, shape) {
			return fmt.Errorf("%w for qr shape, only square, circle and liquid are allowed", ErrInvalidValues)
		}
	}

	if len(s.QRLogo) > 0 {
		_, encoding, err := image.DecodeConfig(bytes.NewReader(s.QRLogo))
		if err != nil {
			return ErrDecodingLogo
		}

		if encoding != "png" && encoding != "jpg" && encoding != "jpeg" {
			return fmt.Errorf("%w only jpg/jpeg and png are allowed", ErrInvalidLogoEncoding)
		}
	}

	return nil
}

func (s Settings) Format() Settings {
	var logo []byte
	if len(s.QRLogo) > 0 {
		img, _, _ := image.Decode(bytes.NewReader(s.QRLogo))

		w := bytes.NewBuffer(nil)
		img = utils.ResizeImage(img, 140, 140)

		// A disaster waiting to happen
		_ = png.Encode(w, img)

		logo = w.Bytes()
	}

	return Settings{
		Init:       s.FormatInit(),
		QRFGColor:  s.QRFGColor,
		QRBGColor:  s.QRBGColor,
		QRShape:    s.QRShape,
		QRLogo:     logo,
		QRShowLogo: s.QRShowLogo,
	}
}
