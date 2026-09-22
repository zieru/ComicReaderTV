package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"comic_reader/pkg/model"
	"comic_reader/server/internal/auth"
	"comic_reader/server/internal/manga"
	"comic_reader/server/internal/pdfengine"
	"comic_reader/server/internal/store"
)

// createDummyPDFBytes menghasilkan minimal valid PDF
func createDummyPDFBytes() []byte {
	return []byte("%PDF-1.4\n1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj\n2 0 obj<</Type/Pages/Count 1/Kids[3 0 R]>>endobj\n3 0 obj<</Type/Page/MediaBox[0 0 300 300]/Parent 2 0 R/Resources<<>>>>endobj\nxref\n0 4\n0000000000 65535 f \n0000000009 00000 n \n0000000052 00000 n \n0000000108 00000 n \ntrailer<</Size 4/Root 1 0 R>>\nstartxref\n185\n%%EOF\n")
}

func setupTestServer(t *testing.T) (*Handler, http.Handler, string) {
	tempDir, err := os.MkdirTemp("", "cr_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	st, err := store.New(filepath.Join(tempDir, "store"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	pe, err := pdfengine.New(filepath.Join(tempDir, "pdf_cache"))
	if err != nil {
		t.Fatalf("Failed to create pdfengine: %v", err)
	}

	handler := NewHandler(st, pe, nil, "1.0.5")
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return handler, mux, tempDir
}

func TestUploadComic(t *testing.T) {
	_, mux, tempDir := setupTestServer(t)
	defer os.RemoveAll(tempDir)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Tambahkan fields
	_ = writer.WriteField("title", "Test Manga Upload")
	_ = writer.WriteField("description", "Komik hasil test upload")

	part, err := writer.CreateFormFile("file", "test_manga.pdf")
	if err != nil {
		t.Fatalf("CreateFormFile error: %v", err)
	}

	dummyPDF := createDummyPDFBytes()
	if _, err := io.Copy(part, bytes.NewReader(dummyPDF)); err != nil {
		t.Fatalf("Copy PDF error: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/comics/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 Created, got %d: %s", rec.Code, rec.Body.String())
	}

	var created model.Comic
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if created.Title != "Test Manga Upload" {
		t.Errorf("Expected title 'Test Manga Upload', got '%s'", created.Title)
	}

	if created.SourceType != "uploaded_pdf" {
		t.Errorf("Expected source_type 'uploaded_pdf', got '%s'", created.SourceType)
	}

	// Cek endpoint download PDF
	reqPdf := httptest.NewRequest(http.MethodGet, "/api/comics/"+created.ID+"/pdf", nil)
	recPdf := httptest.NewRecorder()
	mux.ServeHTTP(recPdf, reqPdf)

	if recPdf.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK for /pdf, got %d: %s", recPdf.Code, recPdf.Body.String())
	}

	if recPdf.Header().Get("Content-Type") != "application/pdf" {
		t.Errorf("Expected application/pdf content type, got: %s", recPdf.Header().Get("Content-Type"))
	}
}

func TestUploadOptionsCORS(t *testing.T) {
	_, mux, tempDir := setupTestServer(t)
	defer os.RemoveAll(tempDir)

	req := httptest.NewRequest(http.MethodOptions, "/api/comics/upload", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected CORS origin *, got: %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestMangaSearch(t *testing.T) {
	_, mux, tempDir := setupTestServer(t)
	defer os.RemoveAll(tempDir)

	req := httptest.NewRequest(http.MethodGet, "/api/manga/search?q=One+Piece", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rec.Code)
	}

	var results []manga.Item
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatalf("Failed to decode results: %v", err)
	}

	if len(results) == 0 {
		t.Logf("Warning: no results returned (might be network/rate limit)")
	} else {
		t.Logf("Found %d manga results for 'One Piece', first: %s", len(results), results[0].Title)
	}
}

func TestAuthProtection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cr_auth_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	st, err := store.New(filepath.Join(tempDir, "store"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	pe, err := pdfengine.New(filepath.Join(tempDir, "pdf_cache"))
	if err != nil {
		t.Fatalf("Failed to create pdfengine: %v", err)
	}

	authMgr := auth.NewManager("mock_token", filepath.Join(tempDir, "auth"), 123456)
	handler := NewHandler(st, pe, authMgr, "1.0.5")
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// 1. Check status
	reqStatus := httptest.NewRequest(http.MethodGet, "/api/auth/status", nil)
	recStatus := httptest.NewRecorder()
	mux.ServeHTTP(recStatus, reqStatus)
	if recStatus.Code != http.StatusOK {
		t.Fatalf("Expected 200 for auth status, got %d", recStatus.Code)
	}

	// 2. Unauthenticated POST /api/comics should be 401 Unauthorized
	comicReq := model.CreateComicRequest{Title: "Unauth Comic", SourceURL: "https://example.com/test.pdf"}
	body, _ := json.Marshal(comicReq)
	reqPost := httptest.NewRequest(http.MethodPost, "/api/comics", bytes.NewReader(body))
	recPost := httptest.NewRecorder()
	mux.ServeHTTP(recPost, reqPost)
	if recPost.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized, got %d", recPost.Code)
	}

	// 3. Unauthenticated GET /api/comics should still be 200 OK (accessible for TV app)
	reqGet := httptest.NewRequest(http.MethodGet, "/api/comics", nil)
	recGet := httptest.NewRecorder()
	mux.ServeHTTP(recGet, reqGet)
	if recGet.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for GET /api/comics, got %d", recGet.Code)
	}

	// 4. Authenticated request using session cookie
	sessionToken := authMgr.CreateSessionForTest()
	reqPostAuth := httptest.NewRequest(http.MethodPost, "/api/comics", bytes.NewReader(body))
	reqPostAuth.AddCookie(&http.Cookie{
		Name:  "comic_session",
		Value: sessionToken,
	})
	recPostAuth := httptest.NewRecorder()
	mux.ServeHTTP(recPostAuth, reqPostAuth)
	if recPostAuth.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created with valid session cookie, got %d: %s", recPostAuth.Code, recPostAuth.Body.String())
	}
}
