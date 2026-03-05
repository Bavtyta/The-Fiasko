package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// BalanceBarLayer отображает полоску баланса игрока.
type BalanceBarLayer struct {
	// Для доступа к данным игрока можно передать функции обратного вызова
	// или интерфейс. Проще всего передать ссылку на игрока, но это создаст
	// циклическую зависимость, если игрок определён в пакете entity.
	// Поэтому будем передавать функции получения текущего баланса и статуса падения.
	getBalance    func() float64
	getMaxBalance func() float64
	isFalling     func() bool
	screenWidth   float64
	screenHeight  float64
}

// NewBalanceBarLayer создаёт слой баланса.
// Принимает функции для получения данных от игрока.
func NewBalanceBarLayer(getBalance, getMaxBalance func() float64, isFalling func() bool, screenW, screenH float64) *BalanceBarLayer {
	return &BalanceBarLayer{
		getBalance:    getBalance,
		getMaxBalance: getMaxBalance,
		isFalling:     isFalling,
		screenWidth:   screenW,
		screenHeight:  screenH,
	}
}

// Update обновляет состояние UI (можно добавить анимацию в будущем).
func (b *BalanceBarLayer) Update() {
	// можно добавить анимацию в будущем
}

// Draw рисует полоску баланса.
func (b *BalanceBarLayer) Draw(screen *ebiten.Image) {
	if b.isFalling() {
		return
	}

	barY := b.screenHeight/2 - 100
	barHeight := 20.0
	maxBarLength := 200.0
	centerX := b.screenWidth / 2

	balance := b.getBalance()
	maxBal := b.getMaxBalance()

	if balance >= 0 {
		length := (balance / maxBal) * maxBarLength
		ebitenutil.DrawRect(screen, centerX, barY, length, barHeight, color.RGBA{0, 255, 0, 255})
	} else {
		length := (-balance / maxBal) * maxBarLength
		ebitenutil.DrawRect(screen, centerX-length, barY, length, barHeight, color.RGBA{255, 0, 0, 255})
	}
}
