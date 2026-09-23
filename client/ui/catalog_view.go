package ui

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
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
	"gioui.org/widget"
	"gioui.org/widget/material"

	"comic_reader/pkg/model"
)

// CatalogView mengelola antarmuka katalog bergaya Leanback Launcher Android TV
type CatalogView struct {
	ServerURL    string
	Comics       []model.Comic
	FocusedIndex int
	HeroFocused  bool // Apakah fokus TV sedang berada di tombol CTA Hero
	Columns      int
	Loading      bool
	ErrorMessage string
	OnSelect     func(comic model.Comic)
	Invalidate   func()

	clicks   []widget.Clickable
	retryBtn widget.Clickable
	heroBtn  widget.Clickable

	coversMu sync.RWMutex
	covers   map[string]paint.ImageOp
}

func NewCatalogView(serverURL string, onSelect func(c model.Comic), invalidate func()) *CatalogView {
	cv := &CatalogView{
		ServerURL:    serverURL,
		Columns:      4,
		FocusedIndex: 0,
		OnSelect:     onSelect,
		Invalidate:   invalidate,
		covers:       make(map[string]paint.ImageOp),
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
		log.Printf("[ComicTV] Memuat katalog dari: %s", endpoint)
		resp, err := client.Get(endpoint)
		if err != nil {
			cv.ErrorMessage = fmt.Sprintf("Gagal koneksi ke %s\nError: %v", endpoint, err)
			log.Printf("[ComicTV] Gagal fetch katalog: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			cv.ErrorMessage = fmt.Sprintf("Server mengembalikan status: %s", resp.Status)
			log.Printf("[ComicTV] Status katalog server tidak valid: %s", resp.Status)
			return
		}

		var list []model.Comic
		if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
			cv.ErrorMessage = fmt.Sprintf("Gagal parsing data komik: %v", err)
			log.Printf("[ComicTV] Decode JSON katalog error: %v", err)
			return
		}

		cv.Comics = list
		if cv.FocusedIndex >= len(list) {
			cv.FocusedIndex = 0
		}
		log.Printf("[ComicTV] Berhasil memuat %d komik", len(list))

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

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Get(coverURL)
	if err != nil {
		log.Printf("[ComicTV] Gagal download cover komik %s: %v", id, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("[ComicTV] Gagal decode cover gambar %s: %v", id, err)
		return
	}

	imgOp := paint.NewImageOp(img)
	cv.coversMu.Lock()
	cv.covers[id] = imgOp
	cv.coversMu.Unlock()

	cv.triggerInvalidate()
}

// HandleKey menangani navigasi D-pad TV di katalog Leanback
func (cv *CatalogView) HandleKey(k RemoteKey) {
	if cv.ErrorMessage != "" || len(cv.Comics) == 0 {
		if k == KeySelect {
			cv.FetchCatalog()
		}
		return
	}

	switch k {
	case KeyRight:
		if !cv.HeroFocused && cv.FocusedIndex < len(cv.Comics)-1 {
			cv.FocusedIndex++
			cv.triggerInvalidate()
		}
	case KeyLeft:
		if !cv.HeroFocused && cv.FocusedIndex > 0 {
			cv.FocusedIndex--
			cv.triggerInvalidate()
		}
	case KeyUp:
		if !cv.HeroFocused {
			cv.HeroFocused = true
			cv.triggerInvalidate()
		}
	case KeyDown:
		if cv.HeroFocused {
			cv.HeroFocused = false
			cv.triggerInvalidate()
		}
	case KeySelect:
		if cv.FocusedIndex >= 0 && cv.FocusedIndex < len(cv.Comics) && cv.OnSelect != nil {
			log.Printf("[ComicTV] Membuka komik terpilih: %s (ID: %s)", cv.Comics[cv.FocusedIndex].Title, cv.Comics[cv.FocusedIndex].ID)
			cv.OnSelect(cv.Comics[cv.FocusedIndex])
		}
	}
}

// Layout merender seluruh antarmuka katalog Android TV model Leanback Launcher
func (cv *CatalogView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.Fill(gtx.Ops, color.NRGBA{R: 11, G: 14, B: 20, A: 255})

	if cv.Loading {
		return cv.renderLoading(gtx, th)
	}

	if cv.ErrorMessage != "" {
		return cv.renderError(gtx, th)
	}

	if len(cv.Comics) == 0 {
		return cv.renderEmpty(gtx, th)
	}

	if len(cv.clicks) < len(cv.Comics) {
		cv.clicks = make([]widget.Clickable, len(cv.Comics))
	}
	for i, c := range cv.Comics {
		if cv.clicks[i].Clicked(gtx) {
			cv.FocusedIndex = i
			cv.HeroFocused = false
			if cv.OnSelect != nil {
				cv.OnSelect(c)
			}
		}
	}
	if cv.heroBtn.Clicked(gtx) {
		if cv.FocusedIndex >= 0 && cv.FocusedIndex < len(cv.Comics) && cv.OnSelect != nil {
			cv.OnSelect(cv.Comics[cv.FocusedIndex])
		}
	}

	screenW := gtx.Constraints.Max.X
	screenH := gtx.Constraints.Max.Y

	margin := 36
	contentW := screenW - (margin * 2)

	// 1. Top Bar
	topBarH := 36
	cv.renderTopBar(gtx, th, margin, contentW, topBarH)

	// 2. Hero Spotlight
	heroY := topBarH + 12
	heroH := screenH * 32 / 100 // 32% layar untuk Hero (sekitar 230px di 720p)
	if heroH < 190 {
		heroH = 190
	}
	if heroH > 250 {
		heroH = 250
	}
	activeComic := cv.Comics[cv.FocusedIndex]
	cv.renderHeroSpotlight(gtx, th, activeComic, margin, heroY, contentW, heroH)

	// 3. Shelf
	shelfY := heroY + heroH + 16
	shelfH := screenH - shelfY - 10
	cv.renderHorizontalShelf(gtx, th, screenW, margin, shelfY, shelfH)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// renderTopBar merender top bar branding + status menggunakan layout.Flex
func (cv *CatalogView) renderTopBar(gtx layout.Context, th *material.Theme, marginX, contentW, barH int) {
	barOffset := op.Offset(image.Pt(marginX, 12)).Push(gtx.Ops)
	barGtx := gtx
	barGtx.Constraints.Min = image.Pt(0, 0)
	barGtx.Constraints.Max = image.Pt(contentW, barH)

	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(barGtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.Body1(th, "COMIC READER TV")
			title.Color = color.NRGBA{R: 0, G: 229, B: 255, A: 255}
			title.TextSize = unit.Sp(15)
			return title.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{}
			m := op.Record(gtx.Ops)
			statusLbl := material.Caption(th, "● ONLINE")
			statusLbl.Color = color.NRGBA{R: 0, G: 230, B: 118, A: 255}
			statusLbl.TextSize = unit.Sp(10)
			dims := layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, statusLbl.Layout)
			call := m.Stop()

			bgRect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
			rrect := clip.UniformRRect(bgRect, dims.Size.Y/2)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 16, G: 38, B: 28, A: 220}, rrect.Op(gtx.Ops))
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 230, B: 118, A: 120}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())
			call.Add(gtx.Ops)
			return dims
		}),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			countText := fmt.Sprintf("%d Komik", len(cv.Comics))
			countLbl := material.Caption(th, countText)
			countLbl.Color = color.NRGBA{R: 140, G: 150, B: 170, A: 255}
			return countLbl.Layout(gtx)
		}),
	)
	barOffset.Pop()
}

// renderHeroSpotlight merender panel info komik terpilih di area atas
func (cv *CatalogView) renderHeroSpotlight(gtx layout.Context, th *material.Theme, c model.Comic, marginX, heroY, contentW, heroH int) {
	heroOffset := op.Offset(image.Pt(marginX, heroY)).Push(gtx.Ops)

	// Background panel Hero: rounded 12, sleek glassmorphism
	heroRect := image.Rect(0, 0, contentW, heroH)
	rrect := clip.UniformRRect(heroRect, 12)
	paint.FillShape(gtx.Ops, color.NRGBA{R: 16, G: 21, B: 32, A: 220}, rrect.Op(gtx.Ops))
	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 18}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())

	// Poster Besar di Sebelah Kanan (proporsional)
	posterH := heroH - 24
	posterW := int(float32(posterH) / 1.45)
	if posterW > contentW/3 {
		posterW = contentW / 3
		posterH = int(float32(posterW) * 1.45)
	}
	posterX := contentW - posterW - 16
	posterY := (heroH - posterH) / 2
	cv.renderPoster(gtx, th, c, posterX, posterY, posterW, posterH, true)

	// Kolom Kiri: Metadata & Sinopsis
	leftPad := 20
	leftW := posterX - leftPad - 16
	contentOffset := op.Offset(image.Pt(leftPad, 14)).Push(gtx.Ops)
	leftGtx := gtx
	leftGtx.Constraints.Min = image.Pt(0, 0)
	leftGtx.Constraints.Max = image.Pt(leftW, heroH-60)

	layout.Flex{Axis: layout.Vertical}.Layout(leftGtx,
		// Badge Label Row
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = image.Point{}
					m := op.Record(gtx.Ops)
					tag := material.Caption(th, "KOMIK PILIHAN")
					tag.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
					tag.TextSize = unit.Sp(9)
					dims := layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, tag.Layout)
					call := m.Stop()

					bgRect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
					r := clip.UniformRRect(bgRect, 4)
					paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 120, B: 215, A: 220}, r.Op(gtx.Ops))
					call.Add(gtx.Ops)
					return dims
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					pageStr := fmt.Sprintf("%d Halaman", c.TotalPages)
					lbl := material.Caption(th, pageStr)
					lbl.Color = color.NRGBA{R: 0, G: 229, B: 255, A: 255}
					lbl.TextSize = unit.Sp(10)
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Caption(th, "PDF HD")
					lbl.Color = color.NRGBA{R: 255, G: 180, B: 0, A: 255}
					lbl.TextSize = unit.Sp(10)
					return lbl.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
		// Judul Komik
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H6(th, c.Title)
			title.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			title.TextSize = unit.Sp(18)
			title.MaxLines = 1
			return title.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
		// Sinopsis
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			desc := c.Description
			if desc == "" {
				desc = "Nikmati pembacaan komik resolusi tinggi langsung dengan remote Android TV."
			}
			desc = strings.ReplaceAll(desc, "\r\n", " ")
			desc = strings.ReplaceAll(desc, "\n", " ")
			desc = strings.Join(strings.Fields(desc), " ")
			if len(desc) > 120 {
				desc = desc[:117] + "..."
			}
			descLbl := material.Caption(th, desc)
			descLbl.Color = color.NRGBA{R: 165, G: 175, B: 195, A: 255}
			descLbl.TextSize = unit.Sp(11)
			descLbl.MaxLines = 2
			return descLbl.Layout(gtx)
		}),
	)
	contentOffset.Pop()

	// Tombol CTA: dinamis di bawah kiri Hero
	btnY := heroH - 42
	btnOffset := op.Offset(image.Pt(leftPad, btnY)).Push(gtx.Ops)
	btnGtx := gtx
	btnGtx.Constraints.Min = image.Point{}
	btnGtx.Constraints.Max = image.Pt(leftW, 34)

	btnColor := color.NRGBA{R: 18, G: 32, B: 50, A: 240}
	borderColor := color.NRGBA{R: 0, G: 200, B: 255, A: 180}
	textColor := color.NRGBA{R: 0, G: 229, B: 255, A: 255}
	btnText := "OK  •  BACA SEKARANG"
	if cv.HeroFocused {
		btnColor = color.NRGBA{R: 0, G: 195, B: 235, A: 255}
		borderColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		textColor = color.NRGBA{R: 10, G: 15, B: 25, A: 255}
		btnText = "▶  BACA SEKARANG (Tekan OK)"
	}

	m := op.Record(btnGtx.Ops)
	lbl := material.Body2(th, btnText)
	lbl.Color = textColor
	lbl.TextSize = unit.Sp(12)
	btnDims := layout.Inset{Top: unit.Dp(6), Bottom: unit.Dp(6), Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(btnGtx, lbl.Layout)
	btnCall := m.Stop()

	btnRect := image.Rect(0, 0, btnDims.Size.X, btnDims.Size.Y)
	br := clip.UniformRRect(btnRect, 6)
	paint.FillShape(btnGtx.Ops, btnColor, br.Op(btnGtx.Ops))
	paint.FillShape(btnGtx.Ops, borderColor, clip.Stroke{Path: br.Path(btnGtx.Ops), Width: 1.5}.Op())
	btnCall.Add(btnGtx.Ops)

	btnOffset.Pop()
	heroOffset.Pop()
}

func (cv *CatalogView) renderPoster(gtx layout.Context, th *material.Theme, c model.Comic, x, y, w, h int, isHero bool) {
	posterOffset := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
	rect := image.Rect(0, 0, w, h)
	rrect := clip.UniformRRect(rect, 10)

	paint.FillShape(gtx.Ops, color.NRGBA{R: 25, G: 30, B: 42, A: 255}, rrect.Op(gtx.Ops))

	cv.coversMu.RLock()
	imgOp, hasCover := cv.covers[c.ID]
	cv.coversMu.RUnlock()

	clipArea := rrect.Push(gtx.Ops)
	if hasCover {
		imgSize := imgOp.Size()
		if imgSize.X > 0 && imgSize.Y > 0 {
			scaleX := float32(w) / float32(imgSize.X)
			scaleY := float32(h) / float32(imgSize.Y)
			scale := scaleX
			if scaleY > scale {
				scale = scaleY
			}
			trans := f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))
			transStack := op.Affine(trans).Push(gtx.Ops)
			imgOp.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			transStack.Pop()
		}
	} else {
		subText := material.Caption(th, "Cover")
		subText.Color = color.NRGBA{R: 100, G: 110, B: 130, A: 255}
		tOffset := op.Offset(image.Pt(w/2-20, h/2-8)).Push(gtx.Ops)
		subText.Layout(gtx)
		tOffset.Pop()
	}
	clipArea.Pop()

	if isHero {
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 229, B: 255, A: 140}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 2}.Op())
	}

	posterOffset.Pop()
}

func (cv *CatalogView) renderHorizontalShelf(gtx layout.Context, th *material.Theme, screenW, margin, shelfY, shelfH int) {
	shelfOffset := op.Offset(image.Pt(0, shelfY)).Push(gtx.Ops)

	// Judul Section (Header Bar)
	titleOffset := op.Offset(image.Pt(margin, 0)).Push(gtx.Ops)
	titleGtx := gtx
	titleGtx.Constraints.Min = image.Point{}
	titleGtx.Constraints.Max = image.Pt(screenW-margin*2, 24)

	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(titleGtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			shelfTitle := material.Caption(th, "KOLEKSI KOMIK")
			shelfTitle.Color = color.NRGBA{R: 215, G: 225, B: 245, A: 255}
			shelfTitle.TextSize = unit.Sp(12)
			return shelfTitle.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			hint := material.Caption(th, "Gunakan Kiri / Kanan pada Remote TV")
			hint.Color = color.NRGBA{R: 120, G: 130, B: 155, A: 255}
			hint.TextSize = unit.Sp(10)
			return hint.Layout(gtx)
		}),
	)
	titleOffset.Pop()

	// Posisi kartu: cardsY diatur ke 34 agar ada jarak aman dari header (0..24),
	// sehingga saat scale 1.05 kartu tidak akan pernah menimpa teks "KOLEKSI KOMIK"!
	cardsY := 34
	availH := shelfH - cardsY - 32
	cardH := availH
	if cardH > 245 {
		cardH = 245
	}
	if cardH < 120 {
		cardH = 120
	}
	cardW := int(float32(cardH) / 1.45)
	cardGap := 18

	// Auto-Scroll Horizontal
	cardStride := cardW + cardGap
	centerTargetX := (screenW / 2) - (cardW / 2)
	currentCardX := margin + (cv.FocusedIndex * cardStride)

	scrollOffsetX := 0
	if currentCardX > centerTargetX {
		scrollOffsetX = centerTargetX - currentCardX
	}

	for i, c := range cv.Comics {
		cardX := margin + (i * cardStride) + scrollOffsetX
		if cardX+cardW+cardGap < 0 || cardX > screenW+50 {
			continue
		}

		isFocused := (i == cv.FocusedIndex && !cv.HeroFocused)
		cv.renderShelfCard(gtx, th, c, i, cardX, cardsY, cardW, cardH, isFocused)
	}

	shelfOffset.Pop()
}

func (cv *CatalogView) renderShelfCard(gtx layout.Context, th *material.Theme, c model.Comic, idx, x, y, w, h int, isFocused bool) {
	scale := float32(1.0)
	if isFocused {
		scale = 1.05
	}

	wScaled := int(float32(w) * scale)
	hScaled := int(float32(h) * scale)
	xScaled := x - (wScaled-w)/2
	yScaled := y - (hScaled-h)/2

	cardOffset := op.Offset(image.Pt(xScaled, yScaled)).Push(gtx.Ops)
	cardGtx := gtx
	cardGtx.Constraints.Min = image.Pt(wScaled, hScaled)
	cardGtx.Constraints.Max = image.Pt(wScaled, hScaled)

	renderCardContent := func(gtx layout.Context) layout.Dimensions {
		cardRect := image.Rect(0, 0, wScaled, hScaled)
		rrect := clip.UniformRRect(cardRect, 8)

		paint.FillShape(gtx.Ops, color.NRGBA{R: 20, G: 26, B: 38, A: 255}, rrect.Op(gtx.Ops))

		cv.coversMu.RLock()
		imgOp, hasCover := cv.covers[c.ID]
		cv.coversMu.RUnlock()

		clipCard := rrect.Push(gtx.Ops)
		if hasCover {
			imgSize := imgOp.Size()
			if imgSize.X > 0 && imgSize.Y > 0 {
				scaleX := float32(wScaled) / float32(imgSize.X)
				scaleY := float32(hScaled) / float32(imgSize.Y)
				scaleImg := scaleX
				if scaleY > scaleImg {
					scaleImg = scaleY
				}
				trans := f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scaleImg, scaleImg))
				transStack := op.Affine(trans).Push(gtx.Ops)
				imgOp.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				transStack.Pop()
			}
		} else {
			subText := material.Caption(th, "...")
			subText.Color = color.NRGBA{R: 100, G: 110, B: 130, A: 255}
			tOffset := op.Offset(image.Pt(wScaled/2-8, hScaled/2-6)).Push(gtx.Ops)
			subText.Layout(gtx)
			tOffset.Pop()
		}
		clipCard.Pop()

		// Border fokus
		if isFocused {
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 229, B: 255, A: 255}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 3}.Op())
		} else {
			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 20}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 0.5}.Op())
		}

		return layout.Dimensions{Size: image.Pt(wScaled, hScaled)}
	}

	if idx < len(cv.clicks) {
		cv.clicks[idx].Layout(cardGtx, renderCardContent)
	} else {
		renderCardContent(cardGtx)
	}
	cardOffset.Pop()

	// Label Judul di Bawah Kartu
	titleOffset := op.Offset(image.Pt(xScaled, yScaled+hScaled+8)).Push(gtx.Ops)
	titleGtx := gtx
	titleGtx.Constraints.Min = image.Point{}
	titleGtx.Constraints.Max.X = wScaled

	titleText := material.Caption(th, c.Title)
	titleText.TextSize = unit.Sp(10)
	titleText.MaxLines = 1
	if isFocused {
		titleText.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	} else {
		titleText.Color = color.NRGBA{R: 140, G: 150, B: 170, A: 255}
	}
	titleText.Layout(titleGtx)
	titleOffset.Pop()
}

func (cv *CatalogView) renderLoading(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H5(th, "Comic Reader TV")
				title.Color = color.NRGBA{R: 0, G: 229, B: 255, A: 255}
				return title.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, "Memuat Katalog Komik...")
				lbl.Color = color.NRGBA{R: 230, G: 230, B: 240, A: 255}
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				sub := material.Caption(th, fmt.Sprintf("Menghubungkan ke %s", cv.ServerURL))
				sub.Color = color.NRGBA{R: 130, G: 140, B: 160, A: 255}
				return sub.Layout(gtx)
			}),
		)
	})
}

func (cv *CatalogView) renderError(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if cv.retryBtn.Clicked(gtx) {
		cv.FetchCatalog()
	}
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				errTitle := material.H5(th, "Gagal Memuat Katalog")
				errTitle.Color = color.NRGBA{R: 255, G: 82, B: 82, A: 255}
				return errTitle.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				errMsg := material.Body1(th, cv.ErrorMessage)
				errMsg.Color = color.NRGBA{R: 210, G: 210, B: 220, A: 255}
				return errMsg.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(18)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &cv.retryBtn, "Coba Lagi (Tekan OK)")
				btn.Background = color.NRGBA{R: 0, G: 140, B: 235, A: 255}
				return btn.Layout(gtx)
			}),
		)
	})
}

func (cv *CatalogView) renderEmpty(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if cv.retryBtn.Clicked(gtx) {
		cv.FetchCatalog()
	}
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, "Belum ada komik. Masukkan komik melalui web admin server.")
				lbl.Color = color.NRGBA{R: 190, G: 195, B: 210, A: 255}
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &cv.retryBtn, "Muat Ulang (Tekan OK)")
				btn.Background = color.NRGBA{R: 0, G: 140, B: 235, A: 255}
				return btn.Layout(gtx)
			}),
		)
	})
}
