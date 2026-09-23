package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"comic_reader/pkg/model"
	"comic_reader/server/internal/auth"
	"comic_reader/server/internal/gdrive"
	"comic_reader/server/internal/gscraper"
	"comic_reader/server/internal/manga"
	"comic_reader/server/internal/pdfengine"
	"comic_reader/server/internal/store"
	"comic_reader/server/internal/updater"
)

type Handler struct {
	store     *store.Store
	pdfEngine *pdfengine.Engine
	authMgr   *auth.Manager
	version   string
}

func NewHandler(st *store.Store, pdfEng *pdfengine.Engine, authMgr *auth.Manager, version string) *Handler {
	return &Handler{
		store:     st,
		pdfEngine: pdfEng,
		authMgr:   authMgr,
		version:   version,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/comics", h.handleComics)
	mux.HandleFunc("/api/comics/upload", h.handleUploadComic)
	mux.HandleFunc("/api/comics/", h.handleComicItem)
	mux.HandleFunc("/api/manga/search", h.handleMangaSearch)
	mux.HandleFunc("/api/auth/status", h.handleAuthStatus)
	mux.HandleFunc("/api/auth/request-otp", h.handleRequestOTP)
	mux.HandleFunc("/api/auth/verify-otp", h.handleVerifyOTP)
	mux.HandleFunc("/api/auth/logout", h.handleLogout)
	mux.HandleFunc("/api/admin/version", h.handleVersion)
	mux.HandleFunc("/api/admin/update", h.handleUpdate)
}

func (h *Handler) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	res, err := updater.CheckLatestRelease("zieru/ComicReaderTV", h.version)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal cek rilis GitHub: %v", err), http.StatusInternalServerError)
		return
	}

	type versionResponse struct {
		updater.VersionCheckResult
		HasBrowser  bool   `json:"has_browser"`
		BrowserPath string `json:"browser_path"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versionResponse{
		VersionCheckResult: *res,
		HasBrowser:         gscraper.HasBrowser(),
		BrowserPath:        gscraper.BrowserPath(),
	})
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !h.checkAuth(r) {
		http.Error(w, "Unauthorized: Silakan login terlebih dahulu via Telegram OTP", http.StatusUnauthorized)
		return
	}

	res, err := updater.CheckLatestRelease("zieru/ComicReaderTV", h.version)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal cek rilis GitHub: %v", err), http.StatusInternalServerError)
		return
	}
	if !res.HasUpdate {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Server sudah menggunakan versi terbaru",
		})
		return
	}

	if err := updater.PerformDebUpdate(res.DownloadURL); err != nil {
		http.Error(w, fmt.Sprintf("Gagal update paket: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"message":        fmt.Sprintf("Pembaruan ke %s berhasil dipasang! Service sedang restart...", res.LatestVersion),
		"latest_version": res.LatestVersion,
	})
}

func (h *Handler) handleComics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	switch r.Method {
	case http.MethodGet:
		comics := h.store.GetAll()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comics)

	case http.MethodPost:
		if !h.checkAuth(r) {
			http.Error(w, "Unauthorized: Silakan login terlebih dahulu via Telegram OTP", http.StatusUnauthorized)
			return
		}

		var req model.CreateComicRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if req.Title == "" || req.SourceURL == "" {
			http.Error(w, "Title dan SourceURL wajib diisi", http.StatusBadRequest)
			return
		}

		idBytes := make([]byte, 6)
		rand.Read(idBytes)
		comicID := hex.EncodeToString(idBytes)

		sourceType := "direct_pdf"
		if gdrive.IsGDriveURL(req.SourceURL) {
			sourceType = "gdrive_pdf"
		}

		// Download / cache PDF & hitung halaman di background
		comic := model.Comic{
			ID:          comicID,
			Title:       req.Title,
			Description: req.Description,
			CoverURL:    req.CoverURL,
			SourceType:  sourceType,
			SourceURL:   req.SourceURL,
			TotalPages:  0,
			Status:      "processing",
		}

		if err := h.store.Upsert(comic); err != nil {
			http.Error(w, fmt.Sprintf("Gagal menyimpan komik: %v", err), http.StatusInternalServerError)
			return
		}

		go h.prefetchComic(comicID, req.SourceURL)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(comic)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleComicItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	// Path parsing: /api/comics/{id} atau /api/comics/{id}/page/{num}
	path := strings.TrimPrefix(r.URL.Path, "/api/comics/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	comicID := parts[0]
	comic, ok := h.store.Get(comicID)
	if !ok {
		http.Error(w, "Komik tidak ditemukan", http.StatusNotFound)
		return
	}

	if len(parts) == 1 {
		// GET /api/comics/{id}
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(comic)
			return
		}
		if r.Method == http.MethodDelete {
			if !h.checkAuth(r) {
				http.Error(w, "Unauthorized: Silakan login terlebih dahulu via Telegram OTP", http.StatusUnauthorized)
				return
			}
			h.store.Delete(comicID)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Sub-resource: /api/comics/{id}/pdf
	if len(parts) >= 2 && parts[1] == "pdf" {
		localPDF, err := h.pdfEngine.EnsureLocalPDF(comic.ID, comic.SourceURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Gagal memuat PDF: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s.pdf\"", comic.Title))
		http.ServeFile(w, r, localPDF)
		return
	}

	// Sub-resource: /api/comics/{id}/page/{num}
	if len(parts) >= 3 && parts[1] == "page" {
		pageNum, err := strconv.Atoi(parts[2])
		if err != nil || pageNum < 1 {
			http.Error(w, "Nomor halaman tidak valid", http.StatusBadRequest)
			return
		}

		localPDF, err := h.pdfEngine.EnsureLocalPDF(comic.ID, comic.SourceURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Gagal memuat PDF: %v", err), http.StatusInternalServerError)
			return
		}

		// Update total pages jika belum terisi
		if comic.TotalPages == 0 {
			count, _ := h.pdfEngine.GetPageCount(localPDF)
			if count > 0 {
				comic.TotalPages = count
				h.store.Upsert(comic)
			}
		}

		imagePath, err := h.pdfEngine.ExtractPageImage(localPDF, pageNum)
		if err != nil {
			http.Error(w, fmt.Sprintf("Gagal mengekstrak halaman %d: %v", pageNum, err), http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, imagePath)
		return
	}

	// Sub-resource: POST /api/comics/{id}/retry
	if len(parts) >= 2 && parts[1] == "retry" {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !h.checkAuth(r) {
			http.Error(w, "Unauthorized: Silakan login terlebih dahulu via Telegram OTP", http.StatusUnauthorized)
			return
		}
		comic.Status = "processing"
		comic.ErrorMsg = ""
		h.store.Upsert(comic)
		go h.prefetchComic(comic.ID, comic.SourceURL)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comic)
		return
	}

	http.NotFound(w, r)
}

func (h *Handler) prefetchComic(cid, surl string) {
	localPath, err := h.pdfEngine.EnsureLocalPDF(cid, surl)
	if err != nil {
		log.Printf("[prefetch] Gagal prefetch PDF %s: %v", cid, err)
		if c, ok := h.store.Get(cid); ok {
			c.Status = "error"
			c.ErrorMsg = err.Error()
			h.store.Upsert(c)
		}
		return
	}
	count, err := h.pdfEngine.GetPageCount(localPath)
	if err != nil {
		log.Printf("[prefetch] Gagal hitung halaman %s: %v", cid, err)
		if c, ok := h.store.Get(cid); ok {
			c.Status = "error"
			c.ErrorMsg = fmt.Sprintf("Gagal membaca halaman PDF: %v", err)
			h.store.Upsert(c)
		}
		return
	}
	if c, ok := h.store.Get(cid); ok {
		c.TotalPages = count
		c.Status = "ready"
		c.ErrorMsg = ""
		h.store.Upsert(c)
		log.Printf("[prefetch] Berhasil memproses komik %s (%d halaman)", cid, count)
	}
}

func (h *Handler) handleUploadComic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !h.checkAuth(r) {
		http.Error(w, "Unauthorized: Silakan login terlebih dahulu via Telegram OTP", http.StatusUnauthorized)
		return
	}

	// Batasi ukuran upload (maks 500MB)
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Gagal memproses form upload: %v", err), http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	description := strings.TrimSpace(r.FormValue("description"))
	coverURL := strings.TrimSpace(r.FormValue("cover_url"))

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File PDF wajib disertakan pada field 'file'", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if title == "" {
		title = header.Filename
		if ext := filepath.Ext(title); ext != "" {
			title = strings.TrimSuffix(title, ext)
		}
	}

	idBytes := make([]byte, 6)
	rand.Read(idBytes)
	comicID := hex.EncodeToString(idBytes)

	localPath, err := h.pdfEngine.SaveUploadedPDF(comicID, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal menyimpan file PDF: %v", err), http.StatusBadRequest)
		return
	}

	totalPages, err := h.pdfEngine.GetPageCount(localPath)
	if err != nil {
		log.Printf("Peringatan: gagal hitung halaman komik %s: %v", comicID, err)
		totalPages = 0
	}

	status := "ready"
	errMsg := ""
	if totalPages == 0 {
		status = "error"
		errMsg = "Gagal menghitung halaman file PDF yang diunggah"
	}

	comic := model.Comic{
		ID:          comicID,
		Title:       title,
		Description: description,
		CoverURL:    coverURL,
		SourceType:  "uploaded_pdf",
		SourceURL:   fmt.Sprintf("local:%s", comicID),
		TotalPages:  totalPages,
		Status:      status,
		ErrorMsg:    errMsg,
	}

	if err := h.store.Upsert(comic); err != nil {
		http.Error(w, fmt.Sprintf("Gagal menyimpan data komik: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comic)
}

func (h *Handler) handleMangaSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]manga.Item{})
		return
	}

	results, err := manga.Search(query)
	if err != nil {
		log.Printf("Gagal mencari manga '%s': %v", query, err)
		http.Error(w, fmt.Sprintf("Gagal mencari manga: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *Handler) checkAuth(r *http.Request) bool {
	if h.authMgr == nil {
		return true
	}
	// Periksa cookie "comic_session"
	if cookie, err := r.Cookie("comic_session"); err == nil && cookie.Value != "" {
		if h.authMgr.ValidateSession(cookie.Value) {
			return true
		}
	}
	// Periksa header Authorization: Bearer <token>
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if h.authMgr.ValidateSession(token) {
			return true
		}
	}
	return false
}

func (h *Handler) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	hasAdmin, botUser, adminUser := false, "", ""
	if h.authMgr != nil {
		hasAdmin, botUser, adminUser = h.authMgr.GetStatus()
	}
	isAuth := h.checkAuth(r)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated":    isAuth,
		"is_authenticated": isAuth,
		"has_admin":        hasAdmin,
		"bot_username":     botUser,
		"admin_username":   adminUser,
	})
}

func (h *Handler) handleRequestOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.authMgr == nil {
		http.Error(w, "Auth manager tidak aktif", http.StatusInternalServerError)
		return
	}

	msg, err := h.authMgr.RequestOTP()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": msg,
	})
}

func (h *Handler) handleVerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	sessionToken, err := h.authMgr.VerifyOTP(req.Code)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Set HTTP Cookie untuk web browser (berlaku 1 hari / 24 jam)
	http.SetCookie(w, &http.Cookie{
		Name:     "comic_session",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   24 * 3600, // 1 hari (24 jam)
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"token":   sessionToken,
		"message": "Login berhasil!",
	})
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	if cookie, err := r.Cookie("comic_session"); err == nil && cookie.Value != "" {
		if h.authMgr != nil {
			h.authMgr.RevokeSession(cookie.Value)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "comic_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Berhasil logout",
	})
}

