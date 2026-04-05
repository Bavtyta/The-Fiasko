package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/common"
	"TheFiaskoTest/internal/render"
)

type SkyLayer struct {
	offsetZ        float64
	parallaxFactor float64
	width          int
	height         int
	color          color.Color
}

// Getters
func (s *SkyLayer) OffsetZ() float64 {
	return s.offsetZ
}

func (s *SkyLayer) ParallaxFactor() float64 {
	return s.parallaxFactor
}

func (s *SkyLayer) Width() int {
	return s.width
}

func (s *SkyLayer) Height() int {
	return s.height
}

func (s *SkyLayer) Color() color.Color {
	return s.color
}

// Setters
func (s *SkyLayer) SetOffsetZ(offsetZ float64) {
	s.offsetZ = offsetZ
}

func (s *SkyLayer) SetParallaxFactor(parallaxFactor float64) {
	if parallaxFactor < 0 {
		parallaxFactor = 0.0
	}
	s.parallaxFactor = parallaxFactor
}

func (s *SkyLayer) SetWidth(width int) {
	s.width = width
}

func (s *SkyLayer) SetHeight(height int) {
	s.height = height
}

func (s *SkyLayer) SetColor(c color.Color) {
	s.color = c
}

func NewSkyLayer(width, height int, parallaxFactor float64) *SkyLayer {
	return &SkyLayer{
		offsetZ:        0,
		parallaxFactor: parallaxFactor,
		width:          width,
		height:         height,
		color:          color.RGBA{135, 206, 235, 255},
	}
}

func (s *SkyLayer) Update(ctx common.WorldContext, delta float64) {
	effectiveSpeed := ctx.GetSpeed() * s.parallaxFactor
	s.offsetZ += effectiveSpeed * delta
	if s.offsetZ > 1000 {
		s.offsetZ -= 1000
	}
}

func (s *SkyLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx common.WorldContext) {
	ebitenutil.DrawRect(screen, 0, 0, float64(s.width), float64(s.height), s.color)
}
