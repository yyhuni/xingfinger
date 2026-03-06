package pkg

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func gzipBody(t *testing.T, body string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write([]byte(body)); err != nil {
		t.Fatalf("gzip write: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}
	return buf.Bytes()
}

func deflateBody(t *testing.T, body string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write([]byte(body)); err != nil {
		t.Fatalf("zlib write: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zlib close: %v", err)
	}
	return buf.Bytes()
}

func TestFetchUsesFinalURLAfterRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/start":
			http.Redirect(w, r, "/final", http.StatusFound)
		case "/final":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = io.WriteString(w, `<html><head><title>final</title></head><body>ok</body></html>`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL + "/start", "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	if resp.URL != server.URL+"/final" {
		t.Fatalf("fetch() URL = %q, want %q", resp.URL, server.URL+"/final")
	}
	if resp.Title != "final" {
		t.Fatalf("fetch() title = %q, want %q", resp.Title, "final")
	}
}

func TestFetchDecodesGzipHTMLResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(gzipBody(t, `<html><head><title>gzip</title></head><body>ok</body></html>`))
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	if resp.Title != "gzip" {
		t.Fatalf("fetch() title = %q, want %q", resp.Title, "gzip")
	}
}

func TestFetchDecodesDeflateHTMLResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Encoding", "deflate")
		_, _ = w.Write(deflateBody(t, `<html><head><title>deflate</title></head><body>ok</body></html>`))
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	if resp.Title != "deflate" {
		t.Fatalf("fetch() title = %q, want %q", resp.Title, "deflate")
	}
}

func TestFetchKeepsJSONBodyWithoutHTMLParsing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = io.WriteString(w, `{"message":"ok","title":"not-html"}`)
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	if resp.Title != "" {
		t.Fatalf("fetch() title = %q, want empty", resp.Title)
	}
	if resp.Body != `{"message":"ok","title":"not-html"}` {
		t.Fatalf("fetch() body = %q, want exact json body", resp.Body)
	}
}

func TestFetchFollowsMultipleRedirectsAndResolvesJSRedirectFromFinalURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/start":
			http.Redirect(w, r, "/middle", http.StatusFound)
		case "/middle":
			http.Redirect(w, r, "/app/login", http.StatusFound)
		case "/app/login":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = io.WriteString(w, `<html><head><title>login</title></head><body><script>window.location.href = "../dashboard"</script></body></html>`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL + "/start", "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	if resp.URL != server.URL+"/app/login" {
		t.Fatalf("fetch() URL = %q, want %q", resp.URL, server.URL+"/app/login")
	}
	if len(resp.JsURLs) != 1 {
		t.Fatalf("fetch() JsURLs len = %d, want 1", len(resp.JsURLs))
	}
	if resp.JsURLs[0] != server.URL+"/dashboard" {
		t.Fatalf("fetch() JsURLs[0] = %q, want %q", resp.JsURLs[0], server.URL+"/dashboard")
	}
}

func TestFetchKeepsHTMLWhenDeflateHeaderIsIncorrect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Encoding", "deflate")
		_, _ = io.WriteString(w, `<html><head><title>plain-html</title></head><body>ok</body></html>`)
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	if resp.Title != "plain-html" {
		t.Fatalf("fetch() title = %q, want %q", resp.Title, "plain-html")
	}
}

func TestFetchLeavesBodyUnchangedWhenContentTypeMissing(t *testing.T) {
	body := []byte{'c', 'a', 'f', 0xe9}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	if resp.Body != string(body) {
		t.Fatalf("fetch() body = %q, want raw bytes %q", resp.Body, string(body))
	}
	if resp.Title != "" {
		t.Fatalf("fetch() title = %q, want empty", resp.Title)
	}
}

func TestFetchHandlesNoContentResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("fetch() status = %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
	if resp.Length != 0 {
		t.Fatalf("fetch() length = %d, want 0", resp.Length)
	}
	if resp.Title != "" {
		t.Fatalf("fetch() title = %q, want empty", resp.Title)
	}
}
