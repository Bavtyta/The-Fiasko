// cmd/game/main.go
package main

import (
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/game"
	"TheFiaskoTest/internal/pools"
	"TheFiaskoTest/internal/state"
)

func main() {
	// Создаём конфигурацию игры
	gameCfg := config.DefaultGameConfig()

	// Создаём менеджер состояний с начальным состоянием MainMenu
	manager := state.NewManager(nil, gameCfg) // временно nil
	mainMenuState := state.NewMainMenuState(manager, gameCfg)
	manager.ChangeState(mainMenuState, nil)

	gameInstance := &Game{manager: manager}

	ebiten.SetWindowSize(gameCfg.ScreenWidth, gameCfg.ScreenHeight)
	ebiten.SetWindowTitle("The Fiasko")
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}

// GameStateType определяет состояние игры как enum.
// КРИТИЧНО (FIX #1): Добавлен GameStateType enum для управления состояниями игры.
type GameStateType int

const (
	// StateMenu - главное меню
	StateMenu GameStateType = iota
	// StatePlaying - игра в процессе
	StatePlaying
	// StateGameOver - игра окончена
	StateGameOver
	// StatePaused - игра на паузе
	StatePaused
)

// Game - главная структура игры.
// КРИТИЧНО (FIX #1): Game владеет ОБОИМИ: World (logic) И Renderer (presentation).
// World = game logic, Renderer = presentation layer.
type Game struct {
	// world - игровой мир (логика)
	// КРИТИЧНО: Game владеет World (game logic)
	world *game.World

	// renderer - рендерер (отрисовка)
	// КРИТИЧНО: Game владеет Renderer (presentation layer)
	// TODO: Будет добавлен в следующих задачах
	// renderer *render.Renderer

	// state - текущее состояние игры
	// КРИТИЧНО: Добавлен GameStateType enum (StateMenu, StatePlaying, StateGameOver, StatePaused)
	state GameStateType

	// lastTime - время последнего кадра для вычисления dt
	lastTime time.Time

	// manager - временный менеджер состояний (для обратной совместимости)
	// TODO: Будет удалён после полного рефакторинга
	manager *state.Manager
}

// NewGame создаёт новый экземпляр игры.
// КРИТИЧНО: Game владеет World И Renderer (разделение logic и presentation).
func NewGame(cfg *config.Config) *Game {
	// Создаём пул препятствий
	pool := pools.NewObstaclePool(cfg.MaxObstacles, cfg)

	// Создаём World (game logic)
	world := game.NewWorld(cfg, pool)

	// TODO: Создать Renderer (presentation layer) в следующих задачах
	// renderer := render.NewRenderer(cfg, resMgr)

	return &Game{
		world: world,
		// renderer: renderer,  // TODO: Будет добавлен позже
		state:    StateMenu,
		lastTime: time.Now(),
		manager:  nil, // Временно nil, будет удалён
	}
}

// Update обновляет состояние игры.
// КРИТИЧНО (FIX #5): Использует math.Min(dt, 0.05) для ограничения dt.
func (g *Game) Update() error {
	// Вычисляем delta time
	now := time.Now()
	dt := now.Sub(g.lastTime).Seconds()
	g.lastTime = now

	// КРИТИЧНО (FIX #5): Ограничиваем dt для стабильности (0.05 вместо 0.1)
	dt = math.Min(dt, 0.05)

	// Обработка полноэкранного режима
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Обрабатываем в зависимости от состояния
	switch g.state {
	case StateMenu:
		// TODO: menu logic
		// Временно используем старый менеджер состояний
		if g.manager != nil {
			return g.manager.Update()
		}
	case StatePlaying:
		g.world.Update(dt)
		if g.world.IsGameOver() {
			g.state = StateGameOver
		}
	case StateGameOver:
		// TODO: game over logic
		// Временно используем старый менеджер состояний
		if g.manager != nil {
			return g.manager.Update()
		}
	case StatePaused:
		// Ничего не обновляем
	}

	return nil
}

// Draw отрисовывает игру.
// КРИТИЧНО (FIX #1): Game владеет renderer и вызывает его напрямую.
// КРИТИЧНО (FIX #4): Отрисовывает score overlay используя world.score.
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateMenu:
		// TODO: draw menu
		// Временно используем старый менеджер состояний
		if g.manager != nil {
			g.manager.Draw(screen)
		}
	case StatePlaying, StateGameOver:
		// КРИТИЧНО (FIX #1): Game владеет renderer и вызывает его напрямую
		// TODO: Uncomment when Renderer is created (task 17.1)
		// g.renderer.DrawWorld(screen, g.world.player, g.world.obstacles)

		// КРИТИЧНО (FIX #4): Draw score overlay using g.world.score
		// TODO: Uncomment when Renderer is created (task 17.1)
		// ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %.0f", g.world.score))

		// Временно используем старый менеджер состояний
		if g.manager != nil {
			g.manager.Draw(screen)
		}

		if g.state == StateGameOver {
			// TODO: draw game over overlay
		}
	case StatePaused:
		// КРИТИЧНО (FIX #1): Game владеет renderer и вызывает его напрямую
		// TODO: Uncomment when Renderer is created (task 17.1)
		// g.renderer.DrawWorld(screen, g.world.player, g.world.obstacles)

		// TODO: draw pause overlay
		// Временно используем старый менеджер состояний
		if g.manager != nil {
			g.manager.Draw(screen)
		}
	}
}

// Layout возвращает размеры экрана.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// TODO: Использовать config из Game после полного рефакторинга
	if g.manager != nil {
		cfg := g.manager.GameConfig()
		return cfg.ScreenWidth, cfg.ScreenHeight
	}
	return 800, 600 // Дефолтные значения
}
