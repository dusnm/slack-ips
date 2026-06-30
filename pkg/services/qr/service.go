package qr

import (
	"bytes"
	"image/png"
	"io"

	"github.com/dusnm/slack-ips/pkg/imgutil"
	"github.com/dusnm/slack-ips/pkg/models"
	"github.com/dusnm/slack-ips/pkg/services/qrcaption"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"github.com/yeqown/go-qrcode/writer/standard/shapes"
)

type (
	Service struct {
		qrCaptionService *qrcaption.Service
	}

	nopCloser struct {
		io.Writer
	}
)

const (
	ShapeSquare = "square"
	ShapeCircle = "circle"
	ShapeLiquid = "liquid"
)

func (nopCloser) Close() error { return nil }

func New(qrCaptionService *qrcaption.Service) *Service {
	return &Service{
		qrCaptionService: qrCaptionService,
	}
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
	case ShapeSquare:
		shape = shapes.Assemble(shapes.SquareFinder(), shapes.SquareBlocks(1))
	case ShapeCircle:
		shape = shapes.Assemble(shapes.RoundedFinder(), shapes.CircleBlocks(0.8))
	case ShapeLiquid:
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

	if caption := user.Settings.GetQRCaption(); caption != "" {
		img, err := png.Decode(bytes.NewReader(b.Bytes()))
		if err != nil {
			return models.QR{}, err
		}

		bgColor, err := imgutil.HexToRGBA(user.Settings.GetQRBGColor())
		if err != nil {
			return models.QR{}, err
		}

		fontColor, err := imgutil.HexToRGBA(user.Settings.GetQRFGColor())
		if err != nil {
			return models.QR{}, err
		}

		img, err = s.qrCaptionService.Do(
			img,
			qrcaption.Style{
				BGColor:   bgColor,
				FontColor: fontColor,
			},
			caption,
		)

		if err != nil {
			return models.QR{}, err
		}

		b.Reset()
		if err = png.Encode(b, img); err != nil {
			return models.QR{}, err
		}
	}

	return models.NewQR(b.Bytes()), nil
}
