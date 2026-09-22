package model

import "time"

// Comic merepresentasikan sebuah komik / manga dalam katalog
type Comic struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CoverURL    string    `json:"cover_url"`
	SourceType  string    `json:"source_type"` // "gdrive_pdf", "direct_pdf"
	SourceURL   string    `json:"source_url"`  // URL Google Drive embed/preview atau direct link
	TotalPages  int       `json:"total_pages"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateComicRequest payload saat menambahkan komik baru
type CreateComicRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	SourceURL   string `json:"source_url"`
}
