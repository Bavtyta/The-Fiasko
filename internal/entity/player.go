package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/common"
	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/render"
	"TheFiaskoTest/internal/world"
)

// PlayerConfig содержит параметры для создания игрока.
type PlayerConfig struct {
	StartX, StartZ, Width, Height float64
	BalanceSpeed                  float64
	Physics                       config.PhysicsConfig
}

// Player представляет игрового персонажа в трёхмерном мире.
type Player struct {
	world        *world.World
	position     core.Vec3
	width        float64
	height       float64
	color        color.Color
	isJumping    bool
	jumpVelocity float64
	groundY      float64
	balance      float64
	maxBalance   float64
	balanceSpeed float64
	isFalling    bool
	physics      config.PhysicsConfig
}

// NewPlayer создаёт игрока на основе конфигурации.
func NewPlayer(world *world.World, cfg PlayerConfig) *Player {
	// Начальная высота будет определена в первом Update через GetSurfaceAt,
	// поэтому здесь можно временно установить 0.
	return &Player{
		world:        world,
		position:     core.Vec3{X: cfg.StartX, Y: 0, Z: cfg.StartZ},
		width:        cfg.Width,
		height:       cfg.Height,
		color:        color.White,
		groundY:      0,
		balance:      0,
		maxBalance:   10,
		balanceSpeed: cfg.BalanceSpeed,
		isFalling:    false,
		physics:      cfg.Physics,
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
func (p *Player) Update(ctx common.WorldContext) {
	if p.isFalling {
		return
	}

	// Получаем информацию о поверхности под игроком
	z := p.position.Z
	if height, surfaceType, ok := p.world.GetSurfaceAt(z); ok {
		p.groundY = height
		if surfaceType == world.SurfaceLiquid {
			p.isFalling = true // падение в воду
		}
	} else {
		// Если нет поверхности (например, игрок улетел за пределы) — падение
		p.isFalling = true
	}

	if p.isJumping {
		p.position.Y += p.jumpVelocity
		p.jumpVelocity -= p.physics.Gravity
		if p.position.Y <= p.groundY {
			p.position.Y = p.groundY
			p.isJumping = false
			p.jumpVelocity = 0
		}
	} else {
		p.position.Y = p.groundY
	}
}

// Jump инициирует прыжок, если игрок на земле.
func (p *Player) Jump(initialVelocity float64) {
	if !p.isJumping && p.position.Y <= p.groundY+0.01 {
		p.isJumping = true
		p.jumpVelocity = initialVelocity
	}
}

// Draw отрисовывает игрока с учётом камеры.
func (p *Player) Draw(screen *ebiten.Image, cam *render.Camera, ctx common.WorldContext) {
	// Корректируем позицию с учётом смещения мира
	relative := p.position

	bottom := core.Vec3{X: relative.X, Y: relative.Y, Z: relative.Z}
	top := core.Vec3{X: relative.X, Y: relative.Y + p.height, Z: relative.Z}

	bx, by, bScale := cam.Project(bottom)
	_, ty, tScale := cam.Project(top)

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
