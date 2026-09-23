package ui

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"sync"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"comic_reader/pkg/model"
)

type ReaderView struct {
	ServerURL   string
	Comic       model.Comic
	CurrentPage int
	TotalPages  int
	ShowHUD     bool
	Loupe       *FloatingLoupe
	OnBack      func()
	Invalidate  func()

	imageMu     sync.RWMutex
	currentPage paint.ImageOp
	pageBounds  image.Rectangle
	isLoading   bool
	lastError   string
}

func NewReaderView(serverURL string, comic model.Comic, onBack func(), invalidate func()) *ReaderView {
	tot := comic.TotalPages
	if tot == 0 {
		tot = 1
	}
	rv := &ReaderView{
		ServerURL:   serverURL,
		Comic:       comic,
		CurrentPage: 1,
		TotalPages:  tot,
		Loupe:       NewFloatingLoupe(),
		OnBack:      onBack,
		Invalidate:  invalidate,
	}
	rv.loadPage(1)
	return rv
}

func (rv *ReaderView) triggerInvalidate() {
	if rv.Invalidate != nil {
		rv.Invalidate()
	}
}

func (rv *ReaderView) loadPage(num int) {
	rv.isLoading = true
	rv.lastError = ""
	rv.triggerInvalidate()

	go func() {
		defer func() {
			rv.triggerInvalidate()
		}()

		url := fmt.Sprintf("%s/api/comics/%s/page/%d", rv.ServerURL, rv.Comic.ID, num)
		log.Printf("[ComicTV] Memuat gambar halaman: %s", url)

		client := &http.Client{Timeout: 18 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			rv.imageMu.Lock()
			rv.lastError = fmt.Sprintf("Gagal download halaman %d: %v", num, err)
			rv.isLoading = false
			rv.imageMu.Unlock()
			log.Printf("[ComicTV] Network error: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			rv.imageMu.Lock()
			rv.lastError = fmt.Sprintf("Server HTTP error: %s", resp.Status)
			rv.isLoading = false
			rv.imageMu.Unlock()
			log.Printf("[ComicTV] Server status error: %s", resp.Status)
			return
		}

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			rv.imageMu.Lock()
			rv.lastError = fmt.Sprintf("Gagal decode gambar: %v", err)
			rv.isLoading = false
			rv.imageMu.Unlock()
			log.Printf("[ComicTV] Image decode error: %v", err)
			return
		}

		rv.imageMu.Lock()
		rv.currentPage = paint.NewImageOp(img)
		rv.pageBounds = img.Bounds()
		rv.isLoading = false
		rv.imageMu.Unlock()
		log.Printf("[ComicTV] Berhasil memuat halaman %d (%dx%d)", num, rv.pageBounds.Dx(), rv.pageBounds.Dy())
	}()
}

// HandleKey menangani input remote D-pad saat membaca
func (rv *ReaderView) HandleKey(k RemoteKey, screenWidth, screenHeight float32) {
	// 1. Jika Floating Loupe aktif, D-pad digunakan untuk navigasi lensa
	if rv.Loupe.Active {
		switch k {
		case KeyUp:
			rv.Loupe.Move(0, -1, screenWidth, screenHeight)
			return
		case KeyDown:
			rv.Loupe.Move(0, 1, screenWidth, screenHeight)
			return
		case KeyLeft:
			rv.Loupe.Move(-1, 0, screenWidth, screenHeight)
			return
		case KeyRight:
			rv.Loupe.Move(1, 0, screenWidth, screenHeight)
			return
		case KeySelect:
			// Tombol OK saat mode loupe: Siklus tingkat Zoom (2.0x -> 2.8x -> 3.8x)
			newZoom := rv.Loupe.CycleZoom()
			log.Printf("[ComicTV] Loupe zoom cycle: %.1fx", newZoom)
			return
		case KeyBack, KeyZoom:
			rv.Loupe.Active = false // Keluar dari mode kaca pembesar
			log.Printf("[ComicTV] Loupe mode dimatikan")
			return
		}
	}

	// 2. Mode Membaca Normal
	switch k {
	case KeyRight:
		if rv.CurrentPage < rv.TotalPages {
			rv.CurrentPage++
			rv.loadPage(rv.CurrentPage)
		}
	case KeyLeft:
		if rv.CurrentPage > 1 {
			rv.CurrentPage--
			rv.loadPage(rv.CurrentPage)
		}
	case KeySelect:
		rv.ShowHUD = !rv.ShowHUD
	case KeyUp:
		// Jika HUD sedang tampil, tekan D-pad Atas dapat mengaktifkan Loupe langsung
		if rv.ShowHUD {
			rv.ShowHUD = false
			rv.Loupe.Active = true
			rv.Loupe.Position = f32.Pt(screenWidth/2, screenHeight/2)
			log.Printf("[ComicTV] Loupe diaktifkan via HUD")
		}
	case KeyZoom:
		rv.Loupe.Toggle()
		if rv.Loupe.Active {
			rv.ShowHUD = false
			rv.Loupe.Position = f32.Pt(screenWidth/2, screenHeight/2)
			log.Printf("[ComicTV] Loupe diaktifkan via tombol Zoom/Play di posisi: (%.1f, %.1f)", rv.Loupe.Position.X, rv.Loupe.Position.Y)
		}
	case KeyBack:
		if rv.ShowHUD {
			rv.ShowHUD = false
		} else if rv.OnBack != nil {
			rv.OnBack()
		}
	}
}

// Layout merender gambar komik fullscreen beserta Loupe presisi dan HUD
func (rv *ReaderView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Background Hitam pekat sinematik agar nyaman di mata saat di layar TV
	paint.Fill(gtx.Ops, color.NRGBA{R: 5, G: 7, B: 10, A: 255})

	rv.imageMu.RLock()
	hasImage := rv.pageBounds.Dx() > 0
	imgOp := rv.currentPage
	bounds := rv.pageBounds
	loading := rv.isLoading
	errMsg := rv.lastError
	rv.imageMu.RUnlock()

	// 1. Render Gambar Komik (Fit to Screen Proportional)
	if hasImage {
		screenW := float32(gtx.Constraints.Max.X)
		screenH := float32(gtx.Constraints.Max.Y)
		imgW := float32(bounds.Dx())
		imgH := float32(bounds.Dy())

		// Scale fit proportional
		scaleX := screenW / imgW
		scaleY := screenH / imgH
		scale := scaleX
		if scaleY < scale {
			scale = scaleY
		}

		destW := imgW * scale
		destH := imgH * scale
		offsetX := (screenW - destW) / 2
		offsetY := (screenH - destH) / 2

		trans := f32.NewAffine2D(scale, 0, offsetX, 0, scale, offsetY)
		transStack := op.Affine(trans).Push(gtx.Ops)
		imgOp.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		transStack.Pop()

		// 2. Render Kaca Pembesar Melayang (Floating Loupe) dengan Continuous Optical Projection
		if rv.Loupe.Active {
			rv.Loupe.Layout(gtx, th, imgOp, bounds, offsetX, offsetY, scale)
		}
	}

	// Loading state
	if loading {
		layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				cardRect := image.Rect(0, 0, 360, 90)
				rrect := clip.UniformRRect(cardRect, 16)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 20, G: 26, B: 38, A: 235}, rrect.Op(gtx.Ops))

				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body1(th, fmt.Sprintf("📖 Memuat Halaman %d / %d...", rv.CurrentPage, rv.TotalPages))
							lbl.Color = color.NRGBA{R: 0, G: 229, B: 255, A: 255}
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							sub := material.Caption(th, "Mengunduh resolusi HD...")
							sub.Color = color.NRGBA{R: 160, G: 160, B: 180, A: 255}
							return sub.Layout(gtx)
						}),
					)
				})
			})
		})
	}

	// Error state
	if errMsg != "" {
		layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				cardRect := image.Rect(0, 0, 480, 110)
				rrect := clip.UniformRRect(cardRect, 16)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 40, G: 15, B: 20, A: 240}, rrect.Op(gtx.Ops))

				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body1(th, "⚠️ "+errMsg)
							lbl.Color = color.NRGBA{R: 255, G: 90, B: 90, A: 255}
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							sub := material.Caption(th, "Tekan ◄ / ► untuk mencoba halaman lain atau Back untuk keluar")
							sub.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
							return sub.Layout(gtx)
						}),
					)
				})
			})
		})
	}

	// 3. Render HUD Overlay Menu (Muncul saat tombol OK ditekan dalam mode normal)
	if rv.ShowHUD && !rv.Loupe.Active {
		rv.renderHUD(gtx, th)
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func (rv *ReaderView) renderHUD(gtx layout.Context, th *material.Theme) {
	// Header Bar Atas
	layout.NW.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			barRect := image.Rect(0, 0, 680, 52)
			rrect := clip.UniformRRect(barRect, 14)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 15, G: 20, B: 30, A: 235}, rrect.Op(gtx.Ops))
			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 40}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())

			return layout.UniformInset(unit.Dp(14)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						title := material.H6(th, rv.Comic.Title)
						title.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						title.TextSize = unit.Sp(16)
						return title.Layout(gtx)
					}),
					layout.Flexed(1, layout.Spacer{}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						pageInfo := material.Body2(th, fmt.Sprintf("Halaman %d dari %d", rv.CurrentPage, rv.TotalPages))
						pageInfo.Color = color.NRGBA{R: 0, G: 229, B: 255, A: 255}
						return pageInfo.Layout(gtx)
					}),
				)
			})
		})
	})

	// Bottom Action Navigation Bar
	layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			barRect := image.Rect(0, 0, 780, 52)
			rrect := clip.UniformRRect(barRect, 14)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 15, G: 20, B: 30, A: 235}, rrect.Op(gtx.Ops))
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 229, B: 255, A: 120}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())

			return layout.UniformInset(unit.Dp(14)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				guide := material.Body2(th, "🔍 [▲ / Z] Kaca Pembesar  |  [◄ / ►] Ganti Halaman  |  [OK] Tutup Menu  |  [BACK] Keluar")
				guide.Color = color.NRGBA{R: 230, G: 230, B: 230, A: 255}
				return layout.Center.Layout(gtx, guide.Layout)
			})
		})
	})
}
