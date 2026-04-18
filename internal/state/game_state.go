package state

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"TheFiaskoTest/internal/asset"
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
	logTexture *ebiten.Image // Текстура бревна
}

func NewGameState(manager *Manager, gameCfg config.GameConfig, cameraCfg config.CameraConfig, physicsCfg config.PhysicsConfig) *GameState {
	w := world.New(50)

	skyLayer := world.NewSkyLayer(gameCfg.ScreenWidth, 300, 0.1)
	w.AddLayer(skyLayer)

	// Загружаем текстуру фона и создаём слой дальнего берега
	backgroundTexture := asset.LoadBackgroundTexture()
	farBankLayer := world.NewFarBankLayer(0, gameCfg.ScreenHeight, 0.2) // от 0 до высоты экрана
	farBankLayer.SetTexture(backgroundTexture)
	w.AddLayer(farBankLayer)

	// Загружаем текстуру реки
	riverTexture := asset.LoadRiverTexture()

	riverLayer := world.NewSegmentLayer(
		0, -25, 2000, 100, 0.3, 20,
		0.0, 0.25,
		color.RGBA{0, 100, 255, 255}, world.SurfaceLiquid,
	)
	riverLayer.SetTexture(riverTexture) // Устанавливаем текстуру для реки
	w.AddLayer(riverLayer)

	// Загружаем текстуру бревна
	logTexture := asset.LoadLogTexture()

	logLayer := world.NewSegmentLayer(0, -20, 10, 40, 1.0, 20, 0.0, 0.30, color.RGBA{139, 69, 19, 255}, world.SurfaceSolid)
	logLayer.SetTexture(logTexture) // Устанавливаем текстуру для слоя
	for _, seg := range logLayer.Segments() {
		seg.SetHeight(seg.Width())
		seg.SetRadialSegments(16)
	}
	w.AddLayer(logLayer)

	player := entity.NewPlayer(config.DefaultConfig())

	// Загружаем текстуры игрока (временно не используются в новой архитектуре)
	_ = asset.LoadPlayerTexture()
	_ = asset.LoadPlayerTextureRight()
	_ = asset.LoadPlayerTextureJump()

	balanceBar := ui.NewBalanceBarLayer(
		func() float64 { return player.Balance },
		func() float64 { return config.DefaultConfig().MaxBalance },
		func() bool { return false }, // TODO: implement IsFalling check
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
		logTexture: logTexture,
	}
}

func (g *GameState) Update() error {
	delta := 1.0 / ebiten.ActualTPS()
	if delta == 0 {
		delta = 1.0 / 60
	}
	g.world.Update(delta)

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
	// TODO: Apply balance input using new architecture
	_ = effectiveDrift

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		// TODO: Jump using new architecture
	}

	// TODO: Update player using new architecture
	// g.player.Update(g.world)

	// Проверка столкновений с препятствиями

	// TODO: Check if player is falling using new architecture
	isFalling := false
	if !isFalling {
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
	// TODO: Draw player using new architecture
	// g.player.Draw(screen, g.camera, g.world)

	// TODO: Draw balance bar using new architecture
	// g.balanceBar.Draw(screen, g.camera, g.player.TiltedUpperWorldPos())

	// Отладочная информация
	balanceText := fmt.Sprintf("Balance: %.2f / %.0f", g.player.Balance, config.DefaultConfig().MaxBalance)
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
