package world

import (
	"TheFiaskoTest/internal/entity"
	"TheFiaskoTest/internal/render"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type LogLayer struct {
	Segments      []*entity.Segment
	Speed         float64
	SegmentLength float64
	SegmentCount  int
	X             float64
	Y             float64
	Width         float64
}

// NewLogLayer принимает дополнительно slopeX и slopeY для наклона сегментов
func NewLogLayer(x, y, width, segmentLength float64, speed float64, count int, slopeX, slopeY float64) *LogLayer {
	segments := make([]*entity.Segment, count)
	for i := 0; i < count; i++ {
		nearZ := float64(i) * segmentLength
		seg := entity.NewSegment(x, y, nearZ, width, segmentLength)
		seg.SetSlope(slopeX, slopeY)
		seg.Color = color.RGBA{139, 69, 19, 255} // коричневый
		segments[i] = seg
	}
	return &LogLayer{
		Segments:      segments,
		Speed:         speed,
		SegmentLength: segmentLength,
		SegmentCount:  count,
		X:             x,
		Y:             y,
		Width:         width,
	}
}

func (l *LogLayer) Update(ctx WorldContext) {
	totalLength := l.SegmentLength * float64(l.SegmentCount)
	for _, seg := range l.Segments {
		seg.Update(l.Speed)
	}

	if l.Segments[0].IsBehindCamera() {
		first := l.Segments[0]
		first.Wrap(totalLength)
		l.Segments = append(l.Segments[1:], first)
	}
}

func (l *LogLayer) Draw(screen *ebiten.Image, cam *render.Camera, ctx WorldContext) {
	for _, seg := range l.Segments {
		seg.Draw(screen, cam)
	}
}
