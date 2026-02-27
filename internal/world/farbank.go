package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/render"
)

type FarBankLayer struct {
	Y      int
	Height int
	Color  color.Color
}

// NewFarBankLayer создаёт слой дальнего берега.
// Параметры: y – вертикальная позиция на экране (верхняя граница), height – высота слоя.
func NewFarBankLayer(y, height int) *FarBankLayer {
	return &FarBankLayer{
		Y:      y,
		Height: height,
		Color:  color.RGBA{34, 139, 34, 255}, // зелёный
	}
}

func (f *FarBankLayer) Update(ctx WorldContext) {
	// Дальний берег статичен, обновление не требуется.
}

func (f *FarBankLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx WorldContext) {
	// Рисуем прямоугольник на весь экран по ширине, начиная с Y, высотой Height.
	ebitenutil.DrawRect(screen, 0, float64(f.Y), cam.ScreenW, float64(f.Height), f.Color)
}
