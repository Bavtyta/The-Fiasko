package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/render"
)

type SkyLayer struct {
	OffsetZ float64
	Speed   float64
	Width   int
	Height  int
	Color   color.Color
}

func NewSkyLayer(width, height int, speed float64) *SkyLayer {
	return &SkyLayer{
		OffsetZ: 0,
		Speed:   speed,
		Width:   width,
		Height:  height,
		Color:   color.RGBA{135, 206, 235, 255},
	}
}

func (s *SkyLayer) Update(ctx WorldContext) {
	s.OffsetZ += s.Speed
	if s.OffsetZ > 1000 {
		s.OffsetZ -= 1000
	}
}

func (s *SkyLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx WorldContext) {
	ebitenutil.DrawRect(screen, 0, 0, float64(s.Width), float64(s.Height), s.Color)
}
