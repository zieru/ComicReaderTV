package gscraper

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
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
// mengekstrak semua halaman dokumen, dan menyimpannya sebagai file PDF di outputPath.
func DownloadViewOnlyPDF(ctx context.Context, fileID string, outputPath string) error {
	binPath, has := launcher.LookPath()
	if !has {
		return errors.New("browser Chromium/Chrome belum terpasang di sistem server. Di Debian/Ubuntu, jalankan: 'sudo apt install -y chromium'")
	}

	log.Printf("[gscraper] Memulai headless Chromium (%s) untuk Google Drive file ID: %s", binPath, fileID)

	// Konfigurasi launcher dengan opsi aman untuk environment server (Debian/Linux VPS)
	l := launcher.New().
		Bin(binPath).
		Headless(true).
		NoSandbox(true).
		Set("disable-dev-shm-usage").
		Set("disable-gpu").
		Set("disable-setuid-sandbox").
		Set("no-first-run").
		Set("no-default-browser-check")

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

	// Beri jeda 2 detik agar runtime scripts Google Drive siap
	time.Sleep(2 * time.Second)

	// Injeksi skrip ekstraksi PDF ke dalam headless page
	extractorJS := `
	async () => {
		// 1. Muat jsPDF
		async function loadJsPDF() {
			if (window.jspdf && window.jspdf.jsPDF) return window.jspdf.jsPDF;
			const cdnUrl = 'https://cdnjs.cloudflare.com/ajax/libs/jspdf/2.5.1/jspdf.umd.min.js';
			let trustedURL = cdnUrl;
			if (window.trustedTypes && trustedTypes.createPolicy) {
				try {
					const policy = trustedTypes.createPolicy('crScraperPolicy', { createScriptURL: (i) => i });
					trustedURL = policy.createScriptURL(cdnUrl);
				} catch(e) {}
			}
			await new Promise((resolve, reject) => {
				const s = document.createElement("script");
				s.src = trustedURL;
				s.onload = resolve;
				s.onerror = reject;
				document.head.appendChild(s);
			});
			return window.jspdf.jsPDF;
		}

		const jsPDF = await loadJsPDF();

		// 2. Temukan scroll container
		const candidates = [
			document.querySelector('.drive-viewer-paginated-scrollable'),
			document.querySelector('.ndfHFb-c4YZDc-Wrql6b'),
			document.querySelector('[role="document"]')?.parentElement,
			document.documentElement,
			document.body
		];
		let scroller = document.scrollingElement || document.documentElement;
		for (const el of candidates) {
			if (el && (el.scrollHeight > el.clientHeight || el === document.body)) {
				scroller = el;
				break;
			}
		}

		// 3. Auto-scroll dan tangkap halaman
		const capturedPages = new Map();
		const scrollStep = Math.max(300, Math.floor(scroller.clientHeight * 0.8));
		let currentScroll = 0;
		const maxScroll = scroller.scrollHeight;

		function getPageNumberFromElement(img) {
			let curr = img;
			while (curr && curr !== document.body) {
				for (const attr of ["data-page-number", "data-page-index", "aria-label"]) {
					const val = curr.getAttribute(attr);
					if (val) {
						const match = val.match(/\b(\d+)\b/);
						if (match) return parseInt(match[1], 10);
					}
				}
				curr = curr.parentElement;
			}
			return null;
		}

		async function captureVisible() {
			const imgs = Array.from(document.querySelectorAll("img")).filter(img => {
				return /^blob:/.test(img.src) || (img.src && img.src.includes("googleusercontent.com"));
			});
			imgs.sort((a, b) => a.getBoundingClientRect().top - b.getBoundingClientRect().top);

			for (const img of imgs) {
				if (!img.complete || img.naturalWidth === 0) {
					await new Promise(r => {
						img.addEventListener("load", r, { once: true });
						setTimeout(r, 1200);
					});
				}
				const w = img.naturalWidth || img.width;
				const h = img.naturalHeight || img.height;
				if (w < 100 || h < 100) continue;

				let pageNum = getPageNumberFromElement(img);
				if (!pageNum) pageNum = capturedPages.size + 1;

				const existing = capturedPages.get(pageNum);
				if (!existing || (w * h > existing.width * existing.height)) {
					try {
						const canvas = document.createElement("canvas");
						canvas.width = w;
						canvas.height = h;
						const ctx = canvas.getContext("2d");
						ctx.drawImage(img, 0, 0, w, h);
						capturedPages.set(pageNum, {
							dataUrl: canvas.toDataURL("image/jpeg", 0.94),
							width: w,
							height: h
						});
					} catch(e) {}
				}
			}
		}

		while (currentScroll <= maxScroll + scrollStep) {
			scroller.scrollTop = currentScroll;
			window.scrollTo(0, currentScroll);
			await captureVisible();
			await new Promise(r => setTimeout(r, 500));
			currentScroll += scrollStep;
			if (currentScroll > scroller.scrollHeight) break;
		}

		if (capturedPages.size === 0) {
			throw new Error("Tidak ada gambar halaman yang terdeteksi di dokumen Google Drive ini");
		}

		// 4. Rekonstruksi PDF
		const sortedKeys = Array.from(capturedPages.keys()).sort((a, b) => a - b);
		let pdf = null;
		for (const pNum of sortedKeys) {
			const page = capturedPages.get(pNum);
			const w = page.width;
			const h = page.height;
			const orientation = w >= h ? "landscape" : "portrait";
			if (!pdf) {
				pdf = new jsPDF({ orientation, unit: "pt", format: [w, h], compress: true });
			} else {
				pdf.addPage([w, h], orientation);
			}
			pdf.addImage(page.dataUrl, "JPEG", 0, 0, w, h);
		}

		// Kembalikan base64
		const rawBase64 = pdf.output("datauristring").split(",")[1];
		return {
			pageCount: sortedKeys.length,
			base64: rawBase64
		};
	}
	`

	log.Printf("[gscraper] Mengeksekusi auto-scroll dan ekstraksi halaman di headless page...")
	evalRes, err := pageWithCtx.Eval(extractorJS)
	if err != nil {
		return fmt.Errorf("ekstraksi headless gagal: %w", err)
	}

	base64Data := evalRes.Value.Get("base64").String()
	pageCount := evalRes.Value.Get("pageCount").Int()
	if base64Data == "" {
		return errors.New("tidak ada data PDF yang dihasilkan dari headless scraper")
	}

	log.Printf("[gscraper] Berhasil mengekstrak %d halaman, mendecode base64 PDF...", pageCount)
	pdfBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("gagal mendecode data base64 PDF: %w", err)
	}

	// Validasi magic PDF
	if len(pdfBytes) < 4 || string(pdfBytes[:4]) != "%PDF" {
		return errors.New("data hasil headless scraper bukan dokumen PDF valid")
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

	log.Printf("[gscraper] PDF berhasil disimpan ke: %s (%d bytes)", outputPath, len(pdfBytes))
	return nil
}
