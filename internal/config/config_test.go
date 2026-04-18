package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Проверка игровых настроек
	if cfg.ScreenWidth <= 0 {
		t.Errorf("ScreenWidth должен быть положительным, получено: %d", cfg.ScreenWidth)
	}
	if cfg.ScreenHeight <= 0 {
		t.Errorf("ScreenHeight должен быть положительным, получено: %d", cfg.ScreenHeight)
	}
	if cfg.TargetFPS <= 0 {
		t.Errorf("TargetFPS должен быть положительным, получено: %d", cfg.TargetFPS)
	}

	// Проверка физики
	if cfg.Gravity <= 0 {
		t.Errorf("Gravity должна быть положительной, получено: %f", cfg.Gravity)
	}
	if cfg.JumpVelocity <= 0 {
		t.Errorf("JumpVelocity должна быть положительной, получено: %f", cfg.JumpVelocity)
	}

	// Проверка игрока
	if cfg.BalanceSpeed <= 0 {
		t.Errorf("BalanceSpeed должна быть положительной, получено: %f", cfg.BalanceSpeed)
	}
	if cfg.MaxBalance <= 0 {
		t.Errorf("MaxBalance должен быть положительным, получено: %f", cfg.MaxBalance)
	}

	// Проверка спавна
	if cfg.SpawnInterval <= 0 {
		t.Errorf("SpawnInterval должен быть положительным, получено: %f", cfg.SpawnInterval)
	}
	if cfg.SpawnRangeX <= 0 {
		t.Errorf("SpawnRangeX должен быть положительным, получено: %f", cfg.SpawnRangeX)
	}
	if cfg.SpawnZ <= 0 {
		t.Errorf("SpawnZ должен быть положительным, получено: %f", cfg.SpawnZ)
	}
	if cfg.ObstacleSpeed <= 0 {
		t.Errorf("ObstacleSpeed должна быть положительной, получено: %f", cfg.ObstacleSpeed)
	}
	if cfg.MaxObstacles <= 0 {
		t.Errorf("MaxObstacles должен быть положительным, получено: %d", cfg.MaxObstacles)
	}
	if cfg.DespawnZ >= 0 {
		t.Errorf("DespawnZ должен быть отрицательным (позади игрока), получено: %f", cfg.DespawnZ)
	}

	// Проверка камеры
	if cfg.CameraFocalLength <= 0 {
		t.Errorf("CameraFocalLength должна быть положительной, получено: %f", cfg.CameraFocalLength)
	}
	if cfg.CameraHorizonY <= 0 {
		t.Errorf("CameraHorizonY должна быть положительной, получено: %f", cfg.CameraHorizonY)
	}
}
