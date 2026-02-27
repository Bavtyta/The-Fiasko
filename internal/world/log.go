package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/entity"
	"TheFiaskoTest/internal/render"
)

type LogLayer struct {
	Entities []*entity.TestPoint
	X, Y     float64
	Width    float64
	Height   float64
	Color    color.Color
}

func NewLogLayer(x, y, width, height float64) *LogLayer {
	return &LogLayer{
		X:        x,
		Y:        y,
		Width:    width,
		Height:   height,
		Color:    color.RGBA{139, 69, 19, 255},
		Entities: []*entity.TestPoint{},
	}
}

func (l *LogLayer) AddEntity(e *entity.TestPoint) {
	l.Entities = append(l.Entities, e)
}

func (l *LogLayer) Update(ctx WorldContext) {
	for _, e := range l.Entities {
		e.Update(ctx)
	}
}

func (l *LogLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx WorldContext) {
	// рисуем бревно
	ebitenutil.DrawRect(screen, l.X, l.Y, l.Width, l.Height, l.Color)

	// рисуем сущности на бревне
	for _, e := range l.Entities {
		e.Draw(screen, cam, ctx)
	}
}
