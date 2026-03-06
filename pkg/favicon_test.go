package pkg

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchFaviconReturnsErrorOnNon200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	_, err := fetchFavicon(server.URL, "")
	if err == nil {
		t.Fatalf("fetchFavicon() error = nil, want non-nil")
	}
}

func TestGetFaviconHashMatchesCalcFaviconHash(t *testing.T) {
	iconData := []byte("favicon-bytes")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(`<html><head><link rel="icon" href="/favicon.ico"></head></html>`))
		case "/favicon.ico":
			w.Header().Set("Content-Type", "image/x-icon")
			_, _ = w.Write(iconData)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	scanner := &Scanner{}
	got := scanner.getFaviconHash(`<html><head><link rel="icon" href="/favicon.ico"></head></html>`, server.URL)
	want := calcFaviconHash(iconData)

	if got != want {
		t.Fatalf("getFaviconHash() = %q, want %q", got, want)
	}
}

func TestFetchFaviconReturnsContentOn200(t *testing.T) {
	iconData := []byte("ok-favicon")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		_, _ = w.Write(iconData)
	}))
	defer server.Close()

	got, err := fetchFavicon(server.URL, "")
	if err != nil {
		t.Fatalf("fetchFavicon() error = %v", err)
	}
	if string(got) != string(iconData) {
		t.Fatalf("fetchFavicon() = %q, want %q", got, iconData)
	}
}
