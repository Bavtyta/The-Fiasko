package entity

import (
	"TheFiaskoTest/internal/core"
	"testing"
)

// TestUpdateObstacle проверяет обновление позиции препятствия
func TestUpdateObstacle(t *testing.T) {
	// Создаём препятствие с начальной позицией и скоростью
	obstacle := &Obstacle{
		Position: core.Vec3{X: 0, Y: 0, Z: 50},
		Velocity: core.Vec3{X: 0, Y: 0, Z: 5},
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
		Type:     ObstacleTypeLog,
	}

	// Обновляем препятствие с dt = 1.0 секунда
	dt := 1.0
	UpdateObstacle(obstacle, dt)

	// Проверяем, что Z уменьшилась (препятствие движется к игроку)
	expectedZ := 50.0 - 5.0*dt // 50 - 5 = 45
	if obstacle.Position.Z != expectedZ {
		t.Errorf("Expected Z = %f, got %f", expectedZ, obstacle.Position.Z)
	}

	// Проверяем, что X и Y не изменились (скорость по этим осям = 0)
	if obstacle.Position.X != 0 {
		t.Errorf("Expected X = 0, got %f", obstacle.Position.X)
	}
	if obstacle.Position.Y != 0 {
		t.Errorf("Expected Y = 0, got %f", obstacle.Position.Y)
	}
}

// TestUpdateObstacleWithXYVelocity проверяет обновление с ненулевой скоростью по X и Y
func TestUpdateObstacleWithXYVelocity(t *testing.T) {
	obstacle := &Obstacle{
		Position: core.Vec3{X: 10, Y: 5, Z: 50},
		Velocity: core.Vec3{X: 2, Y: -1, Z: 5},
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
		Type:     ObstacleTypeLog,
	}

	dt := 0.5
	UpdateObstacle(obstacle, dt)

	// Проверяем все координаты
	expectedX := 10.0 + 2.0*dt   // 10 + 1 = 11
	expectedY := 5.0 + (-1.0)*dt // 5 - 0.5 = 4.5
	expectedZ := 50.0 - 5.0*dt   // 50 - 2.5 = 47.5

	if obstacle.Position.X != expectedX {
		t.Errorf("Expected X = %f, got %f", expectedX, obstacle.Position.X)
	}
	if obstacle.Position.Y != expectedY {
		t.Errorf("Expected Y = %f, got %f", expectedY, obstacle.Position.Y)
	}
	if obstacle.Position.Z != expectedZ {
		t.Errorf("Expected Z = %f, got %f", expectedZ, obstacle.Position.Z)
	}
}

// TestUpdateObstacleMultipleFrames проверяет обновление за несколько кадров
func TestUpdateObstacleMultipleFrames(t *testing.T) {
	obstacle := &Obstacle{
		Position: core.Vec3{X: 0, Y: 0, Z: 50},
		Velocity: core.Vec3{X: 0, Y: 0, Z: 10},
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
		Type:     ObstacleTypeLog,
	}

	// Симулируем 5 кадров по 0.1 секунды
	dt := 0.1
	for i := 0; i < 5; i++ {
		UpdateObstacle(obstacle, dt)
	}

	// После 5 кадров: Z = 50 - 10 * 0.1 * 5 = 50 - 5 = 45
	expectedZ := 50.0 - 10.0*dt*5
	if obstacle.Position.Z != expectedZ {
		t.Errorf("Expected Z = %f, got %f", expectedZ, obstacle.Position.Z)
	}
}
