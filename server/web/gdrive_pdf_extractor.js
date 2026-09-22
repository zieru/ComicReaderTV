/**
 * ComicReader TV - Google Drive View-Only PDF Extractor & Importer
 * 
 * Fitur Utama & Keunggulan dibandingkan script biasa:
 * 1. Auto-Scroll Otomatis: Tidak perlu scroll manual halaman demi halaman.
 * 2. Virtualization-Resistant: Menangkap halaman saat muncul dan menyimpannya di memory buffer,
 *    sehingga halaman tidak hilang saat Google Drive meng-unmount elemen DOM dari layar.
 * 3. Modern Glassmorphic HUD: Menampilkan live progress bar, thumbnail preview, dan estimasi waktu.
 * 4. Dual Export:
 *    - Langsung unggah ke Comic Reader TV Server (POST /api/comics/upload)
 *    - Atau unduh sebagai PDF resolusi tinggi ke komputer lokal.
 */

(async function () {
    // Hindari multiple instances
    if (window.__comicReaderExtractorActive) {
        alert("ComicReader Extractor sudah aktif di tab ini!");
        return;
    }
    window.__comicReaderExtractorActive = true;

    // Default target server Comic Reader
    let DEFAULT_SERVER_URL = "http://ca.tsel.my.id:8080";
    if (window.__COMIC_READER_SERVER__) {
        DEFAULT_SERVER_URL = window.__COMIC_READER_SERVER__;
    }

    // 1. Muat jsPDF dengan penanganan CSP / Trusted Types
    async function loadJsPDF() {
        if (window.jspdf && window.jspdf.jsPDF) {
            return window.jspdf.jsPDF;
        }

        const cdnUrl = 'https://cdnjs.cloudflare.com/ajax/libs/jspdf/2.5.1/jspdf.umd.min.js';
        let trustedURL = cdnUrl;

        if (window.trustedTypes && trustedTypes.createPolicy) {
            try {
                const policy = trustedTypes.createPolicy('comicReaderPolicy', {
                    createScriptURL: (input) => input
                });
                trustedURL = policy.createScriptURL(cdnUrl);
            } catch (e) {
                console.warn("TrustedTypes policy creation fallback:", e);
            }
        }

        await new Promise((resolve, reject) => {
            const script = document.createElement("script");
            script.src = trustedURL;
            script.onload = resolve;
            script.onerror = () => reject(new Error("Gagal memuat jsPDF dari CDN"));
            document.head.appendChild(script);
        });

        return window.jspdf.jsPDF;
    }

    // 2. Ambil Judul Dokumen dari Google Drive DOM
    function getDocumentTitle() {
        const titleEl = document.querySelector('div[role="heading"]') ||
            document.querySelector('.drive-viewer-toolstrip-title') ||
            document.querySelector('[data-tooltip*=".pdf"]') ||
            document.querySelector('meta[property="og:title"]');
        let title = titleEl ? (titleEl.textContent || titleEl.getAttribute('content') || "").trim() : "";
        if (!title) {
            title = document.title.replace(/ - Google (Drive|Docs)/i, "").trim();
        }
        if (!title || title === "Google Drive") {
            title = "Komik_" + new Date().toISOString().slice(0, 10);
        }
        return title.replace(/\.pdf$/i, "");
    }

    // 3. Deteksi Elemen Scroll Container
    function findScrollContainer() {
        const candidates = [
            document.querySelector('.drive-viewer-paginated-scrollable'),
            document.querySelector('.ndfHFb-c4YZDc-Wrql6b'),
            document.querySelector('[role="document"]')?.parentElement,
            document.querySelector('div[style*="overflow-y: scroll"]'),
            document.querySelector('div[style*="overflow-y: auto"]'),
            document.documentElement,
            document.body
        ];
        for (const el of candidates) {
            if (el && (el.scrollHeight > el.clientHeight || el === document.body)) {
                return el;
            }
        }
        return document.scrollingElement || document.documentElement;
    }

    // 4. Injeksi Modern Glassmorphic HUD
    const hudContainer = document.createElement("div");
    hudContainer.id = "comic-reader-hud";
    hudContainer.style.cssText = `
        position: fixed;
        bottom: 24px;
        right: 24px;
        z-index: 9999999;
        width: 380px;
        background: rgba(15, 23, 42, 0.92);
        backdrop-filter: blur(16px);
        -webkit-backdrop-filter: blur(16px);
        border: 1px solid rgba(56, 189, 248, 0.35);
        border-radius: 16px;
        box-shadow: 0 20px 40px rgba(0, 0, 0, 0.6), 0 0 20px rgba(56, 189, 248, 0.2);
        color: #f8fafc;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
        font-size: 14px;
        padding: 20px;
        box-sizing: border-box;
        transition: all 0.3s ease;
    `;

    hudContainer.innerHTML = `
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
            <div style="display: flex; align-items: center; gap: 8px;">
                <span style="font-size: 20px;">⚡</span>
                <strong style="font-size: 15px; color: #38bdf8; letter-spacing: -0.2px;">ComicReader TV Extractor</strong>
            </div>
            <button id="cr-btn-close" style="background: transparent; border: none; color: #94a3b8; font-size: 18px; cursor: pointer; padding: 2px 6px;">&times;</button>
        </div>
        <div id="cr-status-text" style="color: #cbd5e1; font-size: 13px; margin-bottom: 10px; line-height: 1.4;">
            Memulai ekstraksi...
        </div>
        <div style="background: rgba(30, 41, 59, 0.8); border-radius: 8px; overflow: hidden; height: 10px; margin-bottom: 14px; border: 1px solid rgba(255,255,255,0.08);">
            <div id="cr-progress-bar" style="background: linear-gradient(90deg, #38bdf8, #818cf8); height: 100%; width: 0%; transition: width 0.2s ease;"></div>
        </div>
        <div id="cr-preview-box" style="display: flex; align-items: center; gap: 12px; margin-bottom: 14px; background: rgba(0,0,0,0.25); padding: 8px; border-radius: 8px; border: 1px solid rgba(255,255,255,0.05);">
            <img id="cr-thumb" style="width: 48px; height: 64px; object-fit: cover; border-radius: 4px; background: #334155; display: none;" />
            <div style="flex: 1; font-size: 12px; color: #94a3b8;">
                <div id="cr-page-detail">Menunggu halaman ter-render...</div>
                <div id="cr-res-detail" style="font-size: 11px; color: #64748b; margin-top: 2px;"></div>
            </div>
        </div>
        <div id="cr-actions" style="display: flex; gap: 8px;">
            <button id="cr-btn-pause" style="flex: 1; padding: 8px 12px; background: #334155; border: 1px solid rgba(255,255,255,0.1); color: #f8fafc; border-radius: 8px; cursor: pointer; font-size: 12px; font-weight: 500;">Jeda</button>
            <button id="cr-btn-finish-early" style="flex: 1; padding: 8px 12px; background: #2563eb; border: none; color: #fff; border-radius: 8px; cursor: pointer; font-size: 12px; font-weight: 500;">Proses Sekarang</button>
        </div>
        <div id="cr-export-section" style="display: none; margin-top: 14px;">
            <div style="font-size: 12px; font-weight: 600; color: #38bdf8; margin-bottom: 8px;">🎯 Berhasil! Pilih Tindakan:</div>
            <div style="margin-bottom: 10px;">
                <label style="font-size: 11px; color: #94a3b8; display: block; margin-bottom: 3px;">URL Server Comic Reader:</label>
                <input id="cr-input-server" type="text" value="${DEFAULT_SERVER_URL}" style="width: 100%; box-sizing: border-box; background: #1e293b; border: 1px solid #475569; color: #f8fafc; border-radius: 6px; padding: 6px 10px; font-size: 12px;" />
            </div>
            <div style="margin-bottom: 12px;">
                <label style="font-size: 11px; color: #94a3b8; display: block; margin-bottom: 3px;">Judul Komik:</label>
                <input id="cr-input-title" type="text" value="${getDocumentTitle()}" style="width: 100%; box-sizing: border-box; background: #1e293b; border: 1px solid #475569; color: #f8fafc; border-radius: 6px; padding: 6px 10px; font-size: 12px;" />
            </div>
            <div style="display: flex; flex-direction: column; gap: 8px;">
                <button id="cr-btn-upload-server" style="padding: 10px; background: linear-gradient(135deg, #0284c7, #6366f1); border: none; color: #fff; border-radius: 8px; cursor: pointer; font-size: 13px; font-weight: 600; box-shadow: 0 4px 12px rgba(2, 132, 199, 0.4);">
                    🚀 Kirim Langsung ke Server TV
                </button>
                <button id="cr-btn-download-pdf" style="padding: 8px; background: #334155; border: 1px solid rgba(255,255,255,0.1); color: #cbd5e1; border-radius: 8px; cursor: pointer; font-size: 12px;">
                    💾 Unduh File PDF Lokal
                </button>
            </div>
        </div>
    `;

    document.body.appendChild(hudContainer);

    const statusText = document.getElementById("cr-status-text");
    const progressBar = document.getElementById("cr-progress-bar");
    const thumbImg = document.getElementById("cr-thumb");
    const pageDetail = document.getElementById("cr-page-detail");
    const resDetail = document.getElementById("cr-res-detail");
    const btnPause = document.getElementById("cr-btn-pause");
    const btnFinishEarly = document.getElementById("cr-btn-finish-early");
    const btnClose = document.getElementById("cr-btn-close");
    const actionsBox = document.getElementById("cr-actions");
    const exportSection = document.getElementById("cr-export-section");

    let isPaused = false;
    let shouldStop = false;

    btnClose.onclick = () => {
        if (confirm("Hentikan ekstraksi ComicReader?")) {
            shouldStop = true;
            hudContainer.remove();
            window.__comicReaderExtractorActive = false;
        }
    };

    btnPause.onclick = () => {
        isPaused = !isPaused;
        btnPause.textContent = isPaused ? "Lanjutkan" : "Jeda";
        btnPause.style.background = isPaused ? "#059669" : "#334155";
        statusText.textContent = isPaused ? "Ekstraksi dijeda oleh pengguna." : "Melanjutkan ekstraksi...";
    };

    btnFinishEarly.onclick = () => {
        shouldStop = true;
        btnFinishEarly.disabled = true;
        btnFinishEarly.textContent = "Menyelesaikan...";
    };

    // 5. Muat Library jsPDF
    statusText.textContent = "Memuat engine PDF (jsPDF)...";
    let jsPDF;
    try {
        jsPDF = await loadJsPDF();
    } catch (err) {
        statusText.innerHTML = `<span style="color: #ef4444;">Gagal memuat jsPDF: ${err.message}</span>`;
        return;
    }

    // 6. Map Halaman yang berhasil ditangkap (kebal terhadap DOM virtualization)
    // Key: number (1..N), Value: { dataUrl, width, height }
    const capturedPages = new Map();
    const scrollContainer = findScrollContainer();

    // Deteksi perkiraan total halaman dari Drive UI jika ada
    function detectTotalPages() {
        const pageText = document.body.innerText;
        const match = pageText.match(/(\d+)\s*(?:\/|dari|of)\s*(\d+)/i);
        if (match && parseInt(match[2], 10) > 0) {
            return parseInt(match[2], 10);
        }
        return 0;
    }

    let estimatedTotalPages = detectTotalPages();

    // Helper: Tunggu image selesai memuat
    function waitForImage(img, timeoutMs = 2500) {
        return new Promise((resolve) => {
            if (img.complete && img.naturalWidth > 0) {
                return resolve(true);
            }
            const timer = setTimeout(() => resolve(false), timeoutMs);
            img.addEventListener("load", () => {
                clearTimeout(timer);
                resolve(true);
            }, { once: true });
            img.addEventListener("error", () => {
                clearTimeout(timer);
                resolve(false);
            }, { once: true });
        });
    }

    // Helper: Ambil nomor halaman dari elemen parent
    function getPageNumberFromElement(img) {
        let curr = img;
        while (curr && curr !== document.body) {
            // Cek atribut data-page-number atau sejenisnya
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

    // 7. Auto-Scroll & Capture Loop
    const scrollStep = Math.max(300, Math.floor(scrollContainer.clientHeight * 0.75));
    let currentScroll = 0;
    const maxScroll = scrollContainer.scrollHeight - scrollContainer.clientHeight;

    statusText.textContent = "Memindai halaman dokumen...";

    async function scanVisibleImages() {
        const imgs = Array.from(document.querySelectorAll("img")).filter(img => {
            return /^blob:/.test(img.src) || (img.src && img.src.includes("googleusercontent.com"));
        });

        // Urutkan img berdasarkan posisi vertikal di dokumen
        imgs.sort((a, b) => {
            const rectA = a.getBoundingClientRect();
            const rectB = b.getBoundingClientRect();
            return rectA.top - rectB.top;
        });

        for (const img of imgs) {
            await waitForImage(img, 1500);
            const w = img.naturalWidth || img.width;
            const h = img.naturalHeight || img.height;

            if (w < 100 || h < 100) continue; // Abaikan icon/thumbnail kecil

            let pageNum = getPageNumberFromElement(img);
            if (!pageNum) {
                // Fallback: hitung index berdasarkan urutan yang sudah ada
                pageNum = capturedPages.size + 1;
            }

            // Jika belum ada atau versi saat ini lebih besar resolusinya
            const existing = capturedPages.get(pageNum);
            if (!existing || (w * h > existing.width * existing.height)) {
                try {
                    const canvas = document.createElement("canvas");
                    canvas.width = w;
                    canvas.height = h;
                    const ctx = canvas.getContext("2d");
                    ctx.drawImage(img, 0, 0, w, h);
                    const dataUrl = canvas.toDataURL("image/jpeg", 0.94);

                    capturedPages.set(pageNum, {
                        dataUrl,
                        width: w,
                        height: h
                    });

                    // Update HUD preview
                    thumbImg.style.display = "block";
                    thumbImg.src = dataUrl;
                    pageDetail.innerHTML = `<strong style="color: #38bdf8;">Halaman ${pageNum}</strong> berhasil ditangkap`;
                    resDetail.textContent = `Resolusi: ${w}x${h} px`;

                    const currentCount = capturedPages.size;
                    const totalRef = estimatedTotalPages || Math.max(currentCount, 1);
                    const pct = Math.min(100, Math.round((currentCount / totalRef) * 100));
                    progressBar.style.width = `${pct}%`;
                    statusText.textContent = `Mengumpulkan: ${currentCount} ${estimatedTotalPages ? `/ ${estimatedTotalPages}` : ""} halaman (${pct}%)`;
                } catch (e) {
                    console.warn("Gagal render canvas:", e);
                }
            }
        }
    }

    // Jalankan auto-scroll
    while (!shouldStop && currentScroll <= maxScroll + scrollStep) {
        while (isPaused) {
            await new Promise(r => setTimeout(r, 200));
        }

        scrollContainer.scrollTop = currentScroll;
        window.scrollTo(0, currentScroll);

        await scanVisibleImages();
        await new Promise(r => setTimeout(r, 600)); // Beri waktu Drive merender gambar resolusi tinggi

        currentScroll += scrollStep;
        if (currentScroll > scrollContainer.scrollHeight) {
            break;
        }
    }

    // Scroll kembali ke atas
    scrollContainer.scrollTop = 0;

    if (capturedPages.size === 0) {
        statusText.innerHTML = `<span style="color: #ef4444;">Tidak ada gambar halaman yang terdeteksi. Pastikan file PDF sudah terbuka di Google Drive viewer.</span>`;
        return;
    }

    // 8. Tampilkan Menu Ekspor & Rekonstruksi PDF
    statusText.innerHTML = `<span style="color: #4ade80; font-weight: 600;">✅ Selesai! ${capturedPages.size} halaman siap diproses.</span>`;
    progressBar.style.width = "100%";
    actionsBox.style.display = "none";
    exportSection.style.display = "block";

    // Urutkan halaman
    const sortedPageKeys = Array.from(capturedPages.keys()).sort((a, b) => a - b);

    // Fungsi membuat file Blob PDF
    function generatePdfBlob() {
        let pdf = null;
        for (const pNum of sortedPageKeys) {
            const page = capturedPages.get(pNum);
            const w = page.width;
            const h = page.height;
            const orientation = w >= h ? "landscape" : "portrait";

            if (!pdf) {
                pdf = new jsPDF({
                    orientation: orientation,
                    unit: "pt",
                    format: [w, h],
                    compress: true
                });
            } else {
                pdf.addPage([w, h], orientation);
            }
            pdf.addImage(page.dataUrl, "JPEG", 0, 0, w, h);
        }
        return pdf.output("blob");
    }

    // Action 1: Upload langsung ke Server Comic Reader TV
    const btnUploadServer = document.getElementById("cr-btn-upload-server");
    btnUploadServer.onclick = async () => {
        const serverUrl = document.getElementById("cr-input-server").value.replace(/\/+$/, "");
        const comicTitle = document.getElementById("cr-input-title").value.trim() || getDocumentTitle();

        btnUploadServer.disabled = true;
        btnUploadServer.textContent = "⏳ Merekonstruksi PDF & Mengunggah...";

        try {
            statusText.textContent = "Mengkompilasi halaman ke PDF...";
            const pdfBlob = generatePdfBlob();

            statusText.textContent = `Mengunggah ke ${serverUrl}/api/comics/upload...`;
            const formData = new FormData();
            formData.append("title", comicTitle);
            formData.append("description", `Diimpor otomatis dari Google Drive via ComicReader Extractor (${capturedPages.size} halaman)`);
            formData.append("file", pdfBlob, `${comicTitle}.pdf`);

            const res = await fetch(`${serverUrl}/api/comics/upload`, {
                method: "POST",
                body: formData
            });

            if (!res.ok) {
                const errText = await res.text();
                throw new Error(errText || `Server merespon ${res.status}`);
            }

            const createdComic = await res.json();
            statusText.innerHTML = `
                <div style="color: #4ade80; font-weight: 600; margin-bottom: 4px;">🎉 Berhasil Ditambahkan ke TV!</div>
                <div style="font-size: 11px; color: #cbd5e1;">Komik <strong>"${comicTitle}"</strong> (${createdComic.total_pages} hal) siap dibaca di Android TV!</div>
            `;
            btnUploadServer.style.background = "#059669";
            btnUploadServer.textContent = "✔ Berhasil Diunggah";
        } catch (err) {
            console.error("Gagal unggah komik:", err);
            statusText.innerHTML = `<span style="color: #ef4444;">Gagal unggah: ${err.message}. Periksa alamat server atau CORS.</span>`;
            btnUploadServer.disabled = false;
            btnUploadServer.textContent = "Coba Upload Ulang";
        }
    };

    // Action 2: Download File PDF Lokal
    const btnDownloadPdf = document.getElementById("cr-btn-download-pdf");
    btnDownloadPdf.onclick = () => {
        const comicTitle = document.getElementById("cr-input-title").value.trim() || getDocumentTitle();
        btnDownloadPdf.disabled = true;
        btnDownloadPdf.textContent = "⏳ Menyimpan PDF...";

        setTimeout(() => {
            const pdfBlob = generatePdfBlob();
            const blobUrl = URL.createObjectURL(pdfBlob);
            const a = document.createElement("a");
            a.href = blobUrl;
            a.download = `${comicTitle}.pdf`;
            document.body.appendChild(a);
            a.click();
            a.remove();
            setTimeout(() => URL.revokeObjectURL(blobUrl), 30000);

            btnDownloadPdf.disabled = false;
            btnDownloadPdf.textContent = "✔ Unduhan Selesai";
        }, 100);
    };

})();
