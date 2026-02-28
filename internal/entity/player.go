package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/render"
)

// PlayerConfig содержит параметры для создания игрока.
type PlayerConfig struct {
	StartX, StartZ, Width, Height float64
	SurfaceBaseY, SurfaceSlopeY   float64
	BalanceSpeed                  float64
}

// Player представляет игрового персонажа в трёхмерном мире.
type Player struct {
	position      core.Vec3
	width         float64
	height        float64
	color         color.Color
	surfaceBaseY  float64
	surfaceSlopeY float64
	isJumping     bool
	jumpVelocity  float64
	groundY       float64
	balance       float64
	maxBalance    float64
	balanceSpeed  float64
	isFalling     bool
}

// NewPlayer создаёт игрока на основе конфигурации.
func NewPlayer(cfg PlayerConfig) *Player {
	groundY := cfg.SurfaceBaseY + cfg.SurfaceSlopeY*cfg.StartZ
	return &Player{
		position:      core.Vec3{X: cfg.StartX, Y: groundY, Z: cfg.StartZ},
		width:         cfg.Width,
		height:        cfg.Height,
		color:         color.White,
		surfaceBaseY:  cfg.SurfaceBaseY,
		surfaceSlopeY: cfg.SurfaceSlopeY,
		groundY:       groundY,
		balance:       0,
		maxBalance:    10,
		balanceSpeed:  cfg.BalanceSpeed,
		isFalling:     false,
	}
}

// Геттеры
func (p *Player) Position() core.Vec3 { return p.position }
func (p *Player) Width() float64      { return p.width }
func (p *Player) Height() float64     { return p.height }
func (p *Player) Color() color.Color  { return p.color }
func (p *Player) Balance() float64    { return p.balance }
func (p *Player) MaxBalance() float64 { return p.maxBalance }
func (p *Player) IsFalling() bool     { return p.isFalling }
func (p *Player) IsJumping() bool     { return p.isJumping }
func (p *Player) GroundY() float64    { return p.groundY }

// Update обновляет состояние игрока (прыжок, привязку к поверхности).
func (p *Player) Update(ctx WorldContext) {
	if p.isFalling {
		return
	}
	// Вычисляем высоту поверхности в текущей позиции Z
	newGroundY := p.surfaceBaseY + p.surfaceSlopeY*p.position.Z

	const gravity = 0.2
	if p.isJumping {
		p.position.Y += p.jumpVelocity
		p.jumpVelocity -= gravity
		if p.position.Y <= newGroundY {
			p.position.Y = newGroundY
			p.isJumping = false
			p.jumpVelocity = 0
		}
	} else {
		p.position.Y = newGroundY
	}
	p.groundY = newGroundY
}

// Jump инициирует прыжок, если игрок на земле.
func (p *Player) Jump(initialVelocity float64) {
	if !p.isJumping && p.position.Y <= p.groundY+0.01 {
		p.isJumping = true
		p.jumpVelocity = initialVelocity
	}
}

// Draw отрисовывает игрока с учётом камеры.
func (p *Player) Draw(screen *ebiten.Image, cam *render.Camera, ctx WorldContext) {
	// Корректируем позицию с учётом смещения мира
	relative := p.position

	bottom := core.Vec3{X: relative.X, Y: relative.Y, Z: relative.Z}
	top := core.Vec3{X: relative.X, Y: relative.Y + p.height, Z: relative.Z}

	bx, by, bScale := render.Project(bottom, cam)
	_, ty, tScale := render.Project(top, cam)

	if bScale <= 0 || tScale <= 0 {
		return
	}

	scale := (bScale + tScale) / 2
	screenWidth := p.width * scale
	halfW := screenWidth / 2

	ebitenutil.DrawRect(screen, bx-halfW, ty, screenWidth, by-ty, p.color)
}

// ApplyBalanceInput изменяет баланс в зависимости от направления дрейфа.
func (p *Player) ApplyBalanceInput(driftDir int) {
	if p.isFalling {
		return
	}
	delta := p.balanceSpeed * float64(driftDir)
	p.balance += delta
	if p.balance > p.maxBalance {
		p.balance = p.maxBalance
		p.isFalling = true
	} else if p.balance < -p.maxBalance {
		p.balance = -p.maxBalance
		p.isFalling = true
	}
}

// GetZ возвращает глубину для сортировки.
func (p *Player) GetZ() float64 {
	return p.position.Z
}

// SetZ устанавливает глубину.
func (p *Player) SetZ(z float64) {
	p.position.Z = z
}
