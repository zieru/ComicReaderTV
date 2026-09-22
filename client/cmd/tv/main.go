package main

import (
	"flag"
	"log"
	"os"


	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"comic_reader/client/updater"


	"comic_reader/client/ui"
	"comic_reader/pkg/model"
)

type AppState int

const (
	StateCatalog AppState = iota
	StateReader
)

type TVApp struct {
	state       AppState
	catalogView *ui.CatalogView
	readerView  *ui.ReaderView
	theme       *material.Theme
	serverURL   string
}

func main() {
	serverURL := flag.String("server", "http://127.0.0.1:8080", "Alamat server katalog")
	flag.Parse()

	go func() {
		// Buat jendela Fullscreen standar TV (1080p / rasio 16:9)
		w := new(app.Window)
		w.Option(
			app.Title("Comic Reader TV"),
			app.Size(unit.Dp(1280), unit.Dp(720)),
			app.MinSize(unit.Dp(800), unit.Dp(450)),
		)

		// Inisialisasi Auto Updater dengan repo zieru/ComicReaderTV
		appUpdater := updater.NewUpdater("zieru", "ComicReaderTV", "v1.0.0")
		go func() {
			if rel, hasUpdate, err := appUpdater.CheckUpdate(); err == nil && hasUpdate {
				log.Printf("Pembaruan baru tersedia: %s (%s)", rel.Name, rel.TagName)
			}
		}()

		if err := run(w, *serverURL); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func run(w *app.Window, serverURL string) error {
	th := material.NewTheme()
	// material.NewTheme() di Gio versi terbaru sudah otomatis mengikutsertakan collection font default

	tvApp := &TVApp{
		state:     StateCatalog,
		theme:     th,
		serverURL: serverURL,
	}

	tvApp.catalogView = ui.NewCatalogView(serverURL, func(selected model.Comic) {
		tvApp.readerView = ui.NewReaderView(serverURL, selected, func() {
			tvApp.state = StateCatalog
			w.Invalidate()
		})
		tvApp.state = StateReader
		w.Invalidate()
	})

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

			e.Frame(gtx.Ops)
		}
	}
}
