package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ReleaseInfo struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	PublishedAt time.Time      `json:"published_at"`
	Assets      []ReleaseAsset `json:"assets"`
}

type VersionCheckResult struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	HasUpdate      bool   `json:"has_update"`
	DownloadURL    string `json:"download_url"`
	ReleaseNotes   string `json:"release_notes"`
}

// CheckLatestRelease memeriksa rilis terbaru dari GitHub repo
func CheckLatestRelease(repo string, currentVersion string) (*VersionCheckResult, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ComicReader-Server-Updater")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi GitHub API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API merespon status: %s", resp.Status)
	}

	var rel ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("gagal parse JSON rilis GitHub: %w", err)
	}

	var debURL string
	for _, asset := range rel.Assets {
		if strings.HasSuffix(asset.Name, ".deb") {
			debURL = asset.BrowserDownloadURL
			break
		}
	}

	curr := strings.TrimPrefix(currentVersion, "v")
	latest := strings.TrimPrefix(rel.TagName, "v")

	hasUpdate := latest != "" && latest != curr

	return &VersionCheckResult{
		CurrentVersion: currentVersion,
		LatestVersion:  rel.TagName,
		HasUpdate:      hasUpdate,
		DownloadURL:    debURL,
		ReleaseNotes:   rel.Body,
	}, nil
}

// PerformDebUpdate mengunduh file .deb terbaru dan menjalankannya melalui dpkg -i
func PerformDebUpdate(debURL string) error {
	if debURL == "" {
		return fmt.Errorf("URL download paket .deb kosong")
	}

	tmpFile, err := os.CreateTemp("", "comic-reader-update-*.deb")
	if err != nil {
		return fmt.Errorf("gagal membuat file temporary: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	log.Printf("[Updater] Mengunduh paket pembaruan dari %s...", debURL)
	resp, err := http.Get(debURL)
	if err != nil {
		tmpFile.Close()
		return fmt.Errorf("gagal mengunduh file .deb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tmpFile.Close()
		return fmt.Errorf("gagal download paket, status: %s", resp.Status)
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return fmt.Errorf("gagal menulis file .deb: %w", err)
	}
	tmpFile.Close()

	log.Printf("[Updater] Memasang paket .deb menggunakan dpkg -i %s...", tmpFile.Name())
	cmd := exec.Command("dpkg", "-i", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dpkg -i gagal (%w): %s", err, string(output))
	}

	log.Printf("[Updater] Berhasil memasang paket. Menjadwalkan restart service...")
	go func() {
		time.Sleep(1 * time.Second)
		_ = exec.Command("systemctl", "restart", "comic-reader-server").Run()
	}()

	return nil
}
