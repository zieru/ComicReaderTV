package gdrive

import (
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	patterns = []*regexp.Regexp{
		regexp.MustCompile(`/file/d/([a-zA-Z0-9_-]+)`),
		regexp.MustCompile(`[?&]id=([a-zA-Z0-9_-]+)`),
	}
	errorCaptionRegex    = regexp.MustCompile(`class="uc-error-caption">([^<]+)`)
	errorSubcaptionRegex = regexp.MustCompile(`class="uc-error-subcaption">([^<]+)`)
	confirmRegex         = regexp.MustCompile(`[?&]confirm=([a-zA-Z0-9_-]+)`)
	confirmInputRegex    = regexp.MustCompile(`name="confirm"[^>]+value="([a-zA-Z0-9_-]+)"`)
	uuidInputRegex       = regexp.MustCompile(`name="uuid"[^>]+value="([a-zA-Z0-9_-]+)"`)
	actionRegex          = regexp.MustCompile(`id="download-form"[^>]+action="([^"]+)"`)
	downloadLinkRegex    = regexp.MustCompile(`id="uc-download-link"[^>]+href="([^"]+)"`)
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
	jar, _ := cookiejar.New(nil)
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		DisableKeepAlives:     false,
	}

	return &Client{
		httpClient: &http.Client{
			Jar:       jar,
			Transport: transport,
			// Timeout tidak diset pada http.Client agar streaming file besar tidak terputus di tengah jalan
		},
	}
}

// DownloadFile mengunduh stream file dari link publik Google Drive
// Menangani bypass Google Drive virus scan warning untuk file berukuran besar
// dan menangkap error izin / privasi Google Drive secara jelas
func (c *Client) DownloadFile(fileID string) (io.ReadCloser, error) {
	downloadURL := fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID)
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan request ke Google Drive: %w", err)
	}

	// Cek jika Google Drive mengembalikan HTML (peringatan konfirmasi virus/file besar atau pesan error izin)
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, fmt.Errorf("gagal membaca respon HTML Google Drive: %w", readErr)
		}

		bodyStr := string(bodyBytes)

		// 1. Cek apakah ada pesan error dari Google Drive
		if match := errorCaptionRegex.FindStringSubmatch(bodyStr); len(match) > 1 {
			errMsg := html.UnescapeString(strings.TrimSpace(match[1]))
			subMsg := ""
			if subMatch := errorSubcaptionRegex.FindStringSubmatch(bodyStr); len(subMatch) > 1 {
				subMsg = html.UnescapeString(strings.TrimSpace(subMatch[1]))
			}

			if strings.Contains(strings.ToLower(errMsg), "permission to download") ||
				strings.Contains(strings.ToLower(subMsg), "only the owner and editors can download") {
				return nil, fmt.Errorf("izin download file Google Drive ini dinonaktifkan oleh pemilik: '%s'. Buka file di Google Drive > Bagikan (Share) > Ikon Pengaturan/Gerigi > Centang 'Lihat, komentator, dan editor dapat melihat opsi untuk mendownload'", errMsg)
			}
			return nil, fmt.Errorf("Google Drive menolak download: %s (%s)", errMsg, subMsg)
		}

		// 2. Cek apakah ada link / form konfirmasi untuk file berukuran besar
		var confirmURL string

		if linkMatch := downloadLinkRegex.FindStringSubmatch(bodyStr); len(linkMatch) > 1 {
			rawLink := html.UnescapeString(linkMatch[1])
			if strings.HasPrefix(rawLink, "http") {
				confirmURL = rawLink
			} else {
				confirmURL = "https://drive.google.com" + rawLink
			}
		} else if formMatch := actionRegex.FindStringSubmatch(bodyStr); len(formMatch) > 1 {
			actionURL := html.UnescapeString(formMatch[1])
			confirmVal := "t"
			if cMatch := confirmInputRegex.FindStringSubmatch(bodyStr); len(cMatch) > 1 {
				confirmVal = cMatch[1]
			}
			uuidVal := ""
			if uMatch := uuidInputRegex.FindStringSubmatch(bodyStr); len(uMatch) > 1 {
				uuidVal = uMatch[1]
			}
			confirmURL = fmt.Sprintf("%s?id=%s&export=download&confirm=%s&uuid=%s", actionURL, fileID, confirmVal, uuidVal)
		} else if cMatch := confirmRegex.FindStringSubmatch(bodyStr); len(cMatch) > 1 {
			confirmCode := cMatch[1]
			confirmURL = fmt.Sprintf("https://drive.usercontent.google.com/download?id=%s&export=download&confirm=%s", fileID, confirmCode)
		}

		if confirmURL != "" {
			req2, err := http.NewRequest(http.MethodGet, confirmURL, nil)
			if err != nil {
				return nil, err
			}
			req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
			req2.Header.Set("Accept", "*/*")

			resp2, err := c.httpClient.Do(req2)
			if err != nil {
				return nil, fmt.Errorf("gagal request konfirmasi download: %w", err)
			}

			if strings.Contains(resp2.Header.Get("Content-Type"), "text/html") {
				resp2Body, _ := io.ReadAll(resp2.Body)
				resp2.Body.Close()
				preview := string(resp2Body)
				if len(preview) > 200 {
					preview = preview[:200]
				}
				return nil, fmt.Errorf("Google Drive tetap mengembalikan halaman HTML setelah konfirmasi: %s", preview)
			}

			if resp2.StatusCode != http.StatusOK {
				resp2.Body.Close()
				return nil, fmt.Errorf("Google Drive konfirmasi merespon status: %s", resp2.Status)
			}

			return resp2.Body, nil
		}

		// 3. Fallback jika HTML tidak ada tombol konfirmasi dan tidak ada error caption
		title := "Halaman Google Drive"
		titleRegex := regexp.MustCompile(`<title>([^<]+)</title>`)
		if tMatch := titleRegex.FindStringSubmatch(bodyStr); len(tMatch) > 1 {
			title = html.UnescapeString(tMatch[1])
		}
		return nil, fmt.Errorf("Google Drive tidak menyediakan stream file (merespon: '%s'). Pastikan link file bersifat publik dan dapat diunduh", title)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Google Drive merespon status: %s", resp.Status)
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
