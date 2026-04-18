// Пакет game содержит игровую логику и оркестрацию.
package game

import (
	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/entity"
	"TheFiaskoTest/internal/pools"
	"math/rand/v2"
)

// Spawner генерирует препятствия с таймером.
// КРИТИЧНО: Правильная обработка таймера (timer -= interval, НЕ timer = 0).
type Spawner struct {
	timer    float64
	interval float64
	pool     *pools.ObstaclePool
	config   *config.Config
}

// NewSpawner создаёт новый Spawner с заданной конфигурацией и пулом.
func NewSpawner(cfg *config.Config, pool *pools.ObstaclePool) *Spawner {
	return &Spawner{
		timer:    0,
		interval: cfg.SpawnInterval,
		pool:     pool,
		config:   cfg,
	}
}

// Update обновляет таймер спавнера и возвращает новые препятствия.
// КРИТИЧНО: Использует цикл for и timer -= interval (НЕ timer = 0!)
// для сохранения точности таймера и поддержки burst spawn при больших dt.
// КРИТИЧНО: Возвращает []*Obstacle (НЕ *Obstacle) для поддержки множественного спавна.
func (s *Spawner) Update(dt float64) []*entity.Obstacle {
	s.timer += dt

	spawned := make([]*entity.Obstacle, 0, 2)

	// КРИТИЧНО: Используем цикл for и timer -= interval (НЕ timer = 0!)
	// Это сохраняет точность таймера и поддерживает burst spawn при лагах
	for s.timer >= s.interval {
		s.timer -= s.interval // Вычитаем интервал, НЕ обнуляем!

		// Получаем препятствие из пула
		obs := s.pool.Get()

		// КРИТИЧНО (FIX #6): Инициализируем позицию из Config (НЕ хардкод)
		obs.Position = core.Vec3{
			X: (rand.Float64() - 0.5) * s.config.SpawnRangeX,
			Y: 0,
			Z: s.config.SpawnZ,
		}

		// Инициализируем скорость из config
		obs.Velocity = core.Vec3{Z: s.config.ObstacleSpeed}

		// КРИТИЧНО (FIX #3): Устанавливаем Type = ObstacleTypeLog (enum, НЕ string "log")
		obs.Type = entity.ObstacleTypeLog

		spawned = append(spawned, obs)
	}

	return spawned // КРИТИЧНО: Возвращаем []*Obstacle, НЕ *Obstacle
}
