package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/common"
	"TheFiaskoTest/internal/render"
)

type SkyLayer struct {
	offsetZ float64
	speed   float64
	width   int
	height  int
	color   color.Color
}

// Getters
func (s *SkyLayer) OffsetZ() float64 {
	return s.offsetZ
}

func (s *SkyLayer) Speed() float64 {
	return s.speed
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

func (s *SkyLayer) SetSpeed(speed float64) {
	s.speed = speed
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

func NewSkyLayer(width, height int, speed float64) *SkyLayer {
	return &SkyLayer{
		offsetZ: 0,
		speed:   speed,
		width:   width,
		height:  height,
		color:   color.RGBA{135, 206, 235, 255},
	}
}

func (s *SkyLayer) Update(ctx common.WorldContext) {
	s.offsetZ += s.speed
	if s.offsetZ > 1000 {
		s.offsetZ -= 1000
	}
}

func (s *SkyLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx common.WorldContext) {
	ebitenutil.DrawRect(screen, 0, 0, float64(s.width), float64(s.height), s.color)
}
