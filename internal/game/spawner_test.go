package game

import (
	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/entity"
	"TheFiaskoTest/internal/pools"
	"testing"
)

// TestSpawnerUpdate проверяет базовую функциональность Spawner.Update
func TestSpawnerUpdate(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.SpawnInterval = 1.0 // Устанавливаем интервал 1 секунда для теста
	pool := pools.NewObstaclePool(10, cfg)
	spawner := NewSpawner(cfg, pool)

	// Тест 1: До интервала не должно быть спавна
	result := spawner.Update(0.5)
	if len(result) != 0 {
		t.Errorf("Expected 0 obstacles before interval, got %d", len(result))
	}

	// Тест 2: После интервала должен быть спавн
	result = spawner.Update(0.6) // Всего 1.1 секунды
	if len(result) != 1 {
		t.Errorf("Expected 1 obstacle after interval, got %d", len(result))
	}

	// Тест 3: Проверяем, что препятствие правильно инициализировано
	if len(result) > 0 {
		obs := result[0]
		if obs.Type != entity.ObstacleTypeLog {
			t.Errorf("Expected ObstacleTypeLog, got %v", obs.Type)
		}
		if obs.Position.Z != cfg.SpawnZ {
			t.Errorf("Expected Z position %f, got %f", cfg.SpawnZ, obs.Position.Z)
		}
		if obs.Velocity.Z != cfg.ObstacleSpeed {
			t.Errorf("Expected velocity Z %f, got %f", cfg.ObstacleSpeed, obs.Velocity.Z)
		}
	}
}

// TestSpawnerBurstSpawn проверяет burst spawn при больших dt
func TestSpawnerBurstSpawn(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.SpawnInterval = 1.0
	pool := pools.NewObstaclePool(10, cfg)
	spawner := NewSpawner(cfg, pool)

	// Большой dt должен вызвать множественный спавн
	result := spawner.Update(2.5) // Должно заспавниться 2 препятствия
	if len(result) != 2 {
		t.Errorf("Expected 2 obstacles for dt=2.5s, got %d", len(result))
	}
}

// TestSpawnerTimerAccuracy проверяет точность таймера (timer -= interval)
func TestSpawnerTimerAccuracy(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.SpawnInterval = 1.0
	pool := pools.NewObstaclePool(10, cfg)
	spawner := NewSpawner(cfg, pool)

	// Первый спавн: 1.1 секунды (должен заспавниться 1, остаток 0.1)
	result := spawner.Update(1.1)
	if len(result) != 1 {
		t.Errorf("Expected 1 obstacle, got %d", len(result))
	}

	// Второй спавн: 0.95 секунды (всего 1.05, должен заспавниться 1)
	result = spawner.Update(0.95)
	if len(result) != 1 {
		t.Errorf("Expected 1 obstacle (timer accuracy test), got %d", len(result))
	}
}

// TestSpawnerPositionRange проверяет, что позиция X в правильном диапазоне
func TestSpawnerPositionRange(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.SpawnInterval = 0.1 // Быстрый спавн для теста
	cfg.SpawnRangeX = 10.0
	pool := pools.NewObstaclePool(50, cfg)
	spawner := NewSpawner(cfg, pool)

	// Спавним несколько препятствий
	result := spawner.Update(1.0) // Должно заспавниться ~10 препятствий

	// Проверяем, что все X позиции в диапазоне [-5, 5]
	for _, obs := range result {
		if obs.Position.X < -5.0 || obs.Position.X > 5.0 {
			t.Errorf("Position X %f is out of range [-5, 5]", obs.Position.X)
		}
	}
}
