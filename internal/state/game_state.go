package state

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"TheFiaskoTest/internal/entity"
	"TheFiaskoTest/internal/render"
	"TheFiaskoTest/internal/world"
)

const (
	screenWidth    = 1266
	screenHeight   = 768
	driftThreshold = 50.0
)

type GameState struct {
	manager  *Manager
	world    *world.World
	camera   *render.Camera
	player   *entity.Player
	score    float64
	driftDir int // 1 - вправо, -1 - влево
}

func NewGameState(manager *Manager) *GameState {
	w := world.New(0.5)

	// Создаём слои
	skyLayer := world.NewSkyLayer(screenWidth, 300, 0.1)
	w.AddLayer(skyLayer)

	riverLayer := world.NewRiverLayer(2000.0, 0.3, -25, 40.0, 20, 0.0, 0.25)
	w.AddLayer(riverLayer)

	farBankLayer := world.NewFarBankLayer(300, 50)
	w.AddLayer(farBankLayer)

	logLayer := world.NewLogLayer(0, -20, 10, 40, 2.0, 20, 0.0, 0.30)
	w.AddLayer(logLayer)

	// Игрок
	player := entity.NewPlayer(entity.PlayerConfig{
		StartX:        0,
		StartZ:        50,
		Width:         4.0,
		Height:        8.0,
		SurfaceBaseY:  -20,
		SurfaceSlopeY: 0.30,
		BalanceSpeed:  0.2,
	})

	// Слой баланса (отдельный, как ранее)
	balanceLayer := world.NewBalanceBarLayer(
		func() float64 { return player.Balance() },
		func() float64 { return player.MaxBalance() },
		func() bool { return player.IsFalling() },
		screenWidth, screenHeight,
	)
	w.AddLayer(balanceLayer)

	return &GameState{
		manager:  manager,
		world:    w,
		camera:   render.NewCamera(screenWidth, screenHeight),
		player:   player,
		score:    0,
		driftDir: 1,
	}
}

func (g *GameState) Update() error {
	g.world.Update()

	// Обработка дрейфа после порога
	if g.score >= driftThreshold {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.driftDir = -1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.driftDir = 1
		}
	}

	effectiveDrift := 0
	if g.score >= driftThreshold {
		effectiveDrift = g.driftDir
	}
	g.player.ApplyBalanceInput(effectiveDrift)

	// Прыжок
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.player.Jump(2.5)
	}

	g.player.Update(g.world)

	// Начисление очков
	if !g.player.IsFalling() {
		tps := ebiten.ActualTPS()
		if tps == 0 {
			tps = 60
		}
		g.score += 10.0 / tps
	} else {
		// Если игрок упал, переходим в GameOver
		gameOver := NewGameOverState(g.manager, g.score)
		g.manager.ChangeState(gameOver, nil)
	}

	return nil
}

func (g *GameState) Draw(screen *ebiten.Image) {
	g.world.Draw(screen, g.camera)
	g.player.Draw(screen, g.camera, g.world)

	// Отладочная информация (можно оставить)
	balanceText := fmt.Sprintf("Balance: %.2f / %.0f", g.player.Balance(), g.player.MaxBalance())
	ebitenutil.DebugPrintAt(screen, balanceText, 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %.0f", g.score), 10, 30)

	if g.score >= driftThreshold {
		dirStr := "RIGHT"
		if g.driftDir == -1 {
			dirStr = "LEFT"
		}
		ebitenutil.DebugPrintAt(screen, "Drift direction: "+dirStr+" (A/D to change)", 10, 50)
	} else {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Stable (need %.0f points)", driftThreshold), 10, 50)
	}
}

func (g *GameState) Enter(prevState State, data interface{}) {}
func (g *GameState) Exit()                                   {}
