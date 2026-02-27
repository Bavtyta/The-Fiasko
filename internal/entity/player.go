package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	X, Y   float64
	Width  float64
	Height float64
	Color  color.Color
}

func NewPlayer(x, y, w, h float64) *Player {
	return &Player{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
		Color:  color.White,
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, p.X, p.Y, p.Width, p.Height, p.Color)
}
