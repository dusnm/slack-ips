package imgutil

import (
	"encoding/hex"
	"errors"
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/disintegration/imaging"
)

type ResizeMode uint8

const (
	ResizeFill    ResizeMode = iota // crop to fit
	ResizeStretch                   // stretch to fit
	ResizeFit                       // preserve aspect ratio
)

func ResizeImage(
	img image.Image,
	width int,
	height int,
	resizeMode ResizeMode,
) image.Image {
	bounds := img.Bounds()
	if width == bounds.Dx() && height == bounds.Dy() {
		// avoid unnecessary processing
		return img
	}

	switch resizeMode {
	case ResizeFill:
		return imaging.CropCenter(
			imaging.Fit(img, width, height, imaging.Lanczos),
			width,
			height,
		)
	case ResizeStretch:
		return imaging.Resize(img, width, height, imaging.Lanczos)
	case ResizeFit:
		scale := min(
			float64(width)/float64(bounds.Dx()),
			float64(height)/float64(bounds.Dy()),
		)

		scaledWidth := int(math.Round(float64(width) * scale))
		scaledHeight := int(math.Round(float64(height) * scale))

		return imaging.Resize(img, scaledWidth, scaledHeight, imaging.Lanczos)
	default:
		panic("unsupported resize mode")
	}
}

func HexToRGBA(s string) (color.Color, error) {
	s, _ = strings.CutPrefix(strings.TrimSpace(s), "#")
	if len(s) == 3 {
		// Expand short representation into long
		s = string([]byte{
			s[0], s[0],
			s[1], s[1],
			s[2], s[2],
		})
	}

	if len(s) != 6 {
		return color.RGBA{}, errors.New("invalid hex color")
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return color.RGBA{}, err
	}

	// Each hex string is three bytes,
	// one byte for each color value
	// in range [0, 255]
	return color.RGBA{
		R: b[0], // RED value
		G: b[1], // GREEN value
		B: b[2], // BLUE value
		A: 255,  // Fully opaque
	}, nil
}
