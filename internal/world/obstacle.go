package world

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/render"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Инициализация генератора случайных чисел (вызывается один раз при запуске)
func init() {
	rand.Seed(time.Now().UnixNano())
}

type Obstacle struct {
	segment       *Segment
	offsetZ       float64
	width         float64
	height        float64
	rotation      float64
	rotationSpeed float64
	color         color.Color
	radius        float64
}

func NewObstacle(segment *Segment, offsetZ, width, height, rotSpeed float64) *Obstacle {

	// Случайный выбор направления: 1 (по часовой) или -1 (против)
	dir := 1.0
	if rand.Intn(2) == 0 {
		dir = -1.0
	}

	// Случайный множитель скорости от 0.5 до 2.0
	speedFactor := 0.5 + rand.Float64()*1.5
	finalSpeed := rotSpeed * speedFactor * dir

	return &Obstacle{
		segment:       segment,
		offsetZ:       offsetZ,
		width:         width,
		height:        height,
		rotation:      0,
		rotationSpeed: finalSpeed,
		color:         color.RGBA{255, 0, 0, 255},
		radius:        segment.Width() / 2,
	}
}

// Update увеличивает угол поворота на основе прошедшего времени
func (o *Obstacle) Update(delta float64) {
	o.rotation += o.rotationSpeed * delta
}

// WorldPos returns the center position of the obstacle in world coordinates.
// The Z position is calculated as segment.nearZ + offsetZ, ensuring the obstacle
// is positioned within its parent segment's bounds and synchronized with World_Offset_Z.
//
// Requirements: 7.3 (coordinate consistency)
func (o *Obstacle) WorldPos() core.Vec3 {
	z := o.segment.NearZ() + o.offsetZ
	x := o.segment.X() + o.segment.SlopeX()*z
	y := o.segment.BaseY() + o.segment.SlopeY()*z // center of log
	return core.Vec3{X: x, Y: y, Z: z}
}

// Draw отрисовывает препятствие как повёрнутый прямоугольник
func (o *Obstacle) Draw(screen *ebiten.Image, cam *render.Camera) {
	center := o.WorldPos()
	halfW := o.width / 2
	h := o.height
	r := o.radius

	// Локальные координаты: центр вращения в центре бревна, но объект смещён вверх на r
	local := [][2]float64{
		{-halfW, r},     // левый нижний
		{halfW, r},      // правый нижний
		{-halfW, r + h}, // левый верхний
		{halfW, r + h},  // правый верхний
	}

	cosT := math.Cos(o.rotation)
	sinT := math.Sin(o.rotation)

	var worldPts [4]core.Vec3
	for i, l := range local {
		rx := l[0]*cosT + l[1]*sinT
		ry := -l[0]*sinT + l[1]*cosT
		worldPts[i] = core.Vec3{
			X: center.X + rx,
			Y: center.Y + ry,
			Z: center.Z,
		}
	}

	var screenPts [4][2]float64
	for i, wp := range worldPts {
		sx, sy, scale := cam.Project(wp)
		if scale <= 0 {
			return
		}
		screenPts[i] = [2]float64{sx, sy}
	}

	col := o.color
	ebitenutil.DrawLine(screen, screenPts[0][0], screenPts[0][1], screenPts[1][0], screenPts[1][1], col)
	ebitenutil.DrawLine(screen, screenPts[1][0], screenPts[1][1], screenPts[3][0], screenPts[3][1], col)
	ebitenutil.DrawLine(screen, screenPts[3][0], screenPts[3][1], screenPts[2][0], screenPts[2][1], col)
	ebitenutil.DrawLine(screen, screenPts[2][0], screenPts[2][1], screenPts[0][0], screenPts[0][1], col)
}

// Radius возвращает приблизительный радиус для обнаружения столкновений
func (o *Obstacle) Radius() float64 {
	return math.Sqrt(math.Pow(o.width/2, 2) + math.Pow(o.height/2, 2))
}
