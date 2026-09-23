package gscraper

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// HasBrowser mengembalikan true jika biner Chromium/Chrome ditemukan di sistem
func HasBrowser() bool {
	_, has := launcher.LookPath()
	return has
}

// BrowserPath mengembalikan path biner browser jika ditemukan
func BrowserPath() string {
	path, _ := launcher.LookPath()
	return path
}

// DownloadViewOnlyPDF membuka file Google Drive view-only menggunakan headless browser,
// mengekstrak semua halaman dokumen dari DOM canvas/img, dan menyimpannya sebagai file PDF di outputPath.
func DownloadViewOnlyPDF(ctx context.Context, fileID string, outputPath string) error {
	binPath, has := launcher.LookPath()
	if !has {
		return errors.New("browser Chromium/Chrome belum terpasang di sistem server. Di Debian/Ubuntu, jalankan: 'sudo apt install -y chromium'")
	}

	log.Printf("[gscraper] Memulai headless Chromium (%s) untuk Google Drive file ID: %s", binPath, fileID)

	tempUserDir, err := os.MkdirTemp("", "cr-rod-user-*")
	if err == nil {
		defer os.RemoveAll(tempUserDir)
	}

	// Konfigurasi launcher dengan opsi aman untuk environment server (Debian/Linux VPS & Windows)
	l := launcher.New().
		Bin(binPath).
		Headless(true).
		NoSandbox(true).
		Leakless(false).
		Set("disable-dev-shm-usage").
		Set("disable-gpu").
		Set("disable-setuid-sandbox").
		Set("no-first-run").
		Set("no-default-browser-check").
		Set("window-size", "1920,1080").
		Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	if tempUserDir != "" {
		l = l.UserDataDir(tempUserDir)
	}

	controlURL, err := l.Launch()
	if err != nil {
		return fmt.Errorf("gagal meluncurkan browser headless: %w", err)
	}

	browser := rod.New().ControlURL(controlURL)
	if err := browser.Connect(); err != nil {
		return fmt.Errorf("gagal terhubung ke browser: %w", err)
	}
	defer func() {
		_ = browser.Close()
		l.Kill()
	}()

	previewURL := fmt.Sprintf("https://drive.google.com/file/d/%s/preview", fileID)
	log.Printf("[gscraper] Membuka URL preview: %s", previewURL)

	page, err := browser.Page(proto.TargetCreateTarget{URL: previewURL})
	if err != nil {
		return fmt.Errorf("gagal membuat tab halaman: %w", err)
	}
	defer page.Close()

	pageWithCtx := page.Context(ctx)

	// Tunggu halaman selesai dimuat
	if err := pageWithCtx.WaitLoad(); err != nil {
		return fmt.Errorf("gagal memuat halaman preview: %w", err)
	}

	// Beri jeda agar runtime scripts Google Drive siap
	time.Sleep(3 * time.Second)

	// Injeksi skrip ekstraksi PDF ke dalam headless page tanpa library eksternal (CSP-safe)
	extractorJS := `
	async () => {
		// 1. Temukan scroll container Google Drive Viewer
		const scroller = document.querySelector('.ndfHFb-c4YZDc-s2gQvd') ||
			Array.from(document.querySelectorAll('*')).find(el => {
				const s = window.getComputedStyle(el);
				return (s.overflowY === 'auto' || s.overflowY === 'scroll') && el.scrollHeight > el.clientHeight;
			}) || document.scrollingElement || document.documentElement;

		const capturedPages = new Map();
		let pageHeight = 1147;
		let foundAny = false;

		function captureVisible() {
			const imgs = Array.from(document.querySelectorAll("img")).filter(i => {
				return i.src && (i.src.includes("drive-viewer") || i.src.includes("googleusercontent.com") || i.src.startsWith("blob:"));
			});

			for (const img of imgs) {
				const w = img.naturalWidth || img.width;
				const h = img.naturalHeight || img.height;
				if (w < 150 || h < 150) continue;

				if (!foundAny && h > 400) {
					pageHeight = h + 16;
					foundAny = true;
				}

				const rect = img.getBoundingClientRect();
				const absoluteTop = scroller.scrollTop + rect.top;
				const pageNum = Math.max(1, Math.round(absoluteTop / pageHeight) + 1);

				if (!capturedPages.has(pageNum)) {
					try {
						const canvas = document.createElement("canvas");
						canvas.width = w;
						canvas.height = h;
						const ctx = canvas.getContext("2d");
						ctx.drawImage(img, 0, 0, w, h);
						const dataUrl = canvas.toDataURL("image/jpeg", 0.88);
						capturedPages.set(pageNum, dataUrl);
					} catch(e) {}
				}
			}
		}

		// Initial capture
		captureVisible();

		// Auto-scroll loop
		const scrollStep = Math.max(1500, Math.floor(scroller.clientHeight * 1.5));
		let currentScroll = 0;
		const maxScroll = scroller.scrollHeight;
		let lastCount = 0;
		let stallCount = 0;

		while (currentScroll <= maxScroll + scrollStep) {
			scroller.scrollTop = currentScroll;
			await new Promise(r => setTimeout(r, 250));
			captureVisible();

			if (capturedPages.size === lastCount) {
				stallCount++;
			} else {
				stallCount = 0;
				lastCount = capturedPages.size;
			}

			currentScroll += scrollStep;
			if (currentScroll > scroller.scrollHeight) {
				if (stallCount >= 3) break;
			}
		}

		// Pastikan posisi paling bawah tertangkap
		scroller.scrollTop = scroller.scrollHeight;
		await new Promise(r => setTimeout(r, 400));
		captureVisible();

		if (capturedPages.size === 0) {
			throw new Error("Tidak ada gambar halaman yang terdeteksi di dokumen Google Drive ini");
		}

		// Kembalikan daftar halaman terurut
		const sortedKeys = Array.from(capturedPages.keys()).sort((a, b) => a - b);
		const pages = [];
		for (const k of sortedKeys) {
			pages.push({
				pageNum: k,
				dataUrl: capturedPages.get(k)
			});
		}

		return {
			count: pages.length,
			pages: pages
		};
	}
	`

	log.Printf("[gscraper] Mengeksekusi auto-scroll dan ekstraksi canvas halaman...")
	evalRes, err := pageWithCtx.Eval(extractorJS)
	if err != nil {
		return fmt.Errorf("ekstraksi headless gagal: %w", err)
	}

	pagesVal := evalRes.Value.Get("pages")
	pageCount := evalRes.Value.Get("count").Int()
	if pageCount == 0 {
		return errors.New("tidak ada halaman komik yang dihasilkan dari headless scraper")
	}

	log.Printf("[gscraper] Berhasil mengumpulkan %d halaman, mendekode JPEG...", pageCount)
	var jpegs [][]byte
	for i := 0; i < pageCount; i++ {
		p := pagesVal.Get(fmt.Sprintf("%d", i))
		dataUrl := p.Get("dataUrl").String()
		parts := strings.Split(dataUrl, ",")
		if len(parts) == 2 {
			raw, decodeErr := base64.StdEncoding.DecodeString(parts[1])
			if decodeErr == nil && len(raw) > 0 {
				jpegs = append(jpegs, raw)
			}
		}
	}

	if len(jpegs) == 0 {
		return errors.New("gagal mendekode data gambar JPEG dari headless scraper")
	}

	log.Printf("[gscraper] Menyusun dokumen PDF dari %d halaman JPEG...", len(jpegs))
	pdfBytes, err := BuildPDFFromJPEGs(jpegs)
	if err != nil {
		return fmt.Errorf("gagal menyusun PDF dari JPEG: %w", err)
	}

	// Simpan ke file tujuan
	tempOut := outputPath + ".tmp"
	if err := os.WriteFile(tempOut, pdfBytes, 0644); err != nil {
		return fmt.Errorf("gagal menulis file PDF: %w", err)
	}

	if err := os.Rename(tempOut, outputPath); err != nil {
		_ = os.Remove(tempOut)
		return fmt.Errorf("gagal rename file PDF: %w", err)
	}

	log.Printf("[gscraper] PDF berhasil disimpan ke: %s (%d bytes, %d halaman)", outputPath, len(pdfBytes), len(jpegs))
	return nil
}

// BuildPDFFromJPEGs membuat dokumen PDF standar dari kumpulan byte gambar JPEG
func BuildPDFFromJPEGs(jpegs [][]byte) ([]byte, error) {
	if len(jpegs) == 0 {
		return nil, fmt.Errorf("tidak ada gambar JPEG untuk dibuatkan PDF")
	}

	var buf bytes.Buffer
	buf.WriteString("%PDF-1.4\n%\xe2\xe3\xcf\xd3\n")

	type objOffset struct {
		id     int
		offset int
	}
	var offsets []objOffset

	numPages := len(jpegs)

	writeObj := func(id int, data string) {
		offsets = append(offsets, objOffset{id: id, offset: buf.Len()})
		buf.WriteString(fmt.Sprintf("%d 0 obj\n%s\nendobj\n", id, data))
	}

	// 1: Catalog
	writeObj(1, "<< /Type /Catalog /Pages 2 0 R >>")

	// Kids list for Pages
	var kids []string
	for i := 1; i <= numPages; i++ {
		pageObjID := 3 * i
		kids = append(kids, fmt.Sprintf("%d 0 R", pageObjID))
	}
	writeObj(2, fmt.Sprintf("<< /Type /Pages /Kids [ %s ] /Count %d >>", strings.Join(kids, " "), numPages))

	for i, imgData := range jpegs {
		pageIdx := i + 1
		pageObjID := 3 * pageIdx
		imgObjID := 3*pageIdx + 1
		contentObjID := 3*pageIdx + 2

		cfg, err := jpeg.DecodeConfig(bytes.NewReader(imgData))
		w, h := 800, 1131
		if err == nil && cfg.Width > 0 && cfg.Height > 0 {
			w, h = cfg.Width, cfg.Height
		}

		// Page object
		writeObj(pageObjID, fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [ 0 0 %d %d ] /Resources << /XObject << /Im1 %d 0 R >> >> /Contents %d 0 R >>", w, h, imgObjID, contentObjID))

		// Image XObject
		imgHeader := fmt.Sprintf("<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /DCTDecode /Length %d >>\nstream\n", w, h, len(imgData))
		offsets = append(offsets, objOffset{id: imgObjID, offset: buf.Len()})
		buf.WriteString(fmt.Sprintf("%d 0 obj\n%s", imgObjID, imgHeader))
		buf.Write(imgData)
		buf.WriteString("\nendstream\nendobj\n")

		// Content stream
		content := fmt.Sprintf("q %d 0 0 %d 0 0 cm /Im1 Do Q", w, h)
		writeObj(contentObjID, fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(content), content))
	}

	// XRef table
	startXref := buf.Len()
	maxObjID := 3*numPages + 2
	buf.WriteString(fmt.Sprintf("xref\n0 %d\n0000000000 65535 f \n", maxObjID+1))

	offsetMap := make(map[int]int)
	for _, o := range offsets {
		offsetMap[o.id] = o.offset
	}

	for id := 1; id <= maxObjID; id++ {
		off := offsetMap[id]
		buf.WriteString(fmt.Sprintf("%010d 00000 n \n", off))
	}

	buf.WriteString(fmt.Sprintf("trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", maxObjID+1, startXref))

	return buf.Bytes(), nil
}
