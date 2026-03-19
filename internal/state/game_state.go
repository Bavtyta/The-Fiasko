package state

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/entity"
	"TheFiaskoTest/internal/render"
	"TheFiaskoTest/internal/ui"
	"TheFiaskoTest/internal/world"
)

type GameState struct {
	manager    *Manager
	world      *world.World
	camera     *render.Camera
	player     *entity.Player
	balanceBar *ui.BalanceBarLayer
	score      float64
	driftDir   int // 1 - вправо, -1 - влево
	gameConfig config.GameConfig
}

func NewGameState(manager *Manager, gameCfg config.GameConfig, cameraCfg config.CameraConfig, physicsCfg config.PhysicsConfig) *GameState {
	w := world.New(0.5)

	skyLayer := world.NewSkyLayer(gameCfg.ScreenWidth, 300, 0.1)
	w.AddLayer(skyLayer)

	farBankLayer := world.NewFarBankLayer(300, 50)
	w.AddLayer(farBankLayer)

	logLayer := world.NewSegmentLayer(
		0, -20, 7, 40, 2.0, 20,
		0.0, 0.30,
		color.RGBA{139, 69, 19, 255}, world.SurfaceSolid,
	)
	w.AddLayer(logLayer)

	riverLayer := world.NewSegmentLayer(
		0, -25, 2000, 40, 0.3, 20,
		0.0, 0.25,
		color.RGBA{0, 100, 255, 255}, world.SurfaceLiquid,
	)
	w.AddLayer(riverLayer)

	player := entity.NewPlayer(w, entity.PlayerConfig{
		StartX:       0,
		StartZ:       50,
		Width:        4.0,
		Height:       8.0,
		BalanceSpeed: 0.2,
		Physics:      physicsCfg,
		MaxTiltAngle: 0.8,
	})

	balanceBar := ui.NewBalanceBarLayer(
		func() float64 { return player.Balance() },
		func() float64 { return player.MaxBalance() },
		func() bool { return player.IsFalling() },
	)

	return &GameState{
		manager:    manager,
		world:      w,
		camera:     render.NewCamera(float64(gameCfg.ScreenWidth), float64(gameCfg.ScreenHeight), cameraCfg),
		player:     player,
		balanceBar: balanceBar,
		score:      0,
		driftDir:   1,
		gameConfig: gameCfg,
	}
}

func (g *GameState) Update() error {
	g.world.Update()

	if g.score >= g.gameConfig.DriftThreshold {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.driftDir = -1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.driftDir = 1
		}
	}

	effectiveDrift := 0
	if g.score >= g.gameConfig.DriftThreshold {
		effectiveDrift = g.driftDir
	}
	g.player.ApplyBalanceInput(effectiveDrift)

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.player.Jump(2.5)
	}

	g.player.Update(g.world)

	if !g.player.IsFalling() {
		tps := ebiten.ActualTPS()
		if tps == 0 {
			tps = 60
		}
		g.score += 10.0 / tps
	} else {
		gameOver := NewGameOverState(g.manager, g.score, g.gameConfig)
		g.manager.ChangeState(gameOver, nil)
	}

	return nil
}

func (g *GameState) Draw(screen *ebiten.Image) {
	g.world.Draw(screen, g.camera)
	g.player.Draw(screen, g.camera, g.world)

	// Баланс-бар над игроком
	g.balanceBar.Draw(screen, g.camera, g.player.TiltedUpperWorldPos())

	// Отладочная информация
	balanceText := fmt.Sprintf("Balance: %.2f / %.0f", g.player.Balance(), g.player.MaxBalance())
	ebitenutil.DebugPrintAt(screen, balanceText, 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %.0f", g.score), 10, 30)

	if g.score >= g.gameConfig.DriftThreshold {
		dirStr := "RIGHT"
		if g.driftDir == -1 {
			dirStr = "LEFT"
		}
		ebitenutil.DebugPrintAt(screen, "Drift direction: "+dirStr+" (A/D to change)", 10, 50)
	} else {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Stable (need %.0f points)", g.gameConfig.DriftThreshold), 10, 50)
	}
}

func (g *GameState) Enter(prevState State, data interface{}) {}
func (g *GameState) Exit()                                   {}
