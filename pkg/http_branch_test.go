package pkg

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchUsesXPoweredByWhenServerHeaderMissing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Powered-By", "PHP/8.0")
		_, _ = io.WriteString(w, `<html><body>ok</body></html>`)
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "0"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}
	if resp.Server != "PHP/8.0" {
		t.Fatalf("fetch() server = %q, want %q", resp.Server, "PHP/8.0")
	}
}

func TestFetchDoesNotParseJSRedirectForSecondaryTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, `<html><body><script>window.location.href = "/next"</script></body></html>`)
	}))
	defer server.Close()

	resp, err := fetch([]string{server.URL, "1"}, "")
	if err != nil {
		t.Fatalf("fetch() error = %v", err)
	}
	if len(resp.JsURLs) != 0 {
		t.Fatalf("fetch() JsURLs = %v, want empty", resp.JsURLs)
	}
}
