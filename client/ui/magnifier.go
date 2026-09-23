package ui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// FloatingLoupe mengelola status dan render kaca pembesar melayang beresolusi tinggi di atas komik
type FloatingLoupe struct {
	Active     bool
	Position   f32.Point // Posisi kursor/lensa di layar TV (X, Y)
	Radius     float32   // Jari-jari lingkaran lensa (misal: 130 pixel)
	ZoomLevels []float32 // Tingkat perbesaran (2.0x, 2.8x, 3.8x)
	ZoomIndex  int       // Indeks zoom saat ini
	ZoomFactor float32   // Pengali zoom di dalam lensa
	MoveSpeed  float32   // Kecepatan gerak pointer saat tombol D-pad ditekan
}

func NewFloatingLoupe() *FloatingLoupe {
	levels := []float32{2.0, 2.8, 3.8, 6.0, 8.0}
	return &FloatingLoupe{
		Active:     false,
		Position:   f32.Pt(640, 360),
		Radius:     130,
		ZoomLevels: levels,
		ZoomIndex:  1, // Default 2.8x
		ZoomFactor: levels[1],
		MoveSpeed:  32,
	}
}

func formatZoom(z float32) string {
	if z == float32(int(z)) {
		return fmt.Sprintf("%.0fx", z)
	}
	return fmt.Sprintf("%.1fx", z)
}

func (l *FloatingLoupe) Toggle() {
	l.Active = !l.Active
}

// CycleZoom beralih ke level zoom berikutnya saat tombol OK/Select ditekan dalam mode loupe
func (l *FloatingLoupe) CycleZoom() float32 {
	if len(l.ZoomLevels) == 0 {
		return l.ZoomFactor
	}
	l.ZoomIndex = (l.ZoomIndex + 1) % len(l.ZoomLevels)
	l.ZoomFactor = l.ZoomLevels[l.ZoomIndex]
	return l.ZoomFactor
}

// Move menggeser posisi lensa. Lensa bisa bergerak hingga ke ujung layar,
// titik pusat lensa boleh mencapai tepi layar sehingga bisa zoom area ujung-ujung.
func (l *FloatingLoupe) Move(dx, dy float32, maxX, maxY float32) {
	l.Position.X += dx * l.MoveSpeed
	l.Position.Y += dy * l.MoveSpeed

	// Clamp: pusat lensa bisa sampai tepi layar, hanya dijaga agar tidak melewati layar
	minBound := float32(10)
	if l.Position.X < minBound {
		l.Position.X = minBound
	}
	if l.Position.X > maxX-minBound {
		l.Position.X = maxX - minBound
	}
	if l.Position.Y < minBound {
		l.Position.Y = minBound
	}
	if l.Position.Y > maxY-minBound {
		l.Position.Y = maxY - minBound
	}
}

// Layout merender efek kaca pembesar presisi (Continuous Optical Magnification) di atas gambar komik
func (l *FloatingLoupe) Layout(gtx layout.Context, th *material.Theme, imgOp paint.ImageOp, imgBounds image.Rectangle, pageOffsetX, pageOffsetY, pageScale float32) layout.Dimensions {
	if !l.Active {
		return layout.Dimensions{}
	}

	macroStack := op.Record(gtx.Ops)

	// 1. Hitung titik target pada gambar komik yang sedang tepat di bawah lensa
	imgTargetX := (l.Position.X - pageOffsetX) / pageScale
	imgTargetY := (l.Position.Y - pageOffsetY) / pageScale

	// 2. Buat clipping shape lingkaran (lensa kaca pembesar)
	circleBounds := image.Rect(
		int(l.Position.X-l.Radius),
		int(l.Position.Y-l.Radius),
		int(l.Position.X+l.Radius),
		int(l.Position.Y+l.Radius),
	)
	rrect := clip.UniformRRect(circleBounds, int(l.Radius))
	clipArea := rrect.Push(gtx.Ops)

	// Background gelap jika lensa di luar batas komik
	paint.FillShape(gtx.Ops, color.NRGBA{R: 8, G: 10, B: 15, A: 255}, rrect.Op(gtx.Ops))

	// 3. Terapkan transformasi zoom affine
	magScale := pageScale * l.ZoomFactor
	ox := l.Position.X - (imgTargetX * magScale)
	oy := l.Position.Y - (imgTargetY * magScale)

	zoomTransform := f32.NewAffine2D(
		magScale, 0, ox,
		0, magScale, oy,
	)

	transOp := op.Affine(zoomTransform).Push(gtx.Ops)
	imgOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	transOp.Pop()

	clipArea.Pop()

	// 4. Gambar Bezel / Frame Ring Lensa
	drawHollowRing(gtx, circleBounds, int(l.Radius), 3, color.NRGBA{R: 0, G: 229, B: 255, A: 255})
	innerBounds := image.Rect(
		int(l.Position.X-l.Radius+3),
		int(l.Position.Y-l.Radius+3),
		int(l.Position.X+l.Radius-3),
		int(l.Position.Y+l.Radius-3),
	)
	drawHollowRing(gtx, innerBounds, int(l.Radius-3), 1, color.NRGBA{R: 255, G: 255, B: 255, A: 80})

	// 5. Gambar Crosshair Reticle kecil di pusat lensa
	drawReticle(gtx, l.Position)

	// 6. Render Badge Zoom
	drawZoomBadge(gtx, th, l.Position, l.Radius, l.ZoomFactor)

	// 7. Render Bar Panduan Remote TV di bawah layar
	drawTVControlLegend(gtx, th, l.ZoomFactor)

	macro := macroStack.Stop()
	macro.Add(gtx.Ops)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// drawHollowRing menggambar outline lingkaran murni menggunakan clip.Stroke
func drawHollowRing(gtx layout.Context, bounds image.Rectangle, radius int, thickness float32, clr color.NRGBA) {
	rrect := clip.UniformRRect(bounds, radius)
	strokeOp := clip.Stroke{
		Path:  rrect.Path(gtx.Ops),
		Width: thickness,
	}.Op()
	paint.FillShape(gtx.Ops, clr, strokeOp)
}

// drawReticle menggambar bidikan crosshair kecil di pusat lensa
func drawReticle(gtx layout.Context, center f32.Point) {
	reticleColor := color.NRGBA{R: 255, G: 64, B: 129, A: 200}

	leftH := image.Rect(int(center.X)-10, int(center.Y)-1, int(center.X)-3, int(center.Y)+1)
	paint.FillShape(gtx.Ops, reticleColor, clip.Rect(leftH).Op())
	rightH := image.Rect(int(center.X)+3, int(center.Y)-1, int(center.X)+10, int(center.Y)+1)
	paint.FillShape(gtx.Ops, reticleColor, clip.Rect(rightH).Op())
	topV := image.Rect(int(center.X)-1, int(center.Y)-10, int(center.X)+1, int(center.Y)-3)
	paint.FillShape(gtx.Ops, reticleColor, clip.Rect(topV).Op())
	botV := image.Rect(int(center.X)-1, int(center.Y)+3, int(center.X)+1, int(center.Y)+10)
	paint.FillShape(gtx.Ops, reticleColor, clip.Rect(botV).Op())

	centerPt := image.Rect(int(center.X)-1, int(center.Y)-1, int(center.X)+1, int(center.Y)+1)
	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.Rect(centerPt).Op())
}

// drawZoomBadge menampilkan badge zoom pill di bawah/atas lensa, ukuran dinamis membungkus teks
func drawZoomBadge(gtx layout.Context, th *material.Theme, center f32.Point, radius float32, zoom float32) {
	gtx.Constraints.Min = image.Point{}

	m := op.Record(gtx.Ops)
	lbl := material.Caption(th, formatZoom(zoom))
	lbl.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	lbl.TextSize = unit.Sp(12)
	dims := layout.Inset{Top: unit.Dp(3), Bottom: unit.Dp(3), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, lbl.Layout)
	call := m.Stop()

	by := int(center.Y + radius + 6)
	if by+dims.Size.Y > gtx.Constraints.Max.Y-60 {
		by = int(center.Y - radius - float32(dims.Size.Y) - 6)
	}
	bx := int(center.X) - (dims.Size.X / 2)

	offset := op.Offset(image.Pt(bx, by)).Push(gtx.Ops)
	bgRect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
	rrect := clip.UniformRRect(bgRect, dims.Size.Y/2)
	paint.FillShape(gtx.Ops, color.NRGBA{R: 16, G: 22, B: 34, A: 230}, rrect.Op(gtx.Ops))
	paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 229, B: 255, A: 200}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())
	call.Add(gtx.Ops)
	offset.Pop()
}

// drawTVControlLegend menampilkan panduan remote TV di bawah layar, ukuran dinamis pas membungkus teks
func drawTVControlLegend(gtx layout.Context, th *material.Theme, currentZoom float32) {
	layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = image.Point{}
				guideText := fmt.Sprintf("[D-PAD] Geser  •  [OK] Zoom %s  •  [BACK] Keluar Loupe", formatZoom(currentZoom))
				lbl := material.Caption(th, guideText)
				lbl.Color = color.NRGBA{R: 220, G: 225, B: 240, A: 255}
				lbl.TextSize = unit.Sp(12)

				m := op.Record(gtx.Ops)
				dims := layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, lbl.Layout)
				call := m.Stop()

				bgRect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
				rrect := clip.UniformRRect(bgRect, dims.Size.Y/2)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 10, G: 14, B: 22, A: 230}, rrect.Op(gtx.Ops))
				paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 229, B: 255, A: 160}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())
				call.Add(gtx.Ops)
				return dims
			})
		})
	})
}
