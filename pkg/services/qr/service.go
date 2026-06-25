package qr

import (
	"bytes"
	"image/png"
	"io"

	"github.com/dusnm/slack-ips/pkg/models"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"github.com/yeqown/go-qrcode/writer/standard/shapes"
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

	var shape standard.IShape
	switch user.Settings.GetQRShape() {
	case "square":
		shape = shapes.Assemble(shapes.SquareFinder(), shapes.SquareBlocks(1))
	case "circle":
		shape = shapes.Assemble(shapes.RoundedFinder(), shapes.CircleBlocks(0.8))
	case "liquid":
		shape = shapes.Assemble(shapes.RoundedFinder(), shapes.LiquidBlock())
	}

	options = append(options, standard.WithCustomShape(shape))

	if user.Settings.ShouldShowLogo() && len(user.Settings.QRLogo) > 0 {
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
