package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/render"
)

type Segment struct {
	NearZ  float64 // Z ближней (к камере) грани
	BaseY  float64 // базовая высота при Z=0 (Y-пересечение)
	X      float64 // центр по X (при Z=0)
	Width  float64 // ширина в мировых единицах
	Length float64 // длина вдоль Z
	SlopeX float64 // изменение X на единицу Z
	SlopeY float64 // изменение Y на единицу Z
	Color  color.Color
}

// NewSegment создаёт сегмент с заданными параметрами.
// baseY – высота при Z=0; фактическая высота ближней грани = baseY + slopeY * nearZ.
func NewSegment(x, baseY, nearZ, width, length float64) *Segment {
	return &Segment{
		X:      x,
		BaseY:  baseY,
		NearZ:  nearZ,
		Width:  width,
		Length: length,
		SlopeX: 0,
		SlopeY: 0,
		Color:  color.RGBA{255, 255, 255, 255},
	}
}

func (s *Segment) SetSlope(slopeX, slopeY float64) {
	s.SlopeX = slopeX
	s.SlopeY = slopeY
}

func (s *Segment) Update(speed float64) {
	s.NearZ -= speed
}

func (s *Segment) IsBehindCamera() bool {
	return s.NearZ+s.Length <= 0
}

func (s *Segment) Wrap(totalLength float64) {
	s.NearZ += totalLength
}

func (s *Segment) Draw(screen *ebiten.Image, cam *render.Camera) {
	halfW := s.Width / 2
	nearZ := s.NearZ
	farZ := nearZ + s.Length

	// Вычисляем высоты граней с учётом наклона
	yNear := s.BaseY + s.SlopeY*nearZ
	yFar := s.BaseY + s.SlopeY*farZ

	// Координаты углов в мировом пространстве
	corners := []core.Vec3{
		{X: s.X - halfW, Y: yNear, Z: nearZ},                   // левый ближний
		{X: s.X + halfW, Y: yNear, Z: nearZ},                   // правый ближний
		{X: s.X + halfW + s.SlopeX*s.Length, Y: yFar, Z: farZ}, // правый дальний
		{X: s.X - halfW + s.SlopeX*s.Length, Y: yFar, Z: farZ}, // левый дальний
	}

	var screenCorners [4]struct{ x, y float64 }
	for i, c := range corners {
		x, y, scale := render.Project(c, cam)
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
	col := s.Color
	for i := 0; i < 4; i++ {
		next := (i + 1) % 4
		ebitenutil.DrawLine(screen,
			screenCorners[i].x, screenCorners[i].y,
			screenCorners[next].x, screenCorners[next].y,
			col)
	}
}
