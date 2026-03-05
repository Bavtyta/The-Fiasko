package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/render"
)

type Segment struct {
	nearZ  float64 // Z ближней (к камере) грани
	baseY  float64 // базовая высота при Z=0 (Y-пересечение)
	x      float64 // центр по X (при Z=0)
	width  float64 // ширина в мировых единицах
	length float64 // длина вдоль Z
	slopeX float64 // изменение X на единицу Z
	slopeY float64 // изменение Y на единицу Z
	color  color.Color
}

// NewSegment создаёт сегмент с заданными параметрами.
// baseY – высота при Z=0; фактическая высота ближней грани = baseY + slopeY * nearZ.
func NewSegment(x, baseY, nearZ, width, length float64) *Segment {
	return &Segment{
		x:      x,
		baseY:  baseY,
		nearZ:  nearZ,
		width:  width,
		length: length,
		slopeX: 0,
		slopeY: 0,
		color:  color.RGBA{255, 255, 255, 255},
	}
}

func (s *Segment) SetSlope(slopeX, slopeY float64) {
	s.slopeX = slopeX
	s.slopeY = slopeY
}

func (s *Segment) Update(speed float64) {
	s.nearZ -= speed
}

func (s *Segment) IsBehindCamera() bool {
	return s.nearZ+s.length <= 0
}

func (s *Segment) Wrap(totalLength float64) {
	s.nearZ += totalLength
}

func (s *Segment) Draw(screen *ebiten.Image, cam *render.Camera) {
	halfW := s.width / 2
	nearZ := s.nearZ
	farZ := nearZ + s.length

	// Вычисляем высоты граней с учётом наклона
	yNear := s.baseY + s.slopeY*nearZ
	yFar := s.baseY + s.slopeY*farZ

	// Координаты углов в мировом пространстве
	corners := []core.Vec3{
		{X: s.x - halfW, Y: yNear, Z: nearZ},                   // левый ближний
		{X: s.x + halfW, Y: yNear, Z: nearZ},                   // правый ближний
		{X: s.x + halfW + s.slopeX*s.length, Y: yFar, Z: farZ}, // правый дальний
		{X: s.x - halfW + s.slopeX*s.length, Y: yFar, Z: farZ}, // левый дальний
	}

	var screenCorners [4]struct{ x, y float64 }
	for i, c := range corners {
		x, y, scale := cam.Project(c)
		if scale <= 0 {
			// Если ближняя грань за камерой – пропускаем сегмент
			if i < 2 && nearZ <= 0 {
				return
			}
			if scale <= 0 {
				// дальний угол может быть сзади – пропускаем, но продолжаем
				continue
			}
		}
		screenCorners[i].x, screenCorners[i].y = x, y
	}

	// Рисуем контур (позже заменим на текстуру)
	col := s.color
	for i := 0; i < 4; i++ {
		next := (i + 1) % 4
		ebitenutil.DrawLine(screen,
			screenCorners[i].x, screenCorners[i].y,
			screenCorners[next].x, screenCorners[next].y,
			col)
	}
}

// Геттеры для всех полей
func (s *Segment) NearZ() float64 {
	return s.nearZ
}

func (s *Segment) BaseY() float64 {
	return s.baseY
}

func (s *Segment) X() float64 {
	return s.x
}

func (s *Segment) Width() float64 {
	return s.width
}

func (s *Segment) Length() float64 {
	return s.length
}

func (s *Segment) SlopeX() float64 {
	return s.slopeX
}

func (s *Segment) SlopeY() float64 {
	return s.slopeY
}

func (s *Segment) Color() color.Color {
	return s.color
}

// Сеттеры для изменяемых полей
func (s *Segment) SetNearZ(nearZ float64) {
	s.nearZ = nearZ
}

func (s *Segment) SetBaseY(baseY float64) {
	s.baseY = baseY
}

func (s *Segment) SetX(x float64) {
	s.x = x
}

func (s *Segment) SetWidth(width float64) {
	s.width = width
}

func (s *Segment) SetLength(length float64) {
	s.length = length
}

func (s *Segment) SetColor(col color.Color) {
	s.color = col
}
