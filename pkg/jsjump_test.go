package pkg

import (
	"reflect"
	"testing"
)

func TestParseJSRedirectKeepsAbsoluteRedirectURL(t *testing.T) {
	body := `<script>window.location.href = "https://auth.example.com/next"</script>`
	got := parseJSRedirect(body, "https://app.example.com/login")
	want := []string{"https://auth.example.com/next"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseJSRedirect() = %v, want %v", got, want)
	}
}

func TestParseJSRedirectResolvesParentRelativePath(t *testing.T) {
	body := `<script>window.location.href = "../dashboard"</script>`
	got := parseJSRedirect(body, "https://example.com/app/login")
	want := []string{"https://example.com/dashboard"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseJSRedirect() = %v, want %v", got, want)
	}
}

func TestParseJSRedirectSupportsMetaRefresh(t *testing.T) {
	body := `<meta http-equiv="refresh" content="0;url=/next">`
	got := parseJSRedirect(body, "https://example.com/app/login")
	want := []string{"https://example.com/next"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseJSRedirect() = %v, want %v", got, want)
	}
}

func TestParseJSRedirectDeduplicatesTargets(t *testing.T) {
	body := `<script>window.location.href = "/next"</script><script>window.location.href = "/next"</script>`
	got := parseJSRedirect(body, "https://example.com/app/login")
	want := []string{"https://example.com/next"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseJSRedirect() = %v, want %v", got, want)
	}
}

func TestParseJSRedirectSupportsRedirectURLAssignment(t *testing.T) {
	body := `<script>redirectUrl = "/assigned"</script>`
	got := parseJSRedirect(body, "https://example.com/app/login")
	want := []string{"https://example.com/assigned"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseJSRedirect() = %v, want %v", got, want)
	}
}
