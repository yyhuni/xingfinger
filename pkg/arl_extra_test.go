package pkg

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewARLEngineLoadsYAMLFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "arl.yaml")
	content := "- name: test_body\n  rule: body=\"hello\"\n"
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		t.Fatalf("write arl file: %v", err)
	}

	engine, err := NewARLEngine(filename)
	if err != nil {
		t.Fatalf("NewARLEngine() error = %v", err)
	}
	if len(engine.fingerprints) != 1 {
		t.Fatalf("fingerprints len = %d, want 1", len(engine.fingerprints))
	}
	if engine.fingerprints[0].Name != "test_body" {
		t.Fatalf("fingerprint name = %q, want %q", engine.fingerprints[0].Name, "test_body")
	}
}

func TestNewARLEngineReturnsErrorForInvalidYAML(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "arl.yaml")
	if err := os.WriteFile(filename, []byte("- name: [invalid\n"), 0o644); err != nil {
		t.Fatalf("write invalid arl file: %v", err)
	}

	if _, err := NewARLEngine(filename); err == nil {
		t.Fatalf("NewARLEngine() error = nil, want non-nil")
	}
}

func TestMatchConditionSupportsTitleAndIconHash(t *testing.T) {
	if !matchCondition(ARLCondition{Type: "title", Keyword: "Admin"}, "", "", "Super Admin", "") {
		t.Fatalf("title condition should match")
	}
	if !matchCondition(ARLCondition{Type: "icon_hash", Keyword: "123"}, "", "", "", "123") {
		t.Fatalf("icon_hash condition should match")
	}
	if matchCondition(ARLCondition{Type: "unknown", Keyword: "x"}, "", "", "", "") {
		t.Fatalf("unknown condition should not match")
	}
}

func TestExtractARLNameHandlesKnownAndUnknownSuffixes(t *testing.T) {
	cases := map[string]string{
		"demo_body":      "demo",
		"demo_header":    "demo",
		"demo_title":     "demo",
		"demo_icon_hash": "demo",
		"demo_custom":    "demo_custom",
	}

	for input, want := range cases {
		if got := extractARLName(input); got != want {
			t.Fatalf("extractARLName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestMatchReturnsEmptyWhenNoRuleMatches(t *testing.T) {
	engine := &ARLEngine{fingerprints: []ARLFingerprint{{Name: "demo_body", Rule: `body="abc"`}}}
	got := engine.Match("xyz", "", "", "")
	if !reflect.DeepEqual(got, []string(nil)) {
		t.Fatalf("Match() = %v, want nil/empty", got)
	}
}

func TestNewARLEngineReturnsErrorForMissingFile(t *testing.T) {
	if _, err := NewARLEngine(filepath.Join(t.TempDir(), "missing.yaml")); err == nil {
		t.Fatalf("NewARLEngine() error = nil, want non-nil")
	}
}
