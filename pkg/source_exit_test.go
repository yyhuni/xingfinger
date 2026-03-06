package pkg

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestLoadFromFileExitsWhenFileMissing(t *testing.T) {
	if os.Getenv("TEST_LOAD_FROM_FILE_MISSING") == "1" {
		LoadFromFile(filepath.Join(os.TempDir(), "definitely-missing-xingfinger-targets.txt"))
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLoadFromFileExitsWhenFileMissing")
	cmd.Env = append(os.Environ(), "TEST_LOAD_FROM_FILE_MISSING=1")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected subprocess to exit with error")
	}
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.Success() {
		t.Fatalf("expected non-zero exit, got %v", err)
	}
}
