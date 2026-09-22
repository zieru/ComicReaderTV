package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"comic_reader/server/internal/api"
	"comic_reader/server/internal/auth"
	"comic_reader/server/internal/pdfengine"
	"comic_reader/server/internal/store"
	"comic_reader/server/internal/updater"
	"comic_reader/server/web"
)

var Version = "1.0.6"

func main() {
	port := flag.Int("port", 8080, "Port untuk server katalog")
	dataDir := flag.String("data", "./data", "Direktori penyimpanan data dan cache")
	botToken := flag.String("bot-token", "8978586484:AAFfLbux2a-88MLJbplns9Kz4VDfJzdtgi0", "Telegram Bot Token untuk autentikasi OTP")
	chatID := flag.Int64("chat-id", 0, "Telegram Chat ID Admin (opsional, dapat terdeteksi otomatis via /start)")
	showVersion := flag.Bool("v", false, "Tampilkan versi server")
	doUpdate := flag.Bool("update", false, "Periksa dan pasang pembaruan terbaru dari GitHub")
	flag.Parse()

	if envToken := os.Getenv("TELEGRAM_BOT_TOKEN"); envToken != "" && *botToken == "8978586484:AAFfLbux2a-88MLJbplns9Kz4VDfJzdtgi0" {
		*botToken = envToken
	}

	if *showVersion {
		fmt.Printf("Comic Reader Catalog Server v%s\n", Version)
		os.Exit(0)
	}

	if *doUpdate {
		fmt.Printf("Memeriksa pembaruan untuk versi saat ini: v%s...\n", Version)
		res, err := updater.CheckLatestRelease("zieru/ComicReaderTV", Version)
		if err != nil {
			log.Fatalf("Gagal memeriksa pembaruan: %v", err)
		}
		if !res.HasUpdate {
			fmt.Println("Server Anda sudah menggunakan versi terbaru!")
			os.Exit(0)
		}
		fmt.Printf("Versi baru ditemukan: %s. Memulai instalasi...\n", res.LatestVersion)
		if err := updater.PerformDebUpdate(res.DownloadURL); err != nil {
			log.Fatalf("Gagal melakukan pembaruan: %v", err)
		}
		fmt.Printf("Pembaruan ke %s berhasil dipasang! Service sedang di-restart.\n", res.LatestVersion)
		os.Exit(0)
	}

	log.Printf("Memulai Comic Reader Catalog Server v%s di port %d...", Version, *port)

	// 1. Inisialisasi Store
	catalogStore, err := store.New(filepath.Join(*dataDir, "store"))
	if err != nil {
		log.Fatalf("Gagal inisialisasi store: %v", err)
	}

	// 2. Inisialisasi PDF Engine & Cache
	pdfEngine, err := pdfengine.New(filepath.Join(*dataDir, "pdf_cache"))
	if err != nil {
		log.Fatalf("Gagal inisialisasi pdf engine: %v", err)
	}

	// 3. Inisialisasi Auth Manager (Telegram Bot OTP)
	authMgr := auth.NewManager(*botToken, filepath.Join(*dataDir, "auth"), *chatID)

	// 4. Setup HTTP Handler
	mux := http.NewServeMux()
	apiHandler := api.NewHandler(catalogStore, pdfEngine, authMgr, Version)
	apiHandler.RegisterRoutes(mux)

	// 4. Serve Web Admin Dashboard
	mux.Handle("/", http.FileServer(http.FS(web.Files)))

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server siap! Buka http://localhost:%d di browser Anda.", *port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server berhenti dengan error: %v", err)
	}
}
