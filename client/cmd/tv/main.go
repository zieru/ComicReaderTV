package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"comic_reader/client/ui"
	"comic_reader/client/updater"
	"comic_reader/pkg/model"
)

type AppState int

const (
	StateCatalog AppState = iota
	StateReader
)

type UpdateStatus struct {
	mu           sync.Mutex
	available    bool
	releaseTag   string
	downloadURL  string
	downloading  bool
	progress     float32
	downloaded   bool
	savedPath    string
	errorMessage string
}

type TVApp struct {
	state        AppState
	catalogView  *ui.CatalogView
	readerView   *ui.ReaderView
	theme        *material.Theme
	serverURL    string
	updateStatus UpdateStatus
	window       *app.Window
}

const CurrentAppVersion = "v1.0.12"

func main() {
	defaultServer := "http://ca.tsel.my.id:8080"
	if envURL := os.Getenv("COMIC_SERVER_URL"); envURL != "" {
		defaultServer = envURL
	}
	serverURL := flag.String("server", defaultServer, "Alamat server katalog")
	flag.Parse()

	go func() {
		// Buat jendela Fullscreen standar TV (1080p / rasio 16:9)
		w := new(app.Window)
		w.Option(
			app.Title("Comic Reader TV"),
			app.Size(unit.Dp(1280), unit.Dp(720)),
			app.MinSize(unit.Dp(800), unit.Dp(450)),
		)

		if err := run(w, *serverURL); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func run(w *app.Window, serverURL string) error {
	th := material.NewTheme()

	tvApp := &TVApp{
		state:     StateCatalog,
		theme:     th,
		serverURL: serverURL,
		window:    w,
	}

	tvApp.catalogView = ui.NewCatalogView(serverURL, func(selected model.Comic) {
		tvApp.readerView = ui.NewReaderView(serverURL, selected, func() {
			tvApp.state = StateCatalog
			w.Invalidate()
		}, w.Invalidate)
		tvApp.state = StateReader
		w.Invalidate()
	}, w.Invalidate)

	// Inisialisasi Auto Updater di latar belakang
	appUpdater := updater.NewUpdater("zieru", "ComicReaderTV", CurrentAppVersion)
	go tvApp.checkAndApplyUpdate(appUpdater, w)

	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Registrasi penerimaan event Keyboard / Remote TV
			event.Op(gtx.Ops, "main_tag")

			// Handle Hardware Key Events (Remote TV)
			for {
				ev, ok := gtx.Event(key.Filter{Name: ""})
				if !ok {
					break
				}
				if ke, ok := ev.(key.Event); ok {
					rKey := ui.MapKeyEvent(ke)
					if rKey != ui.KeyNone {
						screenWidth := float32(gtx.Constraints.Max.X)
						screenHeight := float32(gtx.Constraints.Max.Y)

						if tvApp.state == StateCatalog {
							tvApp.catalogView.HandleKey(rKey)
						} else if tvApp.state == StateReader && tvApp.readerView != nil {
							tvApp.readerView.HandleKey(rKey, screenWidth, screenHeight)
						}
						w.Invalidate()
					}
				}
			}

			// Render Layar sesuai State
			if tvApp.state == StateCatalog {
				tvApp.catalogView.Layout(gtx, tvApp.theme)
			} else if tvApp.state == StateReader && tvApp.readerView != nil {
				tvApp.readerView.Layout(gtx, tvApp.theme)
			}

			// Render Banner Auto-Update jika ada update
			tvApp.renderUpdateBanner(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

func (app *TVApp) checkAndApplyUpdate(u *updater.Updater, w *app.Window) {
	// Berikan jeda 2 detik saat startup agar TV fokus memuat katalog terlebih dahulu
	time.Sleep(2 * time.Second)

	rel, hasUpdate, err := u.CheckUpdate()
	if err != nil {
		log.Printf("[Updater] Gagal cek update: %v", err)
		return
	}
	if !hasUpdate || rel == nil {
		log.Printf("[Updater] Aplikasi sudah versi terbaru (%s)", CurrentAppVersion)
		return
	}

	apkAsset := rel.FindAPKAsset()
	if apkAsset == nil {
		log.Printf("[Updater] Release %s tidak memiliki aset APK", rel.TagName)
		return
	}

	app.updateStatus.mu.Lock()
	app.updateStatus.available = true
	app.updateStatus.releaseTag = rel.TagName
	app.updateStatus.downloadURL = apkAsset.BrowserDownloadURL
	app.updateStatus.downloading = true
	app.updateStatus.progress = 0
	app.updateStatus.mu.Unlock()
	w.Invalidate()

	log.Printf("[Updater] Mengunduh pembaruan %s dari %s...", rel.TagName, apkAsset.BrowserDownloadURL)
	targetDir := updater.GetWritableUpdateDir()

	savedPath, dErr := u.DownloadAndInstallAPK(apkAsset.BrowserDownloadURL, targetDir, func(pct float32) {
		app.updateStatus.mu.Lock()
		app.updateStatus.progress = pct
		app.updateStatus.mu.Unlock()
		w.Invalidate()
	})

	app.updateStatus.mu.Lock()
	app.updateStatus.downloading = false
	if dErr != nil {
		app.updateStatus.errorMessage = dErr.Error()
		log.Printf("[Updater] Gagal mengunduh APK: %v", dErr)
	} else {
		app.updateStatus.downloaded = true
		app.updateStatus.savedPath = savedPath
		log.Printf("[Updater] Pembaruan tersimpan di %s dan installer dipicu.", savedPath)
	}
	app.updateStatus.mu.Unlock()
	w.Invalidate()
}

func (app *TVApp) renderUpdateBanner(gtx layout.Context) {
	app.updateStatus.mu.Lock()
	avail := app.updateStatus.available
	downloading := app.updateStatus.downloading
	progress := app.updateStatus.progress
	downloaded := app.updateStatus.downloaded
	savedPath := app.updateStatus.savedPath
	errMsg := app.updateStatus.errorMessage
	tag := app.updateStatus.releaseTag
	app.updateStatus.mu.Unlock()

	if !avail && errMsg == "" {
		return
	}

	layout.S.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			bannerWidth := gtx.Constraints.Max.X
			if bannerWidth > 800 {
				bannerWidth = 800
			}
			bannerHeight := 58

			// Center horizontally
			offsetX := (gtx.Constraints.Max.X - bannerWidth) / 2
			op.Offset(image.Pt(offsetX, 0)).Add(gtx.Ops)

			// Background Box
			bgRect := image.Rect(0, 0, bannerWidth, bannerHeight)
			clipArea := clip.RRect{
				Rect: bgRect,
				SE:   10, SW: 10, NW: 10, NE: 10,
			}.Push(gtx.Ops)
			paint.ColorOp{Color: color.NRGBA{R: 24, G: 28, B: 36, A: 245}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			// Border Glow
			borderGlow := clip.Stroke{
				Path:  clip.RRect{Rect: bgRect, SE: 10, SW: 10, NW: 10, NE: 10}.Path(gtx.Ops),
				Width: 2,
			}.Op().Push(gtx.Ops)
			paint.ColorOp{Color: color.NRGBA{R: 0, G: 200, B: 255, A: 220}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			borderGlow.Pop()

			dim := layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Top: unit.Dp(16), Bottom: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				var msg string
				var textColor color.NRGBA

				if errMsg != "" {
					msg = fmt.Sprintf("⚠️ Pembaruan %s gagal: %s", tag, errMsg)
					textColor = color.NRGBA{R: 255, G: 100, B: 100, A: 255}
				} else if downloaded {
					msg = fmt.Sprintf("✅ Pembaruan %s berhasil diunduh (%s). Membuka installer TV...", tag, savedPath)
					textColor = color.NRGBA{R: 76, G: 217, B: 100, A: 255}
				} else if downloading {
					msg = fmt.Sprintf("📥 Mengunduh Pembaruan %s: %.0f%%...", tag, progress*100)
					textColor = color.NRGBA{R: 0, G: 220, B: 255, A: 255}
				} else {
					msg = fmt.Sprintf("🚀 Pembaruan %s tersedia! Menyiapkan unduhan...", tag)
					textColor = color.NRGBA{R: 255, G: 215, B: 0, A: 255}
				}

				lbl := material.Body1(app.theme, msg)
				lbl.Color = textColor
				return lbl.Layout(gtx)
			})

			clipArea.Pop()
			return dim
		})
	})
}

