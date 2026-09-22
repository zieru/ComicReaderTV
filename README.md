# 📺 Comic Reader untuk Android TV & Desktop

Aplikasi pembaca komik/manga performa tinggi berbasis **Go** murni:
- **Client (Android TV & Desktop):** Dibangun dengan **Gio UI** (rendering GPU langsung via Vulkan/OpenGL, navigasi D-pad remote TV, dan fitur **Floating Loupe** / Kaca Pembesar Presisi).
- **Backend (Catalog API Provider):** Menyediakan REST API, otomatis mengurai link **Google Drive (View/Preview/Embed)** atau direct PDF, dan mengekstrak lembaran gambar halaman on-demand.
- **Auto-Update & CI/CD:** Workflow GitHub Actions untuk build APK otomatis dan in-app updater.

---

## 📁 Struktur Proyek

```
comic_reader/
├── server/                     # Catalog Server & PDF Streamer
│   ├── cmd/server/             # Entry point server
│   ├── internal/
│   │   ├── api/                # REST API: /api/comics, /api/comics/{id}/page/{num}
│   │   ├── gdrive/             # Google Drive URL parser & large file downloader
│   │   ├── pdfengine/          # PDF image extractor & page counter
│   │   └── store/              # JSON catalog storage
│   └── web/                    # Web Admin Dashboard (Embed Go)
│
├── client/                     # Client Gio UI (Android TV & Desktop)
│   ├── cmd/tv/                 # Entry point aplikasi TV
│   ├── ui/
│   │   ├── catalog_view.go     # Grid komik dengan D-pad TV focus effect
│   │   ├── reader_view.go      # Reader gambar fullscreen & HUD overlay
│   │   ├── magnifier.go        # Fitur Kaca Pembesar Melayang (Floating Loupe)
│   │   └── dpad.go             # Pemetaan tombol Hardware Remote TV
│   └── updater/                # In-App Updater via GitHub Releases
│
├── pkg/
│   └── model/                  # Model data bersama (Comic entity)
│
└── .github/
    └── workflows/
        └── build-android.yml   # Workflow GitHub Actions Build APK
```

---

## 🚀 Cara Menjalankan

### 1. Menjalankan Server Katalog
```bash
go run ./server/cmd/server --port=8080
```
- Buka dashboard di browser: `http://localhost:8080`
- Anda dapat memasukkan link komik baru berupa link **Google Drive** (misal `https://drive.google.com/file/d/.../preview`) atau link langsung file PDF.

### 2. Menjalankan Client di Komputer (Mode Uji Coba Remote TV)
```bash
go run ./client/cmd/tv --server="http://127.0.0.1:8080"
```

#### Kontrol Navigasi (Remote TV / Keyboard PC):
| Tombol Remote TV | Tombol Keyboard PC | Fungsi di Katalog | Fungsi di Reader |
| :--- | :--- | :--- | :--- |
| **D-Pad Kanan** | Panah Kanan | Pindah kartu ke kanan | Halaman Berikutnya |
| **D-Pad Kiri** | Panah Kiri | Pindah kartu ke kiri | Halaman Sebelumnya |
| **D-Pad Atas** | Panah Atas | Pindah kartu ke atas | Gerakkan Lensa Kaca Pembesar ke Atas |
| **D-Pad Bawah** | Panah Bawah | Pindah kartu ke bawah | Gerakkan Lensa Kaca Pembesar ke Bawah |
| **Tombol OK / Enter** | Enter / Space | Buka komik yang dipilih | Munculkan / Sembunyikan Menu HUD |
| **Tombol Play / Z** | Tombol `Z` | - | **Aktifkan / Matikan Kaca Pembesar (Loupe)** |
| **Tombol Back** | Escape | Keluar aplikasi | Tutup Loupe / Kembali ke Katalog |

---

## 🤖 Build APK untuk Android TV & GitHub Actions

### Kompilasi Lokal ke APK Android
Pastikan Android SDK & NDK sudah terpasang, lalu gunakan tool Gio:
```bash
go install gioui.org/cmd/gogio@latest
gogio -target android -appversion 1.0.0 -appid com.comicreader.tv -o comic_reader_tv.apk ./client/cmd/tv
```

### Auto-Build via GitHub Actions
Cukup lakukan tag rilis baru di git:
```bash
git tag v1.0.0
git push origin v1.0.0
```
GitHub Actions akan secara otomatis mengompilasi APK dan merilisnya di menu **Releases** GitHub Anda.
