package pkg

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chainreactors/fingers/resources"
)

type resourceSnapshot struct {
	eholeData             []byte
	gobyData              []byte
	wappalyzerData        []byte
	fingersHTTPData       []byte
	fingerprinthubWebData []byte
}

func snapshotResources() resourceSnapshot {
	return resourceSnapshot{
		eholeData:             append([]byte(nil), resources.EholeData...),
		gobyData:              append([]byte(nil), resources.GobyData...),
		wappalyzerData:        append([]byte(nil), resources.WappalyzerData...),
		fingersHTTPData:       append([]byte(nil), resources.FingersHTTPData...),
		fingerprinthubWebData: append([]byte(nil), resources.FingerprinthubWebData...),
	}
}

func restoreResources(snapshot resourceSnapshot) {
	resources.EholeData = snapshot.eholeData
	resources.GobyData = snapshot.gobyData
	resources.WappalyzerData = snapshot.wappalyzerData
	resources.FingersHTTPData = snapshot.fingersHTTPData
	resources.FingerprinthubWebData = snapshot.fingerprinthubWebData
}

func writeEHoleFaviconFingerprint(t *testing.T, cmsName, hash string) string {
	t.Helper()

	filename := filepath.Join(t.TempDir(), "ehole.json")
	content := fmt.Sprintf(`{"fingerprint":[{"cms":"%s","method":"faviconhash","location":"body","keyword":["%s"]}]}`+"\n", cmsName, hash)
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("write ehole fingerprint: %v", err)
	}
	return filename
}

func writeARLFingerprint(t *testing.T, rule string) string {
	t.Helper()

	filename := filepath.Join(t.TempDir(), "arl.yaml")
	content := fmt.Sprintf("- name: test_rule\n  rule: %q\n", rule)
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("write arl fingerprint: %v", err)
	}
	return filename
}

func newFaviconTestServer(iconData []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(`<html><head><link rel="icon" href="/favicon.ico"></head><body>ok</body></html>`))
		case "/favicon.ico":
			w.Header().Set("Content-Type", "image/x-icon")
			_, _ = w.Write(iconData)
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestScannerScansPlainHTTPURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><title>ok</title></html>"))
	}))
	defer server.Close()

	oldTimeout := Timeout
	Timeout = 1
	defer func() { Timeout = oldTimeout }()

	scanner := &Scanner{queue: NewQueue(), silent: true}
	scanner.queue.Push([]string{server.URL, "0"})
	scanner.scan()

	if len(scanner.allResults) != 1 {
		t.Fatalf("scan() results = %d, want 1", len(scanner.allResults))
	}
}

func TestScannerDoesNotSilentlyDowngradeHTTPSOnFetchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><title>plain-http</title></html>"))
	}))
	defer server.Close()

	oldTimeout := Timeout
	Timeout = 1
	defer func() { Timeout = oldTimeout }()

	httpsURL := strings.Replace(server.URL, "http://", "https://", 1)
	scanner := &Scanner{queue: NewQueue(), silent: true}
	scanner.queue.Push([]string{httpsURL, "0"})
	scanner.scan()

	if len(scanner.allResults) != 0 {
		t.Fatalf("scan() silently downgraded https target and produced results: %+v", scanner.allResults)
	}
}

func TestScannerDetectsFaviconFingerprintWithoutARL(t *testing.T) {
	snapshot := snapshotResources()
	defer restoreResources(snapshot)

	iconData := []byte("test-favicon")
	server := newFaviconTestServer(iconData)
	defer server.Close()

	eholeFile := writeEHoleFaviconFingerprint(t, "TestFaviconCMS", calcFaviconHash(iconData))
	config := &CustomFingerConfig{EHole: eholeFile, NoDefault: true}
	scanner := NewScanner([]string{server.URL}, 1, "", "", 1, true, false, config)
	scanner.Run()

	if len(scanner.hitResults) != 1 {
		t.Fatalf("hitResults = %d, want 1", len(scanner.hitResults))
	}
	if !strings.Contains(strings.ToLower(scanner.hitResults[0].CMS), strings.ToLower("TestFaviconCMS")) {
		t.Fatalf("CMS = %q, want contains %q", scanner.hitResults[0].CMS, "TestFaviconCMS")
	}
}

func TestScannerStillRunsFaviconFingerprintWhenARLEnginePresent(t *testing.T) {
	snapshot := snapshotResources()
	defer restoreResources(snapshot)

	iconData := []byte("test-favicon")
	server := newFaviconTestServer(iconData)
	defer server.Close()

	eholeFile := writeEHoleFaviconFingerprint(t, "TestFaviconCMS", calcFaviconHash(iconData))
	arlFile := writeARLFingerprint(t, `body="definitely-not-here"`)
	config := &CustomFingerConfig{EHole: eholeFile, ARL: arlFile, NoDefault: true}
	scanner := NewScanner([]string{server.URL}, 1, "", "", 1, true, false, config)
	scanner.Run()

	if len(scanner.hitResults) != 1 {
		t.Fatalf("hitResults = %d, want 1", len(scanner.hitResults))
	}
	if !strings.Contains(strings.ToLower(scanner.hitResults[0].CMS), strings.ToLower("TestFaviconCMS")) {
		t.Fatalf("CMS = %q, want contains %q", scanner.hitResults[0].CMS, "TestFaviconCMS")
	}
}
