package pkg

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadFromFileAddsHTTPSWhenProtocolMissing(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "targets.txt")
	if err := os.WriteFile(filename, []byte("example.com\n"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got := LoadFromFile(filename)
	want := []string{"https://example.com"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadFromFile() = %v, want %v", got, want)
	}
}

func TestLoadFromFileDoesNotMisclassifyDomainContainingHTTP(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "targets.txt")
	if err := os.WriteFile(filename, []byte("myhttp.site\n"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got := LoadFromFile(filename)
	want := []string{"https://myhttp.site"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadFromFile() = %v, want %v", got, want)
	}
}

func TestLoadFromFileSkipsBlankLinesAndPreservesHTTPURL(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "targets.txt")
	content := "\nhttps://example.com\n\nhttp://test.local\n"
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got := LoadFromFile(filename)
	want := []string{"https://example.com", "http://test.local"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadFromFile() = %v, want %v", got, want)
	}
}
