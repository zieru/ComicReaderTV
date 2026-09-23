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
	"gioui.org/io/pointer"
	"gioui.org/io/system"
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

const CurrentAppVersion = "v1.0.15"

var mainTag = new(int)

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

	// Filter lengkap untuk menangkap semua event keyboard/remote & pointer/touch TV
	// PENTING: Pada Android TV, tombol arah (Up/Down/Left/Right) dan Back diperlakukan
	// sebagai SystemEvent oleh Gio. Filter catch-all (Name: "") secara sengaja TIDAK
	// menangkap SystemEvent, sehingga setiap tombol D-pad HARUS didaftarkan secara eksplisit!
	// PENTING: key.FocusFilter{Target: mainTag} wajib ada agar Gio menetapkan status focusable
	// pada mainTag, sehingga AKEYCODE_DPAD_CENTER dapat diteruskan melalui ClickFocus()!
	filters := []event.Filter{
		key.FocusFilter{Target: mainTag},
		key.Filter{Focus: mainTag, Name: key.NameUpArrow},
		key.Filter{Focus: mainTag, Name: key.NameDownArrow},
		key.Filter{Focus: mainTag, Name: key.NameLeftArrow},
		key.Filter{Focus: mainTag, Name: key.NameRightArrow},
		key.Filter{Focus: mainTag, Name: key.NameReturn},
		key.Filter{Focus: mainTag, Name: key.NameEnter},
		key.Filter{Focus: mainTag, Name: key.NameBack},
		key.Filter{Focus: mainTag, Name: key.NameEscape},
		key.Filter{Focus: mainTag, Name: key.NameSpace},
		key.Filter{Focus: mainTag, Name: "DpadCenter"},
		key.Filter{Focus: mainTag, Name: "Select"},
		key.Filter{Focus: mainTag, Name: "Back"},
		key.Filter{Focus: mainTag, Name: "Menu"},
		key.Filter{Focus: mainTag, Name: "Z"},
		key.Filter{Focus: mainTag, Name: "z"},
		key.Filter{Focus: mainTag, Name: "MediaPlayPause"},
		key.Filter{Focus: mainTag, Name: "MediaPlay"},
		key.Filter{Focus: mainTag, Name: ""},
		// Global filters
		key.Filter{Name: key.NameUpArrow},
		key.Filter{Name: key.NameDownArrow},
		key.Filter{Name: key.NameLeftArrow},
		key.Filter{Name: key.NameRightArrow},
		key.Filter{Name: key.NameReturn},
		key.Filter{Name: key.NameEnter},
		key.Filter{Name: key.NameBack},
		key.Filter{Name: key.NameEscape},
		key.Filter{Name: key.NameSpace},
		key.Filter{Name: "DpadCenter"},
		key.Filter{Name: "Select"},
		key.Filter{Name: "Back"},
		key.Filter{Name: "Menu"},
		key.Filter{Name: "Z"},
		key.Filter{Name: "z"},
		key.Filter{Name: "MediaPlayPause"},
		key.Filter{Name: "MediaPlay"},
		key.Filter{Name: ""},
		// Pointer filter: menangkap AKEYCODE_DPAD_CENTER dan click TV
		pointer.Filter{
			Target: mainTag,
			Kinds:  pointer.Press | pointer.Release,
		},
	}

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Daftarkan area input mencakup seluruh layar jendela TV
			inputArea := clip.Rect(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)).Push(gtx.Ops)
			event.Op(gtx.Ops, mainTag)
			inputArea.Pop()

			// Pastikan keyboard focus selalu aktif pada mainTag
			if !gtx.Focused(mainTag) {
				gtx.Execute(key.FocusCmd{Tag: mainTag})
			}

			// Handle Hardware Key Events & Pointer Events (Remote TV)
			for {
				ev, ok := gtx.Event(filters...)
				if !ok {
					break
				}
				switch eventVal := ev.(type) {
				case key.Event:
					rKey := ui.MapKeyEvent(eventVal)
					if rKey != ui.KeyNone {
						log.Printf("[RemoteTV] Key event: %s (State: %v), mapped to: %d", eventVal.Name, eventVal.State, rKey)
						tvApp.handleRemoteKey(rKey, gtx)
						w.Invalidate()
					}
				case pointer.Event:
					if eventVal.Kind == pointer.Release {
						// Pemicuan klik / OK dari DpadCenter TV
						log.Printf("[RemoteTV] Pointer release (DpadCenter/Click) at: (%.1f, %.1f)", eventVal.Position.X, eventVal.Position.Y)
						tvApp.handleRemoteKey(ui.KeySelect, gtx)
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

func (tvApp *TVApp) handleRemoteKey(rKey ui.RemoteKey, gtx layout.Context) {
	screenWidth := float32(gtx.Constraints.Max.X)
	screenHeight := float32(gtx.Constraints.Max.Y)

	// Tangani tombol KELUAR (Back Button)
	if rKey == ui.KeyBack {
		if tvApp.state == StateReader {
			// Saat membaca, tombol Back kembali ke katalog
			tvApp.state = StateCatalog
			return
		}
		if tvApp.state == StateCatalog {
			// Saat di katalog utama, tombol Back keluar dari aplikasi
			tvApp.window.Perform(system.ActionMinimize)
			os.Exit(0)
			return
		}
	}

	if tvApp.state == StateCatalog {
		tvApp.catalogView.HandleKey(rKey)
	} else if tvApp.state == StateReader && tvApp.readerView != nil {
		tvApp.readerView.HandleKey(rKey, screenWidth, screenHeight)
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

