package world

import (
	"github.com/hajimehoshi/ebiten/v2"

	"TheFiaskoTest/internal/render"
)

type Layer interface {
	Update(ctx WorldContext)
	Draw(screen *ebiten.Image, cam *render.Camera, ctx WorldContext)
}

// Интерфейс для сущностей (разрывает цикл)
type WorldContext interface {
	GetSpeed() float64
	GetWorldOffsetZ() float64
}

type World struct {
	WorldOffsetZ float64
	Speed        float64

	Layers []Layer
}

func New(speed float64) *World {
	return &World{
		WorldOffsetZ: 0,
		Speed:        speed,
		Layers:       []Layer{},
	}
}

func (w *World) AddLayer(l Layer) {
	w.Layers = append(w.Layers, l)
}

func (w *World) Update() {
	w.WorldOffsetZ += w.Speed
	for _, layer := range w.Layers {
		layer.Update(w) // World реализует WorldContext
	}
}

func (w *World) Draw(screen *ebiten.Image, cam *render.Camera) {
	for _, layer := range w.Layers {
		layer.Draw(screen, cam, w)
	}
}

// Реализация интерфейса WorldContext
func (w *World) GetSpeed() float64        { return w.Speed }
func (w *World) GetWorldOffsetZ() float64 { return w.WorldOffsetZ }
