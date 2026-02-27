package main

import (
	"TheFiaskoTest/internal/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.New()

	ebiten.SetWindowSize(1266, 768)
	ebiten.SetWindowTitle("Log Runner 3D Prototype")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
