package state

import (
	"TheFiaskoTest/internal/config"
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
		gameCfg := m.manager.GameConfig()
		cameraCfg := config.DefaultCameraConfig()
		physicsCfg := config.DefaultPhysicsConfig()
		gameState := NewGameState(m.manager, gameCfg, cameraCfg, physicsCfg)
		m.manager.ChangeState(gameState, nil)
	}
	return nil
}

func (m *MenuState) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	cfg := m.manager.GameConfig()
	ebitenutil.DebugPrintAt(screen, "THE FIASKO", cfg.ScreenWidth/2-50, cfg.ScreenHeight/2-40)
	ebitenutil.DebugPrintAt(screen, m.startText, cfg.ScreenWidth/2-70, cfg.ScreenHeight/2)
}

func (m *MenuState) Enter(prevState State, data interface{}) {}
func (m *MenuState) Exit()                                   {}
