package qrcaption

import (
	"image"
	"image/color"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type (
	Style struct {
		BGColor   color.Color
		FontColor color.Color
	}

	Service struct {
		font *truetype.Font
	}
)

func New(fontBytes []byte) (*Service, error) {
	ft, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return &Service{
		font: ft,
	}, nil
}

func (s *Service) Do(
	img image.Image,
	style Style,
	caption string,
) (image.Image, error) {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	canvasPadding := int(float64(max(w, h)) * 0.10)
	newW := w + canvasPadding*2
	newH := h + canvasPadding*2

	// The new canvas is expanded by 10% of the original width and height on every side.
	canvas := image.NewRGBA(image.Rect(0, 0, newW, newH))

	// Fill canvas with the background color.
	draw.Draw(
		canvas,
		canvas.Bounds(),
		image.NewUniform(style.BGColor),
		image.Point{},
		draw.Src,
	)

	offset := image.Pt(
		(newW-w)/2,
		(newH-h)/2,
	)

	// Draw the original image in the center of the new canvas.
	draw.Draw(
		canvas,
		image.Rect(offset.X, offset.Y, offset.X+w, offset.Y+h),
		img,
		img.Bounds().Min,
		draw.Over,
	)

	fontFace := truetype.NewFace(s.font, &truetype.Options{
		Size:    float64(canvasPadding) * 0.6,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	defer fontFace.Close()

	drawer := &font.Drawer{
		Dst:  canvas,
		Src:  image.NewUniform(style.FontColor),
		Face: fontFace,
	}

	captionWidth := drawer.MeasureString(caption).Round()
	metrics := fontFace.Metrics()
	textHeight := (metrics.Ascent + metrics.Descent).Round()
	drawer.Dot = fixed.P(
		(newW-captionWidth)/2,
		// Arbitrary padding of 20 serves
		// as a "fix" for the unknown, but consistent
		// padding added by the QR code library.
		canvasPadding/2+metrics.Ascent.Round()-textHeight/2+20,
	)

	drawer.DrawString(caption)

	return canvas, nil
}
