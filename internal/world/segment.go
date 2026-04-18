package world

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"TheFiaskoTest/internal/core"
	"TheFiaskoTest/internal/render"
)

type Segment struct {
	nearZ          float64 // Z ближней (к камере) грани
	baseY          float64 // базовая высота при Z=0 (Y-пересечение)
	x              float64 // центр по X (при Z=0)
	width          float64 // ширина в мировых единицах
	length         float64 // длина вдоль Z
	slopeX         float64 // изменение X на единицу Z
	slopeY         float64 // изменение Y на единицу Z
	color          color.Color
	height         float64 // высота бревна (0 для плоского)
	radialSegments int     // количество граней вокруг оси (минимум 3)
}

// NewSegment создаёт сегмент с заданными параметрами.
// baseY – высота при Z=0; фактическая высота ближней грани = baseY + slopeY * nearZ.
func NewSegment(x, baseY, nearZ, width, length float64) *Segment {
	return &Segment{
		x:              x,
		baseY:          baseY,
		nearZ:          nearZ,
		width:          width,
		length:         length,
		slopeX:         0,
		slopeY:         0,
		color:          color.RGBA{255, 255, 255, 255},
		height:         0,
		radialSegments: 1, // 1 означает плоский прямоугольник
	}
}

func (s *Segment) SetSlope(slopeX, slopeY float64) {
	s.slopeX = slopeX
	s.slopeY = slopeY
}

func (s *Segment) Update(speed float64, delta float64) {
	s.nearZ -= speed * delta
}

func (s *Segment) IsBehindCamera() bool {
	return s.nearZ+s.length <= 0
}

func (s *Segment) Wrap(totalLength float64) {
	s.nearZ += totalLength
}

func (s *Segment) Draw(screen *ebiten.Image, cam *render.Camera, texture *ebiten.Image) {
	if s.height == 0 || s.radialSegments <= 1 {
		s.drawFlat(screen, cam, texture)
	} else {
		s.drawCylinder(screen, cam, texture)
	}
}

func (s *Segment) drawFlat(screen *ebiten.Image, cam *render.Camera, texture *ebiten.Image) {
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

	// Если есть текстура, рисуем с текстурой
	if texture != nil {
		bounds := texture.Bounds()
		texW := float32(bounds.Dx())
		texH := float32(bounds.Dy())

		// Создаём вершины для двух треугольников (прямоугольник)
		vertices := []ebiten.Vertex{
			{DstX: float32(screenCorners[0].x), DstY: float32(screenCorners[0].y), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},       // левый ближний
			{DstX: float32(screenCorners[1].x), DstY: float32(screenCorners[1].y), SrcX: texW, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},    // правый ближний
			{DstX: float32(screenCorners[2].x), DstY: float32(screenCorners[2].y), SrcX: texW, SrcY: texH, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1}, // правый дальний
			{DstX: float32(screenCorners[3].x), DstY: float32(screenCorners[3].y), SrcX: 0, SrcY: texH, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},    // левый дальний
		}
		indices := []uint16{0, 1, 2, 0, 2, 3}
		opts := &ebiten.DrawTrianglesOptions{
			Filter: ebiten.FilterLinear,
		}
		screen.DrawTriangles(vertices, indices, texture, opts)
	} else {
		// Рисуем контур (fallback)
		col := s.color
		for i := 0; i < 4; i++ {
			next := (i + 1) % 4
			ebitenutil.DrawLine(screen,
				screenCorners[i].x, screenCorners[i].y,
				screenCorners[next].x, screenCorners[next].y,
				col)
		}
	}
}

func (s *Segment) drawCylinder(screen *ebiten.Image, cam *render.Camera, texture *ebiten.Image) {
	nearZ := s.nearZ
	farZ := nearZ + s.length
	radius := s.width / 2 // предполагаем круглое сечение: радиус = половина ширины
	// если высота задана отдельно, можно использовать s.height/2, но пока считаем круг

	// Центры окружностей на ближнем и дальнем Z с учётом наклонов
	centerNear := core.Vec3{
		X: s.x + s.slopeX*nearZ,
		Y: s.baseY + s.slopeY*nearZ,
		Z: nearZ,
	}
	centerFar := core.Vec3{
		X: s.x + s.slopeX*farZ,
		Y: s.baseY + s.slopeY*farZ,
		Z: farZ,
	}

	N := s.radialSegments
	if N < 3 {
		N = 3
	}

	// Строим вершины для всех углов (направлений)
	type vertex struct {
		world core.Vec3
		u     float64 // текстурная координата вдоль окружности (0..1)
		v     float64 // вдоль длины (0 на nearZ, 1 на farZ)
	}
	ringNear := make([]vertex, N)
	ringFar := make([]vertex, N)

	for i := 0; i < N; i++ {
		phi := 2 * math.Pi * float64(i) / float64(N)
		u := float64(i) / float64(N)
		dir := core.Vec3{X: math.Cos(phi), Y: math.Sin(phi), Z: 0}

		pNear := centerNear.Add(dir.MulScalar(radius))
		ringNear[i] = vertex{world: pNear, u: u, v: 0}

		pFar := centerFar.Add(dir.MulScalar(radius))
		ringFar[i] = vertex{world: pFar, u: u, v: 1}
	}

	// Проецируем все вершины на экран
	type screenVertex struct {
		x, y float64
		u, v float64
	}
	projNear := make([]screenVertex, N)
	projFar := make([]screenVertex, N)

	for i := 0; i < N; i++ {
		x, y, scale := cam.Project(ringNear[i].world)
		if scale <= 0 {
			return // если ближняя точка не видна, весь сегмент пропадает (упрощение)
		}
		projNear[i] = screenVertex{x: x, y: y, u: ringNear[i].u, v: ringNear[i].v}

		x, y, scale = cam.Project(ringFar[i].world)
		if scale <= 0 {
			return
		}
		projFar[i] = screenVertex{x: x, y: y, u: ringFar[i].u, v: ringFar[i].v}
	}

	// Рисуем каждую грань (между i и i+1)
	for i := 0; i < N; i++ {
		next := (i + 1) % N
		v0 := projNear[i]
		v1 := projNear[next]
		v2 := projFar[next]
		v3 := projFar[i]

		// Получаем размеры текстуры
		var texW, texH float32 = 1, 1
		if texture != nil {
			bounds := texture.Bounds()
			texW = float32(bounds.Dx())
			texH = float32(bounds.Dy())
		}

		// Вершины для двух треугольников с правильными текстурными координатами в пикселях
		vertices := []ebiten.Vertex{
			{DstX: float32(v0.x), DstY: float32(v0.y), SrcX: float32(v0.u) * texW, SrcY: float32(v0.v) * texH, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			{DstX: float32(v1.x), DstY: float32(v1.y), SrcX: float32(v1.u) * texW, SrcY: float32(v1.v) * texH, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			{DstX: float32(v2.x), DstY: float32(v2.y), SrcX: float32(v2.u) * texW, SrcY: float32(v2.v) * texH, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
			{DstX: float32(v3.x), DstY: float32(v3.y), SrcX: float32(v3.u) * texW, SrcY: float32(v3.v) * texH, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		}
		indices := []uint16{0, 1, 2, 0, 2, 3}
		opts := &ebiten.DrawTrianglesOptions{
			Filter: ebiten.FilterLinear,
		}

		if texture != nil {
			screen.DrawTriangles(vertices, indices, texture, opts)
		} else {
			// Если текстуры нет, рисуем контур грани (для отладки)
			col := s.color
			ebitenutil.DrawLine(screen, v0.x, v0.y, v1.x, v1.y, col)
			ebitenutil.DrawLine(screen, v1.x, v1.y, v2.x, v2.y, col)
			ebitenutil.DrawLine(screen, v2.x, v2.y, v3.x, v3.y, col)
			ebitenutil.DrawLine(screen, v3.x, v3.y, v0.x, v0.y, col)
		}
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

func (s *Segment) SetColor(col color.Color) {
	s.color = col
}

func (s *Segment) SetHeight(h float64)     { s.height = h }
func (s *Segment) SetRadialSegments(n int) { s.radialSegments = n }
