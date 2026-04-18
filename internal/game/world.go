// Пакет game содержит игровую логику и оркестрацию.
package game

import (
	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/entity"
	"TheFiaskoTest/internal/pools"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// GameState определяет состояние игры как enum.
// КРИТИЧНО (FIX #4): Добавлен GameState enum для управления состоянием игры.
type GameState int

const (
	// StatePlaying - игра в процессе
	StatePlaying GameState = iota
	// StateGameOver - игра окончена
	StateGameOver
)

// World управляет игровыми объектами как ОРКЕСТРАТОР (НЕ God Object).
// КРИТИЧНО (FIX #1): World НЕ владеет Renderer! Game владеет обоими: World И Renderer.
// World = game logic, Renderer = presentation layer.
type World struct {
	// player - игрок
	player *entity.Player

	// obstacles - активные препятствия
	obstacles []*entity.Obstacle

	// spawner - генератор препятствий
	spawner *Spawner

	// config - конфигурация игры
	config *config.Config

	// pool - пул для переиспользования препятствий
	pool *pools.ObstaclePool

	// state - текущее состояние игры
	// КРИТИЧНО: Добавлен GameState enum (StatePlaying, StateGameOver)
	state GameState

	// score - счёт игрока для отслеживания прогресса
	// КРИТИЧНО (FIX #4): Добавлено поле score для отслеживания прогресса
	score float64

	// КРИТИЧНО (FIX #1): НЕТ поля renderer - World НЕ владеет Renderer!
	// Game владеет обоими: World И Renderer
}

// NewWorld создаёт новый World с заданной конфигурацией и пулом препятствий.
// КРИТИЧНО: НЕТ параметра renderer! Game владеет Renderer, НЕ World.
func NewWorld(cfg *config.Config, pool *pools.ObstaclePool) *World {
	return &World{
		player:    entity.NewPlayer(cfg),
		obstacles: make([]*entity.Obstacle, 0, cfg.MaxObstacles),
		spawner:   NewSpawner(cfg, pool),
		config:    cfg,
		pool:      pool,
		state:     StatePlaying,
		score:     0,
	}
}

// Update обновляет состояние мира (ОРКЕСТРАТОР - делегирует работу отдельным функциям).
// КРИТИЧНО: World ОРКЕСТРИРУЕТ, не выполняет всю логику сам.
func (w *World) Update(dt float64) {
	if w.state != StatePlaying {
		return
	}

	// World ОРКЕСТРИРУЕТ, не выполняет всю логику сам
	// 1. Input handling (КРИТИЧНО: на уровне World, НЕ в entity_logic)
	handleInput(w, dt)

	// 2. Physics updates
	updatePlayer(w, dt)
	updateObstacles(w, dt)

	// 3. Cleanup offscreen obstacles
	removeOffscreenObstacles(w)

	// 4. Spawning
	handleSpawning(w, dt)

	// 5. Collision detection (КРИТИЧНО: с early-out оптимизацией)
	if handleCollisions(w) {
		w.state = StateGameOver
		return
	}

	// 6. Score tracking (КРИТИЧНО: увеличиваем score)
	w.score += dt
}

// handleInput обрабатывает ввод пользователя на уровне World.
// КРИТИЧНО (FIX #2): Input handling на уровне World/Game, НЕ в entity logic.
// Input = внешний источник (клавиатура), НЕ логика entity.
func handleInput(w *World, dt float64) {
	// Балансирование
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		w.player.Balance -= w.config.BalanceSpeed * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		w.player.Balance += w.config.BalanceSpeed * dt
	}

	// Прыжок (только если игрок на земле)
	if ebiten.IsKeyPressed(ebiten.KeyW) && w.player.Velocity.Y == 0 {
		w.player.Velocity.Y = w.config.JumpVelocity
	}
}

// updatePlayer обновляет физику игрока (делегирует UpdatePlayer из entity_logic).
func updatePlayer(w *World, dt float64) {
	entity.UpdatePlayer(w.player, w.config, dt)
}

// updateObstacles обновляет все препятствия (делегирует UpdateObstacle из entity_logic).
func updateObstacles(w *World, dt float64) {
	for _, obs := range w.obstacles {
		entity.UpdateObstacle(obs, dt)
	}
}

// IsGameOver возвращает true, если игра окончена.
func (w *World) IsGameOver() bool {
	return w.state == StateGameOver
}

// handleSpawning обрабатывает генерацию новых препятствий через Spawner.
// КРИТИЧНО: Обрабатывает burst spawn (множественные препятствия при больших dt).
// КРИТИЧНО: Проверяет лимит MaxObstacles для каждого препятствия.
func handleSpawning(w *World, dt float64) {
	// Получаем новые препятствия от spawner (может вернуть несколько при burst spawn)
	newObstacles := w.spawner.Update(dt)

	// КРИТИЧНО: Итерируем по всем возвращённым препятствиям
	for _, obs := range newObstacles {
		// Проверяем лимит MaxObstacles
		if len(w.obstacles) < w.config.MaxObstacles {
			// Лимит не достигнут - добавляем препятствие
			w.obstacles = append(w.obstacles, obs)
		} else {
			// Лимит достигнут - возвращаем препятствие в пул
			w.pool.Put(obs)
		}
	}
}

// removeOffscreenObstacles удаляет препятствия, вышедшие за границы экрана.
// КРИТИЧНО (FIX #6): Проверяет каждое препятствие: если Position.Z < config.DespawnZ, удаляет из списка.
// Возвращает удалённые препятствия в pool.Put().
// Использует эффективный алгоритм удаления (in-place filtering) для избежания аллокаций.
// _Требования: 4.6_
func removeOffscreenObstacles(w *World) {
	// Используем in-place filtering для эффективного удаления
	i := 0
	for _, obs := range w.obstacles {
		// КРИТИЧНО: Используем config.DespawnZ вместо магического числа -10
		if obs.Position.Z < w.config.DespawnZ {
			// Препятствие за границей - возвращаем в пул
			w.pool.Put(obs)
		} else {
			// Препятствие ещё на экране - сохраняем
			w.obstacles[i] = obs
			i++
		}
	}
	// Обрезаем slice до нового размера
	w.obstacles = w.obstacles[:i]
}

// handleCollisions проверяет коллизии между игроком и препятствиями.
// КРИТИЧНО: Использует early-out оптимизацию для пропуска далёких объектов.
// Возвращает true если обнаружена коллизия, false в противном случае.
// _Требования: 4.7_
func handleCollisions(w *World) bool {
	playerZ := w.player.Position.Z

	for _, obs := range w.obstacles {
		// КРИТИЧНО: Early-out оптимизация - пропускаем далёкие объекты
		// Если расстояние по Z больше 5.0, препятствие слишком далеко для коллизии
		if math.Abs(obs.Position.Z-playerZ) > 5.0 {
			continue
		}

		// Проверяем коллизию с близким препятствием
		if CheckAABBCollision(w.player, obs) {
			return true
		}
	}

	return false
}
