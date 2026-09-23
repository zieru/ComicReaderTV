package ui

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"sync"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
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

	hudTimer    *time.Timer
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
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			rv.imageMu.Lock()
			rv.lastError = fmt.Sprintf("Gagal download halaman %d: %v", num, err)
			rv.isLoading = false
			rv.imageMu.Unlock()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			rv.imageMu.Lock()
			rv.lastError = fmt.Sprintf("Server HTTP error: %s", resp.Status)
			rv.isLoading = false
			rv.imageMu.Unlock()
			return
		}

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			rv.imageMu.Lock()
			rv.lastError = fmt.Sprintf("Gagal decode gambar: %v", err)
			rv.isLoading = false
			rv.imageMu.Unlock()
			return
		}

		rv.imageMu.Lock()
		rv.currentPage = paint.NewImageOp(img)
		rv.pageBounds = img.Bounds()
		rv.isLoading = false
		rv.imageMu.Unlock()
	}()
}

// HandleKey menangani input remote D-pad saat membaca
func (rv *ReaderView) HandleKey(k RemoteKey, screenWidth, screenHeight float32) {
	// 1. Jika Floating Loupe aktif, D-pad digunakan untuk mengarahkan kursor pembesar
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
		case KeySelect, KeyBack, KeyZoom:
			rv.Loupe.Active = false // Keluar dari mode kaca pembesar
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
	case KeyZoom:
		rv.Loupe.Toggle()
		if rv.Loupe.Active {
			// Posisikan di tengah layar saat pertama kali diaktifkan
			rv.Loupe.Position = f32.Pt(screenWidth/2, screenHeight/2)
		}
	case KeyBack:
		if rv.ShowHUD {
			rv.ShowHUD = false
		} else if rv.OnBack != nil {
			rv.OnBack()
		}
	}
}

// Layout merender gambar komik fullscreen beserta Loupe dan HUD
func (rv *ReaderView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Background Hitam pekat agar fokus membaca
	paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

	rv.imageMu.RLock()
	hasImage := rv.pageBounds.Dx() > 0
	imgOp := rv.currentPage
	bounds := rv.pageBounds
	loading := rv.isLoading
	errMsg := rv.lastError
	rv.imageMu.RUnlock()

	// 1. Render Gambar Komik (Fit to Screen)
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

		trans := f32.Affine2D{}.
			Offset(f32.Pt(offsetX, offsetY)).
			Scale(f32.Pt(0, 0), f32.Pt(scale, scale))

		macro := op.Record(gtx.Ops)
		op.Affine(trans).Add(gtx.Ops)
		imgOp.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		macro.Stop().Add(gtx.Ops)

		// 2. Render Kaca Pembesar Melayang (Floating Loupe)
		rv.Loupe.Layout(gtx, imgOp, bounds)
	}

	// Loading state
	if loading {
		layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.H6(th, fmt.Sprintf("Memuat Halaman %d...", rv.CurrentPage))
			lbl.Color = color.NRGBA{R: 0, G: 220, B: 255, A: 255}
			return lbl.Layout(gtx)
		})
	}

	// Error state
	if errMsg != "" {
		layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Body1(th, errMsg)
			lbl.Color = color.NRGBA{R: 255, G: 80, B: 80, A: 255}
			return lbl.Layout(gtx)
		})
	}

	// 3. Render HUD Menu Overlay (Muncul saat tombol OK ditekan)
	if rv.ShowHUD {
		rv.renderHUD(gtx, th)
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func (rv *ReaderView) renderHUD(gtx layout.Context, th *material.Theme) {
	// Top Header Bar
	layout.NW.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			title := material.H6(th, rv.Comic.Title)
			title.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			return title.Layout(gtx)
		})
	})

	// Bottom Page Indicator Bar
	layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			info := fmt.Sprintf("Halaman %d / %d  |  Tekan 'Z' / Play untuk Kaca Pembesar", rv.CurrentPage, rv.TotalPages)
			pageText := material.Body1(th, info)
			pageText.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
			return pageText.Layout(gtx)
		})
	})
}
