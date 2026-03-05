package world

import (
	"github.com/hajimehoshi/ebiten/v2"

	"TheFiaskoTest/internal/common"
	"TheFiaskoTest/internal/render"
)

type Layer interface {
	Update(ctx common.WorldContext)
	Draw(screen *ebiten.Image, cam *render.Camera, ctx common.WorldContext)
}

// SurfaceProvider определяет возможность получить информацию о поверхности в точке Z.
type SurfaceProvider interface {
	SurfaceAt(z float64) (height float64, surfaceType SurfaceType, ok bool)
}

type World struct {
	worldOffsetZ float64
	speed        float64

	layers []Layer
}

func New(speed float64) *World {
	return &World{
		worldOffsetZ: 0,
		speed:        speed,
		layers:       []Layer{},
	}
}

func (w *World) AddLayer(l Layer) {
	w.layers = append(w.layers, l)
}

func (w *World) Update() {
	w.worldOffsetZ += w.speed
	for _, layer := range w.layers {
		layer.Update(w) // World реализует WorldContext
	}
}

func (w *World) Draw(screen *ebiten.Image, cam *render.Camera) {
	for _, layer := range w.layers {
		layer.Draw(screen, cam, w)
	}
}

// GetSurfaceAt возвращает высоту и тип поверхности в мировой координате Z.
// Проходит по слоям в порядке их добавления (приоритет: более поздние слои имеют больший приоритет,
// так как они добавляются позже и могут перекрывать ранние).
func (w *World) GetSurfaceAt(z float64) (height float64, surfaceType SurfaceType, ok bool) {
	for _, layer := range w.layers {
		if sp, ok := layer.(SurfaceProvider); ok {
			if h, st, found := sp.SurfaceAt(z); found {
				return h, st, true
			}
		}
	}
	return 0, 0, false
}

// Реализация интерфейса WorldContext
func (w *World) GetSpeed() float64        { return w.speed }
func (w *World) GetWorldOffsetZ() float64 { return w.worldOffsetZ }

// Геттеры для всех полей
func (w *World) Speed() float64        { return w.speed }
func (w *World) WorldOffsetZ() float64 { return w.worldOffsetZ }
func (w *World) Layers() []Layer       { return w.layers }

// Сеттеры для изменяемых полей
func (w *World) SetSpeed(speed float64)         { w.speed = speed }
func (w *World) SetWorldOffsetZ(offset float64) { w.worldOffsetZ = offset }
func (w *World) SetLayers(layers []Layer)       { w.layers = layers }
