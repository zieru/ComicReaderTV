package manga

import (
	"testing"
)

func TestSearchEmpty(t *testing.T) {
	items, err := Search("")
	if err != nil {
		t.Fatalf("Unexpected error for empty query: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(items))
	}
}

func TestSearchManga(t *testing.T) {
	items, err := Search("Detective Conan")
	if err != nil {
		t.Skipf("Network/API unavailable, skipping live test: %v", err)
		return
	}
	if len(items) == 0 {
		t.Errorf("Expected results for 'Detective Conan', got 0")
	}
	first := items[0]
	if first.Title == "" {
		t.Errorf("Expected non-empty title")
	}
}
