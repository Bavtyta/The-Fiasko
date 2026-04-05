package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/common"
	"TheFiaskoTest/internal/render"
)

type FarBankLayer struct {
	y              int
	height         int
	color          color.Color
	parallaxFactor float64
	texture        *ebiten.Image // текстура фона
}

// Getters
func (f *FarBankLayer) Y() int {
	return f.y
}

func (f *FarBankLayer) Height() int {
	return f.height
}

func (f *FarBankLayer) Color() color.Color {
	return f.color
}

func (f *FarBankLayer) ParallaxFactor() float64 {
	return f.parallaxFactor
}

// Setters
func (f *FarBankLayer) SetY(y int) {
	f.y = y
}

func (f *FarBankLayer) SetHeight(height int) {
	f.height = height
}

func (f *FarBankLayer) SetColor(c color.Color) {
	f.color = c
}

func (f *FarBankLayer) SetParallaxFactor(factor float64) {
	if factor < 0 {
		factor = 0
	}
	f.parallaxFactor = factor
}

// SetTexture устанавливает текстуру для дальнего берега
func (f *FarBankLayer) SetTexture(texture *ebiten.Image) {
	f.texture = texture
}

// NewFarBankLayer создаёт слой дальнего берега.
// Параметры: y – вертикальная позиция на экране (верхняя граница), height – высота слоя, parallaxFactor – множитель параллакса.
func NewFarBankLayer(y, height int, parallaxFactor float64) *FarBankLayer {
	if parallaxFactor < 0 {
		parallaxFactor = 0
	}
	return &FarBankLayer{
		y:              y,
		height:         height,
		color:          color.RGBA{34, 139, 34, 255}, // зелёный
		parallaxFactor: parallaxFactor,
	}
}

func (f *FarBankLayer) Update(ctx common.WorldContext, delta float64) {
	// Дальний берег может иметь минимальное движение для эффекта параллакса
	// effectiveSpeed := ctx.GetSpeed() * f.parallaxFactor * delta
	// В текущей реализации дальний берег статичен
}

func (f *FarBankLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx common.WorldContext) {
	// Если есть текстура, рисуем её
	if f.texture != nil {
		opts := &ebiten.DrawImageOptions{}

		// Масштабируем текстуру, чтобы заполнить всю ширину экрана и нужную высоту
		bounds := f.texture.Bounds()
		scaleX := cam.ScreenW() / float64(bounds.Dx())
		scaleY := float64(f.height) / float64(bounds.Dy())

		opts.GeoM.Scale(scaleX, scaleY)
		opts.GeoM.Translate(0, float64(f.y))

		screen.DrawImage(f.texture, opts)
	} else {
		// Fallback: рисуем прямоугольник цветом
		ebitenutil.DrawRect(screen, 0, float64(f.y), cam.ScreenW(), float64(f.height), f.color)
	}
}
