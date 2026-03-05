package common

// WorldContext предоставляет информацию о состоянии игрового мира,
// необходимую сущностям и слоям для обновления и отрисовки.
type WorldContext interface {
	// GetSpeed возвращает текущую скорость движения мира.
	GetSpeed() float64

	// GetWorldOffsetZ возвращает смещение мира по оси Z.
	GetWorldOffsetZ() float64
}
