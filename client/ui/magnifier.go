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
	Radius     float32   // Jari-jari lingkaran lensa (misal: 145 pixel)
	ZoomLevels []float32 // Tingkat perbesaran (2.0x, 2.8x, 3.8x)
	ZoomIndex  int       // Indeks zoom saat ini
	ZoomFactor float32   // Pengali zoom di dalam lensa
	MoveSpeed  float32   // Kecepatan gerak pointer saat tombol D-pad ditekan
}

func NewFloatingLoupe() *FloatingLoupe {
	levels := []float32{2.0, 2.8, 3.8}
	return &FloatingLoupe{
		Active:     false,
		Position:   f32.Pt(640, 360),
		Radius:     145,
		ZoomLevels: levels,
		ZoomIndex:  1, // Default 2.8x
		ZoomFactor: levels[1],
		MoveSpeed:  36,
	}
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

func (l *FloatingLoupe) Move(dx, dy float32, maxX, maxY float32) {
	l.Position.X += dx * l.MoveSpeed
	l.Position.Y += dy * l.MoveSpeed

	// Clamp agar lensa tetap berada di dalam layar TV dengan margin aman
	margin := float32(20)
	if l.Position.X < l.Radius+margin {
		l.Position.X = l.Radius + margin
	}
	if l.Position.X > maxX-l.Radius-margin {
		l.Position.X = maxX - l.Radius - margin
	}
	if l.Position.Y < l.Radius+margin {
		l.Position.Y = l.Radius + margin
	}
	if l.Position.Y > maxY-l.Radius-margin {
		l.Position.Y = maxY - l.Radius - margin
	}
}

// Layout merender efek kaca pembesar presisi (Continuous Optical Magnification) di atas gambar komik
func (l *FloatingLoupe) Layout(gtx layout.Context, th *material.Theme, imgOp paint.ImageOp, imgBounds image.Rectangle, pageOffsetX, pageOffsetY, pageScale float32) layout.Dimensions {
	if !l.Active {
		return layout.Dimensions{}
	}

	macroStack := op.Record(gtx.Ops)

	// 1. Hitung titik target pada gambar komik yang sedang tepat di bawah lensa
	// Proyeksi: Screen -> Image Space
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

	// Isi latar belakang lensa dengan warna gelap jika berada di tepi halaman
	paint.ColorOp{Color: color.NRGBA{R: 15, G: 15, B: 20, A: 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	// 3. Terapkan transformasi zoom affine:
	// Memetakan (imgTargetX, imgTargetY) agar berada tepat di titik pusat kursor l.Position
	// dengan pengali skala total = pageScale * l.ZoomFactor
	magScale := pageScale * l.ZoomFactor
	zoomTransform := f32.Affine2D{}.
		Offset(l.Position).
		Scale(f32.Pt(0, 0), f32.Pt(magScale, magScale)).
		Offset(f32.Pt(-imgTargetX, -imgTargetY))

	transOp := op.Affine(zoomTransform).Push(gtx.Ops)
	imgOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	transOp.Pop()

	// 4. Efek Vignette / Glare halus di tepi dalam lensa
	vignetteColor := color.NRGBA{R: 0, G: 220, B: 255, A: 25}
	paint.ColorOp{Color: vignetteColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	clipArea.Pop()

	// 5. Gambar Bezel / Frame Ring Lensa Kaca Pembesar
	// Ring Luar Neon Cyan Glow
	drawHollowRing(gtx, circleBounds, int(l.Radius), 4, color.NRGBA{R: 0, G: 229, B: 255, A: 255})
	// Ring Dalam Metalik Halus
	innerBounds := image.Rect(
		int(l.Position.X-l.Radius+3),
		int(l.Position.Y-l.Radius+3),
		int(l.Position.X+l.Radius-3),
		int(l.Position.Y+l.Radius-3),
	)
	drawHollowRing(gtx, innerBounds, int(l.Radius-3), 2, color.NRGBA{R: 255, G: 255, B: 255, A: 120})

	// 6. Gambar Crosshair Reticle halus di titik pusat lensa
	drawReticle(gtx, l.Position)

	// 7. Render Badge Indikator Zoom di atas lensa
	drawZoomBadge(gtx, th, l.Position, l.Radius, l.ZoomFactor)

	// 8. Render Bar Panduan Remote TV di bagian bawah layar
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
	reticleColor := color.NRGBA{R: 255, G: 64, B: 129, A: 220} // Neon Pink/Red reticle

	// Garis Horizontal Kiri & Kanan (dengan celah di tengah)
	leftH := image.Rect(int(center.X)-12, int(center.Y)-1, int(center.X)-4, int(center.Y)+1)
	paint.FillShape(gtx.Ops, reticleColor, clip.Rect(leftH).Op())

	rightH := image.Rect(int(center.X)+4, int(center.Y)-1, int(center.X)+12, int(center.Y)+1)
	paint.FillShape(gtx.Ops, reticleColor, clip.Rect(rightH).Op())

	// Garis Vertikal Atas & Bawah
	topV := image.Rect(int(center.X)-1, int(center.Y)-12, int(center.X)+1, int(center.Y)-4)
	paint.FillShape(gtx.Ops, reticleColor, clip.Rect(topV).Op())

	botV := image.Rect(int(center.X)-1, int(center.Y)+4, int(center.X)+1, int(center.Y)+12)
	paint.FillShape(gtx.Ops, reticleColor, clip.Rect(botV).Op())

	// Titik pusat mikro
	centerPt := image.Rect(int(center.X)-1, int(center.Y)-1, int(center.X)+1, int(center.Y)+1)
	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.Rect(centerPt).Op())
}

// drawZoomBadge menampilkan badge floating pill zoom (misal: "🔍 2.8x") tepat di atas lensa
func drawZoomBadge(gtx layout.Context, th *material.Theme, center f32.Point, radius float32, zoom float32) {
	badgeW := 90
	badgeH := 30
	bx := int(center.X) - (badgeW / 2)
	by := int(center.Y - radius - float32(badgeH) - 10)
	if by < 10 {
		by = int(center.Y + radius + 10)
	}

	badgeRect := image.Rect(bx, by, bx+badgeW, by+badgeH)
	rrect := clip.UniformRRect(badgeRect, 15)

	// Latar belakang gelap dengan border cyan
	paint.FillShape(gtx.Ops, color.NRGBA{R: 18, G: 24, B: 38, A: 240}, rrect.Op(gtx.Ops))
	paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 229, B: 255, A: 255}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1.5}.Op())

	// Label Teks Zoom
	offset := op.Offset(image.Pt(bx+14, by+5)).Push(gtx.Ops)
	lbl := material.Caption(th, fmt.Sprintf("🔍 %.1fx", zoom))
	lbl.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	lbl.TextSize = unit.Sp(13)
	lbl.Layout(gtx)
	offset.Pop()
}

// drawTVControlLegend menampilkan panduan remote TV yang elegan di bagian bawah layar saat loupe aktif
func drawTVControlLegend(gtx layout.Context, th *material.Theme, currentZoom float32) {
	layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			pillW := 620
			pillH := 42
			bgRect := image.Rect(0, 0, pillW, pillH)
			rrect := clip.UniformRRect(bgRect, 21)

			// Background Glass Pill
			paint.FillShape(gtx.Ops, color.NRGBA{R: 12, G: 16, B: 24, A: 235}, rrect.Op(gtx.Ops))
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 229, B: 255, A: 180}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1.5}.Op())

			// Layout Teks Petunjuk
			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				guideText := fmt.Sprintf("🔍 LOUPE AKTIF  |  [D-PAD] Geser Lensa  |  [OK] Zoom (%.1fx)  |  [BACK] Selesai", currentZoom)
				lbl := material.Body2(th, guideText)
				lbl.Color = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
				return layout.Center.Layout(gtx, lbl.Layout)
			})
		})
	})
}
