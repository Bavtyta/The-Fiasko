package config

// Config содержит все настройки игры
type Config struct {
	// Игра
	ScreenWidth  int
	ScreenHeight int
	TargetFPS    int

	// Физика
	Gravity      float64
	JumpVelocity float64

	// Игрок
	BalanceSpeed float64
	MaxBalance   float64
	PlayerWidth  float64
	PlayerHeight float64
	PlayerDepth  float64

	// Спавн (КРИТИЧНО: добавлены SpawnRangeX, SpawnZ, DespawnZ)
	SpawnInterval  float64
	SpawnRangeX    float64 // Диапазон X для спавна (НЕ хардкод)
	SpawnZ         float64 // Z позиция спавна (НЕ хардкод)
	ObstacleSpeed  float64
	MaxObstacles   int
	DespawnZ       float64 // Z позиция удаления (вместо магического -10)
	ObstacleWidth  float64 // Ширина препятствия
	ObstacleHeight float64 // Высота препятствия
	ObstacleDepth  float64 // Глубина препятствия

	// Камера
	CameraFocalLength float64
	CameraHorizonY    float64
}

// DefaultConfig возвращает конфигурацию с разумными значениями по умолчанию
func DefaultConfig() *Config {
	return &Config{
		// Игра
		ScreenWidth:  800,
		ScreenHeight: 600,
		TargetFPS:    60,

		// Физика
		Gravity:      20.0,
		JumpVelocity: 10.0,

		// Игрок
		BalanceSpeed: 5.0,
		MaxBalance:   10.0,
		PlayerWidth:  2.0,
		PlayerHeight: 3.0,
		PlayerDepth:  1.0,

		// Спавн
		SpawnInterval:  2.0,
		SpawnRangeX:    10.0, // Диапазон X для спавна
		SpawnZ:         50.0, // Z позиция спавна
		ObstacleSpeed:  5.0,
		MaxObstacles:   20,
		DespawnZ:       -10.0, // Z позиция удаления
		ObstacleWidth:  2.0,
		ObstacleHeight: 2.0,
		ObstacleDepth:  2.0,

		// Камера
		CameraFocalLength: 300.0,
		CameraHorizonY:    300.0,
	}
}
