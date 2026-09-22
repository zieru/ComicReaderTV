package auth

import (
	"os"
	"testing"
)

func TestOTPFlow(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "auth_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mgr := NewManager("", tempDir, 123456789)
	defer mgr.Close()

	// 1. Request OTP
	_, err = mgr.RequestOTP()
	if err != nil {
		t.Fatalf("RequestOTP failed: %v", err)
	}

	mgr.mu.RLock()
	otpCode := mgr.currentOTP
	mgr.mu.RUnlock()

	if len(otpCode) != 6 {
		t.Errorf("Expected 6-digit OTP, got: %s", otpCode)
	}

	// 2. Verify with wrong OTP
	_, err = mgr.VerifyOTP("000000")
	if err == nil {
		t.Errorf("Expected error for wrong OTP, got nil")
	}

	// 3. Verify with correct OTP
	sessionToken, err := mgr.VerifyOTP(otpCode)
	if err != nil {
		t.Fatalf("VerifyOTP failed with valid code: %v", err)
	}

	if len(sessionToken) != 64 {
		t.Errorf("Expected 64-char hex session token, got: %s", sessionToken)
	}

	// 4. Validate session
	if !mgr.ValidateSession(sessionToken) {
		t.Errorf("Session token should be valid")
	}

	// 5. Test reuse of same OTP (should fail)
	_, err = mgr.VerifyOTP(otpCode)
	if err == nil {
		t.Errorf("Expected error on OTP reuse, got nil")
	}

	// 6. Test Revoke Session
	mgr.RevokeSession(sessionToken)
	if mgr.ValidateSession(sessionToken) {
		t.Errorf("Session should be invalid after revocation")
	}
}

func TestRateLimit(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "auth_test_*")
	defer os.RemoveAll(tempDir)

	mgr := NewManager("", tempDir, 123456789)
	defer mgr.Close()

	_, err := mgr.RequestOTP()
	if err != nil {
		t.Fatalf("First RequestOTP failed: %v", err)
	}

	// Permintaan kedua langsung harus kena rate limit
	_, err = mgr.RequestOTP()
	if err == nil {
		t.Errorf("Expected rate limit error on immediate second request")
	}
}

func TestPersistence(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "auth_test_*")
	defer os.RemoveAll(tempDir)

	mgr1 := NewManager("", tempDir, 0)
	mgr1.SetAdminChatID(987654321)
	mgr1.Close()

	// Buat manager baru dengan folder data yang sama
	mgr2 := NewManager("", tempDir, 0)
	defer mgr2.Close()

	hasAdmin, _, _ := mgr2.GetStatus()
	if !hasAdmin || mgr2.adminChatID != 987654321 {
		t.Errorf("Expected adminChatID 987654321 to be persisted, got %d", mgr2.adminChatID)
	}
}
