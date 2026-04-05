package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// BlurConfig defines blur shader parameters
type BlurConfig struct {
	Radius              int     // Blur radius in pixels
	OverlayTransparency float64 // White overlay transparency (0.0-1.0)
}

// DefaultBlurConfig returns default blur parameters
func DefaultBlurConfig() BlurConfig {
	return BlurConfig{
		Radius:              5,
		OverlayTransparency: 0.20, // 20% white overlay
	}
}

// ValidateBlurConfig validates and clamps blur configuration parameters
func ValidateBlurConfig(config BlurConfig) BlurConfig {
	validated := config

	// Clamp radius to 0-20 range
	if validated.Radius < 0 {
		validated.Radius = 0
	}
	if validated.Radius > 20 {
		validated.Radius = 20
	}

	// Clamp transparency to 0.0-1.0 range
	if validated.OverlayTransparency < 0.0 {
		validated.OverlayTransparency = 0.0
	}
	if validated.OverlayTransparency > 1.0 {
		validated.OverlayTransparency = 1.0
	}

	return validated
}

// ApplyBlur applies blur effect and white overlay to an image
func ApplyBlur(src *ebiten.Image, config BlurConfig) *ebiten.Image {
	// Validate configuration
	config = ValidateBlurConfig(config)

	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create output image
	blurred := ebiten.NewImage(width, height)

	// Apply box blur approximation
	temp := applyBoxBlur(src, config.Radius)

	// Draw blurred image
	blurred.DrawImage(temp, nil)

	// Apply white overlay
	overlay := ebiten.NewImage(width, height)
	overlay.Fill(color.RGBA{255, 255, 255, uint8(config.OverlayTransparency * 255)})

	opts := &ebiten.DrawImageOptions{}
	opts.Blend = ebiten.BlendSourceOver
	blurred.DrawImage(overlay, opts)

	return blurred
}

// applyBoxBlur applies a simple box blur (approximation of Gaussian)
func applyBoxBlur(src *ebiten.Image, radius int) *ebiten.Image {
	if radius <= 0 {
		return src
	}

	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create intermediate images for horizontal and vertical passes
	horizontal := ebiten.NewImage(width, height)
	vertical := ebiten.NewImage(width, height)

	// Read source pixels
	srcPixels := make([]byte, width*height*4)
	src.ReadPixels(srcPixels)

	// Horizontal pass
	hPixels := make([]byte, width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b, a int
			count := 0

			for dx := -radius; dx <= radius; dx++ {
				nx := x + dx
				if nx >= 0 && nx < width {
					idx := (y*width + nx) * 4
					r += int(srcPixels[idx])
					g += int(srcPixels[idx+1])
					b += int(srcPixels[idx+2])
					a += int(srcPixels[idx+3])
					count++
				}
			}

			idx := (y*width + x) * 4
			hPixels[idx] = uint8(r / count)
			hPixels[idx+1] = uint8(g / count)
			hPixels[idx+2] = uint8(b / count)
			hPixels[idx+3] = uint8(a / count)
		}
	}

	horizontal.WritePixels(hPixels)

	// Vertical pass
	vPixels := make([]byte, width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var r, g, b, a int
			count := 0

			for dy := -radius; dy <= radius; dy++ {
				ny := y + dy
				if ny >= 0 && ny < height {
					idx := (ny*width + x) * 4
					r += int(hPixels[idx])
					g += int(hPixels[idx+1])
					b += int(hPixels[idx+2])
					a += int(hPixels[idx+3])
					count++
				}
			}

			idx := (y*width + x) * 4
			vPixels[idx] = uint8(r / count)
			vPixels[idx+1] = uint8(g / count)
			vPixels[idx+2] = uint8(b / count)
			vPixels[idx+3] = uint8(a / count)
		}
	}

	vertical.WritePixels(vPixels)
	return vertical
}
