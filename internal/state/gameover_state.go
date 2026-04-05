package state

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"TheFiaskoTest/internal/config"
)

type GameOverState struct {
	manager    *Manager
	score      float64
	gameConfig config.GameConfig
}

func NewGameOverState(manager *Manager, score float64, gameCfg config.GameConfig) *GameOverState {
	return &GameOverState{
		manager:    manager,
		score:      score,
		gameConfig: gameCfg,
	}
}

func (g *GameOverState) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// Перезапуск игры
		cameraCfg := config.DefaultCameraConfig()
		physicsCfg := config.DefaultPhysicsConfig()
		gameState := NewGameState(g.manager, g.gameConfig, cameraCfg, physicsCfg)
		g.manager.ChangeState(gameState, nil)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Выход в главное меню
		mainMenuState := NewMainMenuState(g.manager, g.gameConfig)
		g.manager.ChangeState(mainMenuState, nil)
	}
	return nil
}

func (g *GameOverState) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	msg := fmt.Sprintf("GAME OVER\nScore: %.0f\nPress ENTER to restart\nPress ESC for menu", g.score)
	ebitenutil.DebugPrintAt(screen, msg, g.gameConfig.ScreenWidth/2-100, g.gameConfig.ScreenHeight/2-40)
}

func (g *GameOverState) Enter(prevState State, data interface{}) {}
func (g *GameOverState) Exit()                                   {}
