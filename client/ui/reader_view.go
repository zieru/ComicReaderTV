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
	// 1. Jika Floating Loupe aktif, semua key ditangani di sini
	if rv.Loupe.Active {
		switch k {
		case KeyUp:
			rv.Loupe.Move(0, -1, screenWidth, screenHeight)
		case KeyDown:
			rv.Loupe.Move(0, 1, screenWidth, screenHeight)
		case KeyLeft:
			rv.Loupe.Move(-1, 0, screenWidth, screenHeight)
		case KeyRight:
			rv.Loupe.Move(1, 0, screenWidth, screenHeight)
		case KeySelect:
			newZoom := rv.Loupe.CycleZoom()
			log.Printf("[ComicTV] Loupe zoom cycle: %s", formatZoom(newZoom))
		case KeyBack:
			// BACK saat loupe aktif: matikan loupe, BUKAN keluar dari komik
			rv.Loupe.Active = false
			log.Printf("[ComicTV] Loupe dimatikan, kembali ke mode baca")
		case KeyZoom:
			rv.Loupe.Active = false
			log.Printf("[ComicTV] Loupe dimatikan via tombol zoom")
		}
		return // PENTING: semua key di-consume saat loupe aktif, tidak jatuh ke bawah
	}

	// 2. Mode Membaca Normal (Loupe tidak aktif)
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
		if rv.ShowHUD {
			rv.ShowHUD = false
			rv.Loupe.Active = true
			rv.Loupe.Position = f32.Pt(screenWidth/2, screenHeight/2)
			log.Printf("[ComicTV] Loupe diaktifkan via HUD")
		}
	case KeyZoom:
		rv.Loupe.Active = true
		rv.ShowHUD = false
		rv.Loupe.Position = f32.Pt(screenWidth/2, screenHeight/2)
		log.Printf("[ComicTV] Loupe diaktifkan via tombol Zoom/Play")
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

		// 2. Render Loupe
		if rv.Loupe.Active {
			rv.Loupe.Layout(gtx, th, imgOp, bounds, offsetX, offsetY, scale)
		}
	}

	// Loading state
	if loading {
		rv.renderLoadingOverlay(gtx, th)
	}

	// Error state
	if errMsg != "" {
		rv.renderErrorOverlay(gtx, th, errMsg)
	}

	// 3. Render HUD Overlay Menu
	if rv.ShowHUD && !rv.Loupe.Active {
		rv.renderHUD(gtx, th)
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// renderLoadingOverlay menampilkan indikator loading di tengah layar
func (rv *ReaderView) renderLoadingOverlay(gtx layout.Context, th *material.Theme) {
	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = image.Point{}
		m := op.Record(gtx.Ops)
		dims := layout.Inset{Top: unit.Dp(14), Bottom: unit.Dp(14), Left: unit.Dp(28), Right: unit.Dp(28)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(th, fmt.Sprintf("Memuat Halaman %d / %d...", rv.CurrentPage, rv.TotalPages))
					lbl.Color = color.NRGBA{R: 0, G: 229, B: 255, A: 255}
					return lbl.Layout(gtx)
				}),
			)
		})
		call := m.Stop()

		bgRect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
		rrect := clip.UniformRRect(bgRect, 14)
		paint.FillShape(gtx.Ops, color.NRGBA{R: 16, G: 20, B: 32, A: 230}, rrect.Op(gtx.Ops))
		call.Add(gtx.Ops)
		return dims
	})
}

// renderErrorOverlay menampilkan pesan error di tengah layar
func (rv *ReaderView) renderErrorOverlay(gtx layout.Context, th *material.Theme, errMsg string) {
	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min = image.Point{}
		gtx.Constraints.Max.X = gtx.Constraints.Max.X - 100 // margin kiri-kanan

		m := op.Record(gtx.Ops)
		dims := layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(16), Left: unit.Dp(24), Right: unit.Dp(24)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, errMsg)
					lbl.Color = color.NRGBA{R: 255, G: 110, B: 110, A: 255}
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					sub := material.Caption(th, "◄ / ► Ganti Halaman  •  BACK Keluar")
					sub.Color = color.NRGBA{R: 180, G: 185, B: 200, A: 255}
					return sub.Layout(gtx)
				}),
			)
		})
		call := m.Stop()

		bgRect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
		rrect := clip.UniformRRect(bgRect, 14)
		paint.FillShape(gtx.Ops, color.NRGBA{R: 40, G: 15, B: 20, A: 230}, rrect.Op(gtx.Ops))
		call.Add(gtx.Ops)
		return dims
	})
}

// renderHUD menampilkan overlay HUD info + panduan kontrol saat OK ditekan
func (rv *ReaderView) renderHUD(gtx layout.Context, th *material.Theme) {
	// Dim overlay untuk kesan premium
	dimRect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 120}, clip.Rect(dimRect).Op())

	// Bar Atas: Judul + Halaman (ukuran dinamis)
	layout.N.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: unit.Dp(20), Left: unit.Dp(28), Right: unit.Dp(28)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = image.Point{}
				m := op.Record(gtx.Ops)
				dims := layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							title := material.Body1(th, rv.Comic.Title)
							title.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
							return title.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(24)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							pageInfo := material.Body2(th, fmt.Sprintf("Halaman %d / %d", rv.CurrentPage, rv.TotalPages))
							pageInfo.Color = color.NRGBA{R: 0, G: 229, B: 255, A: 255}
							return pageInfo.Layout(gtx)
						}),
					)
				})
				call := m.Stop()
				bgRect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
				rrect := clip.UniformRRect(bgRect, dims.Size.Y/2)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 12, G: 16, B: 26, A: 230}, rrect.Op(gtx.Ops))
				paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 30}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())
				call.Add(gtx.Ops)
				return dims
			})
		})
	})

	// Bar Bawah: Panduan Kontrol (ukuran dinamis)
	layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = image.Point{}
				m := op.Record(gtx.Ops)
				dims := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					guide := material.Caption(th, "[UP] Kaca Pembesar  •  [KIRI / KANAN] Halaman  •  [OK] Tutup Menu  •  [BACK] Keluar")
					guide.Color = color.NRGBA{R: 210, G: 215, B: 230, A: 255}
					guide.TextSize = unit.Sp(12)
					return guide.Layout(gtx)
				})
				call := m.Stop()
				bgRect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
				rrect := clip.UniformRRect(bgRect, dims.Size.Y/2)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 12, G: 16, B: 26, A: 230}, rrect.Op(gtx.Ops))
				paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 229, B: 255, A: 140}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())
				call.Add(gtx.Ops)
				return dims
			})
		})
	})
}
