package gscraper

import (
	"testing"
)

func TestHasBrowser(t *testing.T) {
	has := HasBrowser()
	t.Logf("Chromium/Chrome installed on system: %v (path: %s)", has, BrowserPath())
}
