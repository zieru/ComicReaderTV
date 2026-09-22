package pdfengine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"sync"

	"comic_reader/server/internal/gdrive"

	"github.com/pdfcpu/pdfcpu/pkg/api"

)

type Engine struct {
	cacheDir    string
	gdriveCli   *gdrive.Client
	downloadMu  sync.Map
}

func New(cacheDir string) (*Engine, error) {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori cache: %w", err)
	}
	return &Engine{
		cacheDir:  cacheDir,
		gdriveCli: gdrive.NewClient(),
	}, nil
}

// EnsureLocalPDF memastikan file PDF diunduh dan tersimpan di cache lokal
func (e *Engine) EnsureLocalPDF(comicID, sourceURL string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(sourceURL))
	urlHash := hex.EncodeToString(hasher.Sum(nil))[:12]

	localPath := filepath.Join(e.cacheDir, fmt.Sprintf("%s_%s.pdf", comicID, urlHash))

	// Jika file sudah ada dan ukuran > 0, langsung gunakan
	if fi, err := os.Stat(localPath); err == nil && fi.Size() > 0 {
		return localPath, nil
	}

	// Gunakan mutex per file agar tidak terjadi download ganda serentak
	muAny, _ := e.downloadMu.LoadOrStore(localPath, &sync.Mutex{})
	mu := muAny.(*sync.Mutex)
	mu.Lock()
	defer mu.Unlock()

	// Double check setelah mendapatkan lock
	if fi, err := os.Stat(localPath); err == nil && fi.Size() > 0 {
		return localPath, nil
	}

	var reader io.ReadCloser
	var err error

	if gdrive.IsGDriveURL(sourceURL) {
		fileID, extractErr := gdrive.ExtractFileID(sourceURL)
		if extractErr != nil {
			return "", extractErr
		}
		reader, err = e.gdriveCli.DownloadFile(fileID)
		if err != nil {
			return "", fmt.Errorf("gagal download gdrive file: %w", err)
		}
	} else {
		// Link download HTTP/HTTPS langsung
		resp, httpErr := http.Get(sourceURL)
		if httpErr != nil {
			return "", fmt.Errorf("gagal download direct url: %w", httpErr)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return "", fmt.Errorf("direct url merespon status: %s", resp.Status)
		}
		reader = resp.Body
	}
	defer reader.Close()

	// Simpan ke temporary file terlebih dahulu
	tempPath := localPath + ".tmp"
	outFile, err := os.Create(tempPath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file temp: %w", err)
	}

	_, err = io.Copy(outFile, reader)
	outFile.Close()
	if err != nil {
		os.Remove(tempPath)
		return "", fmt.Errorf("gagal menyimpan stream PDF: %w", err)
	}

	// Validasi bahwa file yang diunduh adalah file PDF yang valid
	f, err := os.Open(tempPath)
	if err == nil {
		magic := make([]byte, 5)
		n, _ := io.ReadFull(f, magic)
		f.Close()
		if n < 4 || string(magic[:4]) != "%PDF" {
			os.Remove(tempPath)
			return "", fmt.Errorf("file yang diunduh bukan dokumen PDF yang valid (format tidak sesuai atau respon error tersimpan)")
		}
	}

	// Rename atomik
	if err := os.Rename(tempPath, localPath); err != nil {
		return "", fmt.Errorf("gagal rename file PDF: %w", err)
	}

	return localPath, nil
}

// GetPageCount mengembalikan jumlah halaman pada file PDF
func (e *Engine) GetPageCount(pdfPath string) (int, error) {
	return api.PageCountFile(pdfPath)
}

// ExtractPageImage mengekstrak gambar halaman komik dari PDF ke folder cache
// Komik PDF biasanya berisi gambar per halaman penuh.
func (e *Engine) ExtractPageImage(pdfPath string, pageNum int) (string, error) {
	pageImageDir := filepath.Join(e.cacheDir, "extracted", filepath.Base(pdfPath))
	if err := os.MkdirAll(pageImageDir, 0755); err != nil {
		return "", err
	}

	// Cek apakah sudah pernah diekstrak
	expectedPatterns := []string{
		fmt.Sprintf("page_%d.png", pageNum),
		fmt.Sprintf("page_%d.jpg", pageNum),
		fmt.Sprintf("page_%d.webp", pageNum),
	}
	for _, p := range expectedPatterns {
		candidate := filepath.Join(pageImageDir, p)
		if fi, err := os.Stat(candidate); err == nil && fi.Size() > 0 {
			return candidate, nil
		}
	}

	// Ekstrak gambar halaman dari PDF menggunakan pdfcpu api.ExtractImagesFile
	pages := []string{fmt.Sprintf("%d", pageNum)}
	err := api.ExtractImagesFile(pdfPath, pageImageDir, pages, nil)
	if err != nil {
		return "", fmt.Errorf("gagal mengekstrak gambar halaman %d: %w", pageNum, err)
	}

	// Cari file hasil ekstraksi di folder tersebut
	entries, err := os.ReadDir(pageImageDir)
	if err != nil {
		return "", err
	}

	// Format output pdfcpu biasanya menyertakan nomor halaman atau index
	for _, entry := range entries {
		name := entry.Name()
		if strings.Contains(name, fmt.Sprintf("_%d_", pageNum)) || strings.HasPrefix(name, fmt.Sprintf("%d_", pageNum)) {
			return filepath.Join(pageImageDir, name), nil
		}
	}

	// Jika nama file berbeda, ambil file termuda yang dibuat di folder
	if len(entries) > 0 {
		return filepath.Join(pageImageDir, entries[len(entries)-1].Name()), nil
	}

	return "", fmt.Errorf("tidak ditemukan gambar hasil ekstraksi untuk halaman %d", pageNum)
}

