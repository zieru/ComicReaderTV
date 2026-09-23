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
	if l == "" || l == c {
		return false
	}

	var lMaj, lMin, lPatch int
	var cMaj, cMin, cPatch int
	_, _ = fmt.Sscanf(l, "%d.%d.%d", &lMaj, &lMin, &lPatch)
	_, _ = fmt.Sscanf(c, "%d.%d.%d", &cMaj, &cMin, &cPatch)

	if lMaj != cMaj {
		return lMaj > cMaj
	}
	if lMin != cMin {
		return lMin > cMin
	}
	return lPatch > cPatch
}

// FindAPKAsset mencari file .apk di daftar assets release
func (rel *GitHubRelease) FindAPKAsset() *Asset {
	for i := range rel.Assets {
		if strings.HasSuffix(strings.ToLower(rel.Assets[i].Name), ".apk") {
			return &rel.Assets[i]
		}
	}
	return nil
}

// GetWritableUpdateDir mencari folder lokal yang bisa ditulis untuk menyimpan file update
func GetWritableUpdateDir() string {
	candidates := []string{
		"/sdcard/Download",
		"/storage/emulated/0/Download",
		filepath.Join(os.TempDir(), "updates"),
		"/data/data/com.comicreader.tv/cache",
	}
	for _, dir := range candidates {
		if err := os.MkdirAll(dir, 0755); err == nil {
			testFile := filepath.Join(dir, ".test_write")
			if f, err := os.Create(testFile); err == nil {
				_ = f.Close()
				_ = os.Remove(testFile)
				return dir
			}
		}
	}
	return os.TempDir()
}

// DownloadAndInstallAPK mengunduh APK baru dan memanggil installer Android TV
func (u *Updater) DownloadAndInstallAPK(downloadURL, destDir string, onProgress func(percent float32)) (string, error) {
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal mengunduh APK: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gagal mengunduh APK, HTTP %s", resp.Status)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		destDir = os.TempDir()
	}

	apkPath := filepath.Join(destDir, "ComicReaderTV.apk")
	outFile, err := os.Create(apkPath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file update apk: %w", err)
	}
	defer outFile.Close()

	totalBytes := resp.ContentLength
	var downloadedBytes int64
	buf := make([]byte, 32*1024)

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, wErr := outFile.Write(buf[:n]); wErr != nil {
				return "", wErr
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
			return "", err
		}
	}

	if runtime.GOOS == "android" {
		// Memicu intent instalasi di Android TV via am start (Multiple standard actions)
		_ = exec.Command("am", "start", "-a", "android.intent.action.VIEW",
			"-d", "file://"+apkPath,
			"-t", "application/vnd.android.package-archive",
			"-f", "0x10000001",
		).Start()

		_ = exec.Command("am", "start", "-a", "android.intent.action.INSTALL_PACKAGE",
			"-d", "file://"+apkPath,
			"-t", "application/vnd.android.package-archive",
			"-f", "0x10000001",
		).Start()
	}

	return apkPath, nil
}
