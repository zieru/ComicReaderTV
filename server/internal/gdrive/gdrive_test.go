package gdrive

import "testing"

func TestExtractFileID(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
		wantErr  bool
	}{
		{
			name:     "preview link",
			url:      "https://drive.google.com/file/d/1XyZ_12345abcdef/preview",
			expected: "1XyZ_12345abcdef",
			wantErr:  false,
		},
		{
			name:     "view with query params",
			url:      "https://drive.google.com/file/d/2AbCdEfGhIjKlMnOp/view?usp=sharing",
			expected: "2AbCdEfGhIjKlMnOp",
			wantErr:  false,
		},
		{
			name:     "open id link",
			url:      "https://drive.google.com/open?id=3QrStUvWxYz",
			expected: "3QrStUvWxYz",
			wantErr:  false,
		},
		{
			name:     "uc id link",
			url:      "https://drive.google.com/uc?id=4ExampleID&export=download",
			expected: "4ExampleID",
			wantErr:  false,
		},
		{
			name:    "invalid url",
			url:     "https://example.com/not-gdrive",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ExtractFileID(tt.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ExtractFileID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if id != tt.expected {
				t.Errorf("ExtractFileID() = %v, want %v", id, tt.expected)
			}
		})
	}
}

func TestDownloadFile_Restricted(t *testing.T) {
	client := NewClient()
	// ID file yang di-share tanpa izin download publik
	_, err := client.DownloadFile("1Nu-jwvwQIXiggCWR9WG0GaGxTIJ9YKwO")
	if err == nil {
		t.Fatal("harus mengembalikan error untuk file restricted")
	}
	t.Logf("Pesan error yang dihasilkan: %v", err)
}
