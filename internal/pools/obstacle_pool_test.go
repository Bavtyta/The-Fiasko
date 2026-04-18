package pools

import (
	"testing"

	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/entity"
)

// TestObstaclePoolGet проверяет метод Get()
func TestObstaclePoolGet(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := NewObstaclePool(5, cfg)

	// Получаем препятствие из пустого пула (должно создать новое)
	obs := pool.Get()
	if obs == nil {
		t.Fatal("Expected Get() to return a new obstacle, got nil")
	}

	// Проверяем начальные значения из конфигурации
	if obs.Width != cfg.ObstacleWidth || obs.Height != cfg.ObstacleHeight || obs.Depth != cfg.ObstacleDepth {
		t.Errorf("Expected default dimensions (%f, %f, %f), got (%f, %f, %f)",
			cfg.ObstacleWidth, cfg.ObstacleHeight, cfg.ObstacleDepth,
			obs.Width, obs.Height, obs.Depth)
	}
}

// TestObstaclePoolPutAndGet проверяет Put() и повторный Get()
func TestObstaclePoolPutAndGet(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := NewObstaclePool(5, cfg)

	// Создаём препятствие с изменёнными значениями
	obs := &entity.Obstacle{
		Width:  5.0,
		Height: 3.0,
		Depth:  4.0,
	}
	obs.Position.X = 10.0
	obs.Velocity.Y = 5.0

	// Возвращаем в пул
	pool.Put(obs)

	// Получаем обратно
	retrieved := pool.Get()
	if retrieved == nil {
		t.Fatal("Expected Get() to return obstacle from pool, got nil")
	}

	// Проверяем, что состояние было сброшено
	if retrieved.Position.X != 0 || retrieved.Velocity.Y != 0 {
		t.Errorf("Expected Reset() to clear position and velocity, got Position.X=%f, Velocity.Y=%f",
			retrieved.Position.X, retrieved.Velocity.Y)
	}

	if retrieved.Width != 0 || retrieved.Height != 0 || retrieved.Depth != 0 {
		t.Errorf("Expected Reset() to clear dimensions, got (%f, %f, %f)",
			retrieved.Width, retrieved.Height, retrieved.Depth)
	}
}

// TestObstaclePoolMaxSize проверяет ограничение размера пула
func TestObstaclePoolMaxSize(t *testing.T) {
	cfg := config.DefaultConfig()
	maxSize := 3
	pool := NewObstaclePool(maxSize, cfg)

	// Добавляем больше препятствий, чем maxSize
	for i := 0; i < 10; i++ {
		obs := &entity.Obstacle{}
		pool.Put(obs)
	}

	// Проверяем, что в пуле не больше maxSize препятствий
	count := 0
	for i := 0; i < 10; i++ {
		if len(pool.obstacles) == 0 {
			break
		}
		pool.Get()
		count++
	}

	if count > maxSize {
		t.Errorf("Expected pool to contain at most %d obstacles, but got %d", maxSize, count)
	}
}

// TestObstaclePoolReset проверяет, что Reset() вызывается при Put()
func TestObstaclePoolReset(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := NewObstaclePool(5, cfg)

	// Создаём препятствие с установленными значениями
	obs := &entity.Obstacle{
		Width:  10.0,
		Height: 20.0,
		Depth:  30.0,
		Type:   entity.ObstacleTypeRock,
	}
	obs.Position.X = 100.0
	obs.Position.Y = 200.0
	obs.Position.Z = 300.0
	obs.Velocity.X = 1.0
	obs.Velocity.Y = 2.0
	obs.Velocity.Z = 3.0

	// Возвращаем в пул (должен вызвать Reset())
	pool.Put(obs)

	// Получаем обратно
	retrieved := pool.Get()

	// Проверяем, что все поля сброшены
	if retrieved.Position.X != 0 || retrieved.Position.Y != 0 || retrieved.Position.Z != 0 {
		t.Error("Expected Position to be reset to zero")
	}
	if retrieved.Velocity.X != 0 || retrieved.Velocity.Y != 0 || retrieved.Velocity.Z != 0 {
		t.Error("Expected Velocity to be reset to zero")
	}
	if retrieved.Width != 0 || retrieved.Height != 0 || retrieved.Depth != 0 {
		t.Error("Expected dimensions to be reset to zero")
	}
	if retrieved.Type != entity.ObstacleTypeLog {
		t.Errorf("Expected Type to be reset to ObstacleTypeLog, got %v", retrieved.Type)
	}
}
