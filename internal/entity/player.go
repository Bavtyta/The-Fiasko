package entity

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/common"
	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/render"
	"TheFiaskoTest/internal/world"
)

type PlayerConfig struct {
	StartX, StartZ, Width, Height float64
	BalanceSpeed                  float64
	Physics                       config.PhysicsConfig
	MaxTiltAngle                  float64 // максимальный угол наклона в радианах
}

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
	maxTiltAngle float64
}

func NewPlayer(world *world.World, cfg PlayerConfig) *Player {
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
		maxTiltAngle: cfg.MaxTiltAngle,
	}
}

func (p *Player) Update(ctx common.WorldContext) {
	if p.isFalling {
		return
	}
	z := p.position.Z
	if height, surfaceType, ok := p.world.GetSurfaceAt(z); ok {
		p.groundY = height
		if surfaceType == world.SurfaceLiquid {
			p.isFalling = true
		}
	} else {
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

func (p *Player) Jump(initialVelocity float64) {
	if !p.isJumping && p.position.Y <= p.groundY+0.01 {
		p.isJumping = true
		p.jumpVelocity = initialVelocity
	}
}

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

// Draw отрисовывает игрока с учётом наклона (вращение вокруг центра нижней стороны).
func (p *Player) Draw(screen *ebiten.Image, cam *render.Camera, ctx common.WorldContext) {
	if p.isFalling {
		// TODO: анимация падения
		return
	}

	factor := p.balance / p.maxBalance
	if factor < -1 {
		factor = -1
	} else if factor > 1 {
		factor = 1
	}
	theta := factor * p.maxTiltAngle // положительный угол — наклон вправо

	halfW := p.width / 2
	h := p.height

	// Локальные координаты четырёх углов (центр нижней стороны в (0,0))
	local := [][2]float64{
		{-halfW, 0}, // левый нижний
		{halfW, 0},  // правый нижний
		{-halfW, h}, // левый верхний
		{halfW, h},  // правый верхний
	}

	cosT := math.Cos(theta)
	sinT := math.Sin(theta)

	var worldPts [4]core.Vec3
	for i, l := range local {
		// Поворот по часовой стрелке
		rx := l[0]*cosT + l[1]*sinT
		ry := -l[0]*sinT + l[1]*cosT
		worldPts[i] = core.Vec3{
			X: p.position.X + rx,
			Y: p.position.Y + ry,
			Z: p.position.Z,
		}
	}

	// Проецируем на экран
	var screenPts [4][2]float64
	for i, wp := range worldPts {
		sx, sy, scale := cam.Project(wp)
		if scale <= 0 {
			return // если хоть одна точка за камерой, не рисуем
		}
		screenPts[i] = [2]float64{sx, sy}
	}

	col := p.color
	// Левый нижний -> правый нижний
	ebitenutil.DrawLine(screen, screenPts[0][0], screenPts[0][1], screenPts[1][0], screenPts[1][1], col)
	// Правый нижний -> правый верхний
	ebitenutil.DrawLine(screen, screenPts[1][0], screenPts[1][1], screenPts[3][0], screenPts[3][1], col)
	// Правый верхний -> левый верхний
	ebitenutil.DrawLine(screen, screenPts[3][0], screenPts[3][1], screenPts[2][0], screenPts[2][1], col)
	// Левый верхний -> левый нижний
	ebitenutil.DrawLine(screen, screenPts[2][0], screenPts[2][1], screenPts[0][0], screenPts[0][1], col)
}

// TiltedUpperWorldPos возвращает мировые координаты центра верхней грани после наклона.
// Используется для привязки баланс-бара.
func (p *Player) TiltedUpperWorldPos() core.Vec3 {
	factor := p.balance / p.maxBalance
	if factor < -1 {
		factor = -1
	} else if factor > 1 {
		factor = 1
	}
	theta := factor * p.maxTiltAngle
	cosT := math.Cos(theta)
	sinT := math.Sin(theta)
	return core.Vec3{
		X: p.position.X + p.height*sinT,
		Y: p.position.Y + p.height*cosT,
		Z: p.position.Z,
	}
}

// UpperWorldPos возвращает центр верхней грани без наклона (если нужно).
func (p *Player) UpperWorldPos() core.Vec3 {
	return core.Vec3{
		X: p.position.X,
		Y: p.position.Y + p.height,
		Z: p.position.Z,
	}
}

// GetZ возвращает глубину для сортировки.
func (p *Player) GetZ() float64 {
	return p.position.Z
}

func (p *Player) SetZ(z float64) {
	p.position.Z = z
}

func (p *Player) Balance() float64    { return p.balance }
func (p *Player) MaxBalance() float64 { return p.maxBalance }
func (p *Player) IsFalling() bool     { return p.isFalling }
func (p *Player) Position() core.Vec3 { return p.position }
