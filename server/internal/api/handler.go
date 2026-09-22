package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"comic_reader/server/internal/gdrive"
	"comic_reader/pkg/model"
	"comic_reader/server/internal/pdfengine"
	"comic_reader/server/internal/store"
)

type Handler struct {
	store     *store.Store
	pdfEngine *pdfengine.Engine
}

func NewHandler(st *store.Store, pdfEng *pdfengine.Engine) *Handler {
	return &Handler{
		store:     st,
		pdfEngine: pdfEng,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/comics", h.handleComics)
	mux.HandleFunc("/api/comics/", h.handleComicItem)
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

		// Download / cache PDF & hitung halaman di background atau sinkron
		totalPages := 0
		go func(cid, surl string) {
			localPath, err := h.pdfEngine.EnsureLocalPDF(cid, surl)
			if err != nil {
				log.Printf("Gagal prefetch PDF %s: %v", cid, err)
				return
			}
			count, err := h.pdfEngine.GetPageCount(localPath)
			if err != nil {
				log.Printf("Gagal hitung halaman %s: %v", cid, err)
				return
			}
			if comic, ok := h.store.Get(cid); ok {
				comic.TotalPages = count
				h.store.Upsert(comic)
			}
		}(comicID, req.SourceURL)

		comic := model.Comic{
			ID:          comicID,
			Title:       req.Title,
			Description: req.Description,
			CoverURL:    req.CoverURL,
			SourceType:  sourceType,
			SourceURL:   req.SourceURL,
			TotalPages:  totalPages,
		}

		if err := h.store.Upsert(comic); err != nil {
			http.Error(w, fmt.Sprintf("Gagal menyimpan komik: %v", err), http.StatusInternalServerError)
			return
		}

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
			h.store.Delete(comicID)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	http.NotFound(w, r)
}
