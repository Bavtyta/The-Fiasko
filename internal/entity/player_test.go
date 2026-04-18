package entity

import (
	"testing"

	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/core"
)

func TestNewPlayer(t *testing.T) {
	cfg := config.DefaultConfig()
	player := NewPlayer(cfg)

	// Проверяем дефолтные значения позиции
	if player.Position.X != 0 || player.Position.Y != 0 || player.Position.Z != 0 {
		t.Errorf("Expected position (0,0,0), got (%f,%f,%f)",
			player.Position.X, player.Position.Y, player.Position.Z)
	}

	// Проверяем дефолтные значения скорости
	if player.Velocity.X != 0 || player.Velocity.Y != 0 || player.Velocity.Z != 0 {
		t.Errorf("Expected velocity (0,0,0), got (%f,%f,%f)",
			player.Velocity.X, player.Velocity.Y, player.Velocity.Z)
	}

	// Проверяем дефолтный баланс
	if player.Balance != 0 {
		t.Errorf("Expected balance 0, got %f", player.Balance)
	}

	// Проверяем дефолтные размеры из конфигурации
	if player.Width != cfg.PlayerWidth {
		t.Errorf("Expected width %f, got %f", cfg.PlayerWidth, player.Width)
	}
	if player.Height != cfg.PlayerHeight {
		t.Errorf("Expected height %f, got %f", cfg.PlayerHeight, player.Height)
	}
	if player.Depth != cfg.PlayerDepth {
		t.Errorf("Expected depth %f, got %f", cfg.PlayerDepth, player.Depth)
	}
}

func TestPlayerDataOnly(t *testing.T) {
	// Проверяем, что Player - это чистая структура данных
	cfg := config.DefaultConfig()
	player := NewPlayer(cfg)

	// Можем изменять поля напрямую
	player.Position = core.Vec3{X: 10, Y: 5, Z: 20}
	player.Velocity = core.Vec3{X: 1, Y: 2, Z: 3}
	player.Balance = 5.0
	player.Width = 3.0
	player.Height = 4.0
	player.Depth = 2.0

	// Проверяем, что изменения применились
	if player.Position.X != 10 || player.Position.Y != 5 || player.Position.Z != 20 {
		t.Error("Position should be mutable")
	}
	if player.Velocity.X != 1 || player.Velocity.Y != 2 || player.Velocity.Z != 3 {
		t.Error("Velocity should be mutable")
	}
	if player.Balance != 5.0 {
		t.Error("Balance should be mutable")
	}
	if player.Width != 3.0 || player.Height != 4.0 || player.Depth != 2.0 {
		t.Error("Dimensions should be mutable")
	}
}

func TestUpdatePlayer_Gravity(t *testing.T) {
	cfg := config.DefaultConfig()
	player := NewPlayer(cfg)
	cfg.Gravity = 20.0

	// Игрок в воздухе
	player.Position.Y = 5.0
	player.Velocity.Y = 0.0

	UpdatePlayer(player, cfg, 0.1)

	// Гравитация должна уменьшить скорость по Y
	if player.Velocity.Y >= 0 {
		t.Errorf("Expected negative velocity after gravity, got %f", player.Velocity.Y)
	}

	// Позиция должна обновиться
	if player.Position.Y >= 5.0 {
		t.Errorf("Expected position to decrease, got %f", player.Position.Y)
	}
}

func TestUpdatePlayer_GroundLimit(t *testing.T) {
	cfg := config.DefaultConfig()
	player := NewPlayer(cfg)
	cfg.Gravity = 20.0

	// Игрок падает с отрицательной скоростью
	player.Position.Y = 0.5
	player.Velocity.Y = -10.0

	UpdatePlayer(player, cfg, 0.1)

	// Позиция не должна быть ниже 0
	if player.Position.Y < 0 {
		t.Errorf("Expected position Y >= 0, got %f", player.Position.Y)
	}

	// Скорость должна быть сброшена
	if player.Velocity.Y != 0 {
		t.Errorf("Expected velocity Y = 0 on ground, got %f", player.Velocity.Y)
	}
}

func TestUpdatePlayer_NoGravityOnGround(t *testing.T) {
	cfg := config.DefaultConfig()
	player := NewPlayer(cfg)
	cfg.Gravity = 20.0

	// Игрок на земле
	player.Position.Y = 0.0
	player.Velocity.Y = 0.0

	UpdatePlayer(player, cfg, 0.1)

	// Гравитация не должна применяться на земле
	if player.Velocity.Y != 0 {
		t.Errorf("Expected no gravity on ground, got velocity Y = %f", player.Velocity.Y)
	}
}

func TestUpdatePlayer_PositionUpdate(t *testing.T) {
	cfg := config.DefaultConfig()
	player := NewPlayer(cfg)
	cfg.Gravity = 20.0

	// Устанавливаем начальные значения
	player.Position = core.Vec3{X: 0, Y: 10, Z: 0}
	player.Velocity = core.Vec3{X: 5, Y: 10, Z: 3}

	UpdatePlayer(player, cfg, 0.1)

	// Позиция должна обновиться на основе скорости и dt
	expectedX := 0.0 + 5.0*0.1
	expectedZ := 0.0 + 3.0*0.1

	tolerance := 0.0001
	if abs(player.Position.X-expectedX) > tolerance {
		t.Errorf("Expected position X = %f, got %f", expectedX, player.Position.X)
	}
	if abs(player.Position.Z-expectedZ) > tolerance {
		t.Errorf("Expected position Z = %f, got %f", expectedZ, player.Position.Z)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestCheckPlayerFall_WithinBalance(t *testing.T) {
	cfg := config.DefaultConfig()
	player := NewPlayer(cfg)
	cfg.MaxBalance = 10.0

	// Баланс в пределах нормы
	player.Balance = 5.0

	if CheckPlayerFall(player, cfg) {
		t.Error("Player should not fall with balance within limits")
	}

	// Отрицательный баланс в пределах нормы
	player.Balance = -5.0

	if CheckPlayerFall(player, cfg) {
		t.Error("Player should not fall with negative balance within limits")
	}
}

func TestCheckPlayerFall_ExceedsBalance(t *testing.T) {
	cfg := config.DefaultConfig()
	player := NewPlayer(cfg)
	cfg.MaxBalance = 10.0

	// Баланс превышает максимум
	player.Balance = 11.0

	if !CheckPlayerFall(player, cfg) {
		t.Error("Player should fall when balance exceeds max")
	}

	// Отрицательный баланс превышает максимум
	player.Balance = -11.0

	if !CheckPlayerFall(player, cfg) {
		t.Error("Player should fall when negative balance exceeds max")
	}
}

func TestCheckPlayerFall_ExactlyAtLimit(t *testing.T) {
	cfg := config.DefaultConfig()
	player := NewPlayer(cfg)
	cfg.MaxBalance = 10.0

	// Баланс точно на границе
	player.Balance = 10.0

	if CheckPlayerFall(player, cfg) {
		t.Error("Player should not fall when balance equals max")
	}

	// Отрицательный баланс точно на границе
	player.Balance = -10.0

	if CheckPlayerFall(player, cfg) {
		t.Error("Player should not fall when negative balance equals max")
	}
}
