package pkg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveJSONWritesIndentedOutput(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "pretty.json")
	saveJSON(filename, []Result{{URL: "https://example.com"}})

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), "\n  {") {
		t.Fatalf("saveJSON() output is not indented: %q", data)
	}
}
