package pkg

import "testing"

func TestDetectEncodingRecognizesBig5(t *testing.T) {
	if got := detectEncoding("text/html; charset=big5"); got != "big5" {
		t.Fatalf("detectEncoding() = %q, want %q", got, "big5")
	}
}

func TestConvertEncodingReturnsSourceForUnknownDecoder(t *testing.T) {
	src := "plain-text"
	if got := convertEncoding(src, "x-unknown-charset", "utf-8"); got != src {
		t.Fatalf("convertEncoding() = %q, want %q", got, src)
	}
}
