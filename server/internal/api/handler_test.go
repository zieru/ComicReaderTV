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

	handler := NewHandler(st, pe, "1.0.3")
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
