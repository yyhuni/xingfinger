package pkg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chainreactors/fingers/resources"
)

func TestLoadCustomFingerprintsReturnsErrorWhenFileMissing(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.json")
	if err := LoadCustomFingerprints(&CustomFingerConfig{EHole: missing}, true); err == nil {
		t.Fatalf("LoadCustomFingerprints() error = nil, want non-nil")
	}
}

func TestLoadCustomFingerprintsLoadsWappalyzerAndFingers(t *testing.T) {
	snapshot := snapshotResources()
	defer restoreResources(snapshot)

	content := []byte(`{"fingerprint":[]}`)
	wappalyzerFile := filepath.Join(t.TempDir(), "wappalyzer.json")
	fingersFile := filepath.Join(t.TempDir(), "fingers.json")
	if err := os.WriteFile(wappalyzerFile, content, 0o644); err != nil {
		t.Fatalf("write wappalyzer file: %v", err)
	}
	if err := os.WriteFile(fingersFile, content, 0o644); err != nil {
		t.Fatalf("write fingers file: %v", err)
	}

	err := LoadCustomFingerprints(&CustomFingerConfig{Wappalyzer: wappalyzerFile, Fingers: fingersFile}, true)
	if err != nil {
		t.Fatalf("LoadCustomFingerprints() error = %v", err)
	}
	if len(resources.WappalyzerData) == 0 {
		t.Fatalf("expected WappalyzerData to be populated")
	}
	if len(resources.FingersHTTPData) == 0 {
		t.Fatalf("expected FingersHTTPData to be populated")
	}
}
