package game

import (
	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/entity"
	"TheFiaskoTest/internal/pools"
	"testing"
)

// TestHandleCollisions_NoCollision проверяет, что handleCollisions возвращает false
// когда нет коллизий между игроком и препятствиями.
func TestHandleCollisions_NoCollision(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Создаём препятствие далеко от игрока
	obs := &entity.Obstacle{
		Position: core.Vec3{X: 0, Y: 0, Z: 20}, // Далеко по Z
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
	}
	world.obstacles = append(world.obstacles, obs)

	// Проверяем, что коллизии нет
	if handleCollisions(world) {
		t.Error("Expected no collision when obstacle is far away")
	}
}

// TestHandleCollisions_WithCollision проверяет, что handleCollisions возвращает true
// когда есть коллизия между игроком и препятствием.
func TestHandleCollisions_WithCollision(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Создаём препятствие в той же позиции что и игрок
	obs := &entity.Obstacle{
		Position: core.Vec3{X: 0, Y: 0, Z: 0}, // Та же позиция что и игрок
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
	}
	world.obstacles = append(world.obstacles, obs)

	// Проверяем, что есть коллизия
	if !handleCollisions(world) {
		t.Error("Expected collision when obstacle overlaps with player")
	}
}

// TestHandleCollisions_EarlyOutOptimization проверяет, что early-out оптимизация
// работает корректно и пропускает далёкие препятствия.
func TestHandleCollisions_EarlyOutOptimization(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Создаём препятствие на расстоянии > 5.0 по Z (должно быть пропущено)
	farObs := &entity.Obstacle{
		Position: core.Vec3{X: 0, Y: 0, Z: 10}, // Расстояние = 10 > 5.0
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
	}
	world.obstacles = append(world.obstacles, farObs)

	// Создаём препятствие близко, но без коллизии
	nearObs := &entity.Obstacle{
		Position: core.Vec3{X: 10, Y: 0, Z: 2}, // Близко по Z, но далеко по X
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
	}
	world.obstacles = append(world.obstacles, nearObs)

	// Проверяем, что коллизии нет
	if handleCollisions(world) {
		t.Error("Expected no collision with early-out optimization")
	}
}

// TestWorldUpdate_CollisionSetsGameOver проверяет, что World.Update()
// устанавливает состояние StateGameOver при обнаружении коллизии.
func TestWorldUpdate_CollisionSetsGameOver(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Создаём препятствие в той же позиции что и игрок
	obs := &entity.Obstacle{
		Position: core.Vec3{X: 0, Y: 0, Z: 0},
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
	}
	world.obstacles = append(world.obstacles, obs)

	// Проверяем начальное состояние
	if world.state != StatePlaying {
		t.Error("Expected initial state to be StatePlaying")
	}

	// Обновляем мир
	world.Update(0.016)

	// Проверяем, что состояние изменилось на StateGameOver
	if world.state != StateGameOver {
		t.Error("Expected state to be StateGameOver after collision")
	}
}

// TestWorldUpdate_NoCollisionContinuesPlaying проверяет, что World.Update()
// продолжает игру когда нет коллизий.
func TestWorldUpdate_NoCollisionContinuesPlaying(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Создаём препятствие далеко от игрока
	obs := &entity.Obstacle{
		Position: core.Vec3{X: 0, Y: 0, Z: 20},
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
	}
	world.obstacles = append(world.obstacles, obs)

	// Обновляем мир
	world.Update(0.016)

	// Проверяем, что состояние остаётся StatePlaying
	if world.state != StatePlaying {
		t.Error("Expected state to remain StatePlaying when no collision")
	}
}
