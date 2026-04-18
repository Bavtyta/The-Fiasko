package game

import (
	"testing"

	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/pools"
)

// TestWorldUpdate проверяет базовую функциональность World.Update
func TestWorldUpdate(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Проверяем начальное состояние
	if world.score != 0 {
		t.Errorf("Expected initial score to be 0, got %f", world.score)
	}

	// Обновляем мир
	dt := 0.016 // ~60 FPS
	world.Update(dt)

	// Проверяем, что score увеличился
	if world.score != dt {
		t.Errorf("Expected score to be %f, got %f", dt, world.score)
	}

	// Обновляем ещё раз
	world.Update(dt)

	// Проверяем, что score продолжает увеличиваться
	expectedScore := dt * 2
	if world.score != expectedScore {
		t.Errorf("Expected score to be %f, got %f", expectedScore, world.score)
	}
}

// TestWorldUpdateWhenGameOver проверяет, что Update не выполняется когда игра окончена
func TestWorldUpdateWhenGameOver(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Устанавливаем состояние GameOver
	world.state = StateGameOver

	// Обновляем мир
	dt := 0.016
	world.Update(dt)

	// Проверяем, что score НЕ увеличился
	if world.score != 0 {
		t.Errorf("Expected score to remain 0 when game is over, got %f", world.score)
	}
}

// TestWorldIsGameOver проверяет метод IsGameOver
func TestWorldIsGameOver(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Проверяем начальное состояние
	if world.IsGameOver() {
		t.Error("Expected game to not be over initially")
	}

	// Устанавливаем состояние GameOver
	world.state = StateGameOver

	// Проверяем, что игра окончена
	if !world.IsGameOver() {
		t.Error("Expected game to be over after setting state to StateGameOver")
	}
}

// TestWorldSpawning проверяет интеграцию Spawner в World через handleSpawning
func TestWorldSpawning(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.SpawnInterval = 1.0 // Спавн каждую секунду
	cfg.MaxObstacles = 5    // Лимит 5 препятствий
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Проверяем начальное состояние
	if len(world.obstacles) != 0 {
		t.Errorf("Expected 0 obstacles initially, got %d", len(world.obstacles))
	}

	// Обновляем мир с dt < SpawnInterval (не должно заспавниться)
	world.Update(0.5)
	if len(world.obstacles) != 0 {
		t.Errorf("Expected 0 obstacles after 0.5s, got %d", len(world.obstacles))
	}

	// Обновляем мир с dt >= SpawnInterval (должно заспавниться 1 препятствие)
	world.Update(0.6) // Всего 1.1s
	if len(world.obstacles) != 1 {
		t.Errorf("Expected 1 obstacle after 1.1s, got %d", len(world.obstacles))
	}

	// Обновляем мир ещё раз (должно заспавниться ещё 1)
	world.Update(1.0)
	if len(world.obstacles) != 2 {
		t.Errorf("Expected 2 obstacles after another 1.0s, got %d", len(world.obstacles))
	}
}

// TestWorldSpawningMaxObstacles проверяет, что лимит MaxObstacles соблюдается
func TestWorldSpawningMaxObstacles(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.SpawnInterval = 0.1 // Быстрый спавн
	cfg.MaxObstacles = 3    // Лимит 3 препятствия
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Обновляем мир много раз, чтобы заспавнить больше MaxObstacles
	for i := 0; i < 10; i++ {
		world.Update(0.2) // Каждый раз должно спавниться 2 препятствия
	}

	// Проверяем, что количество препятствий не превышает MaxObstacles
	if len(world.obstacles) > cfg.MaxObstacles {
		t.Errorf("Expected at most %d obstacles, got %d", cfg.MaxObstacles, len(world.obstacles))
	}

	// Проверяем, что достигнут лимит
	if len(world.obstacles) != cfg.MaxObstacles {
		t.Errorf("Expected exactly %d obstacles (limit reached), got %d", cfg.MaxObstacles, len(world.obstacles))
	}
}

// TestWorldSpawningBurstSpawn проверяет обработку burst spawn (множественные препятствия за один Update)
func TestWorldSpawningBurstSpawn(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.SpawnInterval = 1.0 // Спавн каждую секунду
	cfg.MaxObstacles = 10   // Достаточный лимит
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Обновляем мир с большим dt (должно заспавниться несколько препятствий сразу)
	world.Update(3.5) // Должно заспавниться 3 препятствия

	if len(world.obstacles) != 3 {
		t.Errorf("Expected 3 obstacles after 3.5s (burst spawn), got %d", len(world.obstacles))
	}
}

// TestRemoveOffscreenObstacles проверяет удаление препятствий за границами экрана
func TestRemoveOffscreenObstacles(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DespawnZ = -10.0 // Граница удаления
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Создаём препятствия вручную с разными Z позициями
	obs1 := pool.Get()
	obs1.Position.Z = 20.0 // На экране
	world.obstacles = append(world.obstacles, obs1)

	obs2 := pool.Get()
	obs2.Position.Z = -15.0 // За границей (< DespawnZ)
	world.obstacles = append(world.obstacles, obs2)

	obs3 := pool.Get()
	obs3.Position.Z = 5.0 // На экране
	world.obstacles = append(world.obstacles, obs3)

	obs4 := pool.Get()
	obs4.Position.Z = -20.0 // За границей (< DespawnZ)
	world.obstacles = append(world.obstacles, obs4)

	// Проверяем начальное состояние
	if len(world.obstacles) != 4 {
		t.Errorf("Expected 4 obstacles initially, got %d", len(world.obstacles))
	}

	// Вызываем removeOffscreenObstacles
	removeOffscreenObstacles(world)

	// Проверяем, что остались только препятствия на экране
	if len(world.obstacles) != 2 {
		t.Errorf("Expected 2 obstacles after removal, got %d", len(world.obstacles))
	}

	// Проверяем, что остались правильные препятствия (с Z >= DespawnZ)
	for _, obs := range world.obstacles {
		if obs.Position.Z < cfg.DespawnZ {
			t.Errorf("Found obstacle with Z=%f (< DespawnZ=%f) after removal", obs.Position.Z, cfg.DespawnZ)
		}
	}

	// Проверяем, что препятствия вернулись в пул (пул должен содержать 2 объекта)
	// Мы не можем напрямую проверить содержимое пула, но можем проверить что Get() возвращает переиспользованные объекты
	pooledObs1 := pool.Get()
	pooledObs2 := pool.Get()

	// Проверяем, что объекты были сброшены (Position должна быть нулевой после Reset)
	if pooledObs1.Position.Z != 0 || pooledObs2.Position.Z != 0 {
		t.Error("Expected pooled obstacles to have reset positions")
	}
}

// TestRemoveOffscreenObstaclesEmpty проверяет removeOffscreenObstacles с пустым списком
func TestRemoveOffscreenObstaclesEmpty(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Проверяем начальное состояние (пустой список)
	if len(world.obstacles) != 0 {
		t.Errorf("Expected 0 obstacles initially, got %d", len(world.obstacles))
	}

	// Вызываем removeOffscreenObstacles (не должно паниковать)
	removeOffscreenObstacles(world)

	// Проверяем, что список остался пустым
	if len(world.obstacles) != 0 {
		t.Errorf("Expected 0 obstacles after removal, got %d", len(world.obstacles))
	}
}

// TestRemoveOffscreenObstaclesAllOnscreen проверяет removeOffscreenObstacles когда все препятствия на экране
func TestRemoveOffscreenObstaclesAllOnscreen(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DespawnZ = -10.0
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Создаём препятствия, все на экране
	for i := 0; i < 5; i++ {
		obs := pool.Get()
		obs.Position.Z = float64(i * 10) // Все Z >= 0 (> DespawnZ)
		world.obstacles = append(world.obstacles, obs)
	}

	// Проверяем начальное состояние
	if len(world.obstacles) != 5 {
		t.Errorf("Expected 5 obstacles initially, got %d", len(world.obstacles))
	}

	// Вызываем removeOffscreenObstacles
	removeOffscreenObstacles(world)

	// Проверяем, что все препятствия остались
	if len(world.obstacles) != 5 {
		t.Errorf("Expected 5 obstacles after removal (all onscreen), got %d", len(world.obstacles))
	}
}

// TestRemoveOffscreenObstaclesAllOffscreen проверяет removeOffscreenObstacles когда все препятствия за границей
func TestRemoveOffscreenObstaclesAllOffscreen(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DespawnZ = -10.0
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := NewWorld(cfg, pool)

	// Создаём препятствия, все за границей
	for i := 0; i < 5; i++ {
		obs := pool.Get()
		obs.Position.Z = -15.0 - float64(i) // Все Z < DespawnZ
		world.obstacles = append(world.obstacles, obs)
	}

	// Проверяем начальное состояние
	if len(world.obstacles) != 5 {
		t.Errorf("Expected 5 obstacles initially, got %d", len(world.obstacles))
	}

	// Вызываем removeOffscreenObstacles
	removeOffscreenObstacles(world)

	// Проверяем, что все препятствия удалены
	if len(world.obstacles) != 0 {
		t.Errorf("Expected 0 obstacles after removal (all offscreen), got %d", len(world.obstacles))
	}
}
