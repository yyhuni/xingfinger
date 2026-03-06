package pkg

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestNewScannerSeedsQueueAndTimeout(t *testing.T) {
	oldTimeout := Timeout
	defer func() { Timeout = oldTimeout }()

	scanner := NewScanner([]string{"https://a.example", "https://b.example"}, 2, "", "", 7, true, false, nil)

	if scanner == nil {
		t.Fatalf("NewScanner() returned nil")
	}
	if Timeout != 7 {
		t.Fatalf("Timeout = %d, want 7", Timeout)
	}
	if scanner.queue.Len() != 2 {
		t.Fatalf("queue.Len() = %d, want 2", scanner.queue.Len())
	}
	if scanner.engine == nil {
		t.Fatalf("expected default engine to be initialized")
	}
	if scanner.customEngine != nil {
		t.Fatalf("expected customEngine to be nil")
	}
}

func TestNewScannerNoDefaultWithoutCustomFingerprints(t *testing.T) {
	oldTimeout := Timeout
	defer func() { Timeout = oldTimeout }()

	scanner := NewScanner([]string{"https://a.example"}, 1, "", "", 3, true, false, &CustomFingerConfig{NoDefault: true})

	if scanner.engine != nil {
		t.Fatalf("expected default engine to be nil when no-default is set")
	}
	if scanner.customEngine != nil {
		t.Fatalf("expected custom engine to be nil without custom fingerprints")
	}
	if scanner.queue.Len() != 1 {
		t.Fatalf("queue.Len() = %d, want 1", scanner.queue.Len())
	}
}

func TestNewScannerLoadsARLEngine(t *testing.T) {
	arlFile := filepath.Join(t.TempDir(), "arl.yaml")
	content := "- name: demo_body\n  rule: body=\"hello\"\n"
	if err := os.WriteFile(arlFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write arl file: %v", err)
	}

	scanner := NewScanner([]string{"https://a.example"}, 1, "", "", 1, true, false, &CustomFingerConfig{ARL: arlFile, NoDefault: true})

	if scanner.arlEngine == nil {
		t.Fatalf("expected arlEngine to be initialized")
	}
}

func TestNewScannerExitsWhenCustomFingerprintLoadFails(t *testing.T) {
	if os.Getenv("TEST_NEW_SCANNER_BAD_CUSTOM") == "1" {
		NewScanner([]string{"https://a.example"}, 1, "", "", 1, true, false, &CustomFingerConfig{EHole: filepath.Join(os.TempDir(), "missing-ehole.json")})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNewScannerExitsWhenCustomFingerprintLoadFails")
	cmd.Env = append(os.Environ(), "TEST_NEW_SCANNER_BAD_CUSTOM=1")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected subprocess to exit with error")
	}
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.Success() {
		t.Fatalf("expected non-zero exit, got %v", err)
	}
}

func TestNewScannerExitsWhenARLLoadFails(t *testing.T) {
	if os.Getenv("TEST_NEW_SCANNER_BAD_ARL") == "1" {
		NewScanner([]string{"https://a.example"}, 1, "", "", 1, true, false, &CustomFingerConfig{ARL: filepath.Join(os.TempDir(), "missing-arl.yaml"), NoDefault: true})
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNewScannerExitsWhenARLLoadFails")
	cmd.Env = append(os.Environ(), "TEST_NEW_SCANNER_BAD_ARL=1")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected subprocess to exit with error")
	}
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.Success() {
		t.Fatalf("expected non-zero exit, got %v", err)
	}
}
