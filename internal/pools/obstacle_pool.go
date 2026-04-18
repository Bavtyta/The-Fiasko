// Пакет pools содержит пулы объектов для переиспользования.
package pools

import (
	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/entity"
)

// ObstaclePool - простой пул для переиспользования препятствий.
// КРИТИЧНО: Ограничение размера (maxSize) для предотвращения бесконечного роста.
type ObstaclePool struct {
	obstacles []*entity.Obstacle
	maxSize   int
	config    *config.Config
}

// NewObstaclePool создаёт новый пул препятствий с заданным максимальным размером и конфигурацией.
func NewObstaclePool(maxSize int, cfg *config.Config) *ObstaclePool {
	return &ObstaclePool{
		obstacles: make([]*entity.Obstacle, 0, maxSize),
		maxSize:   maxSize,
		config:    cfg,
	}
}

// Get возвращает препятствие из пула или создаёт новое, если пул пуст.
func (p *ObstaclePool) Get() *entity.Obstacle {
	if len(p.obstacles) > 0 {
		// Берём из пула
		obs := p.obstacles[len(p.obstacles)-1]
		p.obstacles = p.obstacles[:len(p.obstacles)-1]
		return obs
	}

	// Пул пуст, создаём новое
	return &entity.Obstacle{
		Width:  p.config.ObstacleWidth,
		Height: p.config.ObstacleHeight,
		Depth:  p.config.ObstacleDepth,
	}
}

// Put возвращает препятствие в пул после полного сброса состояния.
// КРИТИЧНО: Полностью сбрасывает состояние через Reset().
func (p *ObstaclePool) Put(obs *entity.Obstacle) {
	// КРИТИЧНО: Полностью сбрасываем состояние
	obs.Reset()

	// КРИТИЧНО: Возвращаем в пул только если не превышен лимит
	if len(p.obstacles) < p.maxSize {
		p.obstacles = append(p.obstacles, obs)
	}
	// Иначе просто отбрасываем (GC соберёт)
}
