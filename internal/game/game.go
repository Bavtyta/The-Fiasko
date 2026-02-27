package game

import (
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

	// -------- Sky layer (Y, height, color) --------
	skyLayer := world.NewSkyLayer(
		1266, // top of screen
		385,  // height
		0.1,
	)
	w.AddLayer(skyLayer)

	// logLayer.AddEntity(entity.NewTestPoint(-15, 0, 50))
	//

	// -------- River layer (Y, height, color) --------
	riverLayer := world.NewRiverLayer(
		2000.0, // ширина (достаточно для заполнения экрана)
		0.3,    // скорость (медленнее бревна)
		-25,    // высота (ниже бревна)
		40.0,   // длина сегмента
		20,     // количество сегментов
		0.0,    // наклон по X (как у бревна, если нужно)
		0.25,   // наклон по Y (параллельно бревну, если у бревна slopeY=0.05)
	)
	w.AddLayer(riverLayer)

	// -------- Far bank layer (Y, height) – color is not passed --------
	farBankLayer := world.NewFarBankLayer(300, 50)
	// If you need to set a specific color, add a method like:
	// farBankLayer.SetColor(color.RGBA{34, 139, 34, 255})
	w.AddLayer(farBankLayer)

	// -------- Log layer (unchanged) --------
	logLayer := world.NewLogLayer(
		0,    // центр по X
		-20,  // базовая высота Y
		40,   // ширина сегмента
		40,   // длина сегмента
		2.0,  // скорость движения
		20,   // количество сегментов
		0.0,  // наклон по X (нет бокового смещения)
		0.30, // наклон по Y: дальняя часть выше на 0.2 * 500 = 100 единиц
	)
	w.AddLayer(logLayer)

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
