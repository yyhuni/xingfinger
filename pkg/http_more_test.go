package pkg

import "testing"

func TestExtractFaviconURLKeepsAbsoluteHTTPURL(t *testing.T) {
	body := `<html><head><link rel="icon" href="http://cdn.example.com/favicon.ico"></head></html>`
	got := extractFaviconURL(body, "https://app.example.com/login")
	want := "http://cdn.example.com/favicon.ico"

	if got != want {
		t.Fatalf("extractFaviconURL() = %q, want %q", got, want)
	}
}
