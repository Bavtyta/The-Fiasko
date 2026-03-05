package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"TheFiaskoTest/internal/common"
	"TheFiaskoTest/internal/render"
)

// SurfaceType определяет тип поверхности сегмента.
type SurfaceType int

const (
	SurfaceSolid  SurfaceType = iota // твёрдая поверхность (бревно, земля)
	SurfaceLiquid                    // жидкость (вода)
)

// SegmentLayer — слой, состоящий из повторяющихся сегментов (река, брёвна и т.п.).
type SegmentLayer struct {
	segments      []*Segment
	speed         float64
	segmentLength float64
	segmentCount  int
	baseX         float64
	baseY         float64
	width         float64
	color         color.Color
	surfaceType   SurfaceType
}

// Getters
func (l *SegmentLayer) Segments() []*Segment {
	return l.segments
}

func (l *SegmentLayer) Speed() float64 {
	return l.speed
}

func (l *SegmentLayer) SegmentLength() float64 {
	return l.segmentLength
}

func (l *SegmentLayer) SegmentCount() int {
	return l.segmentCount
}

func (l *SegmentLayer) BaseX() float64 {
	return l.baseX
}

func (l *SegmentLayer) BaseY() float64 {
	return l.baseY
}

func (l *SegmentLayer) Width() float64 {
	return l.width
}

func (l *SegmentLayer) Color() color.Color {
	return l.color
}

func (l *SegmentLayer) SurfaceType() SurfaceType {
	return l.surfaceType
}

// Setters
func (l *SegmentLayer) SetSegments(segments []*Segment) {
	l.segments = segments
}

func (l *SegmentLayer) SetSpeed(speed float64) {
	l.speed = speed
}

func (l *SegmentLayer) SetSegmentLength(segmentLength float64) {
	l.segmentLength = segmentLength
}

func (l *SegmentLayer) SetSegmentCount(segmentCount int) {
	l.segmentCount = segmentCount
}

func (l *SegmentLayer) SetBaseX(baseX float64) {
	l.baseX = baseX
}

func (l *SegmentLayer) SetBaseY(baseY float64) {
	l.baseY = baseY
}

func (l *SegmentLayer) SetWidth(width float64) {
	l.width = width
}

func (l *SegmentLayer) SetColor(c color.Color) {
	l.color = c
}

func (l *SegmentLayer) SetSurfaceType(surfaceType SurfaceType) {
	l.surfaceType = surfaceType
}

// NewSegmentLayer создаёт новый слой сегментов.
func NewSegmentLayer(baseX, baseY, width, segmentLength float64, speed float64, count int, slopeX, slopeY float64, col color.Color, surfaceType SurfaceType) *SegmentLayer {
	segs := make([]*Segment, count)
	for i := 0; i < count; i++ {
		nearZ := float64(i) * segmentLength
		seg := NewSegment(baseX, baseY, nearZ, width, segmentLength)
		seg.SetSlope(slopeX, slopeY)
		seg.SetColor(col)
		segs[i] = seg
	}
	return &SegmentLayer{
		segments:      segs,
		speed:         speed,
		segmentLength: segmentLength,
		segmentCount:  count,
		baseX:         baseX,
		baseY:         baseY,
		width:         width,
		color:         col,
		surfaceType:   surfaceType,
	}
}

// Update обновляет позиции сегментов и переставляет ушедшие за камеру.
func (l *SegmentLayer) Update(ctx common.WorldContext) {
	totalLength := l.segmentLength * float64(l.segmentCount)
	for _, seg := range l.segments {
		seg.Update(l.speed)
	}
	if len(l.segments) > 0 && l.segments[0].IsBehindCamera() {
		first := l.segments[0]
		first.Wrap(totalLength)
		l.segments = append(l.segments[1:], first)
	}
}

// Draw отрисовывает все сегменты слоя.
func (l *SegmentLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx common.WorldContext) {
	for _, seg := range l.segments {
		seg.Draw(screen, cam)
	}
}

// SurfaceAt возвращает высоту и тип поверхности в заданной мировой координате Z,
// если она находится в пределах какого-либо сегмента слоя.
func (l *SegmentLayer) SurfaceAt(z float64) (height float64, surfaceType SurfaceType, ok bool) {
	for _, seg := range l.segments {
		if z >= seg.NearZ() && z <= seg.NearZ()+seg.Length() {
			// Высота в точке Z: baseY + slopeY * z
			height = seg.BaseY() + seg.SlopeY()*z
			return height, l.surfaceType, true
		}
	}
	return 0, 0, false
}
