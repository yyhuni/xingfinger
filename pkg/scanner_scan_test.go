package pkg

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScanSkipsInvalidQueueItemsAndRequestErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<html><head><title>ok</title></head><body>body</body></html>`))
	}))
	defer server.Close()

	scanner := &Scanner{queue: NewQueue(), silent: true}
	scanner.queue.Push([]string{server.URL, "0"})
	scanner.queue.Push([]string{"http://127.0.0.1:1", "0"})
	scanner.queue.Push("not-a-task")
	scanner.scan()

	if len(scanner.allResults) != 1 {
		t.Fatalf("allResults len = %d, want 1", len(scanner.allResults))
	}
	if scanner.allResults[0].URL != server.URL {
		t.Fatalf("result URL = %q, want %q", scanner.allResults[0].URL, server.URL)
	}
}

func TestScanSecondaryTaskDoesNotAttemptFaviconDetection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<html><head><title>secondary</title></head><body><link rel="icon" href="/favicon.ico"></body></html>`))
	}))
	defer server.Close()

	scanner := &Scanner{queue: NewQueue(), silent: true}
	scanner.queue.Push([]string{server.URL, "1"})
	scanner.scan()

	if len(scanner.allResults) != 1 {
		t.Fatalf("allResults len = %d, want 1", len(scanner.allResults))
	}
	if scanner.allResults[0].Title != "secondary" {
		t.Fatalf("result title = %q, want %q", scanner.allResults[0].Title, "secondary")
	}
}
