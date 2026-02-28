// cmd/game/main.go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"TheFiaskoTest/internal/state"
)

func main() {
	// Создаём менеджер состояний с начальным состоянием Menu
	manager := state.NewManager(nil) // временно nil
	menuState := state.NewMenuState(manager)
	manager.ChangeState(menuState, nil)

	game := &Game{manager: manager}

	ebiten.SetWindowSize(1266, 768)
	ebiten.SetWindowTitle("The Fiasko")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	manager *state.Manager
}

func (g *Game) Update() error {
	return g.manager.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.manager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1266, 768
}
