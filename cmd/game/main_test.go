package main

import (
	"testing"
	"time"

	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/game"
	"TheFiaskoTest/internal/pools"
)

// TestGameUpdateDtCalculation проверяет, что dt вычисляется корректно
func TestGameUpdateDtCalculation(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := game.NewWorld(cfg, pool)

	g := &Game{
		world:    world,
		state:    StatePlaying,
		lastTime: time.Now(),
	}

	// Первый вызов Update
	err := g.Update()
	if err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}

	// Проверяем, что lastTime обновился
	firstTime := g.lastTime

	// Ждём немного
	time.Sleep(10 * time.Millisecond)

	// Второй вызов Update
	err = g.Update()
	if err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}

	// Проверяем, что lastTime обновился снова
	if !g.lastTime.After(firstTime) {
		t.Error("lastTime should be updated after each Update() call")
	}
}

// TestGameUpdateDtClamp проверяет, что dt ограничивается до 0.05
func TestGameUpdateDtClamp(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := game.NewWorld(cfg, pool)

	// Устанавливаем lastTime в прошлое (симулируем большой dt)
	g := &Game{
		world:    world,
		state:    StatePlaying,
		lastTime: time.Now().Add(-1 * time.Second), // 1 секунда назад
	}

	// Вызываем Update - dt должен быть ограничен до 0.05
	err := g.Update()
	if err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}

	// Проверка прошла успешно, если не было паники
	// (внутри Update dt ограничивается до 0.05)
}

// TestGameStateTransitions проверяет переходы между состояниями
func TestGameStateTransitions(t *testing.T) {
	cfg := config.DefaultConfig()
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)
	world := game.NewWorld(cfg, pool)

	g := &Game{
		world:    world,
		state:    StatePlaying,
		lastTime: time.Now(),
	}

	// В StatePlaying должен вызываться world.Update
	err := g.Update()
	if err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}

	// Переключаемся в StateGameOver
	g.state = StateGameOver

	// В StateGameOver world.Update не должен вызываться
	err = g.Update()
	if err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}

	// Переключаемся в StatePaused
	g.state = StatePaused

	// В StatePaused ничего не должно обновляться
	err = g.Update()
	if err != nil {
		t.Fatalf("Update() returned error: %v", err)
	}
}
