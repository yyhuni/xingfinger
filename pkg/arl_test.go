package pkg

import (
	"reflect"
	"testing"
)

func TestParseARLConditionsExtractsAndUnescapesAllTypes(t *testing.T) {
	rule := `body="foo\\\"bar" && header="Server: Apache" && title="Admin" && icon_hash="123"`
	got := parseARLConditions(rule)
	want := []ARLCondition{
		{Type: "body", Keyword: `foo\"bar`},
		{Type: "header", Keyword: "Server: Apache"},
		{Type: "title", Keyword: "Admin"},
		{Type: "icon_hash", Keyword: "123"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseARLConditions() = %#v, want %#v", got, want)
	}
}

func TestARLEngineMatchDeduplicatesSuffixVariants(t *testing.T) {
	engine := &ARLEngine{fingerprints: []ARLFingerprint{
		{Name: "ruoyi_body", Rule: `body="ruoyi"`},
		{Name: "ruoyi_header", Rule: `header="ruoyi"`},
		{Name: "other_title", Rule: `title="missing"`},
	}}

	got := engine.Match("Powered by RuoYi", "X-Powered-By: RuoYi", "", "")
	want := []string{"ruoyi"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Match() = %v, want %v", got, want)
	}
}

func TestARLEngineMatchRequiresAllConditions(t *testing.T) {
	engine := &ARLEngine{fingerprints: []ARLFingerprint{{
		Name: "combo_body",
		Rule: `body="hello" && header="world"`,
	}}}

	got := engine.Match("hello only", "", "", "")
	if len(got) != 0 {
		t.Fatalf("Match() = %v, want empty", got)
	}
}
