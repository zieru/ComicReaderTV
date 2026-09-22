package gdrive

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	// Pola link Google Drive yang sering ditemui:
	// https://drive.google.com/file/d/1aBcDeFgHiJkLmNoPqRsTuVwXyZ/view?usp=sharing
	// https://drive.google.com/file/d/1aBcDeFgHiJkLmNoPqRsTuVwXyZ/preview
	// https://drive.google.com/open?id=1aBcDeFgHiJkLmNoPqRsTuVwXyZ
	// https://drive.google.com/uc?id=1aBcDeFgHiJkLmNoPqRsTuVwXyZ
	patterns = []*regexp.Regexp{
		regexp.MustCompile(`/file/d/([a-zA-Z0-9_-]+)`),
		regexp.MustCompile(`[?&]id=([a-zA-Z0-9_-]+)`),
	}
)

// ExtractFileID mengambil Google Drive File ID dari berbagai macam format URL
func ExtractFileID(rawURL string) (string, error) {
	for _, p := range patterns {
		matches := p.FindStringSubmatch(rawURL)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}
	return "", fmt.Errorf("tidak dapat menemukan Google Drive File ID dari url: %s", rawURL)
}

// Client menangani request download file dari Google Drive
type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 3 * time.Minute,
		},
	}
}

// DownloadFile mengunduh stream file dari link publik Google Drive
// Menangani bypass Google Drive virus scan warning untuk file berukuran besar
func (c *Client) DownloadFile(fileID string) (io.ReadCloser, error) {
	downloadURL := fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID)
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan request ke Google Drive: %w", err)
	}

	// Cek jika Google Drive meminta konfirmasi download untuk file besar ("virus scan warning")
	if resp.StatusCode == http.StatusOK && strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("gagal membaca respon warning: %w", readErr)
		}

		bodyStr := string(bodyBytes)
		confirmRegex := regexp.MustCompile(`confirm=([a-zA-Z0-9_-]+)`)
		match := confirmRegex.FindStringSubmatch(bodyStr)
		if len(match) > 1 {
			confirmCode := match[1]
			// Request kedua dengan query confirm
			secondURL := fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s&confirm=%s", fileID, confirmCode)
			req2, err := http.NewRequest(http.MethodGet, secondURL, nil)
			if err != nil {
				return nil, err
			}
			req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

			// Copy cookies jika ada
			for _, cookie := range resp.Cookies() {
				req2.AddCookie(cookie)
			}

			resp2, err := c.httpClient.Do(req2)
			if err != nil {
				return nil, fmt.Errorf("gagal request konfirmasi file besar: %w", err)
			}
			return resp2.Body, nil
		}
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("google drive merespon dengan status: %s", resp.Status)
	}

	return resp.Body, nil
}

// IsGDriveURL memeriksa apakah URL berasal dari domain Google Drive
func IsGDriveURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return strings.Contains(u.Hostname(), "drive.google.com")
}
