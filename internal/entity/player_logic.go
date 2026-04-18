package entity

import (
	"TheFiaskoTest/internal/config"
	"math"
)

// UpdatePlayer обновляет физику игрока
// Применяет гравитацию, обновляет позицию на основе скорости и dt,
// ограничивает позицию (Position.Y не может быть < 0)
func UpdatePlayer(p *Player, config *config.Config, dt float64) {
	// 1. Применяем гравитацию (только если игрок в воздухе)
	if p.Position.Y > 0 {
		p.Velocity.Y -= config.Gravity * dt
	}

	// 2. Обновляем позицию на основе скорости и dt
	p.Position.X += p.Velocity.X * dt
	p.Position.Y += p.Velocity.Y * dt
	p.Position.Z += p.Velocity.Z * dt

	// 3. Ограничиваем позицию (игрок не может быть ниже земли)
	if p.Position.Y < 0 {
		p.Position.Y = 0
		p.Velocity.Y = 0
	}
}

// CheckPlayerFall проверяет условие падения игрока
// Возвращает true, если баланс игрока превышает максимально допустимый
func CheckPlayerFall(p *Player, config *config.Config) bool {
	return math.Abs(p.Balance) > config.MaxBalance
}
