package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type GitHubRelease struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Body    string  `json:"body"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

type Updater struct {
	RepoOwner      string
	RepoName       string
	CurrentVersion string
	httpClient     *http.Client
}

func NewUpdater(owner, repo, currentVersion string) *Updater {
	return &Updater{
		RepoOwner:      owner,
		RepoName:       repo,
		CurrentVersion: currentVersion,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CheckUpdate memeriksa apakah ada rilis terbaru di GitHub Releases
func (u *Updater) CheckUpdate() (*GitHubRelease, bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", u.RepoOwner, u.RepoName)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("gagal cek update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("github api merespon dengan status: %s", resp.Status)
	}

	var rel GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, false, err
	}

	hasUpdate := isNewer(rel.TagName, u.CurrentVersion)
	return &rel, hasUpdate, nil
}

func isNewer(latest, current string) bool {
	l := strings.TrimPrefix(latest, "v")
	c := strings.TrimPrefix(current, "v")
	return l != "" && l != c
}

// DownloadAndInstallAPK mengunduh APK baru dan memanggil installer Android TV
func (u *Updater) DownloadAndInstallAPK(downloadURL, destDir string, onProgress func(percent float32)) error {
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return err
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("gagal mengunduh APK: %w", err)
	}
	defer resp.Body.Close()

	apkPath := filepath.Join(destDir, "update.apk")
	outFile, err := os.Create(apkPath)
	if err != nil {
		return fmt.Errorf("gagal membuat file update apk: %w", err)
	}
	defer outFile.Close()

	totalBytes := resp.ContentLength
	var downloadedBytes int64
	buf := make([]byte, 32*1024)

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, wErr := outFile.Write(buf[:n]); wErr != nil {
				return wErr
			}
			downloadedBytes += int64(n)
			if totalBytes > 0 && onProgress != nil {
				onProgress(float32(downloadedBytes) / float32(totalBytes))
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	if runtime.GOOS == "android" {
		// Memicu intent instalasi di Android TV via am start
		cmd := exec.Command("am", "start", "-a", "android.intent.action.VIEW", "-d", "file://"+apkPath, "-t", "application/vnd.android.package-archive")
		return cmd.Start()
	}

	return nil
}
