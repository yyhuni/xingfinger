package pkg

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/chainreactors/fingers/resources"
)

func TestLoadFingerFileReturnsRawBytesForUnknownExtension(t *testing.T) {
	content := []byte("raw-data")
	filename := filepath.Join(t.TempDir(), "finger.bin")
	if err := os.WriteFile(filename, content, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	got, err := loadFingerFile(filename)
	if err != nil {
		t.Fatalf("loadFingerFile() error = %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("loadFingerFile() = %q, want %q", got, content)
	}
}

func TestGzipDecompressReturnsErrorForInvalidData(t *testing.T) {
	if _, err := gzipDecompress([]byte("not-gzip")); err == nil {
		t.Fatalf("gzipDecompress() error = nil, want non-nil")
	}
}

func TestLoadCustomFingerprintsLoadsAdditionalBuiltinSources(t *testing.T) {
	snapshot := snapshotResources()
	defer restoreResources(snapshot)

	content := []byte(`{"fingerprint":[]}`)
	gobyFile := filepath.Join(t.TempDir(), "goby.json")
	fingerprintFile := filepath.Join(t.TempDir(), "fingerprinthub.json")
	if err := os.WriteFile(gobyFile, content, 0o644); err != nil {
		t.Fatalf("write goby file: %v", err)
	}
	if err := os.WriteFile(fingerprintFile, content, 0o644); err != nil {
		t.Fatalf("write fingerprinthub file: %v", err)
	}

	err := LoadCustomFingerprints(&CustomFingerConfig{Goby: gobyFile, FingerPrint: fingerprintFile}, true)
	if err != nil {
		t.Fatalf("LoadCustomFingerprints() error = %v", err)
	}

	if len(resources.GobyData) == 0 {
		t.Fatalf("expected GobyData to be populated")
	}
	if len(resources.FingerprinthubWebData) == 0 {
		t.Fatalf("expected FingerprinthubWebData to be populated")
	}
}
