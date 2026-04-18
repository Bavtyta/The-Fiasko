// Пакет game содержит игровую логику и системы.
package game

import (
	"TheFiaskoTest/internal/entity"
	"math"
)

// CheckAABBCollision проверяет столкновение между игроком и препятствием
// используя метод Axis-Aligned Bounding Box (AABB).
//
// AABB коллизия проверяет пересечение по всем трём осям (X, Y, Z).
// Если объекты пересекаются по всем трём осям одновременно, значит есть столкновение.
//
// Параметры:
//   - player: указатель на игрока
//   - obstacle: указатель на препятствие
//
// Возвращает:
//   - true если обнаружено столкновение, false в противном случае
//
// Требования: 4.7, 4.10
func CheckAABBCollision(player *entity.Player, obstacle *entity.Obstacle) bool {
	// Получаем центры объектов
	px, py, pz := player.Position.X, player.Position.Y, player.Position.Z
	ox, oy, oz := obstacle.Position.X, obstacle.Position.Y, obstacle.Position.Z

	// Проверяем пересечение по каждой оси
	// Объекты пересекаются по оси, если расстояние между центрами
	// меньше суммы половин их размеров
	overlapX := math.Abs(px-ox) < (player.Width+obstacle.Width)/2
	overlapY := math.Abs(py-oy) < (player.Height+obstacle.Height)/2
	overlapZ := math.Abs(pz-oz) < (player.Depth+obstacle.Depth)/2

	// Столкновение есть только если пересечение по всем трём осям
	return overlapX && overlapY && overlapZ
}
