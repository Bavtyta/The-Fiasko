package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/render"
)

type BalanceBarLayer struct {
	getBalance    func() float64
	getMaxBalance func() float64
	isFalling     func() bool
}

func NewBalanceBarLayer(getBalance, getMaxBalance func() float64, isFalling func() bool) *BalanceBarLayer {
	return &BalanceBarLayer{
		getBalance:    getBalance,
		getMaxBalance: getMaxBalance,
		isFalling:     isFalling,
	}
}

func (b *BalanceBarLayer) Update() {}

// Draw рисует горизонтальную полоску баланса над центром верхней грани игрока.
// Полоска растёт от центра в сторону наклона.
func (b *BalanceBarLayer) Draw(screen *ebiten.Image, cam *render.Camera, upperCenter core.Vec3) {
	if b.isFalling() {
		return
	}

	// Проецируем центр верхней грани
	sx, sy, scale := cam.Project(upperCenter)
	if scale <= 0 {
		return
	}

	// Отступ вверх от спроецированной точки
	const barOffsetY = 25
	barY := sy - barOffsetY
	barHeight := 20.0
	maxBarLength := 100.0 // половина длины полоски (от центра в одну сторону)

	balance := b.getBalance()
	maxBal := b.getMaxBalance()

	if balance >= 0 {
		length := (balance / maxBal) * maxBarLength
		// Рисуем зелёную часть справа от центра
		ebitenutil.DrawRect(screen, sx, barY, length, barHeight, color.RGBA{0, 255, 0, 255})
	} else {
		length := (-balance / maxBal) * maxBarLength
		// Рисуем красную часть слева от центра
		ebitenutil.DrawRect(screen, sx-length, barY, length, barHeight, color.RGBA{255, 0, 0, 255})
	}
}
