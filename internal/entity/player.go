package entity

import (
	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/core"
)

// Player представляет игрока как чистую структуру данных.
// КРИТИЧНО: Только поля данных, НЕТ методов кроме конструктора.
// НЕТ зависимостей от Config, ResourceManager, Camera в структуре.
type Player struct {
	// Позиция игрока в 3D пространстве
	Position core.Vec3

	// Скорость игрока
	Velocity core.Vec3

	// Баланс игрока (для механики балансирования)
	Balance float64

	// Размеры игрока
	Width  float64
	Height float64
	Depth  float64
}

// NewPlayer создаёт нового игрока с дефолтными значениями из конфигурации.
func NewPlayer(cfg *config.Config) *Player {
	return &Player{
		Position: core.Vec3{X: 0, Y: 0, Z: 0},
		Velocity: core.Vec3{},
		Balance:  0,
		Width:    cfg.PlayerWidth,
		Height:   cfg.PlayerHeight,
		Depth:    cfg.PlayerDepth,
	}
}
