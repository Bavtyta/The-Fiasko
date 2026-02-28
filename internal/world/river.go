package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"TheFiaskoTest/internal/entity"
	"TheFiaskoTest/internal/render"
)

type RiverLayer struct {
	Segments      []*entity.Segment
	Speed         float64
	SegmentLength float64
	SegmentCount  int
	Y             float64 // базовая высота (при Z=0)
	Width         float64
	Color         color.Color
	emptyTex      *ebiten.Image // текстура 1x1 для заливки
}

// NewRiverLayer создаёт слой реки из count сегментов длиной segmentLength,
// шириной width, расположенных на высоте y (ближняя часть), с заданными наклонами.
func NewRiverLayer(width, speed, y, segmentLength float64, count int, slopeX, slopeY float64) *RiverLayer {
	segs := make([]*entity.Segment, count)
	for i := 0; i < count; i++ {
		nearZ := float64(i) * segmentLength
		seg := entity.NewSegment(0, y, nearZ, width, segmentLength)
		seg.SetSlope(slopeX, slopeY)
		seg.Color = color.RGBA{0, 100, 255, 255} // синий (но не используется при заливке)
		segs[i] = seg
	}
	emptyTex := ebiten.NewImage(1, 1)
	emptyTex.Fill(color.White)
	return &RiverLayer{
		Segments:      segs,
		Speed:         speed,
		SegmentLength: segmentLength,
		SegmentCount:  count,
		Y:             y,
		Width:         width,
		Color:         color.RGBA{0, 100, 255, 255},
		emptyTex:      emptyTex,
	}
}

func (r *RiverLayer) Update(ctx WorldContext) {
	totalLength := r.SegmentLength * float64(r.SegmentCount)
	for _, seg := range r.Segments {
		seg.Update(r.Speed)
	}
	if r.Segments[0].IsBehindCamera() {
		first := r.Segments[0]
		first.Wrap(totalLength)
		r.Segments = append(r.Segments[1:], first)
	}
}

func (r *RiverLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx WorldContext) {
	for _, seg := range r.Segments {
		seg.Draw(screen, cam) // рисуем только контур
	}
}
