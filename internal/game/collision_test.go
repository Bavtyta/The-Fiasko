package game

import (
	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/entity"
	"testing"
)

func TestCheckAABBCollision_NoCollision(t *testing.T) {
	player := &entity.Player{
		Position: core.Vec3{X: 0, Y: 0, Z: 0},
		Width:    2.0,
		Height:   3.0,
		Depth:    1.0,
	}

	obstacle := &entity.Obstacle{
		Position: core.Vec3{X: 10, Y: 0, Z: 0}, // Далеко по X
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
	}

	if CheckAABBCollision(player, obstacle) {
		t.Error("Expected no collision when objects are far apart")
	}
}

func TestCheckAABBCollision_WithCollision(t *testing.T) {
	player := &entity.Player{
		Position: core.Vec3{X: 0, Y: 0, Z: 0},
		Width:    2.0,
		Height:   3.0,
		Depth:    1.0,
	}

	obstacle := &entity.Obstacle{
		Position: core.Vec3{X: 0, Y: 0, Z: 0}, // Та же позиция
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
	}

	if !CheckAABBCollision(player, obstacle) {
		t.Error("Expected collision when objects overlap")
	}
}

func TestCheckAABBCollision_EdgeCase(t *testing.T) {
	player := &entity.Player{
		Position: core.Vec3{X: 0, Y: 0, Z: 0},
		Width:    2.0,
		Height:   3.0,
		Depth:    1.0,
	}

	obstacle := &entity.Obstacle{
		Position: core.Vec3{X: 1.5, Y: 0, Z: 0}, // Близко, но касается
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
	}

	// Расстояние по X: 1.5
	// Сумма половин ширин: (2.0 + 2.0) / 2 = 2.0
	// 1.5 < 2.0, значит есть пересечение
	if !CheckAABBCollision(player, obstacle) {
		t.Error("Expected collision when objects are touching")
	}
}

func TestCheckAABBCollision_NoOverlapOnOneAxis(t *testing.T) {
	player := &entity.Player{
		Position: core.Vec3{X: 0, Y: 0, Z: 0},
		Width:    2.0,
		Height:   3.0,
		Depth:    1.0,
	}

	obstacle := &entity.Obstacle{
		Position: core.Vec3{X: 0, Y: 10, Z: 0}, // Далеко по Y
		Width:    2.0,
		Height:   2.0,
		Depth:    2.0,
	}

	// Пересечение по X и Z, но не по Y
	if CheckAABBCollision(player, obstacle) {
		t.Error("Expected no collision when objects don't overlap on Y axis")
	}
}
