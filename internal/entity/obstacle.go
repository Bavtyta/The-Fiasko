// Пакет entity содержит игровые сущности (только данные, без логики).
package entity

import "TheFiaskoTest/internal/core"

// ObstacleType определяет тип препятствия как enum для compile-time проверки типов.
// КРИТИЧНО (FIX #3): Используется enum вместо string для лучшей производительности
// и безопасности типов.
type ObstacleType int

const (
	// ObstacleTypeLog представляет препятствие типа "бревно"
	ObstacleTypeLog ObstacleType = iota
	// ObstacleTypeRock представляет препятствие типа "камень"
	ObstacleTypeRock
	// Будущие типы препятствий можно добавить здесь
)

// Obstacle представляет препятствие в игре (data-only структура).
// КРИТИЧНО: Содержит ТОЛЬКО поля данных, БЕЗ методов (кроме Reset).
// Логика обновления находится в отдельных функциях (obstacle_logic.go).
type Obstacle struct {
	// Position - позиция препятствия в мировых координатах
	Position core.Vec3
	// Velocity - скорость движения препятствия
	Velocity core.Vec3
	// Width - ширина препятствия
	Width float64
	// Height - высота препятствия
	Height float64
	// Depth - глубина препятствия
	Depth float64
	// Type - тип препятствия (enum, НЕ string)
	// КРИТИЧНО (FIX #3): ObstacleType enum для compile-time проверки
	Type ObstacleType
	// КРИТИЧНО: НЕТ active флага - состояние контролируется присутствием в world.obstacles
}

// Reset полностью сбрасывает все поля препятствия к начальным значениям.
// Используется при возврате препятствия в ObjectPool для предотвращения
// утечки состояния между использованиями.
func (o *Obstacle) Reset() {
	o.Position = core.Vec3{}
	o.Velocity = core.Vec3{}
	o.Width = 0
	o.Height = 0
	o.Depth = 0
	o.Type = ObstacleTypeLog // Сброс в дефолтное значение
}
