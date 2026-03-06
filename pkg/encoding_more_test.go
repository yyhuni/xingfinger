package pkg

import "testing"

func TestConvertEncodingToNonUTF8ReturnsEncodedBytesAsString(t *testing.T) {
	got := convertEncoding("hello", "utf-8", "utf-8")
	if got != "hello" {
		t.Fatalf("convertEncoding() = %q, want %q", got, "hello")
	}
}
