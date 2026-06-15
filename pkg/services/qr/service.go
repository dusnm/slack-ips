package qr

import (
	"bytes"
	"image/png"
	"io"

	"github.com/dusnm/slack-ips/pkg/models"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type (
	Service struct{}

	nopCloser struct {
		io.Writer
	}
)

func (nopCloser) Close() error { return nil }

func New() *Service {
	return &Service{}
}

func (s *Service) Generate(user models.User, data string) (models.QR, error) {
	qrc, err := qrcode.New(data)
	if err != nil {
		return models.QR{}, err
	}

	options := []standard.ImageOption{
		standard.WithFgColorRGBHex(user.Settings.GetQRFGColor()),
		standard.WithBgColorRGBHex(user.Settings.GetQRBGColor()),
		standard.WithQRWidth(19),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
	}

	if user.Settings.GetQRShape() == "circle" {
		options = append(options, standard.WithCircleShape())
	}

	if user.Settings.QRShowLogo && len(user.Settings.QRLogo) > 0 {
		buff := bytes.NewReader(user.Settings.QRLogo)
		// Always stored as a PNG
		logo, err := png.Decode(buff)
		if err != nil {
			return models.QR{}, err
		}

		options = append(options, standard.WithLogoImage(logo))
	}

	b := bytes.NewBuffer(nil)
	w := standard.NewWithWriter(
		nopCloser{b},
		options...,
	)

	if err = qrc.Save(w); err != nil {
		return models.QR{}, err
	}

	return models.NewQR(b.Bytes()), nil
}
