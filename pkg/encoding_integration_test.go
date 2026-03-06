package pkg

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetectEncodingReturnsEmptyForUnknownCharset(t *testing.T) {
	got := detectEncoding("text/plain")
	want := ""

	if got != want {
		t.Fatalf("detectEncoding() = %q, want %q", got, want)
	}
}

func TestDecodeToUTF8LeavesContentUnchangedWhenEncodingUnknown(t *testing.T) {
	content := string([]byte{'c', 'a', 'f', 0xe9})
	got := decodeToUTF8(content, "text/plain")
	want := content

	if got != want {
		t.Fatalf("decodeToUTF8() = %q, want raw bytes %q", got, want)
	}
}

func TestFetchDecodesTitleUsingContentTypeCharset(t *testing.T) {
	body := []byte("<html><head><title>caf\xe9</title></head><body>body</body></html>")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=windows-1252")
		_, _ = w.Write(body)
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	if resp.Title != "café" {
		t.Fatalf("fetch() title = %q, want %q", resp.Title, "café")
	}
}

func TestFetchUsesMetaCharsetWhenContentTypeUnknown(t *testing.T) {
	body := []byte("<html><head><meta charset=\"windows-1252\"></head><body>caf\xe9</body></html>")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(body)
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	want := `<html><head><meta charset="windows-1252"></head><body>café</body></html>`
	if resp.Body != want {
		t.Fatalf("fetch() body = %q, want %q", resp.Body, want)
	}
}

func TestFetchUsesTitleEncodingWhenContentTypeUnknown(t *testing.T) {
	body := []byte("<html><head><title>caf\xe9</title></head><body>body</body></html>")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(body)
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	if resp.Title != "café" {
		t.Fatalf("fetch() title = %q, want %q", resp.Title, "café")
	}
}

func TestFetchLeavesBodyUnchangedWhenEncodingUnknown(t *testing.T) {
	body := []byte{'c', 'a', 'f', 0xe9}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write(body)
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	want := string(body)
	if resp.Body != want {
		t.Fatalf("fetch() body = %q, want raw bytes %q", resp.Body, want)
	}
}
