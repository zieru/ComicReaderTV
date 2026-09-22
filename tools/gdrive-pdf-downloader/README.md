# ⚡ Google Drive View-Only PDF Extractor & Importer

Tool canggih untuk mengekstrak dan merekonstruksi dokumen PDF yang terkunci / **View-Only** (pemilik menonaktifkan izin download/cetak) di Google Drive, serta mengimpornya langsung ke **Comic Reader TV**.

---

## 🔍 Mengapa Google Drive View-Only PDF Sulit Diunduh?

Saat pemilik file Google Drive menonaktifkan opsi download ("Pengakses lihat-saja dan pemberi komentar tidak dapat melihat opsi untuk mendownload, mencetak, dan menyalin"):
- Endpoint unduhan Google Drive (`/uc?export=download`) memblokir akses dan mengembalikan pesan error HTML proteksi izin.
- Server backend tidak bisa mendownload biner PDF secara langsung.
- Namun, antarmuka browser Google Drive viewer tetap harus merender tiap halaman sebagai gambar resolusi tinggi (`blob:` image tiles) ke layar agar pengguna bisa membacanya.

---

## 🚀 Peningkatan Signifikan Dibanding Repositori Referensi (`karimelmasry42`)

Repositori referensi (`karimelmasry42/google-drive-view-only-pdf-downloader`) memiliki sejumlah limitasi fatal yang telah kami perbaiki sepenuhnya:

| Fitur | Repositori Referensi (`karimelmasry42`) | **ComicReader Extractor (Improved)** |
|---|---|---|
| **Navigasi Halaman** | ❌ **Manual**: User harus scroll 100+ halaman sendiri sebelum menjalankan script | ✅ **Auto-Scroll Otomatis**: Script men-scroll viewer secara otomatis dari awal hingga akhir |
| **Ketahanan DOM Virtualization** | ❌ **Gagal/Halaman Hilang**: Google Drive menghapus elemen DOM atas saat di-scroll ke bawah, sehingga script lama hanya mengambil beberapa halaman terakhir | ✅ **Virtualization-Resistant**: Tiap halaman langsung disimpan ke memory buffer saat muncul, sehingga tidak ada halaman yang hilang |
| **Progress & UI Feedback** | ❌ **Buta di Konsol**: Tidak ada indikator progress atau status | ✅ **Modern Glassmorphic HUD**: Floating UI di layar dengan progress bar real-time, thumbnail preview, dan estimasi |
| **Kontrol Ekstraksi** | ❌ Tidak ada tombol pause/stop | ✅ Tombol **Jeda**, **Lanjutkan**, dan **Proses Sekarang (Early Finish)** |
| **Integrasi Comic Reader TV** | ❌ Hanya unduh lokal ke PC | ✅ **Direct Upload**: Tombol 1-klik untuk langsung mengirim PDF ke server Comic Reader TV (`ca.tsel.my.id:8080`) |
| **Deteksi Judul Otomatis** | ❌ Default nama `download.pdf` | ✅ Deteksi otomatis judul dokumen dari header Google Drive |

---

## 📖 Cara Penggunaan

### Cara 1: Menggunakan Bookmarklet (Paling Praktis)
1. Buka Web Admin Comic Reader TV (`http://ca.tsel.my.id:8080/` atau lokal).
2. Di bagian **"Bypass Google Drive View-Only"**, tarik tombol **"⚡ Import GDrive PDF"** ke Bookmarks Bar browser Anda.
3. Buka tab file PDF di Google Drive yang terproteksi.
4. Klik bookmarklet tersebut dari Bookmarks Bar.
5. HUD akan muncul di pojok kanan bawah dan secara otomatis memindai seluruh halaman dokumen.
6. Klik **"🚀 Kirim Langsung ke Server TV"** untuk langsung membaca komik di TV!

### Cara 2: Melalui Developer Tools Console
1. Buka file PDF yang terproteksi di Google Drive viewer.
2. Buka DevTools (`F12` atau `Ctrl + Shift + I`), pilih tab **Console**.
3. Buka file [`gdrive_pdf_extractor.js`](./gdrive_pdf_extractor.js), salin seluruh kodenya, dan paste ke Console lalu tekan `Enter`.
4. Ekstraksi otomatis akan berjalan hingga selesai.

---

## 🛠 Endpoint API Server Terkait

Server Comic Reader TV menyediakan endpoint multipart untuk menerima PDF hasil ekstraksi maupun file PDF lokal:

- **Method**: `POST`
- **Path**: `/api/comics/upload`
- **Headers**: `Content-Type: multipart/form-data`
- **Form Fields**:
  - `title`: Judul komik (string)
  - `description`: Deskripsi komik (opsional)
  - `cover_url`: URL cover komik (opsional)
  - `file`: Biner dokumen PDF (wajib)

Endpoint ini mendukung CORS lengkap sehingga aman dipanggil langsung dari origin `drive.google.com`.
