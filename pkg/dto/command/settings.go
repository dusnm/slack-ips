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

	"github.com/dusnm/slack-ips/pkg/imgutil"
)

var (
	ErrInvalidHexColor     = errors.New("invalid hex color")
	ErrInvalidValues       = errors.New("invalid values")
	ErrCaptionTooLong      = errors.New("caption too long")
	ErrDecodingLogo        = errors.New("decoding logo failed")
	ErrInvalidLogoEncoding = errors.New("invalid logo encoding")
)

type (
	Settings struct {
		Init
		QRFGColor  string
		QRBGColor  string
		QRShape    string
		QRCaption  string
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

	if s.QRCaption != "" {
		caption := strings.TrimSpace(s.QRCaption)
		if len(caption) > 50 {
			return ErrCaptionTooLong
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
		img = imgutil.ResizeImage(img, 140, 140, imgutil.ResizeFit)

		// A disaster waiting to happen
		_ = png.Encode(w, img)

		logo = w.Bytes()
	}

	return Settings{
		Init:       s.FormatInit(),
		QRFGColor:  strings.TrimSpace(s.QRFGColor),
		QRBGColor:  strings.TrimSpace(s.QRBGColor),
		QRShape:    strings.TrimSpace(s.QRShape),
		QRCaption:  strings.TrimSpace(s.QRCaption),
		QRLogo:     logo,
		QRShowLogo: s.QRShowLogo,
	}
}
