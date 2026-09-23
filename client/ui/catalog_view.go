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
	// Jika sedang error atau data kosong, tombol OK/Select memicu reload
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
		// Berpindah fokus ke tombol aksi Hero Spotlight
		if !cv.HeroFocused {
			cv.HeroFocused = true
			cv.triggerInvalidate()
		}
	case KeyDown:
		// Kembali ke rak kartu komik
		if cv.HeroFocused {
			cv.HeroFocused = false
			cv.triggerInvalidate()
		}
	case KeySelect:
		if cv.FocusedIndex >= 0 && cv.FocusedIndex < len(cv.Comics) && cv.OnSelect != nil {
			log.Printf("[ComicTV] Memilih komik: %s (ID: %s)", cv.Comics[cv.FocusedIndex].Title, cv.Comics[cv.FocusedIndex].ID)
			cv.OnSelect(cv.Comics[cv.FocusedIndex])
		}
	}
}

// Layout merender seluruh antarmuka katalog Android TV model Leanback Launcher
func (cv *CatalogView) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Latar belakang gelap khas Android TV (Cinematic Deep Slate)
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

	// Pastikan array clickable sesuai jumlah komik
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

	// 1. Render Top Header Bar (Branding & Server Status)
	cv.renderTopBar(gtx, th, screenW)

	// Komik yang sedang aktif disorot
	activeComic := cv.Comics[cv.FocusedIndex]

	// 2. Render Hero Spotlight (Bagian Atas ~46% Layar)
	heroH := int(float32(screenH) * 0.46)
	cv.renderHeroSpotlight(gtx, th, activeComic, screenW, heroH)

	// 3. Render Leanback Shelf / Carousel (Bagian Bawah ~54% Layar)
	shelfY := heroH + 10
	cv.renderHorizontalShelf(gtx, th, screenW, screenH-shelfY, shelfY)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func (cv *CatalogView) renderTopBar(gtx layout.Context, th *material.Theme, width int) {
	barOffset := op.Offset(image.Pt(36, 20)).Push(gtx.Ops)
	barGtx := gtx
	barGtx.Constraints.Max.X = width - 72

	layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(barGtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H6(th, "📺 COMIC READER TV")
			title.Color = color.NRGBA{R: 0, G: 229, B: 255, A: 255}
			title.TextSize = unit.Sp(16)
			return title.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Online Status Pill
			pillRect := image.Rect(0, 0, 110, 24)
			rrect := clip.UniformRRect(pillRect, 12)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 16, G: 38, B: 28, A: 220}, rrect.Op(gtx.Ops))
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 230, B: 118, A: 160}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())

			offset := op.Offset(image.Pt(10, 3)).Push(gtx.Ops)
			statusLbl := material.Caption(th, "● SERVER ON")
			statusLbl.Color = color.NRGBA{R: 0, G: 230, B: 118, A: 255}
			statusLbl.TextSize = unit.Sp(11)
			statusLbl.Layout(gtx)
			offset.Pop()
			return layout.Dimensions{Size: image.Pt(110, 24)}
		}),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			countText := fmt.Sprintf("%d Komik Tersedia", len(cv.Comics))
			countLbl := material.Body2(th, countText)
			countLbl.Color = color.NRGBA{R: 150, G: 160, B: 180, A: 255}
			return countLbl.Layout(gtx)
		}),
	)
	barOffset.Pop()
}

func (cv *CatalogView) renderHeroSpotlight(gtx layout.Context, th *material.Theme, c model.Comic, screenW, heroH int) {
	heroOffset := op.Offset(image.Pt(36, 56)).Push(gtx.Ops)
	heroGtx := gtx
	heroW := screenW - 72
	heroGtx.Constraints.Min = image.Pt(heroW, heroH-60)
	heroGtx.Constraints.Max = image.Pt(heroW, heroH-60)

	// Latar belakang panel Hero transparan
	heroRect := image.Rect(0, 0, heroW, heroH-60)
	rrect := clip.UniformRRect(heroRect, 16)
	paint.FillShape(gtx.Ops, color.NRGBA{R: 18, G: 23, B: 34, A: 160}, rrect.Op(gtx.Ops))
	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 15}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())

	contentOffset := op.Offset(image.Pt(28, 20)).Push(gtx.Ops)
	leftW := heroW - 240 // Lebar kolom teks

	// 1. Kolom Kiri: Metadata & Sinopsis
	leftGtx := heroGtx
	leftGtx.Constraints.Max.X = leftW
	layout.Flex{Axis: layout.Vertical}.Layout(leftGtx,
		// Badge Label "⭐ KOMIK PILIHAN"
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			tagRect := image.Rect(0, 0, 130, 22)
			r := clip.UniformRRect(tagRect, 6)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 120, B: 215, A: 200}, r.Op(gtx.Ops))
			tagOffset := op.Offset(image.Pt(10, 2)).Push(gtx.Ops)
			tag := material.Caption(th, "⭐ KOMIK PILIHAN")
			tag.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			tag.TextSize = unit.Sp(10)
			tag.Layout(gtx)
			tagOffset.Pop()
			return layout.Dimensions{Size: image.Pt(130, 22)}
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
		// Judul Komik Utama
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H5(th, c.Title)
			title.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			return title.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
		// Info Badges: Halaman, Format, Status
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					pageStr := fmt.Sprintf("📄 %d Halaman", c.TotalPages)
					lbl := material.Body2(th, pageStr)
					lbl.Color = color.NRGBA{R: 0, G: 229, B: 255, A: 255}
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, "⚡ Format PDF HD")
					lbl.Color = color.NRGBA{R: 255, G: 180, B: 0, A: 255}
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, "ID: #"+c.ID)
					lbl.Color = color.NRGBA{R: 140, G: 150, B: 170, A: 255}
					return lbl.Layout(gtx)
				}),
			)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		// Sinopsis Singkat (Maks 2 baris)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			desc := c.Description
			if desc == "" {
				desc = "Buka dan nikmati pembacaan komik beresolusi tinggi langsung dengan remote kontrol Android TV."
			}
			// Batasi teks sinopsis agar rapi di TV
			if len(desc) > 160 {
				desc = desc[:157] + "..."
			}
			descLbl := material.Body2(th, desc)
			descLbl.Color = color.NRGBA{R: 190, G: 195, B: 210, A: 255}
			return descLbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(14)}.Layout),
		// Tombol Aksi CTA "▶ Tekan OK untuk Mulai Membaca"
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btnW := 300
			btnH := 40
			btnRect := image.Rect(0, 0, btnW, btnH)
			br := clip.UniformRRect(btnRect, 10)

			btnColor := color.NRGBA{R: 0, G: 150, B: 255, A: 240}
			borderColor := color.NRGBA{R: 0, G: 229, B: 255, A: 255}
			if cv.HeroFocused {
				btnColor = color.NRGBA{R: 0, G: 200, B: 118, A: 255}
				borderColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			}

			paint.FillShape(gtx.Ops, btnColor, br.Op(gtx.Ops))
			paint.FillShape(gtx.Ops, borderColor, clip.Stroke{Path: br.Path(gtx.Ops), Width: 2}.Op())

			btnOffset := op.Offset(image.Pt(24, 9)).Push(gtx.Ops)
			btnText := "▶ Tekan [OK] untuk Membaca"
			if cv.HeroFocused {
				btnText = "▶ [ENTER] BACA SEKARANG"
			}
			lbl := material.Body1(th, btnText)
			lbl.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			lbl.Layout(gtx)
			btnOffset.Pop()

			return layout.Dimensions{Size: image.Pt(btnW, btnH)}
		}),
	)
	contentOffset.Pop()

	// 2. Kolom Kanan: Large Hero Poster Card Preview
	posterW := 150
	posterH := int(float32(posterW) * 1.45)
	posterX := heroW - posterW - 28
	posterY := (heroH - 60 - posterH) / 2

	cv.renderPoster(gtx, th, c, posterX, posterY, posterW, posterH, true)

	heroOffset.Pop()
}

func (cv *CatalogView) renderPoster(gtx layout.Context, th *material.Theme, c model.Comic, x, y, w, h int, isHero bool) {
	posterOffset := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
	rect := image.Rect(0, 0, w, h)
	rrect := clip.UniformRRect(rect, 10)

	// Background poster
	paint.FillShape(gtx.Ops, color.NRGBA{R: 25, G: 30, B: 42, A: 255}, rrect.Op(gtx.Ops))

	// Gambar cover
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
			macro := op.Record(gtx.Ops)
			op.Affine(trans).Add(gtx.Ops)
			imgOp.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			macro.Stop().Add(gtx.Ops)
		}
	} else {
		// Placeholder teks saat loading
		subText := material.Caption(th, "📖 Cover")
		subText.Color = color.NRGBA{R: 120, G: 130, B: 150, A: 255}
		tOffset := op.Offset(image.Pt(w/2-24, h/2-8)).Push(gtx.Ops)
		subText.Layout(gtx)
		tOffset.Pop()
	}
	clipArea.Pop()

	// Glow Border untuk Poster Hero
	if isHero {
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 229, B: 255, A: 160}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 2}.Op())
	}

	posterOffset.Pop()
}

func (cv *CatalogView) renderHorizontalShelf(gtx layout.Context, th *material.Theme, screenW, shelfH, shelfY int) {
	shelfOffset := op.Offset(image.Pt(0, shelfY)).Push(gtx.Ops)

	// Judul Section Rak
	titleOffset := op.Offset(image.Pt(36, 0)).Push(gtx.Ops)
	shelfTitle := material.Body1(th, "📚 KOLEKSI KOMIK")
	shelfTitle.Color = color.NRGBA{R: 220, G: 225, B: 240, A: 255}
	shelfTitle.Layout(gtx)
	titleOffset.Pop()

	cardW := 150
	cardH := int(float32(cardW) * 1.45)
	cardGap := 22
	cardsY := 28

	// Hitung Auto-Scroll Horizontal agar kartu yang sedang disorot selalu terlihat nyaman di tengah layar TV
	cardStride := cardW + cardGap
	centerTargetX := (screenW / 2) - (cardW / 2)
	currentCardX := 36 + (cv.FocusedIndex * cardStride)

	scrollOffsetX := 0
	if currentCardX > centerTargetX {
		scrollOffsetX = centerTargetX - currentCardX
	}

	for i, c := range cv.Comics {
		cardX := 36 + (i * cardStride) + scrollOffsetX
		// Lewati kartu yang berada di luar layar
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
		scale = 1.10 // Efek Pop-Up membesar saat disorot remote TV
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
		// Glow Border menyala jika kartu aktif disorot TV
		cardRect := image.Rect(0, 0, wScaled, hScaled)
		rrect := clip.UniformRRect(cardRect, 10)

		// Latar belakang kartu
		paint.FillShape(gtx.Ops, color.NRGBA{R: 25, G: 32, B: 45, A: 255}, rrect.Op(gtx.Ops))

		// Render cover komik
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
				macro := op.Record(gtx.Ops)
				op.Affine(trans).Add(gtx.Ops)
				imgOp.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				macro.Stop().Add(gtx.Ops)
			}
		} else {
			// Placeholder
			subText := material.Caption(th, "📖 Memuat...")
			subText.Color = color.NRGBA{R: 120, G: 130, B: 150, A: 255}
			tOffset := op.Offset(image.Pt(wScaled/2-32, hScaled/2-8)).Push(gtx.Ops)
			subText.Layout(gtx)
			tOffset.Pop()
		}
		clipCard.Pop()

		// Efek Sorot Remote TV (Vibrant Cyan Border)
		if isFocused {
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 229, B: 255, A: 255}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 3.5}.Op())
		} else {
			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 25}, clip.Stroke{Path: rrect.Path(gtx.Ops), Width: 1}.Op())
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
	titleGtx.Constraints.Max.X = wScaled

	titleText := material.Body2(th, c.Title)
	if isFocused {
		titleText.Color = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	} else {
		titleText.Color = color.NRGBA{R: 160, G: 170, B: 185, A: 255}
	}
	titleText.Layout(titleGtx)
	titleOffset.Pop()
}

func (cv *CatalogView) renderLoading(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H5(th, "📺 Comic Reader TV")
				title.Color = color.NRGBA{R: 0, G: 229, B: 255, A: 255}
				return title.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(18)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(th, "Memuat Katalog Komik...")
				lbl.Color = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				sub := material.Caption(th, fmt.Sprintf("Menghubungkan ke %s", cv.ServerURL))
				sub.Color = color.NRGBA{R: 140, G: 150, B: 170, A: 255}
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
				errTitle := material.H5(th, "⚠️ Gagal Memuat Katalog")
				errTitle.Color = color.NRGBA{R: 255, G: 82, B: 82, A: 255}
				return errTitle.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				errMsg := material.Body1(th, cv.ErrorMessage)
				errMsg.Color = color.NRGBA{R: 220, G: 220, B: 220, A: 255}
				return errMsg.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &cv.retryBtn, "Coba Lagi (Tekan OK)")
				btn.Background = color.NRGBA{R: 0, G: 150, B: 255, A: 255}
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
				lbl := material.Body1(th, "Belum ada komik di katalog. Masukkan komik melalui web admin server.")
				lbl.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &cv.retryBtn, "Muat Ulang (Tekan OK)")
				btn.Background = color.NRGBA{R: 0, G: 150, B: 255, A: 255}
				return btn.Layout(gtx)
			}),
		)
	})
}
