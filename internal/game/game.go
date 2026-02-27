package game

import (
	"TheFiaskoTest/internal/entity"
	"TheFiaskoTest/internal/render"
	"TheFiaskoTest/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	world  *world.World
	camera *render.Camera
}

func New() *Game {
	w := world.New(0.5)

	// -------- Log layer (unchanged) --------
	logLayer := world.NewLogLayer(500, 350, 266, 20)
	logLayer.AddEntity(entity.NewTestPoint(0, 0, 50))
	w.AddLayer(logLayer)

	// -------- River layer (Y, height, color) --------
	riverLayer := world.NewRiverLayer(
		350, // Y position
		100, // height
		0.5, // water color
	)
	w.AddLayer(riverLayer)

	// -------- Far bank layer (Y, height) – color is not passed --------
	farBankLayer := world.NewFarBankLayer(300, 50)
	// If you need to set a specific color, add a method like:
	// farBankLayer.SetColor(color.RGBA{34, 139, 34, 255})
	w.AddLayer(farBankLayer)

	// -------- Sky layer (Y, height, color) --------
	skyLayer := world.NewSkyLayer(
		1266, // top of screen
		300,  // height
		0.1,
	)
	w.AddLayer(skyLayer)

	return &Game{
		world:  w,
		camera: render.NewCamera(1266, 768),
	}
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Draw(screen, g.camera)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1266, 768
}
