package world

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"TheFiaskoTest/internal/common"
	"TheFiaskoTest/internal/render"
)

const maxActiveObstacles = 100

type Layer interface {
	Update(ctx common.WorldContext, delta float64)
	Draw(screen *ebiten.Image, cam *render.Camera, ctx common.WorldContext)
}

// SurfaceProvider определяет возможность получить информацию о поверхности в точке Z.
type SurfaceProvider interface {
	SurfaceAt(z float64) (height float64, surfaceType SurfaceType, ok bool)
}

type World struct {
	worldOffsetZ float64
	speed        float64

	layers        []Layer
	obstacles     []*Obstacle
	baseSpawnDist float64 // базовое расстояние между препятствиями (в единицах Z)
	lastSpawnZ    float64 // последняя позиция Z, где было создано препятствие
}

func (w *World) AddLayer(l Layer) {
	w.layers = append(w.layers, l)
}

func New(speed float64) *World {
	return &World{
		worldOffsetZ:  0,
		speed:         speed,
		baseSpawnDist: 50.0, // появляться каждые 50 единиц пути при базовой скорости
		lastSpawnZ:    0.0,  // инициализация - первый спавн произойдет после прохождения baseSpawnDist
		layers:        []Layer{},
		obstacles:     []*Obstacle{},
	}
}

func (w *World) AddObstacle(obs *Obstacle) {
	w.obstacles = append(w.obstacles, obs)
}

// spawnObstacle creates a new obstacle on the farthest solid surface segment.
// The obstacle placement uses synchronized coordinates: segment.nearZ + offset.
// This ensures obstacles spawn at positions consistent with World_Offset_Z progression.
//
// Requirements: 5.6 (distance-based spawning), 7.3 (coordinate consistency)
func (w *World) spawnObstacle() {
	var farthestSeg *Segment
	maxZ := -math.MaxFloat64

	// Find the farthest solid surface segment across all layers
	// Segments are synchronized with World_Offset_Z through parallaxFactor=1.0
	for _, layer := range w.layers {
		if sl, ok := layer.(*SegmentLayer); ok && sl.SurfaceType() == SurfaceSolid {
			for _, seg := range sl.Segments() {
				segZ := seg.NearZ() + seg.Length() // far end of segment
				if segZ > maxZ {
					maxZ = segZ
					farthestSeg = seg
				}
			}
		}
	}

	if farthestSeg == nil {
		return // no suitable segment found
	}

	// Create obstacle at random position along segment with margin zones
	// to avoid spawning too close to segment boundaries
	const margin = 5.0
	segLength := farthestSeg.Length()

	// If segment is too short for margin zones, use entire segment
	var offsetZ float64
	if segLength > 2*margin {
		offsetZ = margin + rand.Float64()*(segLength-2*margin)
	} else {
		offsetZ = rand.Float64() * segLength
	}

	// Obstacle world position will be: segment.nearZ + offsetZ
	// This ensures the obstacle is within segment bounds [nearZ, nearZ+length]
	obs := NewObstacle(farthestSeg, offsetZ, 3, 5, 2*math.Pi)
	w.AddObstacle(obs)
}

func (w *World) Update(delta float64) {
	w.worldOffsetZ += w.speed * delta // если нужно общее смещение мира

	// Обновляем слои
	for _, layer := range w.layers {
		layer.Update(w, delta)
	}

	// Обновляем существующие препятствия
	for _, obs := range w.obstacles {
		obs.Update(delta)
	}

	// Генерация новых препятствий на основе пройденного расстояния
	// Используем цикл для обработки нескольких спавнов за один кадр при больших delta
	// Добавляем защиту от бесконечного цикла
	maxSpawnsPerFrame := 10
	spawnsThisFrame := 0
	for w.worldOffsetZ-w.lastSpawnZ >= w.baseSpawnDist && spawnsThisFrame < maxSpawnsPerFrame {
		if len(w.obstacles) < maxActiveObstacles {
			w.spawnObstacle()
		}
		w.lastSpawnZ += w.baseSpawnDist
		spawnsThisFrame++
	}

	// Удаляем препятствия, которые полностью позади камеры
	var remaining []*Obstacle
	for _, obs := range w.obstacles {
		if obs.WorldPos().Z > -10 { // допустим, порог
			remaining = append(remaining, obs)
		}
	}
	w.obstacles = remaining
}

func (w *World) Draw(screen *ebiten.Image, cam *render.Camera) {
	for _, layer := range w.layers {
		layer.Draw(screen, cam, w)
	}
	for _, obs := range w.obstacles {
		obs.Draw(screen, cam)
	}
}

// Получение списка препятствий (для проверки столкновений)
func (w *World) Obstacles() []*Obstacle {
	return w.obstacles
}

// GetSurfaceAt возвращает высоту и тип поверхности в мировой координате Z.
// Проходит по слоям в порядке их добавления (приоритет: более поздние слои имеют больший приоритет,
// так как они добавляются позже и могут перекрывать ранние).

// Реализация интерфейса WorldContext
func (w *World) GetSpeed() float64        { return w.speed }
func (w *World) GetWorldOffsetZ() float64 { return w.worldOffsetZ }

// SurfaceInfo содержит информацию о поверхности в точке Z
type SurfaceInfo struct {
	Height  float64
	Type    SurfaceType
	Segment *Segment // указатель на сегмент (для твёрдых поверхностей)
}

// GetSurfaceAt возвращает информацию о самой высокой поверхности в точке Z.
func (w *World) GetSurfaceAt(z float64) (SurfaceInfo, bool) {
	var best SurfaceInfo
	best.Height = -math.MaxFloat64
	found := false

	for _, layer := range w.layers {
		if sp, ok := layer.(SurfaceProvider); ok {
			if h, st, ok := sp.SurfaceAt(z); ok {
				// Пытаемся получить сегмент, если слой является SegmentLayer
				var seg *Segment
				if sl, ok := layer.(*SegmentLayer); ok {
					seg = sl.SegmentAt(z)
				}
				if !found || h > best.Height {
					best.Height = h
					best.Type = st
					best.Segment = seg
					found = true
				}
			}
		}
	}
	return best, found
}
