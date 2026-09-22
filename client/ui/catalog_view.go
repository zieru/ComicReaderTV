package ui

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"sync"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"comic_reader/pkg/model"
)

type CatalogView struct {
	ServerURL    string
	Comics       []model.Comic
	FocusedIndex int
	Columns      int
	Loading      bool
	ErrorMessage string
	OnSelect     func(comic model.Comic)

	coversMu sync.RWMutex
	covers   map[string]paint.ImageOp
}

func NewCatalogView(serverURL string, onSelect func(c model.Comic)) *CatalogView {
	cv := &CatalogView{
		ServerURL:    serverURL,
		Columns:      4, // 4 kolom di layar landscape Android TV
		OnSelect:     onSelect,
		covers:       make(map[string]paint.ImageOp),
	}
	cv.FetchCatalog()
	return cv
}

func (cv *CatalogView) FetchCatalog() {
	cv.Loading = true
	go func() {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(cv.ServerURL + "/api/comics")
		if err != nil {
			cv.ErrorMessage = fmt.Sprintf("Gagal koneksi ke server: %v", err)
			cv.Loading = false
			return
		}
		defer resp.Body.Close()

		var list []model.Comic
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			cv.ErrorMessage = "Gagal parsing respon server"
			cv.Loading = false
			return
		}

		cv.Comics = list
		cv.Loading = false
	}()
}

// HandleKey menangani navigasi D-pad TV di katalog
func (cv *CatalogView) HandleKey(k RemoteKey) {
	if len(cv.Comics) == 0 {
		return
	}

	switch k {
	case KeyRight:
		if cv.FocusedIndex < len(cv.Comics)-1 {
			cv.FocusedIndex++
		}
	case KeyLeft:
		if cv.FocusedIndex > 0 {
			cv.FocusedIndex--
		}
	case KeyDown:
		next := cv.FocusedIndex + cv.Columns
		if next < len(cv.Comics) {
			cv.FocusedIndex = next
		}
	case KeyUp:
		prev := cv.FocusedIndex - cv.Columns
		if prev >= 0 {
			cv.FocusedIndex = prev
		}
	case KeySelect:
		if cv.FocusedIndex >= 0 && cv.FocusedIndex < len(cv.Comics) && cv.OnSelect != nil {
			cv.OnSelect(cv.Comics[cv.FocusedIndex])
		}
	}
}

// Layout merender grid katalog di layar TV
func (cv *CatalogView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Latar belakang gelap khas Android TV
	paint.Fill(gtx.Ops, color.NRGBA{R: 18, G: 18, B: 18, A: 255})

	if cv.Loading {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.H6(th, "Memuat Katalog Komik...").Layout(gtx)
		})
	}

	if cv.ErrorMessage != "" {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, cv.ErrorMessage).Layout(gtx)
		})
	}

	if len(cv.Comics) == 0 {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return material.Body1(th, "Belum ada komik di katalog. Masukkan komik melalui admin web server.").Layout(gtx)
		})
	}

	// Hitung ukuran item grid berdasarkan lebar layar TV
	totalWidth := gtx.Constraints.Max.X
	padding := 32
	availableWidth := totalWidth - (padding * 2)
	gap := 20
	cardWidth := (availableWidth - (cv.Columns-1)*gap) / cv.Columns
	cardHeight := int(float32(cardWidth) * 1.45) // Rasio komik ~1:1.45

	return layout.UniformInset(unit.Dp(float32(padding))).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		for i, c := range cv.Comics {
			row := i / cv.Columns
			col := i % cv.Columns

			x := col * (cardWidth + gap)
			y := 60 + row*(cardHeight+gap+40) // 60px offset header

			isFocused := (i == cv.FocusedIndex)

			// Render kartu komik
			cv.renderComicCard(gtx, th, c, x, y, cardWidth, cardHeight, isFocused)
		}

		// Render Header Title di TV
		titleLabel := material.H5(th, "📺 Comic Reader TV")
		titleLabel.Color = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
		titleLabel.Layout(gtx)

		return layout.Dimensions{Size: gtx.Constraints.Max}
	})
}

func (cv *CatalogView) renderComicCard(gtx layout.Context, th *material.Theme, c model.Comic, x, y, w, h int, isFocused bool) {
	scale := float32(1.0)
	if isFocused {
		scale = 1.06 // Efek pop-up membesar saat disorot remote TV
	}

	wScaled := int(float32(w) * scale)
	hScaled := int(float32(h) * scale)
	xScaled := x - (wScaled-w)/2
	yScaled := y - (hScaled-h)/2

	cardRect := image.Rect(xScaled, yScaled, xScaled+wScaled, yScaled+hScaled)

	// Gambar Border Glow jika fokus
	if isFocused {
		glowBorder := image.Rect(cardRect.Min.X-4, cardRect.Min.Y-4, cardRect.Max.X+4, cardRect.Max.Y+4)
		clipArea := clip.RRect{Rect: glowBorder, SE: 12, SW: 12, NW: 12, NE: 12}.Push(gtx.Ops)
		paint.ColorOp{Color: color.NRGBA{R: 0, G: 220, B: 255, A: 255}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		clipArea.Pop()
	}

	// Body Kartu Komik
	clipCard := clip.RRect{Rect: cardRect, SE: 8, SW: 8, NW: 8, NE: 8}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{R: 35, G: 35, B: 35, A: 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	clipCard.Pop()

	// Label Judul di bawah kartu
	titleMacro := op.Record(gtx.Ops)
	op.Offset(image.Pt(xScaled, yScaled+hScaled+8)).Add(gtx.Ops)
	gtxLimited := gtx
	gtxLimited.Constraints.Max.X = wScaled

	titleText := material.Body2(th, c.Title)
	if isFocused {
		titleText.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	} else {
		titleText.Color = color.NRGBA{R: 180, G: 180, B: 180, A: 255}
	}
	titleText.Layout(gtxLimited)
	titleMacro.Stop().Add(gtx.Ops)
}
