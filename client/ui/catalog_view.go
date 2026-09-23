package ui

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"
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

type CatalogView struct {
	ServerURL    string
	Comics       []model.Comic
	FocusedIndex int
	Columns      int
	Loading      bool
	ErrorMessage string
	OnSelect     func(comic model.Comic)
	Invalidate   func()

	coversMu sync.RWMutex
	covers   map[string]paint.ImageOp
}

func NewCatalogView(serverURL string, onSelect func(c model.Comic), invalidate func()) *CatalogView {
	cv := &CatalogView{
		ServerURL:  serverURL,
		Columns:    4, // 4 kolom di layar landscape Android TV
		OnSelect:   onSelect,
		Invalidate: invalidate,
		covers:     make(map[string]paint.ImageOp),
	}
	cv.FetchCatalog()
	return cv
}

func (cv *CatalogView) triggerInvalidate() {
	if cv.Invalidate != nil {
		cv.Invalidate()
	}
}

func (cv *CatalogView) FetchCatalog() {
	cv.Loading = true
	cv.ErrorMessage = ""
	cv.triggerInvalidate()

	go func() {
		defer func() {
			cv.Loading = false
			cv.triggerInvalidate()
		}()

		client := &http.Client{Timeout: 12 * time.Second}
		endpoint := strings.TrimRight(cv.ServerURL, "/") + "/api/comics"
		resp, err := client.Get(endpoint)
		if err != nil {
			cv.ErrorMessage = fmt.Sprintf("Gagal koneksi ke %s\nError: %v", endpoint, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			cv.ErrorMessage = fmt.Sprintf("Server mengembalikan status: %s", resp.Status)
			return
		}

		var list []model.Comic
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			cv.ErrorMessage = fmt.Sprintf("Gagal parsing data komik: %v", err)
			return
		}

		cv.Comics = list
		cv.FocusedIndex = 0

		// Unduh cover tiap komik di latar belakang
		for _, c := range list {
			if c.CoverURL != "" {
				go cv.fetchCover(c.ID, c.CoverURL)
			}
		}
	}()
}

func (cv *CatalogView) fetchCover(id, coverURL string) {
	cv.coversMu.RLock()
	_, exists := cv.covers[id]
	cv.coversMu.RUnlock()
	if exists {
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(coverURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return
	}

	imgOp := paint.NewImageOp(img)
	cv.coversMu.Lock()
	cv.covers[id] = imgOp
	cv.coversMu.Unlock()

	cv.triggerInvalidate()
}

// HandleKey menangani navigasi D-pad TV di katalog
func (cv *CatalogView) HandleKey(k RemoteKey) {
	// Jika sedang error atau data kosong, tombol OK/Select memicu reload
	if cv.ErrorMessage != "" || len(cv.Comics) == 0 {
		if k == KeySelect {
			cv.FetchCatalog()
		}
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
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := material.H5(th, "📺 Comic Reader TV")
					title.Color = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
					return title.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.H6(th, "Memuat Katalog Komik...")
					lbl.Color = color.NRGBA{R: 0, G: 220, B: 255, A: 255}
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					sub := material.Body2(th, fmt.Sprintf("Menghubungkan ke %s", cv.ServerURL))
					sub.Color = color.NRGBA{R: 160, G: 160, B: 160, A: 255}
					return sub.Layout(gtx)
				}),
			)
		})
	}

	if cv.ErrorMessage != "" {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					errTitle := material.H5(th, "⚠️ Gagal Memuat Katalog")
					errTitle.Color = color.NRGBA{R: 255, G: 82, B: 82, A: 255}
					return errTitle.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					errMsg := material.Body1(th, cv.ErrorMessage)
					errMsg.Color = color.NRGBA{R: 220, G: 220, B: 220, A: 255}
					return errMsg.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					hint := material.Body2(th, "Tekan tombol [OK / Enter] pada remote TV untuk mencoba lagi")
					hint.Color = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
					return hint.Layout(gtx)
				}),
			)
		})
	}

	if len(cv.Comics) == 0 {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(th, "Belum ada komik di katalog. Masukkan komik melalui web admin server.")
					lbl.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					hint := material.Body2(th, "Tekan tombol [OK / Enter] untuk memuat ulang")
					hint.Color = color.NRGBA{R: 76, G: 175, B: 80, A: 255}
					return hint.Layout(gtx)
				}),
			)
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

	// Body Kartu Komik (Background)
	clipCard := clip.RRect{Rect: cardRect, SE: 8, SW: 8, NW: 8, NE: 8}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{R: 35, G: 35, B: 35, A: 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	// Gambar Cover jika sudah selesai diunduh
	cv.coversMu.RLock()
	imgOp, hasCover := cv.covers[c.ID]
	cv.coversMu.RUnlock()

	if hasCover {
		imgSize := imgOp.Size()
		if imgSize.X > 0 && imgSize.Y > 0 {
			scaleX := float32(wScaled) / float32(imgSize.X)
			scaleY := float32(hScaled) / float32(imgSize.Y)

			macro := op.Record(gtx.Ops)
			trans := f32.Affine2D{}.
				Offset(f32.Pt(float32(xScaled), float32(yScaled))).
				Scale(f32.Pt(0, 0), f32.Pt(scaleX, scaleY))
			op.Affine(trans).Add(gtx.Ops)
			imgOp.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			macro.Stop().Add(gtx.Ops)
		}
	} else {
		// Placeholder saat cover sedang diunduh
		subText := material.Caption(th, "📖 Memuat...")
		subText.Color = color.NRGBA{R: 120, G: 120, B: 120, A: 255}
		titleMacro := op.Record(gtx.Ops)
		op.Offset(image.Pt(xScaled+12, yScaled+hScaled/2-10)).Add(gtx.Ops)
		gtxSub := gtx
		gtxSub.Constraints.Max.X = wScaled - 24
		subText.Layout(gtxSub)
		titleMacro.Stop().Add(gtx.Ops)
	}

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
