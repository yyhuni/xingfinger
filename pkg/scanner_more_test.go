package pkg

import (
	"os"
	"strings"
	"testing"
)

func TestDetectFingerprintsReturnsEmptyWhenNoEngines(t *testing.T) {
	scanner := &Scanner{}
	got := scanner.detectFingerprints([]byte("HTTP/1.1 200 OK\r\n\r\n<body>hello</body>"))
	if len(got) != 0 {
		t.Fatalf("detectFingerprints() = %v, want empty", got)
	}
}

func TestDetectFaviconReturnsNilWhenNoFaviconFound(t *testing.T) {
	scanner := &Scanner{}
	got := scanner.detectFavicon("<html><head></head><body>no icon</body></html>", "://bad-url")
	if got != nil {
		t.Fatalf("detectFavicon() = %v, want nil", got)
	}
}

func TestGetFaviconHashReturnsEmptyWhenBaseURLInvalid(t *testing.T) {
	scanner := &Scanner{}
	if got := scanner.getFaviconHash("<html></html>", "://bad-url"); got != "" {
		t.Fatalf("getFaviconHash() = %q, want empty", got)
	}
}

func TestRunSavesResultsToOutputFile(t *testing.T) {
	iconData := []byte("run-favicon")
	server := newFaviconTestServer(iconData)
	defer server.Close()

	eholeFile := writeEHoleFaviconFingerprint(t, "RunCMS", calcFaviconHash(iconData))
	outputFile := t.TempDir() + "/result.json"
	config := &CustomFingerConfig{EHole: eholeFile, NoDefault: true}
	scanner := NewScanner([]string{server.URL}, 1, outputFile, "", 1, true, false, config)
	scanner.Run()

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), server.URL) {
		t.Fatalf("output file missing URL: %q", data)
	}
	if !strings.Contains(strings.ToLower(string(data)), "runcms") {
		t.Fatalf("output file missing CMS: %q", data)
	}
}
