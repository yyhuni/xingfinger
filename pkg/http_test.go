package pkg

import "testing"

func TestExtractFaviconURLPreservesHTTPSForProtocolRelativePath(t *testing.T) {
	body := `<html><head><link rel="icon" href="//cdn.example.com/favicon.ico"></head></html>`
	got := extractFaviconURL(body, "https://app.example.com/login")
	want := "https://cdn.example.com/favicon.ico"

	if got != want {
		t.Fatalf("extractFaviconURL() = %q, want %q", got, want)
	}
}

func TestExtractFaviconURLUsesDefaultWhenTagMissing(t *testing.T) {
	body := `<html><head></head><body>no icon</body></html>`
	got := extractFaviconURL(body, "https://app.example.com/login")
	want := "https://app.example.com/favicon.ico"

	if got != want {
		t.Fatalf("extractFaviconURL() = %q, want %q", got, want)
	}
}

func TestExtractFaviconURLResolvesRelativePath(t *testing.T) {
	body := `<html><head><link rel="icon" href="assets/favicon.ico"></head></html>`
	got := extractFaviconURL(body, "https://app.example.com/app/login")
	want := "https://app.example.com/app/assets/favicon.ico"

	if got != want {
		t.Fatalf("extractFaviconURL() = %q, want %q", got, want)
	}
}
