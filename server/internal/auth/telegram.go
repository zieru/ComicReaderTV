package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	botToken      string
	dataDir       string
	botUsername   string
	adminChatID   int64
	adminUsername string

	mu             sync.RWMutex
	currentOTP     string
	otpExpiry      time.Time
	lastRequestOTP time.Time

	sessions map[string]time.Time // token -> expiry
	stopPoll chan struct{}
}

type tgUser struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type tgChat struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type tgMessage struct {
	MessageID int64  `json:"message_id"`
	From      tgUser `json:"from"`
	Chat      tgChat `json:"chat"`
	Text      string `json:"text"`
}

type tgUpdate struct {
	UpdateID int64     `json:"update_id"`
	Message  *tgMessage `json:"message"`
}

type tgResponse struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

func NewManager(botToken string, dataDir string, initialChatID int64) *Manager {
	m := &Manager{
		botToken:    strings.TrimSpace(botToken),
		dataDir:     dataDir,
		adminChatID: initialChatID,
		sessions:    make(map[string]time.Time),
		stopPoll:    make(chan struct{}),
	}

	// Coba muat chat ID yang tersimpan dari file jika initialChatID belum ditentukan
	if m.adminChatID == 0 {
		savedID := m.loadSavedChatID()
		if savedID != 0 {
			m.adminChatID = savedID
			log.Printf("[auth] Memuat Admin Telegram Chat ID tersimpan: %d", savedID)
		}
	}

	// Dapatkan info bot
	m.initBotInfo()

	// Jalankan background listener untuk menangkap /start dari admin
	if m.botToken != "" {
		go m.pollTelegramUpdates()
	}

	return m
}

func (m *Manager) Close() {
	select {
	case <-m.stopPoll:
	default:
		close(m.stopPoll)
	}
}

func (m *Manager) initBotInfo() {
	if m.botToken == "" {
		return
	}
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", m.botToken)
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("[auth] Gagal getMe bot Telegram: %v", err)
		return
	}
	defer resp.Body.Close()

	var tgRes tgResponse
	if err := json.NewDecoder(resp.Body).Decode(&tgRes); err == nil && tgRes.OK {
		var user tgUser
		if err := json.Unmarshal(tgRes.Result, &user); err == nil {
			m.botUsername = user.Username
			log.Printf("[auth] Bot Telegram terhubung: @%s (ID: %d)", user.Username, user.ID)
		}
	}
}

func (m *Manager) loadSavedChatID() int64 {
	path := filepath.Join(m.dataDir, "admin_chat_id.txt")
	content, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	id, err := strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
	if err != nil {
		return 0
	}
	return id
}

func (m *Manager) saveChatID(id int64) {
	_ = os.MkdirAll(m.dataDir, 0755)
	path := filepath.Join(m.dataDir, "admin_chat_id.txt")
	_ = os.WriteFile(path, []byte(strconv.FormatInt(id, 10)), 0644)
}

// SetAdminChatID menetapkan admin chat ID secara eksplisit
func (m *Manager) SetAdminChatID(id int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.adminChatID = id
	m.saveChatID(id)
}

// GetStatus mengembalikan status bot dan admin
func (m *Manager) GetStatus() (hasAdmin bool, botUsername string, adminUsername string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.adminChatID != 0, m.botUsername, m.adminUsername
}

// RequestOTP menghasilkan kode OTP 6 digit dan mengirimkannya ke Telegram admin
func (m *Manager) RequestOTP() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.adminChatID == 0 {
		return "", errors.New("admin belum terhubung ke bot. Silakan buka bot di Telegram dan kirim pesan /start terlebih dahulu")
	}

	// Rate limiting: minimal jeda 15 detik antar permintaan
	if time.Since(m.lastRequestOTP) < 15*time.Second {
		remain := int((15 * time.Second - time.Since(m.lastRequestOTP)).Seconds())
		return "", fmt.Errorf("harap tunggu %d detik sebelum meminta OTP baru", remain)
	}

	// Generate 6 digit angka acak (100000 - 999999)
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", fmt.Errorf("gagal menghasilkan angka acak: %w", err)
	}
	otp := strconv.Itoa(int(n.Int64() + 100000))

	m.currentOTP = otp
	m.otpExpiry = time.Now().Add(5 * time.Minute)
	m.lastRequestOTP = time.Now()

	// Kirim pesan ke Telegram
	text := fmt.Sprintf(
		"🔐 <b>KODE OTP COMIC READER TV</b>\n\n"+
			"Kode Login Anda: <code>%s</code>\n\n"+
			"⏰ Berlaku selama <b>5 menit</b>.\n"+
			"⚠️ <i>Jangan berikan kode ini kepada siapa pun!</i>",
		otp,
	)

	go m.sendTelegramMessage(m.adminChatID, text)

	log.Printf("[auth] OTP berhasil dibuat dan dikirim ke Telegram Chat ID %d", m.adminChatID)
	return "OTP telah dikirim ke Telegram Anda", nil
}

// VerifyOTP memverifikasi kode OTP yang dimasukkan user
func (m *Manager) VerifyOTP(code string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	code = strings.TrimSpace(code)
	if code == "" {
		return "", errors.New("kode OTP wajib diisi")
	}

	if m.currentOTP == "" || time.Now().After(m.otpExpiry) {
		return "", errors.New("kode OTP sudah kadaluarsa atau belum diminta. Silakan minta kode baru")
	}

	if m.currentOTP != code {
		return "", errors.New("kode OTP salah. Periksa kembali pesan di Telegram Anda")
	}

	// Reset OTP agar tidak dapat digunakan ulang (One-Time)
	m.currentOTP = ""
	m.otpExpiry = time.Time{}

	// Generate session token 32 byte hex (64 karakter)
	tokenBytes := make([]byte, 32)
	_, _ = rand.Read(tokenBytes)
	sessionToken := hex.EncodeToString(tokenBytes)

	// Sesi aktif selama 7 hari
	m.sessions[sessionToken] = time.Now().Add(7 * 24 * time.Hour)

	log.Printf("[auth] Login berhasil via OTP! Sesi dibuat (berlaku 7 hari)")
	return sessionToken, nil
}

// ValidateSession memeriksa apakah session token valid dan belum kadaluarsa
func (m *Manager) ValidateSession(token string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if token == "" {
		return false
	}
	expiry, exists := m.sessions[token]
	if !exists {
		return false
	}
	if time.Now().After(expiry) {
		return false
	}
	return true
}

// RevokeSession menghapus token sesi login (Logout)
func (m *Manager) RevokeSession(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, token)
}

// CreateSessionForTest membuat sesi aktif langsung untuk keperluan automated test
func (m *Manager) CreateSessionForTest() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	tokenBytes := make([]byte, 32)
	_, _ = rand.Read(tokenBytes)
	sessionToken := hex.EncodeToString(tokenBytes)
	m.sessions[sessionToken] = time.Now().Add(7 * 24 * time.Hour)
	return sessionToken
}

func (m *Manager) sendTelegramMessage(chatID int64, htmlText string) {
	if m.botToken == "" || chatID == 0 {
		return
	}
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", m.botToken)
	form := url.Values{}
	form.Set("chat_id", strconv.FormatInt(chatID, 10))
	form.Set("text", htmlText)
	form.Set("parse_mode", "HTML")

	resp, err := http.PostForm(apiURL, form)
	if err != nil {
		log.Printf("[auth] Gagal kirim pesan Telegram: %v", err)
		return
	}
	defer resp.Body.Close()
}

// pollTelegramUpdates mendengarkan pesan dari user untuk auto-pairing admin
func (m *Manager) pollTelegramUpdates() {
	var offset int64 = 0
	client := &http.Client{Timeout: 35 * time.Second}

	for {
		select {
		case <-m.stopPoll:
			return
		default:
		}

		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=25", m.botToken, offset)
		resp, err := client.Get(apiURL)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		var tgRes struct {
			OK     bool       `json:"ok"`
			Result []tgUpdate `json:"result"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&tgRes); err != nil {
			resp.Body.Close()
			time.Sleep(3 * time.Second)
			continue
		}
		resp.Body.Close()

		if tgRes.OK {
			for _, u := range tgRes.Result {
				if u.UpdateID >= offset {
					offset = u.UpdateID + 1
				}
				if u.Message == nil {
					continue
				}

				chatID := u.Message.Chat.ID
				user := u.Message.From
				text := strings.TrimSpace(u.Message.Text)

				// Jika pesan adalah /start atau /login atau /admin
				if strings.HasPrefix(text, "/start") || strings.HasPrefix(text, "/login") || strings.HasPrefix(text, "/admin") || strings.HasPrefix(text, "/otp") {
					m.mu.Lock()
					isNew := (m.adminChatID != chatID)
					m.adminChatID = chatID
					m.adminUsername = user.Username
					m.saveChatID(chatID)
					m.mu.Unlock()

					log.Printf("[auth] Akun Telegram admin terhubung: @%s (Chat ID: %d)", user.Username, chatID)

					welcomeText := fmt.Sprintf(
						"👋 <b>Halo %s!</b>\n\n"+
							"Akun Telegram Anda berhasil dihubungkan sebagai <b>Admin Comic Reader TV</b>!\n\n"+
							"✅ Anda sekarang dapat meminta kode OTP kapan saja melalui web dashboard.",
						user.FirstName,
					)
					if isNew {
						m.sendTelegramMessage(chatID, welcomeText)
					}

					// Jika pesan adalah /otp, langsung kirimkan OTP
					if strings.HasPrefix(text, "/otp") {
						_, _ = m.RequestOTP()
					}
				}
			}
		}

		time.Sleep(1 * time.Second)
	}
}
