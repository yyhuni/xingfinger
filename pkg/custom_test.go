package pkg

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/chainreactors/fingers/resources"
)

func TestLoadFingerFileCompressesJSONInput(t *testing.T) {
	content := []byte(`{"fingerprint":[]}`)
	filename := filepath.Join(t.TempDir(), "finger.json")
	if err := os.WriteFile(filename, content, 0o644); err != nil {
		t.Fatalf("write json: %v", err)
	}

	got, err := loadFingerFile(filename)
	if err != nil {
		t.Fatalf("loadFingerFile() error = %v", err)
	}

	decoded, err := gzipDecompress(got)
	if err != nil {
		t.Fatalf("gzipDecompress() error = %v", err)
	}
	if !bytes.Equal(decoded, content) {
		t.Fatalf("decoded gzip = %q, want %q", decoded, content)
	}
}

func TestLoadFingerFileReturnsGzipInputAsIs(t *testing.T) {
	content := []byte(`{"fingerprint":[]}`)
	compressed, err := gzipCompress(content)
	if err != nil {
		t.Fatalf("gzipCompress() error = %v", err)
	}

	filename := filepath.Join(t.TempDir(), "finger.json.gz")
	if err := os.WriteFile(filename, compressed, 0o644); err != nil {
		t.Fatalf("write gzip: %v", err)
	}

	got, err := loadFingerFile(filename)
	if err != nil {
		t.Fatalf("loadFingerFile() error = %v", err)
	}
	if !bytes.Equal(got, compressed) {
		t.Fatalf("loadFingerFile() = %q, want raw gzip bytes", got)
	}
}

func TestLoadCustomFingerprintsNoDefaultClearsBuiltinResources(t *testing.T) {
	snapshot := snapshotResources()
	defer restoreResources(snapshot)

	content := []byte(`{"fingerprint":[]}`)
	filename := filepath.Join(t.TempDir(), "ehole.json")
	if err := os.WriteFile(filename, content, 0o644); err != nil {
		t.Fatalf("write ehole file: %v", err)
	}

	err := LoadCustomFingerprints(&CustomFingerConfig{EHole: filename, NoDefault: true}, true)
	if err != nil {
		t.Fatalf("LoadCustomFingerprints() error = %v", err)
	}

	decoded, err := gzipDecompress(resources.EholeData)
	if err != nil {
		t.Fatalf("gzipDecompress() error = %v", err)
	}
	if !bytes.Equal(decoded, content) {
		t.Fatalf("resources.EholeData decoded = %q, want %q", decoded, content)
	}
	if len(resources.GobyData) != 0 || len(resources.WappalyzerData) != 0 || len(resources.FingersHTTPData) != 0 || len(resources.FingerprinthubWebData) != 0 {
		t.Fatalf("expected non-EHole builtin resources to be cleared")
	}
}
