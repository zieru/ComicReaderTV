package gscraper

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func TestHasBrowser(t *testing.T) {
	has := HasBrowser()
	t.Logf("Chromium/Chrome installed on system: %v (path: %s)", has, BrowserPath())
}

func TestBuildPDFFromJPEGs(t *testing.T) {
	// Create a dummy 100x100 JPEG in memory
	var imgBuf bytes.Buffer
	m := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			m.Set(x, y, color.RGBA{R: 200, G: 50, B: 50, A: 255})
		}
	}
	if err := jpeg.Encode(&imgBuf, m, nil); err != nil {
		t.Fatalf("jpeg encode failed: %v", err)
	}

	pdfBytes, err := BuildPDFFromJPEGs([][]byte{imgBuf.Bytes(), imgBuf.Bytes(), imgBuf.Bytes()})
	if err != nil {
		t.Fatalf("BuildPDFFromJPEGs failed: %v", err)
	}

	tmpPdf := filepath.Join(t.TempDir(), "test_generated.pdf")
	if err := os.WriteFile(tmpPdf, pdfBytes, 0644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	// Validate with pdfcpu
	count, err := api.PageCountFile(tmpPdf)
	if err != nil {
		t.Fatalf("pdfcpu validate failed: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 pages, got %d", count)
	}
	t.Logf("Success! Generated PDF has %d pages and valid xref table!", count)
}
