package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeEHoleKeywordFingerprint(t *testing.T, cmsName, keyword string) string {
	t.Helper()
	filename := filepath.Join(t.TempDir(), "ehole-keyword.json")
	content := fmt.Sprintf(`{"fingerprint":[{"cms":"%s","method":"keyword","location":"body","keyword":["%s"]}]}`+"\n", cmsName, keyword)
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("write ehole keyword fingerprint: %v", err)
	}
	return filename
}

func TestDetectFingerprintsUsesCustomEngine(t *testing.T) {
	eholeFile := writeEHoleKeywordFingerprint(t, "KeywordCMS", "hello-keyword")
	scanner := NewScanner([]string{"https://example.com"}, 1, "", "", 1, true, false, &CustomFingerConfig{EHole: eholeFile, NoDefault: true})
	raw := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n<html><body>hello-keyword</body></html>")

	got := scanner.detectFingerprints(raw)
	if len(got) == 0 {
		t.Fatalf("detectFingerprints() = empty, want at least one match")
	}
	if !strings.Contains(strings.ToLower(strings.Join(got, ",")), "keywordcms") {
		t.Fatalf("detectFingerprints() = %v, want contains %q", got, "KeywordCMS")
	}
}

func TestDetectFaviconUsesCustomEngine(t *testing.T) {
	iconData := []byte("direct-favicon")
	server := newFaviconTestServer(iconData)
	defer server.Close()

	eholeFile := writeEHoleFaviconFingerprint(t, "DirectFaviconCMS", calcFaviconHash(iconData))
	scanner := NewScanner([]string{server.URL}, 1, "", "", 1, true, false, &CustomFingerConfig{EHole: eholeFile, NoDefault: true})
	body := `<html><head><link rel="icon" href="/favicon.ico"></head></html>`

	got := scanner.detectFavicon(body, server.URL)
	if len(got) == 0 {
		t.Fatalf("detectFavicon() = empty, want at least one match")
	}
	if !strings.Contains(strings.ToLower(strings.Join(got, ",")), "directfaviconcms") {
		t.Fatalf("detectFavicon() = %v, want contains %q", got, "DirectFaviconCMS")
	}
}
