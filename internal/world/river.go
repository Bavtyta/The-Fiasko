package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/render"
)

type RiverLayer struct {
	OffsetZ float64
	Speed   float64
	Width   int
	Height  int
	Color   color.Color
}

func NewRiverLayer(width, height int, speed float64) *RiverLayer {
	return &RiverLayer{
		OffsetZ: 0,
		Speed:   speed,
		Width:   width,
		Height:  height,
		Color:   color.RGBA{0, 100, 255, 255},
	}
}

// Update теперь принимает WorldContext (даже если не используется)
func (r *RiverLayer) Update(ctx WorldContext) {
	r.OffsetZ += r.Speed
	if r.OffsetZ > 1000 {
		r.OffsetZ -= 1000
	}
}

// Draw получил два новых параметра (пока игнорируем их)
func (r *RiverLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx WorldContext) {
	ebitenutil.DrawRect(screen, 0, 0, float64(r.Width), float64(r.Height), r.Color)
}
