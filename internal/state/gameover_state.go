package state

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameOverState struct {
	manager *Manager
	score   float64
}

func NewGameOverState(manager *Manager, score float64) *GameOverState {
	return &GameOverState{
		manager: manager,
		score:   score,
	}
}

func (g *GameOverState) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// Перезапуск игры
		gameState := NewGameState(g.manager)
		g.manager.ChangeState(gameState, nil)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Выход в меню
		menuState := NewMenuState(g.manager)
		g.manager.ChangeState(menuState, nil)
	}
	return nil
}

func (g *GameOverState) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	msg := fmt.Sprintf("GAME OVER\nScore: %.0f\nPress ENTER to restart\nPress ESC for menu", g.score)
	ebitenutil.DebugPrintAt(screen, msg, screenWidth/2-100, screenHeight/2-40)
}

func (g *GameOverState) Enter(prevState State, data interface{}) {}
func (g *GameOverState) Exit()                                   {}
