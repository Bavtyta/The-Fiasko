package state

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MenuState struct {
	manager   *Manager
	startText string
}

func NewMenuState(manager *Manager) *MenuState {
	return &MenuState{
		manager:   manager,
		startText: "Press ENTER to start",
	}
}

func (m *MenuState) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// Переключаемся на игровое состояние
		gameState := NewGameState(m.manager)
		m.manager.ChangeState(gameState, nil)
	}
	return nil
}

func (m *MenuState) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	ebitenutil.DebugPrintAt(screen, "THE FIASKO", 1266/2-50, 768/2-40)
	ebitenutil.DebugPrintAt(screen, m.startText, 1266/2-70, 768/2)
}

func (m *MenuState) Enter(prevState State, data interface{}) {}
func (m *MenuState) Exit()                                   {}
