package ui

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// FloatingLoupe mengelola status dan render kaca pembesar melayang di atas komik
type FloatingLoupe struct {
	Active     bool
	Position   f32.Point // Posisi kursor/lensa di layar TV (X, Y)
	Radius     float32   // Jari-jari lingkaran lensa (misal: 120 pixel)
	ZoomFactor float32   // Pengali zoom di dalam lensa (misal: 2.2x)
	MoveSpeed  float32   // Kecepatan gerak pointer saat tombol D-pad ditekan
}

func NewFloatingLoupe() *FloatingLoupe {
	return &FloatingLoupe{
		Active:     false,
		Position:   f32.Pt(400, 300),
		Radius:     130,
		ZoomFactor: 2.2,
		MoveSpeed:  30,
	}
}

func (l *FloatingLoupe) Toggle() {
	l.Active = !l.Active
}

func (l *FloatingLoupe) Move(dx, dy float32, maxX, maxY float32) {
	l.Position.X += dx * l.MoveSpeed
	l.Position.Y += dy * l.MoveSpeed

	// Clamp agar pointer tidak keluar layar TV
	if l.Position.X < l.Radius {
		l.Position.X = l.Radius
	}
	if l.Position.X > maxX-l.Radius {
		l.Position.X = maxX - l.Radius
	}
	if l.Position.Y < l.Radius {
		l.Position.Y = l.Radius
	}
	if l.Position.Y > maxY-l.Radius {
		l.Position.Y = maxY - l.Radius
	}
}

// Layout merender efek kaca pembesar di atas gambar komik
func (l *FloatingLoupe) Layout(gtx layout.Context, imgOp paint.ImageOp, imgBounds image.Rectangle) layout.Dimensions {
	if !l.Active {
		return layout.Dimensions{}
	}

	// 1. Simpan state op sebelum transformasi
	macroStack := op.Record(gtx.Ops)

	// Buat clipping shape lingkaran (lensa kaca pembesar)
	circleBounds := image.Rect(
		int(l.Position.X-l.Radius),
		int(l.Position.Y-l.Radius),
		int(l.Position.X+l.Radius),
		int(l.Position.Y+l.Radius),
	)
	clipArea := clip.Ellipse(circleBounds).Push(gtx.Ops)

	// 2. Terapkan transformasi Zoom di dalam area clipping
	// Titik tumpu zoom (pivot) tepat di koordinat kursor l.Position
	zoomTransform := f32.Affine2D{}.
		Offset(l.Position).
		Scale(f32.Pt(0, 0), f32.Pt(l.ZoomFactor, l.ZoomFactor)).
		Offset(l.Position.Mul(-1))

	transOp := op.Affine(zoomTransform).Push(gtx.Ops)

	// Render gambar komik yang sudah diperbesar di dalam lingkaran
	imgOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	transOp.Pop()
	clipArea.Pop()

	// 3. Gambar Ring / Frame Border Lensa Kaca Pembesar (Warna Cyan / Putih menyala)
	drawRing(gtx, l.Position, l.Radius, 4, color.NRGBA{R: 0, G: 220, B: 255, A: 255})

	// 4. Gambar Crosshair titik pusat kecil di tengah lensa
	drawCenterPoint(gtx, l.Position, color.NRGBA{R: 255, G: 60, B: 60, A: 200})

	macro := macroStack.Stop()
	macro.Add(gtx.Ops)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func drawRing(gtx layout.Context, center f32.Point, radius float32, thickness float32, clr color.NRGBA) {
	// Gambar lingkaran luar
	rOuter := radius + thickness/2
	rInner := radius - thickness/2
	
	rectOuter := image.Rect(
		int(center.X-rOuter),
		int(center.Y-rOuter),
		int(center.X+rOuter),
		int(center.Y+rOuter),
	)
	
	clipArea := clip.Ellipse(rectOuter).Push(gtx.Ops)
	paint.ColorOp{Color: clr}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	clipArea.Pop()

	// Masking transparan bagian dalam agar menjadi cincin (ring)
	rectInner := image.Rect(
		int(center.X-rInner),
		int(center.Y-rInner),
		int(center.X+rInner),
		int(center.Y+rInner),
	)
	// Kita biarkan simplifikasi visual dengan border outline
	_ = rectInner
}

func drawCenterPoint(gtx layout.Context, center f32.Point, clr color.NRGBA) {
	size := 3
	r := image.Rect(int(center.X)-size, int(center.Y)-size, int(center.X)+size, int(center.Y)+size)
	clipArea := clip.Ellipse(r).Push(gtx.Ops)
	paint.ColorOp{Color: clr}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	clipArea.Pop()
}

// Distance helper
func dist(a, b f32.Point) float32 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}
