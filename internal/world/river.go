package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"TheFiaskoTest/internal/core"
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
		r.drawFilledSegment(screen, seg, cam)
	}
}

func (r *RiverLayer) drawFilledSegment(screen *ebiten.Image, seg *entity.Segment, cam *render.Camera) {
	halfW := seg.Width / 2
	nearZ := seg.NearZ
	farZ := nearZ + seg.Length

	yNear := seg.BaseY + seg.SlopeY*nearZ
	yFar := seg.BaseY + seg.SlopeY*farZ

	corners := []core.Vec3{
		{X: seg.X - halfW, Y: yNear, Z: nearZ},
		{X: seg.X + halfW, Y: yNear, Z: nearZ},
		{X: seg.X + halfW + seg.SlopeX*seg.Length, Y: yFar, Z: farZ},
		{X: seg.X - halfW + seg.SlopeX*seg.Length, Y: yFar, Z: farZ},
	}

	var pts [4]struct{ x, y float64 }
	for i, c := range corners {
		x, y, scale := render.Project(c, cam)
		if scale <= 0 {
			if i < 2 && nearZ <= 0 {
				return
			}
			return
		}
		pts[i].x, pts[i].y = x, y
	}

	rc, gc, bc, ac := r.Color.RGBA()
	rCol := float32(rc) / 65535
	gCol := float32(gc) / 65535
	bCol := float32(bc) / 65535
	aCol := float32(ac) / 65535

	vertices := []ebiten.Vertex{
		{DstX: float32(pts[0].x), DstY: float32(pts[0].y), SrcX: 0, SrcY: 0, ColorR: rCol, ColorG: gCol, ColorB: bCol, ColorA: aCol},
		{DstX: float32(pts[1].x), DstY: float32(pts[1].y), SrcX: 0, SrcY: 0, ColorR: rCol, ColorG: gCol, ColorB: bCol, ColorA: aCol},
		{DstX: float32(pts[2].x), DstY: float32(pts[2].y), SrcX: 0, SrcY: 0, ColorR: rCol, ColorG: gCol, ColorB: bCol, ColorA: aCol},
		{DstX: float32(pts[3].x), DstY: float32(pts[3].y), SrcX: 0, SrcY: 0, ColorR: rCol, ColorG: gCol, ColorB: bCol, ColorA: aCol},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, r.emptyTex, nil)
}
