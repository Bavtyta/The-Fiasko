package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/render"
)

type TestPoint struct {
	Position core.Vec3
}

func NewTestPoint(x, y, z float64) *TestPoint {
	return &TestPoint{
		Position: core.Vec3{X: x, Y: y, Z: z},
	}
}

// Обновление точки с использованием WorldContext
func (t *TestPoint) Update(ctx WorldContext) {
	t.Position.Z -= ctx.GetSpeed()

	const MaxDepth = 500
	if t.Position.Z < ctx.GetWorldOffsetZ() {
		t.Position.Z += MaxDepth
	}
}

func (t *TestPoint) Draw(screen *ebiten.Image, cam *render.Camera, ctx WorldContext) {
	relative := t.Position
	relative.Z -= ctx.GetWorldOffsetZ()

	x, y, scale := render.Project(relative, cam)
	if scale == 0 {
		// Объект за near plane – не рисуем
		return
	}

	const baseSize = 10.0
	const maxSize = 150.0

	size := baseSize * scale
	if size > maxSize {
		size = maxSize
	}

	ebitenutil.DrawRect(screen, x-size/2, y-size/2, size, size, color.White)
}

func (t *TestPoint) GetZ() float64 {
	return t.Position.Z
}

func (t *TestPoint) SetZ(z float64) {
	t.Position.Z = z
}
