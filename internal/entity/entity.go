package entity

import (
	"TheFiaskoTest/internal/render"

	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	Update(ctx WorldContext)
	Draw(screen *ebiten.Image, cam *render.Camera, ctx WorldContext)
	GetZ() float64
	SetZ(z float64)
}

type WorldContext interface {
	GetWorldOffsetZ() float64
	GetSpeed() float64
}
