package pkg

import "testing"

func TestDetectEncodingRecognizesUTF8(t *testing.T) {
	got := detectEncoding("text/html; charset=utf-8")
	want := "utf-8"

	if got != want {
		t.Fatalf("detectEncoding() = %q, want %q", got, want)
	}
}

func TestDetectEncodingPreservesWindows1252(t *testing.T) {
	got := detectEncoding("text/html; charset=windows-1252")
	want := "windows-1252"

	if got != want {
		t.Fatalf("detectEncoding() = %q, want %q", got, want)
	}
}

func TestDecodeToUTF8UsesContentTypeCharset(t *testing.T) {
	content := string([]byte{'c', 'a', 'f', 0xe9})
	got := decodeToUTF8(content, "text/html; charset=windows-1252")
	want := "café"

	if got != want {
		t.Fatalf("decodeToUTF8() = %q, want %q", got, want)
	}
}

func TestDecodeToUTF8MetaCharsetOverridesContentType(t *testing.T) {
	content := string([]byte(
		`<html><head><meta charset="windows-1252"></head><body>caf` + string([]byte{0xe9}) + `</body></html>`,
	))
	got := decodeToUTF8(content, "text/html; charset=utf-8")
	want := `<html><head><meta charset="windows-1252"></head><body>café</body></html>`

	if got != want {
		t.Fatalf("decodeToUTF8() = %q, want %q", got, want)
	}
}

func TestDecodeToUTF8DetectsTitleEncodingWhenHeaderClaimsUTF8(t *testing.T) {
	content := string([]byte(
		`<html><head><title>caf` + string([]byte{0xe9}) + `</title></head><body>body</body></html>`,
	))
	got := decodeToUTF8(content, "text/html; charset=utf-8")
	want := `<html><head><title>café</title></head><body>body</body></html>`

	if got != want {
		t.Fatalf("decodeToUTF8() = %q, want %q", got, want)
	}
}
