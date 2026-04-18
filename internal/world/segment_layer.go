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
	segments       []*Segment
	parallaxFactor float64
	segmentLength  float64
	segmentCount   int
	baseX          float64
	baseY          float64
	width          float64
	color          color.Color
	surfaceType    SurfaceType
	texture        *ebiten.Image
}

// Getters
func (l *SegmentLayer) Segments() []*Segment {
	return l.segments
}

func (l *SegmentLayer) SurfaceType() SurfaceType {
	return l.surfaceType
}

// Setters
func (l *SegmentLayer) SetTexture(tex *ebiten.Image) {
	l.texture = tex
}

// NewSegmentLayer создаёт новый слой сегментов.
func NewSegmentLayer(baseX, baseY, width, segmentLength float64, parallaxFactor float64, count int, slopeX, slopeY float64, col color.Color, surfaceType SurfaceType) *SegmentLayer {
	segs := make([]*Segment, count)
	for i := 0; i < count; i++ {
		nearZ := float64(i) * segmentLength
		seg := NewSegment(baseX, baseY, nearZ, width, segmentLength)
		seg.SetSlope(slopeX, slopeY)
		seg.SetColor(col)
		segs[i] = seg
	}
	return &SegmentLayer{
		segments:       segs,
		parallaxFactor: parallaxFactor,
		segmentLength:  segmentLength,
		segmentCount:   count,
		baseX:          baseX,
		baseY:          baseY,
		width:          width,
		color:          col,
		surfaceType:    surfaceType,
	}
}

// Update обновляет позиции сегментов и переставляет ушедшие за камеру.
func (l *SegmentLayer) Update(ctx common.WorldContext, delta float64) {
	effectiveSpeed := ctx.GetSpeed() * l.parallaxFactor
	totalLength := l.segmentLength * float64(l.segmentCount)
	for _, seg := range l.segments {
		seg.Update(effectiveSpeed, delta)
	}

	// Переставляем первый сегмент только когда второй сегмент тоже ушёл за камеру
	if len(l.segments) > 0 && l.segments[0].IsBehindCamera() {
		first := l.segments[0]
		first.Wrap(totalLength)
		l.segments = append(l.segments[1:], first)
	}
}

// Draw отрисовывает все сегменты слоя.
func (l *SegmentLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx common.WorldContext) {
	for _, seg := range l.segments {
		seg.Draw(screen, cam, l.texture)
	}
}

// SurfaceAt возвращает высоту и тип поверхности в заданной мировой координате Z,
// если она находится в пределах какого-либо сегмента слоя.
func (l *SegmentLayer) SurfaceAt(z float64) (height float64, surfaceType SurfaceType, ok bool) {
	for _, seg := range l.segments {
		if z >= seg.NearZ() && z <= seg.NearZ()+seg.Length() {
			baseHeight := seg.BaseY() + seg.SlopeY()*z
			if l.surfaceType == SurfaceSolid {
				radius := seg.Width() / 2
				return baseHeight + radius, l.surfaceType, true
			}
			return baseHeight, l.surfaceType, true
		}
	}
	return 0, 0, false
}

// SegmentAt возвращает сегмент слоя, содержащий точку z, или nil.
func (l *SegmentLayer) SegmentAt(z float64) *Segment {
	for _, seg := range l.segments {
		if z >= seg.NearZ() && z <= seg.NearZ()+seg.Length() {
			return seg
		}
	}
	return nil
}
