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
