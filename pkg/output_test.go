package pkg

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveResultsWritesJSONArrayToJSONFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "result.json")
	results := []Result{{URL: "https://example.com", CMS: "nginx", StatusCode: 200}}

	saveResults(filename, results)

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	var got []Result
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(got) != 1 || got[0].URL != results[0].URL || got[0].CMS != results[0].CMS || got[0].StatusCode != results[0].StatusCode {
		t.Fatalf("saveResults() wrote %+v, want %+v", got, results)
	}
}

func TestSaveResultsSkipsUnsupportedExtension(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "result.txt")
	results := []Result{{URL: "https://example.com"}}

	saveResults(filename, results)

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		t.Fatalf("expected %s to not be created, stat err = %v", filename, err)
	}
}
