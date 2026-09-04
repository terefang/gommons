package noise

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Default rendering dimensions and scaling constants
const (
	DefaultImageWidth  = 4096
	DefaultImageHeight = 4096
	DefaultImageScale  = 0.002
)

// MapToColor maps a noise value in the range [-1.0, 1.0] to an 8-bit grayscale color.
func MapToColor(val float64) color.Color {
	if val < -1.0 {
		val = -1.0
	} else if val > 1.0 {
		val = 1.0
	}

	normalized := (val + 1.0) * 0.5
	gray := uint8(normalized * 255.0)

	return color.RGBA{R: gray, G: gray, B: gray, A: 255}
}

// RenderConfig holds parameters required for rendering noise samples into an image.
type RenderConfig struct {
	Width  int
	Height int
	Scale  float64
}

// NewDefaultRenderConfig returns a RenderConfig initialized with default constants.
func NewDefaultRenderConfig() RenderConfig {
	return RenderConfig{
		Width:  DefaultImageWidth,
		Height: DefaultImageHeight,
		Scale:  DefaultImageScale,
	}
}

// RenderNoiseToPNG samples a 2D slice from a noise generator function and saves it as a PNG file.
func RenderNoiseToPNG(filePath string, config RenderConfig, sampler func(x, y float64) float64) error {
	width := config.Width
	if width <= 0 {
		width = DefaultImageWidth
	}

	height := config.Height
	if height <= 0 {
		height = DefaultImageHeight
	}

	scale := config.Scale
	if scale <= 0 {
		scale = DefaultImageScale
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			nx := float64(x) * scale
			ny := float64(y) * scale

			noiseValue := sampler(nx, ny)
			img.Set(x, y, MapToColor(noiseValue))
		}
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("failed to encode PNG %s: %w", filePath, err)
	}

	return nil
}
