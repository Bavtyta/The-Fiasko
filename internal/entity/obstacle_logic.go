// Пакет entity содержит игровые сущности и их логику.
package entity

// UpdateObstacle обновляет позицию препятствия на основе скорости и dt.
// Препятствия движутся к игроку (уменьшение Z координаты).
// КРИТИЧНО: Реализует требования 1.1, 1.2, 1.4, 1.5, 1.6, 1.7, 3.5
func UpdateObstacle(o *Obstacle, dt float64) {
	// Препятствия движутся к игроку (уменьшение Z)
	o.Position.Z -= o.Velocity.Z * dt

	// Обновляем X и Y позиции (если препятствие имеет скорость по этим осям)
	o.Position.X += o.Velocity.X * dt
	o.Position.Y += o.Velocity.Y * dt
}
